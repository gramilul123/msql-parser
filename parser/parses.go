package parser

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/jmoiron/sqlx"

	"mysql-parser/config"
	"mysql-parser/models"
	"mysql-parser/mysql"
)

const (
	GET_FIELDS_QUERY       = "SHOW COLUMNS FROM table_name;"
	GET_TABLE_DATA_QUERY   = "SELECT * FROM table_name WHERE where_text;"
	GET_CREATE_TABLE_QUERY = "SHOW CREATE TABLE table_name;"
)

func Run() {
	cnf := config.New()

	db, err := mysql.NewMysql(cnf.Db)
	if err != nil {

		log.Fatalf(err.Error())
	}

	tablesConditions := cnf.Tables
	tableList := &models.TablesList{}

	tableList, err = ParseTables(db, cnf, tableList, tablesConditions)
	if err != nil {

		log.Fatalf(err.Error())
	}
}

func ParseTables(db *sqlx.DB, cnf *config.Config, tablesList *models.TablesList, tables map[string]string) (*models.TablesList, error) {

	var err error

	for table, where := range tables {
		tableData := &models.TableData{
			Table: table,
		}

		if len(where) == 0 {

			return tablesList, fmt.Errorf("where con`t be empty for table, %s", table)
		}

		if tableData.Data, err = GetTableData(db, table, where); err != nil {

			return tablesList, err
		}

		if tableData.Fields, err = GetTableFields(db, table); err != nil {

			return tablesList, err
		}

		err = CreateTableDump(db, cnf, table, tableData)
		if err != nil {

			return tablesList, err
		}

		tablesList.Tables = append(tablesList.Tables, tableData)
	}

	return tablesList, nil
}

func GetTableData(db *sqlx.DB, table, where string) (*models.Data, error) {

	data := &models.Data{}

	query := strings.Replace(GET_TABLE_DATA_QUERY, "table_name", table, -1)
	query = strings.Replace(query, "where_text", where, -1)
	rows, err := db.Query(query)
	if err != nil {

		return data, fmt.Errorf("can't make get table data query, %w", err)
	}

	cols, err := rows.Columns()
	if err != nil {

		return data, fmt.Errorf("can't get rows, %w", err)
	}

	data.Fields = append(data.Fields, cols...)

	for rows.Next() {

		row := make([]*sql.NullString, 0, len(cols))
		columns := make([]sql.NullString, len(cols))
		columnPointers := make([]interface{}, len(cols))

		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return data, fmt.Errorf("can't rows scan table data, %w", err)
		}

		for i := range cols {

			val := columnPointers[i]
			row = append(row, val.(*sql.NullString))
		}

		data.ValuesList = append(data.ValuesList, models.Values{Values: row})
	}

	return data, nil
}

func GetTableFields(db *sqlx.DB, table string) ([]*models.Field, error) {

	fields := []*models.Field{}
	query := strings.Replace(GET_FIELDS_QUERY, "table_name", table, -1)

	err := db.Select(&fields, query)
	if err != nil {

		return fields, fmt.Errorf("can't get table field, %w", err)
	}

	return fields, nil
}

func GetTableDump(db *sqlx.DB, table string) (*models.ShowTable, error) {
	dump := &models.ShowTable{}
	query := strings.Replace(GET_CREATE_TABLE_QUERY, "table_name", table, -1)

	rows, err := db.Query(query)
	if err != nil {

		return dump, fmt.Errorf("can't get show table, %w", err)
	}

	for rows.Next() {
		if err := rows.Scan(
			&dump.Table,
			&dump.Create,
		); err != nil {
			return dump, fmt.Errorf("can't rows scan, %w", err)
		}
	}

	return dump, nil
}

func CreateTableDump(db *sqlx.DB, cnf *config.Config, table string, tableData *models.TableData) error {

	f, err := os.Create("./dumps/" + table + ".sql")
	if err != nil {

		return err
	}

	dump := &models.Dump{
		CompleteTime: time.Now().String(),
		Name:         table,
	}

	createTable, err := GetTableDump(db, table)
	if err != nil {

		return err
	}

	dump.Table = createTable.Create
	dump.Fields, dump.Values = PrepareInsertValues(tableData)

	if dump.ServerVersion, err = GetServerVersion(db); err != nil {

		return err
	}

	tmpl, err := template.New("mysqldump").Parse(models.DUMP_TEMPLATE)
	if err != nil {

		return err
	}

	if err = tmpl.Execute(f, dump); err != nil {

		return err
	}

	return nil
}

func GetServerVersion(db *sqlx.DB) (string, error) {

	var server_version sql.NullString

	if err := db.QueryRow("SELECT version()").Scan(&server_version); err != nil {

		return "", err
	}

	return server_version.String, nil
}

func PrepareInsertValues(tableData *models.TableData) (string, string) {

	var fieldsList, valuesList []string
	var fields, values string

	if len(tableData.Data.ValuesList) > 0 {

		fieldsList = append(fieldsList, tableData.Data.Fields...)
		fields = "(" + strings.Join(fieldsList, ", ") + ")"

		for _, row := range tableData.Data.ValuesList {
			valuesRow := []string{}

			for key, value := range row.Values {

				val := "NULL"
				if value.Valid {
					val = PrepareValue(value.String, key, tableData.Fields)
				}
				valuesRow = append(valuesRow, val)
			}
			valuesList = append(valuesList, "("+strings.Join(valuesRow, ", ")+")")
		}

		values = strings.Join(valuesList, ", ")
	}

	return fields, values
}

func PrepareValue(value string, key int, fields []*models.Field) string {

	if strings.Contains(fields[key].Type, "char") ||
		strings.Contains(fields[key].Type, "text") ||
		strings.Contains(fields[key].Type, "date") {
		value = "'" + value + "'"
	}

	return value
}
