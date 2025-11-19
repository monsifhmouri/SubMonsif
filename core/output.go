package core

import (
	"encoding/json"
	"fmt"
	"os"
)

func Output(results []Result, outputFile string, verbose bool) {
	var output *os.File
	var err error

	if outputFile != "" {
		output, err = os.Create(outputFile)
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			return
		}
		defer output.Close()
	}

	for _, result := range results {
		line := formatResult(result, verbose)
		fmt.Println(line)

		if output != nil {
			output.WriteString(line + "\n")
			
			// JSON output for advanced processing
			jsonData, _ := json.Marshal(result)
			output.WriteString("//JSON// " + string(jsonData) + "\n")
		}
	}
}

func formatResult(result Result, verbose bool) string {
	base := fmt.Sprintf("[%d] %s [%s] [Tech: %v] [Time: %v]",
		result.StatusCode, result.URL, result.Server, 
		result.Technologies, result.ResponseTime)

	if verbose {
		base += fmt.Sprintf(" [Size: %d] [IP: %s]", result.ContentLength, result.IP)
	}

	return base
}