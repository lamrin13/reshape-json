package main

import (
	"fmt"

	"github.com/lamrin13/reshape-json/transform"
)

func main() {
	original := `{
		"name": "Nirmal",
		"street": "ssss",
		"unit": "12",
		"role": "SWE",
		"last4": "1234",
		"expiryMonth": 12,
		"expiryYear": 2023
	}`

	mappingConfig := transform.MappingConfig{
		"user.name":         "name",
		"user.role":         "role",
		"card.lastFour":     "last4",
		"card.expiry.Month": "expiryMonth",
		"card.expiry.Year":  "expiryYear",
		"address":           "unit+street",
	}

	desiredJSON, err := mappingConfig.Reshape([]byte(original))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(desiredJSON))
}
