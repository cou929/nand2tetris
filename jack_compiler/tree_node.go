package main

import (
	"fmt"
	"strings"
)

// TreeNode represents all nodes (both inner node and leaf node) of the parse tree.
type TreeNode interface {
	Name() string
	Type() NodeType
	AppendChild(node TreeNode)
	ChildNodes() []TreeNode
	SetValue(v string)
	Value() string
	SetMeta()
	Xml() string
}

type InnerNode struct {
	Children  []TreeNode
	Typ       NodeType
	N         string
	XMLMarkup bool
}

func NewInnerNode(typ NodeType, name string, x bool) *InnerNode {
	return &InnerNode{
		Typ:       typ,
		N:         name,
		XMLMarkup: x,
	}
}

func (n *InnerNode) Name() string {
	return n.N
}

func (n *InnerNode) Type() NodeType {
	return n.Typ
}

func (n *InnerNode) AppendChild(node TreeNode) {
	n.Children = append(n.Children, node)
}

func (n *InnerNode) ChildNodes() []TreeNode {
	return n.Children
}

func (n *InnerNode) SetValue(v string) {
	// do nothing
}

func (n *InnerNode) Value() string {
	return ""
}

func (n *InnerNode) SetMeta() {
}

func (n *InnerNode) Xml() string {
	res := []string{}
	if n.XMLMarkup {
		res = append(res, fmt.Sprintf("<%s>", n.Name()))
	}
	for _, c := range n.ChildNodes() {
		res = append(res, c.Xml())
	}
	if n.XMLMarkup {
		res = append(res, fmt.Sprintf("</%s>", n.Name()))
	}
	return strings.Join(res, "\n")
}

func NewClassNode() *InnerNode {
	return NewInnerNode(ClassType, "class", true)
}

func NewClassVarDecNode() *InnerNode {
	return NewInnerNode(ClassVarDecType, "classVarDec", true)
}

func NewTypeNode() *InnerNode {
	return NewInnerNode(TypeType, "type", false)
}

func NewSubroutineDecNode() *InnerNode {
	return NewInnerNode(SubroutineDecType, "subroutineDec", true)
}

func NewParameterListNode() *InnerNode {
	return NewInnerNode(ParameterListType, "parameterList", true)
}

func NewSubroutineBodyNode() *InnerNode {
	return NewInnerNode(SubroutineBodyType, "subroutineBody", true)
}

func NewVarDecNode() *InnerNode {
	return NewInnerNode(VarDecType, "varDec", true)
}

func NewClassNameNode() *InnerNode {
	return NewInnerNode(ClassNameType, "className", false)
}

func NewSubroutineNameNode() *InnerNode {
	return NewInnerNode(SubroutineNameType, "subroutineName", false)
}

func NewVarNameNode() *InnerNode {
	return NewInnerNode(VarNameType, "varName", false)
}

func NewStatementsNode() *InnerNode {
	return NewInnerNode(StatementsType, "statements", true)
}

func NewStatementNode() *InnerNode {
	return NewInnerNode(StatementType, "statement", false)
}

func NewLetStatementNode() *InnerNode {
	return NewInnerNode(LetStatementType, "letStatement", true)
}

func NewIfStatementNode() *InnerNode {
	return NewInnerNode(IfStatementType, "ifStatement", true)
}

func NewWhileStatementNode() *InnerNode {
	return NewInnerNode(WhileStatementType, "whileStatement", true)
}

func NewDoStatementNode() *InnerNode {
	return NewInnerNode(DoStatementType, "doStatement", true)
}

func NewReturnStatementNode() *InnerNode {
	return NewInnerNode(ReturnStatementType, "returnStatement", true)
}

func NewExpressionNode() *InnerNode {
	return NewInnerNode(ExpressionType, "expression", true)
}

func NewTermNode() *InnerNode {
	return NewInnerNode(TermType, "term", true)
}

func NewSubroutineCallNode() *InnerNode {
	return NewInnerNode(SubroutineCallType, "subroutineCall", false)
}

func NewExpressionListNode() *InnerNode {
	return NewInnerNode(ExpressionListType, "expressionList", true)
}

func NewOpNode() *InnerNode {
	return NewInnerNode(OpType, "op", false)
}

func NewUnaryOpNode() *InnerNode {
	return NewInnerNode(UnaryOpType, "unaryOp", false)
}

func NewKeywordConstantNode() *InnerNode {
	return NewInnerNode(KeywordConstantType, "keywordConstant", false)
}

type LeafNode struct {
	Typ       NodeType
	N         string
	V         string
	XMLMarkup bool
}

func NewLeafNode(typ NodeType, name string, x bool) *LeafNode {
	return &LeafNode{
		Typ:       typ,
		N:         name,
		XMLMarkup: x,
	}
}

func (n *LeafNode) Name() string {
	return n.N
}

func (n *LeafNode) Type() NodeType {
	return n.Typ
}

func (n *LeafNode) AppendChild(node TreeNode) {
	// do nothing
}

func (n *LeafNode) ChildNodes() []TreeNode {
	return nil
}

func (n *LeafNode) SetValue(v string) {
	n.V = v
}

func (n *LeafNode) Value() string {
	return n.V
}

func (n *LeafNode) SetMeta() {
}

func (n *LeafNode) Xml() string {
	return fmt.Sprintf("<%s>%s</%s>", n.Name(), escapeXml(n.Value()), n.Name())
}

func AdaptTokenToNode(token Token) TreeNode {
	node := NewLeafNode(token.Type(), token.Name(), true)
	node.SetValue(token.String())
	return node
}
