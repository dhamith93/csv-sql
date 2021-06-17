package file

import (
	"csv-sql/internal/table"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gabriel-vasile/mimetype"
)

type File struct {
	Path    string
	Headers []string
	Table   string
	Content [][]string
}

func GetMimeType(path string) (string, error) {
	mtype, err := mimetype.DetectFile(path)
	if err != nil {
		return "", err
	}
	return mtype.String(), nil
}

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

func GetXLSXSheetList(path string) ([]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	return f.GetSheetList(), nil
}

func ReadXLSXFile(path string, sheet string) ([][]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func WriteToCSV(path string, result table.Table) {
	csvFile, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error creating %v : %v\n", path, err.Error())
		return
	}
	csvWriter := csv.NewWriter(csvFile)
	csvWriter.Write(result.Headers)
	csvWriter.WriteAll(result.Data)
	csvWriter.Flush()
	csvFile.Close()
}

func IsFile(path string) bool {
	_, err := os.Open(path)
	return err == nil
}
