package main

import (
	"fmt"
	"strings"
)

// TreeNode is non-terminal symbols.
// Inner node of parse tree, in other words.
// Satisfies Token interface.
type TreeNode interface {
	AppendChild(token Token)
	ChildNodes() []Token
	Token
}

type GeneralNode struct {
	Children  []Token
	Typ       NodeType
	XMLHeader string
}

func NewGeneralNode(Typ NodeType, xh string) *GeneralNode {
	return &GeneralNode{
		Typ:       Typ,
		XMLHeader: xh,
	}
}

func (n *GeneralNode) Type() NodeType {
	return n.Typ
}

func (n *GeneralNode) AppendChild(tkn Token) {
	n.Children = append(n.Children, tkn)
}

func (n *GeneralNode) ChildNodes() []Token {
	return n.Children
}

func (n *GeneralNode) Xml() string {
	res := []string{}
	if n.XMLHeader != "" {
		res = append(res, fmt.Sprintf("<%s>", n.XMLHeader))
	}
	for _, c := range n.Children {
		res = append(res, c.Xml())
	}
	if n.XMLHeader != "" {
		res = append(res, fmt.Sprintf("</%s>", n.XMLHeader))
	}
	return strings.Join(res, "\n")
}

func (n *GeneralNode) String() string {
	switch n.Type() {
	case ClassNameType, SubroutineNameType, VarDecType, OpType, UnaryOpType, KeyConstType:
		return n.ChildNodes()[0].String()
	}
	return ""
}

func NewClassNode() *GeneralNode {
	return NewGeneralNode(ClassType, "class")
}

func NewClassVarDecNode() *GeneralNode {
	return NewGeneralNode(ClassVarDecType, "classVarDec")
}

func NewTypeNode() *GeneralNode {
	return NewGeneralNode(TypeType, "")
}

func NewSubroutineDecNode() *GeneralNode {
	return NewGeneralNode(SubroutineDecType, "subroutineDec")
}

func NewParameterListNode() *GeneralNode {
	return NewGeneralNode(ParameterListType, "parameterList")
}

func NewSubroutineBodyNode() *GeneralNode {
	return NewGeneralNode(SubroutineBodyType, "subroutineBody")
}

func NewVarDecNode() *GeneralNode {
	return NewGeneralNode(VarDecType, "varDec")
}

func NewClassNameNode() *GeneralNode {
	return NewGeneralNode(ClassNameType, "")
}

func NewSubroutineNameNode() *GeneralNode {
	return NewGeneralNode(SubroutineNameType, "")
}

func NewVarNameNode() *GeneralNode {
	return NewGeneralNode(VarNameType, "")
}

func NewStatementsNode() *GeneralNode {
	return NewGeneralNode(StatementsType, "statements")
}

func NewStatementNode() *GeneralNode {
	return NewGeneralNode(StatementType, "")
}

func NewLetStatementNode() *GeneralNode {
	return NewGeneralNode(LetStatementType, "letStatement")
}

func NewIfStatementNode() *GeneralNode {
	return NewGeneralNode(IfStatementType, "ifStatement")
}

func NewWhileStatementNode() *GeneralNode {
	return NewGeneralNode(WhileStatementType, "whileStatement")
}

func NewDoStatementNode() *GeneralNode {
	return NewGeneralNode(DoStatementType, "doStatement")
}

func NewReturnStatementNode() *GeneralNode {
	return NewGeneralNode(ReturnStatementType, "returnStatement")
}

func NewExpressionNode() *GeneralNode {
	return NewGeneralNode(ExpressionType, "expression")
}

func NewTermNode() *GeneralNode {
	return NewGeneralNode(TermType, "term")
}

func NewSubroutineCallNode() *GeneralNode {
	return NewGeneralNode(SubroutineCallType, "")
}

func NewExpressionListNode() *GeneralNode {
	return NewGeneralNode(ExpressionListType, "expressionList")
}

func NewOpNode() *GeneralNode {
	return NewGeneralNode(OpType, "")
}

func NewUnaryOpNode() *GeneralNode {
	return NewGeneralNode(UnaryOpType, "")
}

func NewKeyConstNode() *GeneralNode {
	return NewGeneralNode(KeyConstType, "")
}
