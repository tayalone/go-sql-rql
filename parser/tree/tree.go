package tree

// Node of Query Tree
type Node interface {
	GetParent() Node
	IsLeaf() bool
	AddChild(n Node)
	GetSQLQuery() (string, []interface{})
	MergeChildren(children []Node)
	GetChildren() []Node
}
