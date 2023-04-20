package field

import (
	"fmt"
	"strings"

	Tree "github.com/tayalone/go-sql-rql/parser/tree"
)

// Node of Field
type Node struct {
	name     string
	parent   Tree.Node
	children []Tree.Node
}

// GetParent of fields name
func (n *Node) GetParent() Tree.Node {
	return n.parent
}

// IsLeaf of fields name
func (n *Node) IsLeaf() bool {
	return false
}

// AddChild fields name
func (n *Node) AddChild(c Tree.Node) {
	n.children = append(n.children, c)
}

// GetSQLQuery fields name
func (n *Node) GetSQLQuery() (string, []interface{}) {
	var values []interface{}

	var conditions []string

	for _, child := range n.children {
		c, v := child.GetSQLQuery()

		str := fmt.Sprintf("%s %s", n.name, c)

		conditions = append(conditions, str)
		v = append(values, v...)
	}

	return strings.Join(conditions, " AND "), values
}

// MergeChildren from fields name
func (n *Node) MergeChildren(children []Node) {
}

// GetChildren from fields name
func (n *Node) GetChildren() []Node {
	return nil
}
