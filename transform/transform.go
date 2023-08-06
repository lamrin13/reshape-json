package transform

import (
	"encoding/json"
	"strings"
)

type MappingConfig map[string]string

// transformJSON transforms the input JSON data based on the provided mapping configuration.
func (mappingConfig MappingConfig) Reshape(input []byte) ([]byte, error) {
	var original map[string]interface{}
	err := json.Unmarshal(input, &original)
	if err != nil {
		return nil, err
	}

	desired := make(map[string]interface{})
	for key, value := range mappingConfig {
		var formattedValue any
		if strings.Contains(value, "+") {
			subKeys := strings.Split(value, "+")
			var stringVals []string
			for _, v := range subKeys {
				stringVals = append(stringVals, original[v].(string))
			}
			formattedValue = strings.Join(stringVals, "-")
		} else {
			formattedValue = original[value]
		}
		if strings.Contains(key, ".") {
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
		} else {
			desired[key] = make(map[string]interface{})
			desired[key] = formattedValue
		}
	}
	return json.Marshal(desired)
}
