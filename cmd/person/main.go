package main

import (
    "fmt"
    "log"
    "github.com/Perlishnov/gotrainingproject/internal/models"
)

func main() {
    testJSON := `{"name": "Lucas", "age": 20}`
    
    person, err := models.FromJSON(testJSON)
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    
    fmt.Printf("Parsed: %+v\n", person)
    fmt.Printf("Is adult? %v\n", models.IsAdult(*person))
}