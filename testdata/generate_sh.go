// Generated sample sh code for load testing
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run generate_sh.go <size_in_mb> <output_file>")
		return
	}

	sizeInMB := os.Args[1]
	outputFile := os.Args[2]
	numericPart := strings.TrimSuffix(sizeInMB, "MB")

	content := strings.Repeat("echo 'Fil it up !'\n", 1024*100) // a large file~ 	1,700 KB ~ 1.5mb

	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Failed to create file: %s\n", err)
		return
	}
	defer file.Close()
	limit, err := strconv.Atoi(numericPart)
	if err != nil {
		fmt.Println("Invalid input. Please provide a numeric value followed by 'MB'.")
		return
	}
	for i := 0; i < limit; i++ { // Adjust the multiplier for larger files
		_, err := file.WriteString(content)
		if err != nil {
			fmt.Printf("Failed to write to file: %s\n", err)
			return
		}
	}

	fmt.Printf("Generated shell script of approximately %s: %s\n", sizeInMB, outputFile)
}
