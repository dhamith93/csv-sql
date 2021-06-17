package main

import (
	"csv-sql/internal/completer"
	"csv-sql/internal/database"
	"csv-sql/internal/display"
	"csv-sql/internal/file"
	"csv-sql/internal/stringops"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	defer handleExit()
	fmt.Println("Welcome to CSV-SQL")
	files := make([]file.File, 0)
	tableCount := 0
	dbName := "/tmp/csvql_db_" + stringops.RandSeq(10) + ".db"
	db := database.CreateDB(dbName)
	defer db.Close()

	for {
		response := strings.TrimSpace(prompt.Input("cmd > ", completer.Completer))
		responseArr := strings.Fields(response)
		if len(responseArr) > 0 {
			if response == "SHOW TABLES" {
				display.ShowTables(db)
				continue
			}

			cmd := strings.ToUpper(responseArr[0])

			if cmd == "LOAD" {
				files = loadFile(responseArr, files, db, tableCount)
				continue
			}

			if cmd == "EXIT" {
				break
			}

			if cmd == "DB" {
				db = openDB(responseArr, db, dbName)
				continue
			}

			if cmd == "SELECT" {
				display.PrintTable(database.GetData(db, response))
				continue
			}

			if cmd == "SAVE" {
				saveFile(responseArr, db, files)
				continue
			}

			affectedRows, err := database.RunQuery(db, response)
			if err != nil {
				fmt.Printf("Error in %v : %v\n", response, err.Error())
				continue
			}

			fmt.Printf("%v row(s) affected...\n", affectedRows)
		}
	}
	os.Remove(dbName)
}

func openDB(responseArr []string, db *sql.DB, dbName string) *sql.DB {
	if len(responseArr) > 1 {
		if file.IsFile(responseArr[1]) {
			var err error
			db, err = database.OpenDB(db, responseArr[1])

			if err != nil {
				fmt.Printf("Error opening DB : %v\n", err.Error())
				fmt.Println("Falling back to default DB")
				db, _ = database.OpenDB(db, dbName)
			}

		}
	} else {
		fmt.Println("Not a valid file to open as a DB")
	}
	return db
}

func loadFile(responseArr []string, files []file.File, db *sql.DB, tableCount int) []file.File {
	if len(responseArr) == 3 {
		path := responseArr[1]
		tableName := responseArr[2]
		mimeType, err := file.GetMimeType(path)
		if err != nil {
			fmt.Printf("Error reading mimetype of file %v : %v\n", path, err.Error())
			return files
		}

		if !validFileType(mimeType) {
			fmt.Printf("Error reading file %v : file type not supported\n", path)
			return files
		}

		for i := range files {
			if files[i].Table == tableName {
				fmt.Printf("Error table exists: %v\n", tableName)
				return files
			}
		}

		for {
			fmt.Println("File has a header row (y/n)?")
			response := strings.ToUpper(strings.TrimSpace(prompt.Input("> ", completer.Completer)))
			if response == "Y" || response == "N" {
				var (
					content [][]string
					fileErr error
					headers = make([]string, 0)
				)

				if mimeType == "text/csv" {
					content, fileErr = file.ReadCSVFile(path)
				}

				if mimeType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
					sheets, err := file.GetXLSXSheetList(path)
					if err != nil {
						fmt.Printf("Error reading file %v : %v\n", path, err.Error())
						return files
					}
					fmt.Printf("What is the name of the sheet to be loaded? %v\n", strings.Join(sheets, ", "))
					sheet := strings.TrimSpace(prompt.Choose("> ", sheets))
					content, fileErr = file.ReadXLSXFile(path, sheet)
				}

				if fileErr != nil {
					fmt.Printf("Error reading file %v : %v\n", path, fileErr.Error())
					return files
				}

				if len(content) == 0 {
					fmt.Println("Empty file")
					return files
				}

				if response == "Y" {
					headers = content[0]
					content = content[1:]
				} else {
					for {
						fmt.Println("Enter " + strconv.Itoa(len(content[0])) + " headers separated by commas")
						headers = strings.Split(strings.TrimSpace(prompt.Input("> ", completer.Completer)), ",")
						if len(content[0]) == len(headers) {
							break
						}
					}
				}

				for i := range headers {
					headers[i] = strings.TrimSpace(headers[i])
				}

				file := file.File{
					Path:    path,
					Headers: headers,
					Table:   strings.TrimSpace(tableName),
				}

				files = append(files, file)
				file.Content = content
				database.PopulateTables(db, &file)
				tableCount++
				file.Content = nil
				break
			}
		}

	}
	return files
}

func saveFile(responseArr []string, db *sql.DB, files []file.File) {
	if len(responseArr) == 3 {
		table := strings.TrimSpace(responseArr[1])
		path := strings.TrimSpace(responseArr[2])

		res := database.GetData(db, "SELECT COUNT(*) FROM sqlite_master where type='table' AND tbl_name = '"+table+"'")
		tableCount, _ := strconv.Atoi(res.Data[0][0])
		if tableCount == 0 {
			fmt.Println("Table not found")
			return
		}

		result := database.GetData(db, "SELECT * FROM "+table)
		if len(result.Data) > 0 {
			file.WriteToCSV(path, result)
		}
	} else {
		fmt.Println("Please use SAVE table_name /path/to/file")
	}
}

func validFileType(mimeType string) bool {
	return mimeType == "text/csv" || mimeType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
}

func handleExit() {
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}
