package main

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

var commentPattern = regexp.MustCompile(`\/\/.*?$`)
var spacePattern = regexp.MustCompile(`\s+`)

type Parser struct {
	reader   io.Reader
	fileName string
	curFunc  string
	curLine  int
}

func NewParser(reader io.Reader, n string) *Parser {
	return &Parser{
		reader:   reader,
		fileName: n,
	}
}

func (p *Parser) Parse() ([]*Command, error) {
	var res []*Command

	scanner := bufio.NewScanner(p.reader)
	for scanner.Scan() {
		l := scanner.Text()
		c, err := p.parseLine(l)
		if err != nil {
			return nil, err
		}
		if c == nil {
			continue
		}

		p.curLine++
		if c.Type == CommandFunction {
			p.curFunc = string(c.Arg1)
		}
		c.SetMeta(p.fileName, p.curFunc, p.curLine)
		if c.Type == CommandReturn {
			p.curFunc = ""
		}

		res = append(res, c)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (p *Parser) parseLine(line string) (*Command, error) {
	commentRemoved := commentPattern.ReplaceAllString(line, "")
	trimmed := strings.TrimSpace(commentRemoved)

	if trimmed == "" {
		return nil, nil
	}

	tokens := spacePattern.Split(trimmed, -1)

	c, err := NewCommand(tokens)
	if err != nil {
		return nil, err
	}

	return c, nil
}
