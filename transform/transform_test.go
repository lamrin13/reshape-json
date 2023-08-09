package transform

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func BenchmarkReshape(b *testing.B) {
	m := MappingConfig{
		"user.firstname":   "name",
		"user.email":       "email",
		"user.phoneNumber": "phone",
	}

	original := `{
		"name": "Nirmal",
		"email": "nirmal.patel@abc.com",
		"phone": "6478668421"
	}`

	for i := 0; i < b.N; i++ {
		m.Reshape([]byte(original))
	}
}

func TestReshape(t *testing.T) {
	paths, err := filepath.Glob(filepath.Join("testdata", "*.mapping"))
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range paths {
		_, filename := filepath.Split(path)
		testname := filename[:len(filename)-len(filepath.Ext(path))]

		// Each path turns into a test: the test name is the filename without the
		// extension.
		t.Run(testname, func(t *testing.T) {
			source, err := os.ReadFile(path)
			if err != nil {
				t.Fatal("error reading config file:", err)
			}

			goldenfile := filepath.Join("testdata", testname+".golden")
			golden, err := os.ReadFile(goldenfile)
			if err != nil {
				t.Fatal("error reading golden file:", err)
			}

			originalfile := filepath.Join("testdata", testname+".input")
			original, err := os.ReadFile(originalfile)
			if err != nil {
				t.Fatal("error reading original file:", err)
			}
			var mappingConfig map[string]string
			json.Unmarshal(source, &mappingConfig)

			m := MappingConfig(mappingConfig)
			var expected, returned map[string]interface{}
			json.Unmarshal(golden, &expected)

			returnedJSON, err := m.Reshape(original)
			if err != nil {
				t.Errorf(err.Error())
			}
			json.Unmarshal(returnedJSON, &returned)
			if !reflect.DeepEqual(expected, returned) {
				t.Errorf("Expected %s\n, Got %s", golden, returnedJSON)
			}

		})
	}

}
