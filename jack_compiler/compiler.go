package main

import (
	"fmt"
	"strings"
)

type Compiler struct{}

func NewCompiler() *Compiler {
	return &Compiler{}
}

func (c *Compiler) Compile(pt TreeNode) (string, error) {
	codes, err := c.compile(pt)
	if err != nil {
		return "", err
	}
	return strings.Join(codes, "\n"), nil
}

func (c *Compiler) compile(pt TreeNode) ([]string, error) {
	switch pt.Type() {
	case ClassType:
		return c.compileClass(pt)
	case SubroutineDecType:
		return c.compileSubroutineDec(pt)
	case SubroutineBodyType:
		return c.compileSubroutineBody(pt)
	case StatementsType:
		return c.compileStatements(pt)
	case DoStatementType:
		return c.compileDoStatement(pt)
	case ReturnStatementType:
		return c.compileReturnStatement(pt)
	case IdentifierType:
		return c.compileIdentifier(pt)
	case ExpressionListType:
		return c.compileExpressionList(pt)
	case ExpressionType:
		return c.compileExpression(pt)
	}
	return []string{}, fmt.Errorf("Not supported %v", pt.Type())
}

func (c *Compiler) compileClass(pt TreeNode) ([]string, error) {
	return []string{}, nil
}

func (c *Compiler) compileSubroutineDec(pt TreeNode) ([]string, error) {
	return []string{}, nil
}

func (c *Compiler) compileSubroutineBody(pt TreeNode) ([]string, error) {
	return []string{}, nil
}

func (c *Compiler) compileStatements(pt TreeNode) ([]string, error) {
	return []string{}, nil
}

func (c *Compiler) compileDoStatement(pt TreeNode) ([]string, error) {
	return []string{}, nil
}

func (c *Compiler) compileReturnStatement(pt TreeNode) ([]string, error) {
	return []string{}, nil
}

func (c *Compiler) compileIdentifier(pt TreeNode) ([]string, error) {
	return []string{}, nil
}

func (c *Compiler) compileExpressionList(pt TreeNode) ([]string, error) {
	return []string{}, nil
}

func (c *Compiler) compileExpression(pt TreeNode) ([]string, error) {
	return []string{}, nil
}
