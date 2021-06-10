package helpers

import (
	"csv-sql/entity"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func PopulateTables(db *sql.DB, file entity.File) {
	header := ""
	for i := 0; i < len(file.Headers); i++ {
		header += "\"" + file.Headers[i] + "\" TEXT"
		if i != (len(file.Headers) - 1) {
			header += ","
		}
	}
	query := "CREATE TABLE " + file.Table + " (" + header + ");"
	statement, err := db.Prepare(query)
	if err != nil {
		fmt.Printf("Error in %v : %v\n", query, err.Error())
	}
	statement.Exec()
	tx, _ := db.Begin()
	hasError := false

	for _, row := range file.Content {
		var builder strings.Builder

		for _, c := range row {
			builder.WriteString("\"")
			builder.WriteString(c)
			builder.WriteString("\"")
			builder.WriteString(",")
		}

		data := strings.TrimSuffix(builder.String(), ",")

		insertQuery := "INSERT INTO " + file.Table + " VALUES (" + data + ");"
		statement, err := tx.Prepare(insertQuery)
		if err != nil {
			fmt.Printf("Error in %v : %v\n", insertQuery, err.Error())
		}
		_, err = statement.Exec()
		if err != nil {
			hasError = true
			break
		}
	}
	if hasError {
		tx.Rollback()
	} else {
		defer tx.Commit()
	}
	file.Content = nil
}

func RunQuery(db *sql.DB, query string) (int64, error) {
	tx, _ := db.Begin()
	statement, err := tx.Prepare(query)
	if err != nil {
		fmt.Printf("Error in %v : %v\n", query, err.Error())
		return 0, nil
	}
	res, err := statement.Exec()
	if err != nil {
		tx.Rollback()
		fmt.Printf("Error in %v : %v\n", query, err.Error())
		return 0, nil
	}
	defer tx.Commit()
	return res.RowsAffected()
}

func GetData(db *sql.DB, query string) [][]string {
	row, err := db.Query(query)
	if err != nil {
		fmt.Printf("Error in %v : %v\n", query, err.Error())
		return nil
	}
	defer row.Close()
	columns, err := row.Columns()
	if err != nil {
		fmt.Printf("Error reading columns %v : %v\n", query, err.Error())
		return nil
	}
	output := make([][]string, 0)
	rawResult := make([][]byte, len(columns))
	dest := make([]interface{}, len(columns))
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}
	for row.Next() {
		row.Scan(dest...)
		res := make([]string, 0)
		for _, raw := range rawResult {
			if raw != nil {
				res = append(res, string(raw))
			}
		}
		output = append(output, res)
	}
	return output
}
