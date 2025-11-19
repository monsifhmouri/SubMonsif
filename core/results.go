package core

import (
	"fmt"
	"os"
	"sort"
)

func SaveResults(results []string, outputFile string) {
	// 
	sort.Strings(results)

	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			return
		}
		defer file.Close()

		for _, result := range results {
			file.WriteString(result + "\n")
		}
		fmt.Printf("[+] Results saved to: %s\n", outputFile)
	}

	// 
	for _, result := range results {
		fmt.Println(result)
	}
}