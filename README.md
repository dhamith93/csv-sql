# csv-sql
Command-line tool to load csv and excel (xlsx) files and run sql commands 

[![Go](https://github.com/dhamith93/csv-sql/actions/workflows/go.yml/badge.svg)](https://github.com/dhamith93/csv-sql/actions/workflows/go.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/dhamith93/csv-sql)](https://goreportcard.com/report/github.com/dhamith93/csv-sql)

## Usage

csv-sql supports loading and saving results as CSV and XLSX files with data processing with SQLite compatible sql commands.

Also, this can be used to open existing SQLite DBs and extract data as CSV.

### Loading a file
```
LOAD /path/to/file table_name
```
You can set up headers if the first row is not a header.
For XLSX files, when loading, this will ask to select the sheet of the file to load.

### Opening a existing SQLite DB
```
DB /path/to/db
```

### Creating a new table with a select query
```sql
CREATE TABLE emp_user AS SELECT emp.emp_id, emp.name, user.user_name, user.role FROM emp INNER JOIN user ON emp.user_id = user.id
```

### Saving a table as a csv

This only supports CSV for now

```
SAVE table_name /path/to/save.csv
```
## Screenshots

![](screenshots/screenshot_01.png)
