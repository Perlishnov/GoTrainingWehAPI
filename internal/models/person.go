package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Person struct{
	Name string `json:"name"`
	Age int `json:"age"`
}

func fromJSON(input string) (*Person, error)  {
	var person Person
	decoder := json.NewDecoder(strings.NewReader(input))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&person); err != nil {
		return nil, fmt.Errorf("JSON parsing error: %w", err)
	}

	if strings.TrimSpace(person.Name) == "" {
		return nil, fmt.Errorf("Name can not be empty")
	}
	if person.Age < 0 {
		return nil, fmt.Errorf("Age can not be negative")
	}

	return  &person, nil
	
}