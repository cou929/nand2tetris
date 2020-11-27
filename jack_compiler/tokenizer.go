package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type TokenType int

const (
	KeywordType TokenType = iota + 1
	SymbolType
	IntConstType
	StrConstType
	IdentifierType
)

type Token interface {
	Type() TokenType
	String() string
	Int() int
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

func (t KeywordToken) Type() TokenType {
	return KeywordType
}

func (t KeywordToken) String() string {
	return string(t)
}

func (t KeywordToken) Int() int {
	return 0
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

func (t SymbolToken) Type() TokenType {
	return SymbolType
}

func (t SymbolToken) String() string {
	return string(t)
}

func (t SymbolToken) Int() int {
	return 0
}

type IntConstToken int

func NewIntConstToken(in int) (IntConstToken, bool) {
	if 0 <= in || in <= 32767 {
		return IntConstToken(in), true
	}
	return 0, false
}

func (t IntConstToken) Type() TokenType {
	return IntConstType
}

func (t IntConstToken) String() string {
	return ""
}

func (t IntConstToken) Int() int {
	return int(t)
}

type StrConstToken string

func NewStrConstToken(in string) (StrConstToken, bool) {
	valid := regexp.MustCompile(`^[^"\n]+$`)
	if valid.MatchString(in) {
		return StrConstToken(in), true
	}
	return "", false
}

func (t StrConstToken) Type() TokenType {
	return StrConstType
}

func (t StrConstToken) String() string {
	return string(t)
}

func (t StrConstToken) Int() int {
	return 0
}

type IdentifierToken string

func NewIdentifierToken(in string) (IdentifierToken, bool) {
	valid := regexp.MustCompile(`^[a-zA-Z_]([a-zA-Z0-9_]+)?$`)
	if valid.MatchString(in) {
		return IdentifierToken(in), true
	}
	return "", false
}

func (t IdentifierToken) Type() TokenType {
	return IdentifierType
}

func (t IdentifierToken) String() string {
	return string(t)
}

func (t IdentifierToken) Int() int {
	return 0
}

type tokenizeState int

const (
	ordinal            tokenizeState = iota + 1
	multiCommentOpened               // /* */
	stringOpened                     // "blah"
)

type Tokenizer struct {
	reader io.Reader
	tokens []Token
	state  tokenizeState
	buf    string
}

func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{
		reader: r,
		state:  ordinal,
	}
}

func (t *Tokenizer) Tokenize() ([]Token, error) {
	var res []Token

	scanner := bufio.NewScanner(t.reader)
	for scanner.Scan() {
		l := scanner.Text()
		tokens, err := t.parseLine(l)
		if err != nil {
			return nil, err
		}
		if tokens == nil {
			continue
		}
		res = append(res, tokens...)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	t.tokens = res
	return res, nil
}

func (t *Tokenizer) parseLine(l string) ([]Token, error) {
	var res []Token
	runes := []rune(l)

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if t.state == ordinal {
			if t.singleComment(runes, i) {
				tkn, err := t.bufToToken()
				if err != nil {
					return nil, err
				}
				if tkn != nil {
					res = append(res, tkn)
				}

				break
			}

			if t.multiCommentOpen(runes, i) {
				tkn, err := t.bufToToken()
				if err != nil {
					return nil, err
				}
				if tkn != nil {
					res = append(res, tkn)
				}
				t.clearBuf()

				if err := t.transit(multiCommentOpened); err != nil {
					return nil, fmt.Errorf("Invalid multi comment opening %s. %w", l, err)
				}
				i++
				continue
			}

			if t.stringQuote(r) {
				tkn, err := t.bufToToken()
				if err != nil {
					return nil, err
				}
				if tkn != nil {
					res = append(res, tkn)
				}
				t.clearBuf()

				if err := t.transit(stringOpened); err != nil {
					return nil, fmt.Errorf("Invalid string opening %s. %w", l, err)
				}
				continue
			}

			if t.delim(r) {
				tkn, err := t.bufToToken()
				if err != nil {
					return nil, err
				}
				if tkn != nil {
					res = append(res, tkn)
				}
				t.clearBuf()

				sym, ok := NewSymbolToken(string(r))
				if ok {
					res = append(res, sym)
				}

				continue
			}

			t.appendBuf(r)
		}

		if t.state == multiCommentOpened {
			if t.multiCommentClose(runes, i) {
				if err := t.transit(ordinal); err != nil {
					return nil, fmt.Errorf("Invalid multi comment closing %s. %w", l, err)
				}
				i++
				continue
			}
			continue
		}

		if t.state == stringOpened {
			if t.stringQuote(r) {
				tkn, err := t.bufToToken()
				if err != nil {
					return nil, err
				}
				if tkn != nil {
					res = append(res, tkn)
				}
				t.clearBuf()

				if err := t.transit(ordinal); err != nil {
					return nil, fmt.Errorf("Invalid multi comment closing %s. %w", l, err)
				}
				continue
			}

			t.appendBuf(r)
		}
	}

	return res, nil
}

func (t *Tokenizer) appendBuf(c rune) {
	t.buf += string(c)
}

func (t *Tokenizer) clearBuf() {
	t.buf = ""
}

func (t *Tokenizer) bufToToken() (Token, error) {
	b := t.buf

	if b == "" {
		return nil, nil
	}

	kw, ok := NewKeywordToken(b)
	if ok {
		return kw, nil
	}

	i, err := strconv.Atoi(b)
	if err == nil {
		it, ok := NewIntConstToken(i)
		if ok {
			return it, nil
		}
	}

	if t.state == ordinal {
		id, ok := NewIdentifierToken(b)
		if ok {
			return id, nil
		}
	}

	if t.state == stringOpened {
		st, ok := NewStrConstToken(b)
		if ok {
			return st, nil
		}
	}

	return nil, fmt.Errorf("Invalid token %s", t.buf)
}

func (t *Tokenizer) singleComment(runes []rune, idx int) bool {
	if idx+1 >= len(runes) {
		return false
	}
	if runes[idx] != '/' || runes[idx+1] != '/' {
		return false
	}
	return true
}

func (t *Tokenizer) multiCommentOpen(runes []rune, idx int) bool {
	if idx+1 >= len(runes) {
		return false
	}
	if runes[idx] != '/' || runes[idx+1] != '*' {
		return false
	}
	return true
}

func (t *Tokenizer) multiCommentClose(runes []rune, idx int) bool {
	if idx+1 >= len(runes) {
		return false
	}
	if runes[idx] != '*' || runes[idx+1] != '/' {
		return false
	}
	return true
}

func (t *Tokenizer) stringQuote(r rune) bool {
	if r == '"' {
		return true
	}
	return false
}

func (t *Tokenizer) delim(c rune) bool {
	if t.state != ordinal {
		return false
	}
	if unicode.IsSpace(c) {
		return true
	}
	if _, ok := NewSymbolToken(string(c)); ok {
		return true
	}
	return false
}

func (t *Tokenizer) transit(next tokenizeState) error {
	if !t.canTransit(next) {
		return fmt.Errorf("Invalid status transition from %v to %v", t.state, next)
	}
	t.state = next
	return nil
}

func (t *Tokenizer) canTransit(next tokenizeState) bool {
	cur := t.state
	switch cur {
	case ordinal:
		switch next {
		case multiCommentOpened, stringOpened:
			return true
		}
	case multiCommentOpened:
		switch next {
		case ordinal:
			return true
		}
	case stringOpened:
		switch next {
		case ordinal:
			return true
		}
	}
	return false
}

func (t *Tokenizer) Xml() string {
	if len(t.tokens) <= 0 {
		return ""
	}

	var res []string
	res = append(res, `<tokens>`)
	for _, tkn := range t.tokens {
		switch tkn.Type() {
		case KeywordType:
			res = append(res, fmt.Sprintf("<keyword>%s</keyword>", t.escapeXml(tkn.String())))
		case SymbolType:
			res = append(res, fmt.Sprintf("<symbol>%s</symbol>", t.escapeXml(tkn.String())))
		case IntConstType:
			res = append(res, fmt.Sprintf("<integerConstant>%d</integerConstant>", tkn.Int()))
		case StrConstType:
			res = append(res, fmt.Sprintf("<stringConstant>%s</stringConstant>", t.escapeXml(tkn.String())))
		case IdentifierType:
			res = append(res, fmt.Sprintf("<identifier>%s</identifier>", t.escapeXml(tkn.String())))
		}
	}
	res = append(res, `</tokens>`)

	return strings.Join(res, "\n")
}

func (t *Tokenizer) escapeXml(in string) string {
	a := strings.ReplaceAll(in, "&", "&amp;")
	b := strings.ReplaceAll(a, "<", "&lt;")
	c := strings.ReplaceAll(b, ">", "&gt;")
	return c
}
