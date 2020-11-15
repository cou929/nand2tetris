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
	reader io.Reader
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{
		reader: reader,
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
	tokens := spacePattern.Split(trimmed, -1)

	c, err := NewCommand(tokens)
	if err != nil {
		return nil, err
	}

	return c, nil
}
