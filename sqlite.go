package schemata

import (
	"database/sql"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"h12.me/gengo"
)

type SQLite struct {
	DB *sql.DB
}

func (d SQLite) Schema(source string) (*Schema, error) {
	if isSelectStmt(source) {
		return d.schemaFromSelect(source)
	}
	return d.schema(source)
}

func (d SQLite) schemaFromSelect(stmt string) (*Schema, error) {
	view := "view_" + strconv.Itoa(rand.Int())
	createStmt := fmt.Sprintf("CREATE VIEW %s AS %s", view, fmt.Sprintf(stmt, "TRUE"))
	if _, err := d.DB.Exec(createStmt); err != nil {
		return nil, err
	}
	defer d.DB.Exec(fmt.Sprintf("DROP VIEW %s", view))
	s, err := d.schema(view)
	if err != nil {
		return nil, err
	}
	s.Name = stmt
	s.FromSelect = true
	return s, nil
}

func (d SQLite) schema(table string) (*Schema, error) {
	rows, err := d.DB.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return nil, err
	}
	schema := Schema{Name: table, DB: "sqlite", GoName: gengo.SnakeToUpperCamel(table)}
	for rows.Next() {
		var i, null, key int
		var field, type_ string
		var default_ *string
		if err := rows.Scan(&i, &field, &type_, &null, &default_, &key); err != nil {
			return nil, err
		}
		schema.Fields = append(schema.Fields, d.parseField(field, type_, null, key))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &schema, nil
}

func (d SQLite) parseField(field, type_ string, null, key int) Field {
	return Field{
		Name:     field,
		GoName:   gengo.SnakeToUpperCamel(field),
		Primary:  key == 1,
		Nullable: null == 1,
		Type:     type_,
		GoType:   ParseSQLiteType(type_),
	}
}
func ParseSQLiteType(type_ string) reflect.Type {
	ss := strings.Split(type_, "(")
	switch ss[0] {
	case "INTEGER", "BOOLEAN":
		return reflect.TypeOf(int(0))
	case "REAL":
		return reflect.TypeOf(float64(0))
	case "DATETIME", "TIMESTAMP":
		return reflect.TypeOf("")
	case "TEXT":
		return reflect.TypeOf("")
	case "BLOB":
		return reflect.TypeOf([]byte{})
	}
	panic("unknown type " + type_)
}
