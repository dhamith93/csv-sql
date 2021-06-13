package helpers

import (
	"csv-sql/entity"
	"database/sql"
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

func ShowTables(db *sql.DB) {
	tableNamesQuery := "SELECT name FROM sqlite_master"
	tableNames := GetData(db, tableNamesQuery)
	result := make([][]string, 0)
	resultTable := entity.Table{
		Headers: []string{"table", "columns"},
	}

	for _, table := range tableNames.Data {
		columnNamesQuery := "SELECT name FROM pragma_table_info('" + table[0] + "')"
		columnNames := GetData(db, columnNamesQuery)
		var columns []string

		for _, v1 := range columnNames.Data {
			s := strings.Join(v1, ", ")
			columns = append(columns, s)
		}
		result = append(result, []string{table[0], strings.Join(columns, ", ")})
	}

	resultTable.Data = result
	PrintTable(resultTable)
}
