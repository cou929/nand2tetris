package main

import "fmt"

type AsmCode struct {
	fileName string
	lineNum  int
	line     []string
}

func (a AsmCode) Code() []string {
	return a.line
}

func NewAsmCode(n string, i int, c *Command) (*AsmCode, error) {
	res := &AsmCode{
		fileName: n,
		lineNum:  i,
	}

	if c.Type == CommandPush {
		if c.Arg1 == "constant" {
			res.line = []string{
				fmt.Sprintf("@%d", c.Arg2),
				"D=A", // D=val
				"@SP",
				"A=M",
				"M=D", // M[SP]=val
				"@SP",
				"M=M+1", // increment stack pointer
			}
			return res, nil
		}
	}

	if c.Type == CommandArithmetic {
		switch c.Arg1 {
		case "add":
			res.line = []string{
				// pop y and set to D
				"@SP",
				"A=M-1",
				"D=M",
				"@SP",
				"M=M-1",
				// pop x and x+y
				"@SP",
				"A=M-1",
				"D=M+D",
				"@SP",
				"M=M-1",
				// push result
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			}
			return res, nil
		case "sub":
			res.line = []string{
				// pop y and set to D
				"@SP",
				"A=M-1",
				"D=M",
				"@SP",
				"M=M-1",
				// pop x and x-y
				"@SP",
				"A=M-1",
				"D=M-D",
				"@SP",
				"M=M-1",
				// push result
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			}
			return res, nil
		case "neg":
			res.line = []string{
				// pop y, set -y to D
				"@SP",
				"A=M-1",
				"D=-M",
				"@SP",
				"M=M-1",
				// push result
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			}
			return res, nil
		case "eq":
			res.line = []string{
				// pop y and set to D
				"@SP",
				"A=M-1",
				"D=M",
				"@SP",
				"M=M-1",
				// pop x and compare x-y with 0
				"@SP",
				"A=M-1",
				"D=M-D",
				"@SP",
				"M=M-1",
				fmt.Sprintf("@IS_ZERO.%s.%d", res.fileName, res.lineNum),
				"D;JEQ",
				fmt.Sprintf("@IS_NOT_ZERO.%s.%d", res.fileName, res.lineNum),
				"0;JMP",
				// comparison result to D
				fmt.Sprintf("(IS_ZERO.%s.%d)", res.fileName, res.lineNum),
				"@0",
				"D=!A",
				fmt.Sprintf("@END.%s.%d", res.fileName, res.lineNum),
				"0;JMP",
				fmt.Sprintf("(IS_NOT_ZERO.%s.%d)", res.fileName, res.lineNum),
				"@0",
				"D=A",
				fmt.Sprintf("@END.%s.%d", res.fileName, res.lineNum),
				"0;JMP",
				// push added result
				fmt.Sprintf("(END.%s.%d)", res.fileName, res.lineNum),
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			}
			return res, nil
		case "gt":
			res.line = []string{
				// pop y and set to D
				"@SP",
				"A=M-1",
				"D=M",
				"@SP",
				"M=M-1",
				// pop x and compare x-y with 0
				"@SP",
				"A=M-1",
				"D=M-D",
				"@SP",
				"M=M-1",
				fmt.Sprintf("@IS_ZERO.%s.%d", res.fileName, res.lineNum),
				"D;JGT",
				fmt.Sprintf("@IS_NOT_ZERO.%s.%d", res.fileName, res.lineNum),
				"0;JMP",
				// comparison result to D
				fmt.Sprintf("(IS_ZERO.%s.%d)", res.fileName, res.lineNum),
				"@0",
				"D=!A",
				fmt.Sprintf("@END.%s.%d", res.fileName, res.lineNum),
				"0;JMP",
				fmt.Sprintf("(IS_NOT_ZERO.%s.%d)", res.fileName, res.lineNum),
				"@0",
				"D=A",
				fmt.Sprintf("@END.%s.%d", res.fileName, res.lineNum),
				"0;JMP",
				// push added result
				fmt.Sprintf("(END.%s.%d)", res.fileName, res.lineNum),
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			}
			return res, nil
		case "lt":
			res.line = []string{
				// pop y and set to D
				"@SP",
				"A=M-1",
				"D=M",
				"@SP",
				"M=M-1",
				// pop x and compare x-y with 0
				"@SP",
				"A=M-1",
				"D=M-D",
				"@SP",
				"M=M-1",
				fmt.Sprintf("@IS_ZERO.%s.%d", res.fileName, res.lineNum),
				"D;JLT",
				fmt.Sprintf("@IS_NOT_ZERO.%s.%d", res.fileName, res.lineNum),
				"0;JMP",
				// comparison result to D
				fmt.Sprintf("(IS_ZERO.%s.%d)", res.fileName, res.lineNum),
				"@0",
				"D=!A",
				fmt.Sprintf("@END.%s.%d", res.fileName, res.lineNum),
				"0;JMP",
				fmt.Sprintf("(IS_NOT_ZERO.%s.%d)", res.fileName, res.lineNum),
				"@0",
				"D=A",
				fmt.Sprintf("@END.%s.%d", res.fileName, res.lineNum),
				"0;JMP",
				// push added result
				fmt.Sprintf("(END.%s.%d)", res.fileName, res.lineNum),
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			}
			return res, nil
		case "and":
			res.line = []string{
				// pop y and set to D
				"@SP",
				"A=M-1",
				"D=M",
				"@SP",
				"M=M-1",
				// pop x and operate `x And y`
				"@SP",
				"A=M-1",
				"D=M&D",
				"@SP",
				"M=M-1",
				// push added result
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			}
			return res, nil
		case "or":
			res.line = []string{
				// pop y and set to D
				"@SP",
				"A=M-1",
				"D=M",
				"@SP",
				"M=M-1",
				// pop x and operate `x Or y`
				"@SP",
				"A=M-1",
				"D=M|D",
				"@SP",
				"M=M-1",
				// push added result
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			}
			return res, nil
		case "not":
			res.line = []string{
				// pop y and operate `Not y`
				"@SP",
				"A=M-1",
				"D=!M",
				"@SP",
				"M=M-1",
				// push added result
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			}
			return res, nil
		}
	}
	return nil, nil
}
