package main

import "fmt"

type AsmCode struct {
	Line []string
}

func NewAsmCode(c *Command) (*AsmCode, error) {
	if c.Type == CommandPush {
		if c.Arg1 == "constant" {
			return &AsmCode{
				Line: []string{
					fmt.Sprintf("@%d", c.Arg2),
					"D=A", // D=val
					"@SP",
					"A=M",
					"M=D", // M[SP]=val
					"@SP",
					"M=M+1", // increment stack pointer
				},
			}, nil
		}
	}

	if c.Type == CommandArithmetic {
		if c.Arg1 == "add" {
			return &AsmCode{
				Line: []string{
					// pop and set to D
					"@SP",
					"A=M-1",
					"D=M",
					"@SP",
					"M=M-1", // decrement stack pointer
					// pop and add
					"@SP",
					"A=M-1",
					"D=D+M", // add popped 2 values
					"@SP",
					"M=M-1", // decrement stack pointer
					// push added result
					"@SP",
					"A=M",
					"M=D",
					"@SP",
					"M=M+1",
				},
			}, nil
		}
	}
	return nil, nil
}
