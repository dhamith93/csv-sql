package helpers

import (
	"csv-sql/entity"
	"os"

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
