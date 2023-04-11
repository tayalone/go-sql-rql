package main

import "github.com/tayalone/go-sql-rql/parser"

func main() {
	q := `fields=id%2Cname&sort=id%2C-created_at&skip=0&limit=10`

	p := parser.New(nil)

	p.Parse(q)
}
