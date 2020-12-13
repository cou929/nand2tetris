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
	SetMeta(parent NodeType, grandParent NodeType, s *SymbolTableEntry) error
	Meta() *IDMeta
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

func (n *InnerNode) SetMeta(parent NodeType, grandParent NodeType, s *SymbolTableEntry) error {
	return fmt.Errorf("InnerNode does not supported SetMeta")
}

func (n *InnerNode) Meta() *IDMeta {
	return nil
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

func NewTypeNode() *OneChildNode {
	return NewOneChildNode(TypeType, "type", false)
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

func NewClassNameNode() *OneChildNode {
	return NewOneChildNode(ClassNameType, "className", false)
}

func NewSubroutineNameNode() *OneChildNode {
	return NewOneChildNode(SubroutineNameType, "subroutineName", false)
}

func NewVarNameNode() *OneChildNode {
	return NewOneChildNode(VarNameType, "varName", false)
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

func NewOpNode() *OneChildNode {
	return NewOneChildNode(OpType, "op", false)
}

func NewUnaryOpNode() *OneChildNode {
	return NewOneChildNode(UnaryOpType, "unaryOp", false)
}

func NewKeywordConstantNode() *OneChildNode {
	return NewOneChildNode(KeywordConstantType, "keywordConstant", false)
}

type OneChildNode struct {
	Typ       NodeType
	N         string
	V         string
	Children  []TreeNode
	XMLMarkup bool
}

func NewOneChildNode(typ NodeType, name string, x bool) *OneChildNode {
	return &OneChildNode{
		Typ:       typ,
		N:         name,
		XMLMarkup: x,
	}
}

func (n *OneChildNode) Name() string {
	return n.N
}

func (n *OneChildNode) Type() NodeType {
	return n.Typ
}

func (n *OneChildNode) AppendChild(node TreeNode) {
	n.Children = []TreeNode{node} // must have only one child
}

func (n *OneChildNode) ChildNodes() []TreeNode {
	return n.Children
}

func (n *OneChildNode) SetValue(v string) {
	// do nothing
}

func (n *OneChildNode) Value() string {
	if len(n.Children) < 1 {
		return ""
	}
	return n.Children[0].Value()
}

func (n *OneChildNode) SetMeta(parent NodeType, grandParent NodeType, s *SymbolTableEntry) error {
	return fmt.Errorf("OneChildNode does not supported SetMeta")
}

func (n *OneChildNode) Meta() *IDMeta {
	return n.Children[0].Meta()
}

func (n *OneChildNode) Xml() string {
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

type LeafNode struct {
	Typ       NodeType
	N         string
	V         string
	XMLMarkup bool
	IDMeta    *IDMeta
}

type IDMeta struct {
	Category    IdCategory
	Declaration bool
	SymbolInfo  *SymbolInfo
}

type IdCategory int

const (
	IdCatVar IdCategory = iota + 1
	IdCatArg
	IdCatStatic
	IdCatField
	IdCatClass
	IdCatSub
)

func (i IdCategory) String() string {
	switch i {
	case IdCatVar:
		return "IdCatVar"
	case IdCatArg:
		return "IdCatArg"
	case IdCatStatic:
		return "IdCatStatic"
	case IdCatField:
		return "IdCatField"
	case IdCatClass:
		return "IdCatClass"
	case IdCatSub:
		return "IdCatSub"
	}
	return "Invalid"
}

type SymbolInfo struct {
	Kind  VarKind
	Type  string
	Index int
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

func (n *LeafNode) SetMeta(parent NodeType, grandParent NodeType, s *SymbolTableEntry) error {
	if n.Type() != IdentifierType {
		return fmt.Errorf("SetMeta is only for IdentifierType %v", n.Value())
	}
	if parent != ClassNameType && parent != SubroutineNameType && parent != VarNameType {
		return fmt.Errorf("Invalid parent type %v %v %v", n.Value(), n.Type(), parent)
	}

	meta := &IDMeta{}

	// Category
	switch s == nil {
	case true:
		switch parent {
		case ClassNameType:
			meta.Category = IdCatClass
		case SubroutineNameType:
			meta.Category = IdCatSub
		}
	case false:
		switch s.Kind {
		case Static:
			meta.Category = IdCatStatic
		case Field:
			meta.Category = IdCatField
		case Argument:
			meta.Category = IdCatArg
		case Var:
			meta.Category = IdCatVar
		}
	}

	// Declaration
	meta.Declaration = false
	switch grandParent {
	case ClassType, ClassVarDecType, SubroutineDecType, ParameterListType, VarDecType:
		meta.Declaration = true
	}

	// SymbolInfo
	if s != nil {
		meta.SymbolInfo = &SymbolInfo{
			Kind:  s.Kind,
			Type:  s.Typ,
			Index: s.Index,
		}
	}

	n.IDMeta = meta

	return nil
}

func (n *LeafNode) Meta() *IDMeta {
	return n.IDMeta
}

func (n *LeafNode) Xml() string {
	if idAttr && n.Type() == IdentifierType {
		if n.IDMeta.SymbolInfo != nil {
			return fmt.Sprintf("<%s category=\"%s\" declaration=\"%t\" kind=\"%s\" type=\"%s\" index=\"%d\">%s</%s>", n.Name(), n.IDMeta.Category, n.IDMeta.Declaration, n.IDMeta.SymbolInfo.Kind, n.IDMeta.SymbolInfo.Type, n.IDMeta.SymbolInfo.Index, escapeXml(n.Value()), n.Name())
		}
		return fmt.Sprintf("<%s category=\"%s\" declaration=\"%t\">%s</%s>", n.Name(), n.IDMeta.Category, n.IDMeta.Declaration, escapeXml(n.Value()), n.Name())
	}
	return fmt.Sprintf("<%s>%s</%s>", n.Name(), escapeXml(n.Value()), n.Name())
}

func AdaptTokenToNode(token Token) TreeNode {
	node := NewLeafNode(token.Type(), token.Name(), true)
	node.SetValue(token.String())
	return node
}
