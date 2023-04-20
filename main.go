package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	ParserValid "github.com/tayalone/go-sql-rql/parser/valid"
)

const TimeLayout = "2006-01-02T15:04:05.999Z"

var operation = map[string]interface{}{
	"$eq":         struct{}{},
	"$is":         struct{}{},
	"$gt":         struct{}{},
	"$gte":        struct{}{},
	"$lt":         struct{}{},
	"$lte":        struct{}{},
	"$between":    struct{}{},
	"$notBetween": struct{}{},
	"$in":         struct{}{},
}

var linkOperation = map[string]interface{}{
	"$and": struct{}{},
	"$or":  struct{}{},
}

var notOperation = map[string]interface{}{
	"$not": struct{}{},
}

type node interface {
	GetParent() node
	IsLeaf() bool
	AddChild(n node)
	GetSQLQuery() (string, []interface{})
}

type rootNode struct {
	children []node
}

func (rt *rootNode) GetParent() node {
	return nil
}

func (rt *rootNode) IsLeaf() bool {
	return false
}

func (rt *rootNode) AddChild(n node) {
	rt.children = append(rt.children, n)
}

func (rt *rootNode) isLeaf() bool {
	return false
}

func (rt *rootNode) GetSQLQuery() (string, []interface{}) {
	fmt.Println("root GetSQLQuery")
	var values []interface{}
	var conditions []string

	fmt.Println("root GetSQLQuery", rt.children)

	for _, child := range rt.children {
		c, v := child.GetSQLQuery()
		fmt.Println("root child c ", c, " v ", v)

		conditions = append(conditions, c)
		v = append(values, v...)
	}

	return strings.Join(conditions, " AND "), values
}

// / Operation -----------------------------------------------

type eqNode struct {
	op     string
	value  interface{}
	parent node
	isNot  bool
}

func (eq *eqNode) GetParent() node {
	return eq.parent
}

func (eq *eqNode) IsLeaf() bool {
	return true
}

func (eq *eqNode) AddChild(n node) {
}

func (eq *eqNode) GetSQLQuery() (string, []interface{}) {
	op := eq.op
	if eq.isNot {
		op = "!" + op
	}

	return fmt.Sprintf("%s ?", op), []interface{}{eq.value}
}

// / --------------------------------------------------------

// / ---------------------------------------------------------

type fieldNode struct {
	name     string
	parent   node
	children []node
}

func (fn *fieldNode) GetParent() node {
	return fn.parent
}

func (fn *fieldNode) IsLeaf() bool {
	return false
}

func (fn *fieldNode) AddChild(n node) {
	fn.children = append(fn.children, n)
}

func (fn *fieldNode) isLeaf() bool {
	return false
}

func (fn *fieldNode) GetSQLQuery() (string, []interface{}) {
	var values []interface{}

	var conditions []string

	for _, child := range fn.children {
		c, v := child.GetSQLQuery()

		str := fmt.Sprintf("%s %s", fn.name, c)

		conditions = append(conditions, str)
		v = append(values, v...)
	}

	return strings.Join(conditions, " AND "), values
}

type queryTree struct {
	obj                 interface{}
	filedTypeMap        map[string]reflect.Type
	allFieldMap         map[string]interface{}
	allSelectedFieldMap map[string]interface{}
	allFilteredFieldMap map[string]interface{}
	allSortedMap        map[string]interface{}
}

type QueryTree interface {
	Parse(query string) node
	buildTree(jsonBody string) node
}

func newQueryTree(obj interface{}) QueryTree {
	t := reflect.TypeOf(obj)

	filedTypeMap := map[string]reflect.Type{}
	allFieldMap := map[string]interface{}{}
	allSelectedFieldMap := map[string]interface{}{}
	allFilteredFieldMap := map[string]interface{}{}
	allSortedMap := map[string]interface{}{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// fmt.Printf("Name: %s, Type: %s\n", field.Name, field.Type)
		// fmt.Println(field.Type.Kind())
		filedTypeMap[field.Name] = field.Type
		qTag := field.Tag.Get("q")
		qTagValues := strings.Split(qTag, ";")

		nameIndex := -1
		for i, v := range qTagValues {
			if strings.Contains(v, "name") {
				nameIndex = i
				break
			}
		}
		// / Do When NameIndex > -1
		if nameIndex > -1 {
			// / Set Name To Map
			rawName := qTagValues[nameIndex]
			splitName := strings.Split(rawName, ":")
			if len(splitName) != 2 {
				panic("name must be `name:field_name`")
			}

			fieldName := strings.TrimSpace(splitName[1])

			allFieldMap[fieldName] = struct{}{}
			filedTypeMap[fieldName] = field.Type

			for i, v := range qTagValues {
				if i == nameIndex {
					continue
				}
				switch v {
				case "selected":
					allSelectedFieldMap[fieldName] = struct{}{}
				case "filtered":
					allFilteredFieldMap[fieldName] = struct{}{}
				case "sorted":
					allSortedMap[fieldName] = struct{}{}
				}
			}

		}
	}
	return &queryTree{
		obj:                 obj,
		filedTypeMap:        filedTypeMap,
		allFieldMap:         allFieldMap,
		allSelectedFieldMap: allSelectedFieldMap,
		allFilteredFieldMap: allFilteredFieldMap,
		allSortedMap:        allSortedMap,
	}
}

func (qt *queryTree) buildTree(jsonBody string) node {
	// check
	decoded, err := url.QueryUnescape(jsonBody)
	if err != nil {
		// Handle error
		panic("filter invalid")
	}

	var filter map[string]interface{}
	err = json.Unmarshal([]byte(decoded), &filter)
	if err != nil {
		// Handle error
		panic("filter invalid")
	}

	// create rootNode
	root := rootNode{
		children: []node{},
	}
	fmt.Println("create root node", root)

	for key, value := range filter {
		if _, isNotOP := notOperation[key]; isNotOP {
			fmt.Println("Built NotOP Node")
		} else if _, isOP := operation[key]; isOP {
			fmt.Println("Built OP Node")
		} else if _, isField := qt.allSelectedFieldMap[key]; isField {
			fmt.Println("Built field :", key, value)
			qt.buildFieldNode(&root, key, value)
		}
	}

	return &root
}

func (qt *queryTree) buildFieldNode(
	parent node,
	fieldName string,
	value interface{},
) node {
	// Build Node
	current := &fieldNode{
		parent:   parent,
		name:     fieldName,
		children: []node{},
	}

	// / check value type match with fieldName
	fieldType, _ := qt.filedTypeMap[fieldName]

	if match, reaValue := ParserValid.MatchType(fieldType, value, TimeLayout); match {
		fmt.Println("FieldNode' child is match Type will create eqOp Node ", reaValue)
		op := qt.buildOpNode("$eq", value, current, false)
		current.AddChild(op)
	} else {
		fmt.Println("check if  v == map[string]{}")

		// // /  check if  v == map[string]{}
		// b, err := json.Marshal(value)
		// // var ops map[string]interface{}
		// if err != nil {
		// 	// Handle error
		// 	panic("buildFieldNode Marshal error")
		// }
		// var ops map[string]interface{}
		// err = json.Unmarshal(b, &ops)
		// if err != nil {
		// 	panic("buildFieldNode Unmarshal error")
		// }
	}

	if len(current.children) > 0 {
		parent.AddChild(current)
	}
	return nil
}

func (qt *queryTree) buildOpNode(
	op string,
	value interface{},
	parent node,
	isNot bool,
) node {
	switch op {
	case "$eq":
		{
			node := eqNode{
				op:     "=",
				value:  value,
				parent: parent,
				isNot:  isNot,
			}
			return &node
		}
	}

	return nil
}

func (qt *queryTree) Parse(query string) node {
	fmt.Println("qt Parse query", query)
	splits := strings.Split(query, "&")
	fmt.Println("splits", splits)

	var root node

	for _, key := range splits {
		ext := strings.Split(key, "=")
		if len(ext) == 2 {
			switch ext[0] {
			case "filter":
				{
					fmt.Println("query have filter")
					root = qt.buildTree(ext[1])
				}
			}
		}
	}

	return root
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	type obj struct {
		ID        uint      `json:"id" q:"name:id;selected;filtered;sorted"`
		Email     string    `json:"email" q:"name:email;selected;filtered;sorted"`
		ENUM      string    `json:"enum" q:"name:enum;selected;filtered;sorted"`
		Activate  bool      `json:"active,omitempty" q:"name:active;selected;filtered;sorted"`
		CreatedAt time.Time `json:"created_at,omitempty" q:"name:created_at;selected;filtered;sorted"`
	}
	qt := newQueryTree(obj{})
	// fmt.Println(qt)

	r.GET("/demo", func(c *gin.Context) {
		queryString := c.Request.URL.RawQuery

		root := qt.Parse(queryString)

		conditions, _ := root.GetSQLQuery()

		c.JSON(http.StatusOK, gin.H{
			"message":     "pong",
			"queryString": queryString,
			"conditions":  conditions,
		})
	})

	r.Run(":3000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
