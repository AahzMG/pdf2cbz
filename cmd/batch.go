package cmd

import (
	"fmt"
	"os"
)

func batchPDF2CBZ(source string) (result string) {
	files, err := os.ReadDir(source)
	if err != nil {
		return fmt.Sprintf("Error %f", err)
	}
	if len(files) < 1 {
		return "Error no files found"
	}

	for _, file := range files {
		filepath := source + "/" + file.Name()
		res := convertPDF2CBZ(filepath)
		fmt.Println("Processed ", res)
	}

	return source
}
