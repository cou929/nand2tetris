package main

import (
	"fmt"
	"regexp"
	"strconv"
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
	String() string
	Name() string
}

type KeywordToken string

func NewKeywordToken(in string) (KeywordToken, bool) {
	switch in {
	case "class", "constructor", "function", "method", "field", "static", "var", "int", "char", "boolean",
		"void", "true", "false", "null", "this", "let", "do", "if", "else", "while", "return":
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

func (t KeywordToken) Int() int {
	return 0
}

func (t KeywordToken) Xml() string {
	return fmt.Sprintf("<%s>%s</%s>", t.Name(), escapeXml(t.String()), t.Name())
}

func (t KeywordToken) Name() string {
	return "keyword"
}

type SymbolToken string

func NewSymbolToken(in string) (SymbolToken, bool) {
	switch in {
	case "{", "}", "(", ")", "[", "]", ".", ",", ";", "+", "-", "*", "/", "&", "|", "<", ">", "=", "~":
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

func (t SymbolToken) Int() int {
	return 0
}

func (t SymbolToken) Xml() string {
	return fmt.Sprintf("<%s>%s</%s>", t.Name(), escapeXml(t.String()), t.Name())
}

func (t SymbolToken) Name() string {
	return "symbol"
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

func (t IntConstToken) String() string {
	return strconv.Itoa(int(t))
}

func (t IntConstToken) Int() int {
	return int(t)
}

func (t IntConstToken) Xml() string {
	return fmt.Sprintf("<%s>%d</%s>", t.Name(), t.Int(), t.Name())
}

func (t IntConstToken) Name() string {
	return "integerConstant"
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

func (t StrConstToken) Int() int {
	return 0
}

func (t StrConstToken) Xml() string {
	return fmt.Sprintf("<%s>%s</%s>", t.Name(), escapeXml(t.String()), t.Name())
}

func (t StrConstToken) Name() string {
	return "stringConstant"
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

func (t IdentifierToken) Int() int {
	return 0
}

func (t IdentifierToken) Xml() string {
	return fmt.Sprintf("<%s>%s</%s>", t.Name(), escapeXml(t.String()), t.Name())
}

func (t IdentifierToken) Name() string {
	return "identifier"
}

func escapeXml(in string) string {
	a := strings.ReplaceAll(in, "&", "&amp;")
	b := strings.ReplaceAll(a, "<", "&lt;")
	c := strings.ReplaceAll(b, ">", "&gt;")
	return c
}
