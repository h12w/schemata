package schemata

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"h12.me/gengo"
)

func (s *Schema) JSON(w io.Writer) {
	buf, _ := json.MarshalIndent(s, "", "    ")
	w.Write(buf)
}

func (s *Schema) Struct(tags ...string) *gengo.TypeDecl {
	return &gengo.TypeDecl{
		Name: s.GoName,
		Type: gengo.Type{
			Kind:   gengo.StructKind,
			Fields: s.goFields(tags),
		},
	}
}

func (s *Schema) InsertQuery(w io.Writer) {
	fpn(w, "INSERT INTO %s (", s.Name)
	s.Fields.ToList(w)
	fpn(w, "\n) VALUES (%s?)", strings.Repeat("?,", len(s.Fields)-1))
}

func (s *Schema) InsertIgnoreQuery(w io.Writer) {
	switch s.DB {
	case "mysql":
		fpn(w, "INSERT IGNORE INTO %s (", s.Name)
	case "sqlite":
		fpn(w, "INSERT OR IGNORE INTO %s (", s.Name)
	}
	s.Fields.ToList(w)
	fpn(w, "\n) VALUES (%s?)", strings.Repeat("?,", len(s.Fields)-1))
}

func (s *Schema) goFields(tags []string) (fields []*gengo.Field) {
	for i := range s.Fields {
		fields = append(fields, s.Fields[i].goField(tags))
	}
	return
}

func (f *Field) goField(tags []string) *gengo.Field {
	var tag gengo.Tag
	if len(tags) > 0 {
		tag = make(gengo.Tag, len(tags))
		for i := range tag {
			tag[i] = &gengo.TagPart{
				Encoding:  tags[i],
				Name:      f.Name,
				OmitEmpty: true,
			}
		}
	}
	return &gengo.Field{
		Name: f.GoName,
		Type: gengo.Type{
			Kind:  gengo.IdentKind,
			Ident: GoType{f.GoType}.String(),
		},
		Tag: tag,
	}
}

func (s *Schema) Select(w io.Writer) {
	if s.FromSelect {
		fmt.Fprint(w, s.Name)
		return
	}
	s.Fields.Select(w)
	fpn(w, "FROM\n    %s", s.Name)
	w.Write([]byte("WHERE\n    %s"))
}

func (s *Schema) Scan(w io.Writer, name string) {
	fpn(w, "var v %s", name)
	fpn(w, "if err := rows.Scan(")
	for _, f := range s.Fields {
		fpn(w, "    &v.%s,", gengo.SnakeToUpperCamel(f.Name))
	}
	fpn(w, "); err != nil {")
	fpn(w, "    return err")
	fpn(w, "}")
}

func (fs Fields) Select(w io.Writer) {
	fpn(w, "SELECT")
	fs.ToList(w)
	fpn(w, "")
}

func (fs Fields) ToList(w io.Writer) {
	for i, f := range fs {
		if i > 0 {
			fpn(w, ",")
		}
		fp(w, "    %s", f.Name)
	}
}

func fp(w io.Writer, format string, v ...interface{}) {
	fmt.Fprintf(w, format, v...)
}

func fpn(w io.Writer, format string, v ...interface{}) {
	fmt.Fprintf(w, format, v...)
	fmt.Fprintln(w)
}
