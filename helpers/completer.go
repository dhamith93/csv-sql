package helpers

import "github.com/c-bata/go-prompt"

func Completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "LOAD", Description: "LOAD /path/to/file table_name"},
		{Text: "SAVE", Description: "SAVE table_name /path/to/file"},
		{Text: "DB", Description: "DB /path/to/sqlite/db"},
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
