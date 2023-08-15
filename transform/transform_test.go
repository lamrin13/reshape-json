package transform

import (
	"encoding/json"
	"log"
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
	entries, err := os.ReadDir("./testdata")
	if err != nil {
		log.Fatal(err)
	}
	for _, testDir := range entries {

		testname := testDir.Name()

		// Each path turns into a test: the test name is the folder name
		t.Run(testname, func(t *testing.T) {
			path := filepath.Join("testdata", testname, testname+".mapping")
			source, err := os.ReadFile(path)
			if err != nil {
				t.Fatal("error reading config file:", err)
			}

			goldenfile := filepath.Join("testdata", testname, testname+".golden")
			golden, err := os.ReadFile(goldenfile)
			if err != nil {
				t.Fatal("error reading golden file:", err)
			}

			originalfile := filepath.Join("testdata", testname, testname+".input")
			original, err := os.ReadFile(originalfile)
			if err != nil {
				t.Fatal("error reading original file:", err)
			}
			var mappingConfig map[string]string
			json.Unmarshal(source, &mappingConfig)

			m := MappingConfig(mappingConfig)
			var expected, returned any
			json.Unmarshal(golden, &expected)

			returnedJSON, err := m.Reshape(original)
			if err != nil {
				t.Errorf(err.Error())
			}
			json.Unmarshal(returnedJSON, &returned)
			log.Println(returned, expected)
			if !reflect.DeepEqual(expected, returned) {
				t.Errorf("Expected %s\n, Got %s", golden, returnedJSON)
			}

		})
	}

}
