// cmd/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"smuggr.xyz/arrango/common/models/input"
	"smuggr.xyz/arrango/core/solver"
)

func main() {
	solver := solver.Solver{
		PopulationSize: 50,
		Generations:    1000,
		MutationRate:   0.1,
	}
	result := solver.Solve(input.ExampleInputData)

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatalf("Error converting result to JSON: %v", err)
	}

	fmt.Println("Result:", string(jsonResult))
	filePath := "/home/karol/Documents/Repositories/Arrango/web/mock/public/timetables.json"
	err = os.WriteFile(filePath, jsonResult, 0644)
	if err != nil {
		log.Fatalf("Error writing result to file: %v", err)
	}
}
