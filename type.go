package schemata

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"
)

type DB interface {
	ParseType(string) reflect.Type
}

type Schema struct {
	DB         string `json:"db"`
	Name       string `json:"name"`
	Fields     Fields `json:"fields"`
	FromSelect bool   `json:"from_select"`
}

type Field struct {
	Name     string `json:"name"`
	Primary  bool   `json:"primary,omitempty"`
	Nullable bool   `json:"nullable,omitemtpy"`
	Type     string `json:"type"`
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

func (s *Schema) String() string {
	buf, _ := json.MarshalIndent(s, "", "\t")
	return string(buf)
}

func LoadSchema(file string) (*Schema, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	var s Schema
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}
