package schemata

import (
	"fmt"
	"reflect"
	"time"
)

type DB interface {
	ParseType(string) reflect.Type
}

type Schema struct {
	Name   string
	Fields Fields
	DB     DB
}

type Field struct {
	Name     string
	Primary  bool
	Nullable bool
	Type     string
}

type Fields []Field

type GoType struct {
	reflect.Type
}

func (t *GoType) MarshalText() (text []byte, err error) {
	return []byte(t.Type.String()), nil
}

func (t *GoType) UnmarshalText(text []byte) error {
	switch string(text) {
	case "string":
		t.Type = reflect.TypeOf("")
	case "int":
		t.Type = reflect.TypeOf(int(0))
	case "bool":
		t.Type = reflect.TypeOf(bool(false))
	case "time.Time":
		t.Type = reflect.TypeOf(time.Time{})
	case "float32":
		t.Type = reflect.TypeOf(float32(0))
	case "float64":
		t.Type = reflect.TypeOf(float64(0))
	}
	if t.Type == nil {
		return fmt.Errorf("fail to parse type %s", string(text))
	}
	return nil
}

func (f *Field) GoType(db DB) GoType {
	return GoType{db.ParseType(f.Type)}
}
