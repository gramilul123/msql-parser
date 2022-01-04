package models

import "database/sql"

type Field struct {
	Field   string         `db:"Field"`
	Type    string         `db:"Type"`
	Null    string         `db:"Null"`
	Key     string         `db:"Key"`
	Default sql.NullString `db:"Default"`
	Extra   string         `db:"Extra"`
}

type Key struct {
}

type Values struct {
	Values []*sql.NullString
}

type Data struct {
	Fields     []string
	ValuesList []Values
}

type TableData struct {
	Table  string
	Data   *Data
	Fields []*Field
	Level  int
	Keys   []*Key
}

type TablesList struct {
	Tables []*TableData
}

type ShowTable struct {
	Table  string
	Create string
}

type Dump struct {
	ServerVersion string
	CompleteTime  string
	Name          string
	Table         string
	Fields        string
	Values        string
}
