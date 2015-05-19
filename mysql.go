package schemata

import (
	"database/sql"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	db *sql.DB
}

func (d Mysql) Schema(name, source string) (*Schema, error) {
	if isSelectStmt(source) {
		return d.schemaFromSelect(name, source)
	}
	return d.schema(name, source)
}
func isSelectStmt(source string) bool {
	return strings.ToUpper(strings.Fields(source)[0]) == "SELECT"
}

func (d Mysql) schemaFromSelect(name, stmt string) (*Schema, error) {
	view := "view_" + strconv.Itoa(rand.Int())
	createStmt := fmt.Sprintf("CREATE VIEW %s AS %s", view, stmt)
	if _, err := d.db.Exec(createStmt); err != nil {
		return nil, err
	}
	defer d.db.Exec(fmt.Sprintf("DROP VIEW %s", view))
	return d.schema(name, view)
}

func (d Mysql) schema(name, table string) (*Schema, error) {
	rows, err := d.db.Query(fmt.Sprintf("SHOW COLUMNS FROM %s", table))
	if err != nil {
		return nil, err
	}
	schema := Schema{Name: name}
	for rows.Next() {
		var field, type_, null, key, extra string
		var default_ sql.NullString
		if err := rows.Scan(&field, &type_, &null, &key, &default_, &extra); err != nil {
			return nil, err
		}
		schema.Fields = append(schema.Fields, d.parseField(field, type_, null, key))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &schema, nil
}

func (d Mysql) parseField(field, type_, null, key string) Field {
	return Field{
		Name:     field,
		Primary:  key == "PRI",
		Nullable: null == "YES",
		DbType:   type_,
		GoType:   d.parseType(type_),
	}
}
func (d Mysql) parseType(type_ string) reflect.Type {
	ss := strings.Split(type_, "(")
	switch ss[0] {
	case "tinyint", "int", "integer", "smallint", "mediumint", "bigint":
		return reflect.TypeOf(int(0))
	case "boolean", "bool":
		return reflect.TypeOf(bool(false))
	case "decimal", "float":
		return reflect.TypeOf(float32(0))
	case "double":
		return reflect.TypeOf(float64(0))
	case "date", "datetime", "timestamp", "time", "year":
		return reflect.TypeOf(time.Time{})
	case "char", "varchar", "text", "tinytext":
		return reflect.TypeOf("")
	}
	return nil
}
