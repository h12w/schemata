package schemata

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func (s *Schema) StructName() string {
	return lowerCamel(s.Name)
}

func (s *Schema) JSON(w io.Writer) {
	buf, _ := json.MarshalIndent(s, "", "    ")
	w.Write(buf)
}

func (s *Schema) Struct(w io.Writer) {
	fp(w, "type %s struct {\n", s.StructName())
	for i := range s.Fields {
		s.Fields[i].Struct(w, s.DB)
	}
	fp(w, "}\n")
}

func (f *Field) Struct(w io.Writer, db DB) {
	fp(w, "    %-25s *%-10s    `json:\"%s,omitempty\"`\n", upperCamel(f.Name), f.GoType(db).String(), f.Name)
}

func (s *Schema) Select(w io.Writer) {
	s.Fields.Select(w)
}

func (s *Schema) Scan(w io.Writer) {
	fp(w, "var v %s\n", s.StructName())
	fp(w, "rows.Scan(\n")
	for _, f := range s.Fields {
		fp(w, "    &v.%s,\n", upperCamel(f.Name))
	}
	fp(w, ")\n")
}

func (fs Fields) Select(w io.Writer) {
	fp(w, "SELECT\n")
	for i, f := range fs {
		if i > 0 {
			fp(w, ",\n")
		}
		fp(w, "    %s", f.Name)
	}
	fp(w, "\n")
}

func (s *Schema) From(w io.Writer) {
	fp(w, "FROM\n    %s\n", s.Name)
}

func upperCamel(s string) string {
	ss := strings.Split(s, "_")
	for i := range ss {
		ss[i] = strings.Title(ss[i])
	}
	return strings.Join(ss, "")
}

func lowerCamel(s string) string {
	ss := strings.Split(s, "_")
	for i := 1; i < len(ss); i++ {
		ss[i] = strings.Title(ss[i])
	}
	return strings.Join(ss, "")
}

func fp(w io.Writer, format string, v ...interface{}) {
	fmt.Fprintf(w, format, v...)
}
