package schemata

import (
	"fmt"
	"io"
	"strings"
)

func (s *Schema) ToStruct(w io.Writer) {
	fp(w, "type struct %s {\n", goName(s.Name))
	for i := range s.Fields {
		s.Fields[i].ToStruct(w)
	}
	fp(w, "}\n")
}

func (f *Field) ToStruct(w io.Writer) {
	fp(w, "    %-25s *%-10s    `json:\"%s,omitempty\"`\n", goName(f.Name), f.GoType.String(), f.Name)
}

func goName(s string) string {
	ss := strings.Split(s, "_")
	for i := range ss {
		ss[i] = strings.Title(ss[i])
	}
	return strings.Join(ss, "")
}

func fp(w io.Writer, format string, v ...interface{}) {
	fmt.Fprintf(w, format, v...)
}
