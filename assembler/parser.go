package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"unicode"
)

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

type lineParsingStatus int

const (
	initialized = iota + 1
	openedACommand
	openedCCommand
	closedDest
	closedComp
	openedLCommand
	closedLCommand
	finished
)

type lineParsingCommentStatus int

const (
	notStartedComment = iota + 1
	openedComment
	closedComment
)

type lineParsingState struct {
	status        lineParsingStatus
	commentStatus lineParsingCommentStatus
	buf           string
}

func newLineParsingState() lineParsingState {
	return lineParsingState{
		status:        initialized,
		commentStatus: notStartedComment,
		buf:           "",
	}
}

func (l *lineParsingState) transit(next lineParsingStatus) error {
	if next == openedACommand {
		if l.status != initialized {
			return fmt.Errorf("Invalid lineParsingStatus transition. %d to %d", l.status, next)
		}
	}

	if next == openedCCommand {
		if l.status != initialized {
			return fmt.Errorf("Invalid lineParsingStatus transition. %d to %d", l.status, next)
		}
	}

	if next == closedDest {
		if l.status != openedCCommand {
			return fmt.Errorf("Invalid lineParsingStatus transition. %d to %d", l.status, next)
		}
	}

	l.status = next
	return nil
}

func (l *lineParsingState) transitComment() error {
	if l.commentStatus == notStartedComment {
		l.commentStatus = openedComment
		return nil
	}
	if l.commentStatus == openedComment {
		l.commentStatus = closedComment
		return nil
	}
	return fmt.Errorf("Unexpected call of transitComment. Maybe you should stop parsing before")
}

func (l *lineParsingState) appendBuf(r rune) {
	l.buf += string(r)
}

func (l *lineParsingState) resetBuf() {
	l.buf = ""
}

func (p *Parser) parseLine(line string) (*Command, error) {
	var res *Command
	state := newLineParsingState()

	log.Printf("Start line parsing `%s`\n", line)

	for _, r := range line {
		log.Println(string(r), state.status)

		// Spaces
		if unicode.IsSpace(r) {
			continue
		}

		// Comment
		if r == '/' {
			if err := state.transitComment(); err != nil {
				return nil, err
			}
			if state.commentStatus == closedComment {
				break
			}
			continue
		}

		// L Command
		if r == '(' {
			if err := state.transit(openedLCommand); err != nil {
				return nil, err
			}
			res = NewCommand(LCommand)
			continue
		}

		if r == ')' {
			if err := state.transit(closedLCommand); err != nil {
				return nil, err
			}
			continue
		}

		// A Command
		if r == '@' {
			if err := state.transit(openedACommand); err != nil {
				return nil, err
			}
			res = NewCommand(ACommand)
			continue
		}

		// C Command
		if state.status == initialized {
			if err := state.transit(openedCCommand); err != nil {
				return nil, err
			}
			res = NewCommand(CCommand)
		}

		if r == '=' {
			if err := state.transit(closedDest); err != nil {
				return nil, err
			}
			if err := res.SetDest(state.buf); err != nil {
				return nil, err
			}
			state.resetBuf()
			continue
		}

		if r == ';' {
			if err := state.transit(closedComp); err != nil {
				return nil, err
			}
			if err := res.SetComp(state.buf); err != nil {
				return nil, err
			}
			state.resetBuf()
			continue
		}

		state.appendBuf(r)
	}

	if res == nil || res.Type == 0 {
		return nil, nil
	}

	if res.Type == ACommand || res.Type == LCommand {
		if err := res.SetSymbol(state.buf); err != nil {
			return nil, err
		}
	}
	if res.Type == CCommand {
		if state.status == openedCCommand || state.status == closedDest {
			if err := res.SetComp(state.buf); err != nil {
				return nil, err
			}
		}
		if state.status == closedComp {
			if err := res.SetJump(state.buf); err != nil {
				return nil, err
			}
		}
	}

	log.Println("Type", res.Type, "Symbol", res.Symbol, "Dest", res.Dest, "Comp", res.Comp, "Jump", res.Jump)
	return res, nil
}
