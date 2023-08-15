package main

import (
	"fmt"

	"github.com/lamrin13/reshape-json/transform"
)

func main() {
	original := `{
		"name": "John",
		"street": "Dongg Street",
		"unit": "12",
		"role": "SWE",
		"last4": "1234",
		"expiryMonth": 12,
		"expiryYear": 2023,
		"fruits": [
			{
				"t1": "John",
				"t2": "john.a@abc.com"
			},
			{
				"t1": "Bob",
				"t2": "bob.b@abc.com"
			}
		]
	}`

	mappingConfig := transform.MappingConfig{
		"user.name":          "name",
		"user.role":          "role",
		"card.lastFour":      "last4",
		"card.expiry":        "expiryMonth+expiryYear/-",
		"address":            "unit+street",
		"user.favoriteFoods": "fruits",
		"o1":                 "t1",
		"o2":                 "t2",
	}

	desiredJSON, err := mappingConfig.Reshape([]byte(original))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(desiredJSON))
}
