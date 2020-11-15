package main

import "fmt"

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

func NewCommand(typ CommandType) *Command {
	return &Command{Type: typ}
}

func (c *Command) SetArg1(arg string) error {
	if c.Type == CommandReturn {
		return fmt.Errorf("Invalid usage of Arg1. Type is %d", c.Type)
	}
	if c.Type == CommandArithmetic {
		op := ArithmeticOp(arg)
		if !op.Valid() {
			return fmt.Errorf("Invalid arithmetic command %s", arg)
		}
	}
	c.Arg1 = CommandArg1(arg)
	return nil
}

func (c *Command) SetArg2(arg int) error {
	if c.Type != CommandPush && c.Type != CommandPop && c.Type != CommandFunction && c.Type != CommandCall {
		return fmt.Errorf("Invalid useage of Arg2. Type is %d", c.Type)
	}
	c.Arg2 = CommandArg2(arg)
	return nil
}
