package transform

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type MappingConfig map[string]string

// transformJSON transforms the input JSON data based on the provided mapping configuration.
func (mappingConfig MappingConfig) Reshape(input []byte) ([]byte, error) {
	var original any
	err := json.Unmarshal(input, &original)
	if err != nil {
		return nil, err
	}

	switch v := original.(type) {
	case []interface{}:
		return reshapeArray(v, mappingConfig)
	case map[string]interface{}:
		return reshapeObject(v, mappingConfig)
	default:
		return nil, errors.New("Invalid JSON input")
	}
}

func getSeparator(value string) (string, error) {
	elems := strings.Split(value, "/")
	switch l := len(elems); l {
	case 1:
		return " ", nil
	case 2:
		return elems[1], nil
	default:
		return "", errors.New("multiple / found in config for " + value)
	}
}

func getStringVal(value interface{}) (string, error) {
	switch i := value.(type) {
	case float64:
		return strconv.FormatFloat(i, 'f', -1, 64), nil
	case int:
		return strconv.Itoa(i), nil
	case string:
		return i, nil
	default:
		return "", errors.New("Only allowed string, float and ints for merging fields")
	}
}

func reshapeObject(original map[string]interface{}, mappingConfig MappingConfig) ([]byte, error) {
	desired := make(map[string]interface{})
	var flag bool
	for key, value := range mappingConfig {
		var (
			formattedValue any
			err            error
			separator      string
		)
		flag = false
		if strings.Contains(value, "+") {
			separator, err = getSeparator(value)
			if err != nil {
				return nil, err
			}
			subKeys := strings.Split(strings.Split(value, "/")[0], "+")
			var stringVals []string
			for _, v := range subKeys {
				if _, ok := original[v]; !ok {
					flag = true
					break
				}
				strValue, err := getStringVal(original[v])
				if err != nil {
					return nil, err
				}
				stringVals = append(stringVals, strValue)
			}
			if !flag {
				formattedValue = strings.Join(stringVals, separator)
			}
		} else {
			if _, ok := original[value]; !ok {
				continue
			}
			if reflect.TypeOf(original[value]).String() == "[]interface {}" {
				arrObject := original[value].([]interface{})
				if reflect.TypeOf(arrObject[0]).String() != "map[string]interface {}" {
					formattedValue = original[value]
				} else {
					formattedBytes, err := reshapeArray(arrObject, mappingConfig)
					if err != nil {
						return nil, err
					}
					json.Unmarshal(formattedBytes, &formattedValue)
				}
			} else {
				formattedValue = original[value]
			}
		}
		if !flag && strings.Contains(key, ".") {
			nests := strings.Split(key, ".")
			if _, ok := desired[nests[0]]; !ok {
				desired[nests[0]] = make(map[string]interface{})
			}
			prev := desired
			for i := 1; i < len(nests); i++ {
				temp := prev[nests[i-1]].(map[string]interface{})
				if _, ok := temp[nests[i]]; !ok {
					temp[nests[i]] = make(map[string]interface{})
				}
				prev = temp
			}
			prev[nests[len(nests)-1]] = formattedValue
		} else if !flag {
			desired[key] = make(map[string]interface{})
			desired[key] = formattedValue
		}
	}
	return json.Marshal(desired)
}

func reshapeArray(original []interface{}, mappingConfig MappingConfig) ([]byte, error) {
	desired := make([][]byte, len(original))
	for i, v := range original {
		obj := v.(map[string]interface{})
		tmp, err := reshapeObject(obj, mappingConfig)
		desired[i] = make([]byte, len(tmp))
		if err != nil {
			return nil, err
		}
		desired[i] = tmp
	}
	var reshapedArray []byte
	reshapedArray = append(reshapedArray, byte('['))
	reshapedArray = append(reshapedArray, bytes.Join(desired, []byte(","))...)
	reshapedArray = append(reshapedArray, byte(']'))
	return reshapedArray, nil
}
