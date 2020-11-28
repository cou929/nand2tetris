package main

import (
	"fmt"
	"regexp"
	"strings"
)

type Tokens []Token

func NewTokens(tokens []Token) Tokens {
	return Tokens(tokens)
}

func (t *Tokens) Xml() string {
	res := []string{`<tokens>`}
	for _, tkn := range *t {
		res = append(res, tkn.Xml())
	}
	res = append(res, `</tokens>`)
	return strings.Join(res, "\n")
}

type Token interface {
	Type() NodeType
	Xml() string
}

type KeywordToken string

func NewKeywordToken(in string) (KeywordToken, bool) {
	switch in {
	case "class":
		fallthrough
	case "constructor":
		fallthrough
	case "function":
		fallthrough
	case "method":
		fallthrough
	case "field":
		fallthrough
	case "static":
		fallthrough
	case "var":
		fallthrough
	case "int":
		fallthrough
	case "char":
		fallthrough
	case "boolean":
		fallthrough
	case "void":
		fallthrough
	case "true":
		fallthrough
	case "false":
		fallthrough
	case "null":
		fallthrough
	case "this":
		fallthrough
	case "let":
		fallthrough
	case "do":
		fallthrough
	case "if":
		fallthrough
	case "else":
		fallthrough
	case "while":
		fallthrough
	case "return":
		return KeywordToken(in), true
	}

	return "", false
}

func (t KeywordToken) Type() NodeType {
	return KeywordType
}

func (t KeywordToken) String() string {
	return string(t)
}

func (t KeywordToken) Xml() string {
	return fmt.Sprintf("<keyword>%s</keyword>", escapeXml(t.String()))
}

type SymbolToken string

func NewSymbolToken(in string) (SymbolToken, bool) {
	switch in {
	case "{":
		fallthrough
	case "}":
		fallthrough
	case "(":
		fallthrough
	case ")":
		fallthrough
	case "[":
		fallthrough
	case "]":
		fallthrough
	case ".":
		fallthrough
	case ",":
		fallthrough
	case ";":
		fallthrough
	case "+":
		fallthrough
	case "-":
		fallthrough
	case "*":
		fallthrough
	case "/":
		fallthrough
	case "&":
		fallthrough
	case "|":
		fallthrough
	case "<":
		fallthrough
	case ">":
		fallthrough
	case "=":
		fallthrough
	case "~":
		return SymbolToken(in), true
	}

	return "", false
}

func (t SymbolToken) Type() NodeType {
	return SymbolType
}

func (t SymbolToken) String() string {
	return string(t)
}

func (t SymbolToken) Xml() string {
	return fmt.Sprintf("<symbol>%s</symbol>", escapeXml(t.String()))
}

type IntConstToken int

func NewIntConstToken(in int) (IntConstToken, bool) {
	if 0 <= in || in <= 32767 {
		return IntConstToken(in), true
	}
	return 0, false
}

func (t IntConstToken) Type() NodeType {
	return IntConstType
}

func (t IntConstToken) Int() int {
	return int(t)
}

func (t IntConstToken) Xml() string {
	return fmt.Sprintf("<integerConstant>%d</integerConstant>", t.Int())
}

type StrConstToken string

func NewStrConstToken(in string) (StrConstToken, bool) {
	valid := regexp.MustCompile(`^[^"\n]+$`)
	if valid.MatchString(in) {
		return StrConstToken(in), true
	}
	return "", false
}

func (t StrConstToken) Type() NodeType {
	return StrConstType
}

func (t StrConstToken) String() string {
	return string(t)
}

func (t StrConstToken) Xml() string {
	return fmt.Sprintf("<stringConstant>%s</stringConstant>", escapeXml(t.String()))
}

type IdentifierToken string

func NewIdentifierToken(in string) (IdentifierToken, bool) {
	valid := regexp.MustCompile(`^[a-zA-Z_]([a-zA-Z0-9_]+)?$`)
	if valid.MatchString(in) {
		return IdentifierToken(in), true
	}
	return "", false
}

func (t IdentifierToken) Type() NodeType {
	return IdentifierType
}

func (t IdentifierToken) String() string {
	return string(t)
}

func (t IdentifierToken) Xml() string {
	return fmt.Sprintf("<identifier>%s</identifier>", escapeXml(t.String()))
}

func escapeXml(in string) string {
	a := strings.ReplaceAll(in, "&", "&amp;")
	b := strings.ReplaceAll(a, "<", "&lt;")
	c := strings.ReplaceAll(b, ">", "&gt;")
	return c
}
