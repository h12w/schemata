package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"h12.me/schemata"

	"github.com/docopt/docopt-go"
)

func main() {
	usage := `Schemata
Usage:
  schemata extract <db> <conn-str> <table>
  schemata generate struct <struct-name> <schema-json>
  schemata generate select <schema-json>
  schemata generate scan <struct-name> <schema-json>

`

	arg, _ := docopt.Parse(usage, nil, true, "Schemata", false)
	if arg["extract"].(bool) {
		db, conn, table := arg["<db>"].(string), arg["<conn-str>"].(string), arg["<table>"].(string)
		x, err := sql.Open(db, conn)
		if err != nil {
			log.Fatal(err)
		}

		switch db {
		case "mysql":
			s, _ := schemata.MySQL{DB: x}.Schema(table)
			fmt.Println(s)
		default:
			fmt.Println(arg)
		}
	} else if arg["generate"].(bool) {
		file := arg["<schema-json>"].(string)
		if arg["struct"].(bool) {
			structName := arg["<struct-name>"].(string)
			s, err := schemata.LoadSchema(file)
			if err != nil {
				log.Fatal(err)
			}
			s.Struct(os.Stdout, structName)
		} else if arg["select"].(bool) {
			s, err := schemata.LoadSchema(file)
			if err != nil {
				log.Fatal(err)
			}
			s.Select(os.Stdout)
			s.From(os.Stdout)
		} else if arg["scan"].(bool) {
			s, err := schemata.LoadSchema(file)
			if err != nil {
				log.Fatal(err)
			}
			s.Scan(os.Stdout, arg["<struct-name>"].(string))
		}
	}
}
