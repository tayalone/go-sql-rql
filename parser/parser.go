package parser

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/tayalone/go-sql-rql/parser/param"
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
	Parse(q string) param.QueryParams
}

type parser struct {
	defaultLimit   *int
	allFields      map[string]interface{}
	disabledFields map[string]interface{}
}

// New Query Parser
func New(defaultLimit *int) Parser {
	return &parser{
		defaultLimit:   defaultLimit,
		allFields:      mockAllFields,
		disabledFields: map[string]interface{}{},
	}
}

func (p *parser) getFields(str string) []string {
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

func (p *parser) getSkip(str string) int {
	tmp, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	skip := int(tmp)
	if skip < 0 {
		return 0
	}
	return int(skip)
}

func (p *parser) getLimit(str string) *int {
	tmp, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return p.defaultLimit
	}
	limit := int(tmp)
	if limit < 0 {
		return p.defaultLimit
	}
	if p.defaultLimit == nil {
		return &limit
	} else if limit > *p.defaultLimit {
		return p.defaultLimit
	}
	return &limit
}

func (p parser) Parse(q string) param.QueryParams {
	params, err := url.ParseQuery(q)
	if err != nil {
		panic("query format error")
	}

	qp := param.New()

	if v := params.Get("fields"); v != "" {
		fields := p.getFields(v)
		qp.SetFields(fields)
	}

	if v := params.Get("skip"); v != "" {
		skip := p.getSkip(v)
		qp.SetSkip(skip)
	}

	if v := params.Get("limit"); v != "" {
		limit := p.getLimit(v)
		qp.SetLimit(limit)
	}

	return qp
}
