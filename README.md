# csv-sql
Command-line tool to load csv files and run sql commands 

[![Go](https://github.com/dhamith93/csv-sql/actions/workflows/go.yml/badge.svg)](https://github.com/dhamith93/csv-sql/actions/workflows/go.yml)

## Usage

csv-sql supports loading and saving results as CSV files with data processing with SQLite compatible sql commands

### Loading a file
```
LOAD /path/to/file.csv table_name
```

### Creating a new table with a select query
```sql
CREATE TABLE emp_user AS SELECT emp.emp_id, emp.name, user.user_name, user.role FROM emp INNER JOIN user ON emp.user_id = user.id
```

### Saving a table as a csv
```
SAVE table_name /path/to/save.csv
```
## Screenshots

![](screenshots/screenshot_01.png)
