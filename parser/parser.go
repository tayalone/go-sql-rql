package parser

import (
	"fmt"
	"net/url"
	"strings"
)

var mockAllFields = map[string]interface{}{
	"id":   struct{}{},
	"name": struct{}{},
}

var mockDisabledFields = map[string]interface{}{
	"preventField": struct{}{},
}

// Parser interface
type Parser interface {
	Parse(q string)
}

type parser struct {
	defaultLimit   *int
	allFields      map[string]interface{}
	disabledFields map[string]interface{}
}

// New Query Parser
func New(defaultLimit *int) Parser {
	// myMap := map[string]bool{}

	// myMap["1"] = 1

	return parser{
		defaultLimit:   defaultLimit,
		allFields:      mockAllFields,
		disabledFields: map[string]interface{}{},
	}
}

func (p parser) getFields(str string) []string {
	expectFields := strings.Split(str, ",")

	allFields := []string{}
	if str == "" || len(expectFields) == 0 {
		for key := range p.allFields {
			allFields = append(allFields, key)
		}
		return allFields
	}

	for _, v := range expectFields {
		_, inAllField := p.allFields[v]

		inDisabledFields := false
		if inAllField && !inDisabledFields {
			allFields = append(allFields, v)
		}
	}

	return allFields
}

func (p parser) Parse(q string) {
	params, err := url.ParseQuery(q)
	if err != nil {
		panic("query format error")
	}

	if v := params.Get("fields"); v != "" {
		fmt.Println("action: fields", v)
	}

	// fmt.Println("q is", q)
	// actions := strings.Split(q, "&")
	// for i, a := range actions {
	// 	fmt.Printf("key no: %d, action: %s\n", i, a)

	// 	payloads := strings.Split(a, "=")
	// 	if len(payloads) != 2 {
	// 		panic("query action invalid format")
	// 	}
	// 	key := payloads[0]
	// 	value := payloads[1]

	// }
}
