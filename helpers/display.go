package helpers

import (
	"csv-sql/entity"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func PrintTable(resultTable entity.Table) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(resultTable.Headers)
	table.SetAutoFormatHeaders(false)
	for _, v := range resultTable.Data {
		table.Append(v)
	}
	table.Render()
}

func PrintFiles(files []entity.File) {
	fmt.Println("---")
	for _, file := range files {
		fmt.Printf("Table: %v\nHeaders: %v\nPath: %v\n---\n", file.Table, strings.Join(file.Headers, ", "), file.Path)
	}
}
