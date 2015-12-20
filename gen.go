package schemata

import (
	"encoding/json"
	"fmt"
	"io"

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
	fp(w, "FROM\n    %s\n", s.Name)
	w.Write([]byte("WHERE\n    %s"))
}

func (s *Schema) Scan(w io.Writer, name string) {
	fp(w, "var v %s\n", name)
	fp(w, "if err := rows.Scan(\n")
	for _, f := range s.Fields {
		fp(w, "    &v.%s,\n", gengo.GoUpperName(f.Name))
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

func fp(w io.Writer, format string, v ...interface{}) {
	fmt.Fprintf(w, format, v...)
}
