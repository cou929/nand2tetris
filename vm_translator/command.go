package main

import (
	"fmt"
	"strconv"
)

type CommandType int

const (
	CommandArithmetic CommandType = iota + 1
	CommandPush
	CommandPop
	CommandLabel
	CommandGoto
	CommandIf
	CommandFunction
	CommandReturn
	CommandCall
)

type CommandArg1 string

type CommandArg2 int

type Command struct {
	Type CommandType
	Arg1 CommandArg1
	Arg2 CommandArg2
}

type ArithmeticOp string

func (a ArithmeticOp) Valid() bool {
	switch a {
	case "add":
		fallthrough
	case "sub":
		fallthrough
	case "neg":
		fallthrough
	case "eq":
		fallthrough
	case "gt":
		fallthrough
	case "lt":
		fallthrough
	case "and":
		fallthrough
	case "or":
		fallthrough
	case "not":
		return true
	}

	return false
}

type MemorySegment string

func (m MemorySegment) Valid() bool {
	switch m {
	case "argument":
		fallthrough
	case "local":
		fallthrough
	case "static":
		fallthrough
	case "constant":
		fallthrough
	case "this":
		fallthrough
	case "that":
		fallthrough
	case "pointer":
		fallthrough
	case "temp":
		return true
	}
	return false
}

func NewCommand(tokens []string) (*Command, error) {
	if len(tokens) > 3 {
		return nil, fmt.Errorf("Too many tokens %v", tokens)
	}
	c := &Command{}
	requiredTokenLen := 0

	switch tokens[0] {
	case "push":
		c.Type = CommandPush
		requiredTokenLen = 3
	case "pop":
		c.Type = CommandPop
		requiredTokenLen = 3
	case "label":
		c.Type = CommandLabel
		requiredTokenLen = 2
	case "goto":
		c.Type = CommandGoto
		requiredTokenLen = 2
	case "if-goto":
		c.Type = CommandIf
		requiredTokenLen = 2
	case "function":
		c.Type = CommandFunction
		requiredTokenLen = 3
	case "call":
		c.Type = CommandCall
		requiredTokenLen = 3
	case "return":
		c.Type = CommandReturn
		requiredTokenLen = 1
	default:
		c.Type = CommandArithmetic
		requiredTokenLen = 1
	}

	if len(tokens) != requiredTokenLen {
		return nil, fmt.Errorf("Invalid number of tokens %v", tokens)
	}

	if len(tokens) > 1 {
		if err := c.setArg1(tokens[1]); err != nil {
			return nil, err
		}
	}

	if c.Type == CommandArithmetic {
		if err := c.setArg1(tokens[0]); err != nil {
			return nil, err
		}
	}

	if len(tokens) > 2 {
		i, err := strconv.Atoi(tokens[2])
		if err != nil {
			return nil, fmt.Errorf("Invalid arg2 format %s", tokens[2])
		}
		if err := c.setArg2(i); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Command) setArg1(arg string) error {
	if c.Type == CommandReturn {
		return fmt.Errorf("Invalid usage of Arg1. Type is %d", c.Type)
	}
	if c.Type == CommandArithmetic {
		op := ArithmeticOp(arg)
		if !op.Valid() {
			return fmt.Errorf("Invalid arithmetic command %s", arg)
		}
	}
	if c.Type == CommandPush || c.Type == CommandPop {
		m := MemorySegment(arg)
		if !m.Valid() {
			return fmt.Errorf("Invalid memory segment %s", arg)
		}
		if c.Type == CommandPop && m == "constant" {
			return fmt.Errorf("Invalid command argument. pop constant is not supported. %s", arg)
		}
	}
	c.Arg1 = CommandArg1(arg)
	return nil
}

func (c *Command) setArg2(arg int) error {
	if c.Type != CommandPush && c.Type != CommandPop && c.Type != CommandFunction && c.Type != CommandCall {
		return fmt.Errorf("Invalid useage of Arg2. Type is %d", c.Type)
	}
	c.Arg2 = CommandArg2(arg)
	return nil
}
