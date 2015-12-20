package schemata

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"h12.me/gengo"
)

func (s *Schema) StructName() string {
	return lowerCamel(s.Name)
}

func (s *Schema) JSON(w io.Writer) {
	buf, _ := json.MarshalIndent(s, "", "    ")
	w.Write(buf)
}

func (s *Schema) Struct(w io.Writer, name string) error {
	decl := &gengo.TypeDecl{
		Name: name,
		Type: gengo.Type{
			Kind:   gengo.StructKind,
			Fields: s.fields(),
		},
	}
	return decl.Marshal(w)
}

func (s *Schema) fields() (fields []*gengo.Field) {
	for i := range s.Fields {
		fields = append(fields, s.Fields[i].goField(s.DB))
	}
	return
}

func (f *Field) goField(db string) *gengo.Field {
	goType := GoType{parseType(db, f.Type)}
	return &gengo.Field{
		Name: upperCamel(f.Name),
		Type: gengo.Type{
			Kind:  gengo.IdentKind,
			Ident: goType.String(),
		},
		Tag: &gengo.Tag{
			Parts: []*gengo.TagPart{
				{
					Encoding:  "json",
					Name:      f.Name,
					OmitEmpty: true,
				},
			},
		},
	}
}
func parseType(db string, t string) reflect.Type {
	switch db {
	case "mysql":
		return ParseMySQLType(t)
	case "sqlite":
		return ParseSQLiteType(t)
	}
	return nil
}

func (s *Schema) Select(w io.Writer) {
	if s.FromSelect {
		fmt.Fprint(w, s.Name)
		return
	}
	s.Fields.Select(w)
	fp(w, "FROM\n    %s\n", s.Name)
	fmt.Fprint(w, "WHERE\n    %s")
}

func (s *Schema) Scan(w io.Writer, name string) {
	fp(w, "var v %s\n", name)
	fp(w, "if err := rows.Scan(\n")
	for _, f := range s.Fields {
		fp(w, "    &v.%s,\n", upperCamel(f.Name))
	}
	fp(w, "); err != nil {\n")
	fp(w, "    return err\n")
	fp(w, "}\n")
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
