package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"unicode"
)

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

func (t *Tokenizer) Tokenize() (Tokens, error) {
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
	return NewTokens(res), nil
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

	i, err := strconv.Atoi(b)
	if err == nil {
		it, ok := NewIntConstToken(i)
		if ok {
			return it, nil
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
