package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ValidationResult struct {
	FileName   string `json:"templateName"`
	Result     string `json:"result"`
	Exceptions string `json:"exceptions"`
}

// Directory containing the templates
const kTemplateDir = "test"

func main() {
	// Result slice to store the validation results
	var results []ValidationResult

	// Find all the template files in the directory
	err := filepath.Walk(kTemplateDir, func(path string, info os.FileInfo, err error) error {
		// Skip directories and process only regular files with ".yaml" or ".yml" extensions
		if err != nil || info.IsDir() || (!isYAMLfile(path) && !isJSONfile(path)) {
			return nil
		}

		// Run the SAM CLI command for each template file
		result := "fail"
		exceptions := ""
		if err := runSAMCommand("validate", "-t", path, "--region", "us-west-2"); err == nil {
			result = "pass"
		} else {
			exceptions = err.Error()
		}

		// Add the result to the slice
		results = append(results, ValidationResult{
			FileName:   filepath.Base(path),
			Result:     result,
			Exceptions: exceptions,
		})

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the directory: %s\n", err)
		os.Exit(1)
	}

	// Write the results to a JSON file
	jsonData, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		fmt.Printf("Error encoding JSON: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
	// jsonFileName := filepath.Join(kResultDirectory, "validation_results.json")
	// if err := os.WriteFile(jsonFileName, jsonData, 0644); err != nil {
	// 	fmt.Printf("Error writing JSON file: %s\n", err)
	// 	os.Exit(1)
	// }

}

// Helper function to check if a file has a ".yaml" or ".yml" extension
func isYAMLfile(file string) bool {
	return strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml")
}

// Helper function to check if a file has a ".json" extension
func isJSONfile(file string) bool {
	return strings.HasSuffix(file, ".json")
}

// Helper function to run the SAM CLI command with given arguments
func runSAMCommand(args ...string) error {
	cmd := exec.Command("sam", args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
