package main

import (
	"database/sql"
	"fmt"
	"log"

	"h12.me/schemata"

	"github.com/docopt/docopt-go"
)

func main() {
	usage := `Schemata
Usage:
  schemata extract <db> <conn_str> <table>
  schemata generate <schema_json>

`

	arg, _ := docopt.Parse(usage, nil, true, "Schemata", false)
	if arg["extract"].(bool) {
		db, conn, table := arg["<db>"].(string), arg["<conn_str>"].(string), arg["<table>"].(string)
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
	}
}
