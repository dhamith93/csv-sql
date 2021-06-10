package entity

type File struct {
	Path    string
	Headers []string
	Table   string
	Content [][]string
}
