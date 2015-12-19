package schemata

import (
	"strings"
)

func isSelectStmt(source string) bool {
	return strings.ToUpper(strings.Fields(source)[0]) == "SELECT"
}
