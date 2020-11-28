package main

type ParseTree struct {
	Node     Token
	Children []ParseTree
}

func (pt *ParseTree) Xml() string {
	return ""
}
