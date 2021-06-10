package main

import (
	"csv-sql/entity"
	"csv-sql/helpers"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	_ "github.com/mattn/go-sqlite3"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "LOAD", Description: "LOAD /path/to/file table_name"},
		{Text: "SAVE", Description: "SAVE table_name /path/to/file"},
		{Text: "SHOW TABLES", Description: ""},
		{Text: "EXIT", Description: ""},
		{Text: "SELECT", Description: ""},
		{Text: "INSERT", Description: ""},
		{Text: "INTO", Description: ""},
		{Text: "VALUES", Description: ""},
		{Text: "UPDATE", Description: ""},
		{Text: "DELETE", Description: ""},
		{Text: "FROM", Description: ""},
		{Text: "WHERE", Description: ""},
		{Text: "AND", Description: ""},
		{Text: "INNER", Description: ""},
		{Text: "LEFT", Description: ""},
		{Text: "RIGHT", Description: ""},
		{Text: "FULL", Description: ""},
		{Text: "JOIN", Description: ""},
		{Text: "ON", Description: ""},
		{Text: "SET", Description: ""},
		{Text: "LIMIT", Description: ""},
		{Text: "ORDER", Description: ""},
		{Text: "ASC", Description: ""},
		{Text: "DESC", Description: ""},
		{Text: "NULL", Description: ""},
		{Text: "LIKE", Description: ""},
		{Text: "IS", Description: ""},
		{Text: "NOT", Description: ""},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func main() {
	defer handleExit()
	fmt.Println("Welcome to csvql")
	files := make([]entity.File, 0)
	tableCount := 0
	dbName := "/tmp/csvql_db_" + helpers.RandSeq(10) + ".db"
	fmt.Println(dbName)
	db := createDB(dbName)
	defer db.Close()

	for {
		response := strings.TrimSpace(prompt.Input("cmd > ", completer))
		responseArr := strings.Fields(response)
		if len(responseArr) > 0 {
			cmd := strings.ToUpper(responseArr[0])

			if cmd == "LOAD" {
				files = loadFile(responseArr, files, db, tableCount)
				continue
			}

			if cmd == "EXIT" {
				break
			}

			if cmd == "SELECT" {
				helpers.PrintMemUsage()
				fmt.Println(helpers.GetData(db, response))
				continue
			}

			if cmd == "SHOW" && response == "SHOW TABLES" {
				printFiles(files)
				continue
			}

			if cmd == "SAVE" {
				saveFile(responseArr, db, files)
				continue
			}

			affectedRows, err := helpers.RunQuery(db, response)
			if err != nil {
				fmt.Printf("Error in %v : %v\n", response, err.Error())
				continue
			}

			fmt.Printf("%v row(s) affected...\n", affectedRows)
		} else {
			continue
		}
	}

	helpers.PrintMemUsage()
	os.Remove(dbName)
}

func loadFile(responseArr []string, files []entity.File, db *sql.DB, tableCount int) []entity.File {
	if len(responseArr) == 3 {
		path := responseArr[1]
		tableName := responseArr[2]

		for i := range files {
			if files[i].Table == tableName {
				fmt.Printf("Error table exists: %v\n", tableName)
				return files
			}
		}

		if _, err := os.Stat(path); err == nil || os.IsExist(err) {
			for {
				fmt.Println("File has a header row (y/n)?")
				response := strings.ToUpper(strings.TrimSpace(prompt.Input("> ", completer)))
				if response == "Y" || response == "N" {
					content, fileErr := helpers.ReadCSVFile(path)
					if fileErr != nil {
						fmt.Printf("Error reading file %v : %v\n", path, fileErr.Error())
						break
					}

					if len(content) == 0 {
						fmt.Println("Empty file")
						break
					}

					headers := make([]string, 0)

					if response == "Y" {
						headers = content[0]
						content = content[1:]
					} else {
						for {
							fmt.Println("Enter " + strconv.Itoa(len(content[0])) + " headers seperated by commas")
							headers = strings.Split(strings.TrimSpace(prompt.Input("> ", completer)), ",")
							if len(content[0]) == len(headers) {
								break
							}
						}
					}

					file := entity.File{
						Path:    path,
						Headers: headers,
						Table:   strings.TrimSpace(tableName),
					}

					files = append(files, file)
					file.Content = content
					helpers.PopulateTables(db, file)
					tableCount++
					break
				}
			}
		} else {
			fmt.Printf("File doesn't exists: %v\n", path)
		}
	} else {
		fmt.Println("Please use LOAD /path/to/file.csv table_name")
	}
	return files
}

func saveFile(responseArr []string, db *sql.DB, files []entity.File) {
	if len(responseArr) == 3 {
		table := strings.TrimSpace(responseArr[1])
		path := strings.TrimSpace(responseArr[2])
		file := entity.File{}
		found := false

		for i := range files {
			if files[i].Table == table {
				file = files[i]
				found = true
				break
			}
		}

		if !found {
			fmt.Println("Table not found")
			return
		}

		result := helpers.GetData(db, "SELECT * FROM "+table)
		if len(result) > 0 {
			helpers.WriteToCSV(path, file.Headers, result)
		}
	} else {
		fmt.Println("Please use SAVE table_name /path/to/file")
	}
}

func createDB(dbName string) *sql.DB {
	file, err := os.Create(dbName)
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	db, _ := sql.Open("sqlite3", dbName)
	return db
}

func printFiles(files []entity.File) {
	fmt.Println("---")
	for _, file := range files {
		fmt.Printf("Table: %v\nHeaders: %v\nPath: %v\n---\n", file.Table, strings.Join(file.Headers, ", "), file.Path)
	}
}

func handleExit() {
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}
