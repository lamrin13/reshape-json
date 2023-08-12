package main

import (
	"encoding/json"
	"fmt"

	"github.com/lamrin13/reshape-json/transform"
)

func main() {
	original := `[
		{
			"name": "John",
			"email": "john.a@abc.com"
		},
		{
			"name": "Bob",
			"email": "bob.b@abc.com"
		}
	]`

	mappingConfig := transform.MappingConfig{
		"usr.userName":     "name",
		"usr.emailAddress": "email",
	}

	desiredJSON, err := mappingConfig.Reshape([]byte(original))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(json.Valid(desiredJSON))
	fmt.Println(string(desiredJSON))
}
