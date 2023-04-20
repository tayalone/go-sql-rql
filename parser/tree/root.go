package tree

import (
	"fmt"
	"strings"
)

// Root of query tree
type Root struct {
	children []Node
}

// GetParent of root
func (rt *Root) GetParent() Node {
	return nil
}

// IsLeaf of root
func (rt *Root) IsLeaf() bool {
	return false
}

// AddChild loot
func (rt *Root) AddChild(n Node) {
	rt.children = append(rt.children, n)
}

// GetSQLQuery from root
func (rt *Root) GetSQLQuery() (string, []interface{}) {
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

// MergeChildren from root
func (rt *Root) MergeChildren(children []Node) {
	rt.children = append(rt.children, children...)
}

// GetChildren from root
func (rt *Root) GetChildren() []Node {
	return rt.children
}
