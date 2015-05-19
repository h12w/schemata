package schemata

import (
	"reflect"
)

type Schema struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name     string
	Primary  bool
	Nullable bool
	DbType   string
	GoType   reflect.Type
}
