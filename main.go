package main

import (
	"fmt"

	"github.com/tayalone/go-sql-rql/parser"
)

func main() {
	q := `fields=id%2Cname&sort=id%2C-created_at&skip=112&limit=10`

	p := parser.New(nil)
	qp := p.Parse(q)

	fmt.Println("fields", qp.GetFields())
	fmt.Println("skip", qp.GetSkip())
}
