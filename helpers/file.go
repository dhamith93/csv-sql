package helpers

import (
	"encoding/csv"
	"fmt"
	"os"
)

// ReadCSVFile read from given CSV file
func ReadCSVFile(path string) ([][]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return [][]string{}, err
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true
	records, err := reader.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}

func WriteToCSV(path string, headers []string, data [][]string) {
	csvFile, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error creating %v : %v\n", path, err.Error())
		return
	}
	csvWriter := csv.NewWriter(csvFile)
	csvWriter.Write(headers)
	csvWriter.WriteAll(data)
	csvWriter.Flush()
	csvFile.Close()
}
