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
		base := ""
		switch c.Arg1 {
		case "local":
			base = fmt.Sprintf("@LCL")
		case "argument":
			base = fmt.Sprintf("@ARG")
		case "this":
			base = fmt.Sprintf("@THIS")
		case "that":
			base = fmt.Sprintf("@THAT")
		case "pointer":
			base = fmt.Sprintf("@R3")
		case "temp":
			base = fmt.Sprintf("@R5")
		case "constant":
			base = fmt.Sprintf("@%d", c.Arg2)
		case "static":
			base = fmt.Sprintf("@%s.%d", res.fileName, c.Arg2)
		}

		switch c.Arg1 {
		case "constant":
			res.line = []string{
				base,
				"D=A",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			}
			return res, nil
		case "static":
			res.line = []string{
				base,
				"D=M",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			}
			return res, nil
		case "pointer":
			fallthrough
		case "temp":
			res.line = []string{
				base,
				"D=A",
				fmt.Sprintf("@%d", c.Arg2),
				"A=D+A",
				"D=M",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			}
			return res, nil
		default: // local, argument, this, that
			res.line = []string{
				base,
				"D=M",
				fmt.Sprintf("@%d", c.Arg2),
				"A=D+A",
				"D=M",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			}
			return res, nil
		}
	}

	if c.Type == CommandPop {
		if c.Arg1 == "constant" {
			return nil, fmt.Errorf("pop command not supported constant")
		}
		base := ""
		switch c.Arg1 {
		case "local":
			base = fmt.Sprintf("@LCL")
		case "argument":
			base = fmt.Sprintf("@ARG")
		case "this":
			base = fmt.Sprintf("@THIS")
		case "that":
			base = fmt.Sprintf("@THAT")
		case "pointer":
			base = fmt.Sprintf("@R3")
		case "temp":
			base = fmt.Sprintf("@R5")
		case "static":
			base = fmt.Sprintf("@%s.%d", res.fileName, c.Arg2)
		}
		switch c.Arg1 {
		case "static":
			res.line = []string{
				// address to set
				base,
				"D=A",
				"@POP_DEST",
				"M=D",
				// pop and set
				"@SP",
				"A=M-1",
				"D=M",
				"@SP",
				"M=M-1",
				"@POP_DEST",
				"A=M",
				"M=D",
			}
			return res, nil
		case "pointer":
			fallthrough
		case "temp":
			res.line = []string{
				// address to set
				base,
				"D=A",
				fmt.Sprintf("@%d", c.Arg2),
				"D=D+A",
				"@POP_DEST",
				"M=D",
				// pop and set
				"@SP",
				"A=M-1",
				"D=M",
				"@SP",
				"M=M-1",
				"@POP_DEST",
				"A=M",
				"M=D",
			}
			return res, nil
		default: // local, argument, this, that
			res.line = []string{
				// address to set
				base,
				"D=M",
				fmt.Sprintf("@%d", c.Arg2),
				"D=D+A",
				"@POP_DEST",
				"M=D",
				// pop and set
				"@SP",
				"A=M-1",
				"D=M",
				"@SP",
				"M=M-1",
				"@POP_DEST",
				"A=M",
				"M=D",
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

	if c.Type == CommandLabel {
		// TODO function name
		res.line = []string{
			fmt.Sprintf("(%s)", c.Arg1),
		}
		return res, nil
	}

	if c.Type == CommandGoto {
		// TODO function name
		res.line = []string{
			fmt.Sprintf("@%s", c.Arg1),
			"0;JMP",
		}
		return res, nil
	}

	if c.Type == CommandIf {
		// TODO function name
		res.line = []string{
			"@SP",
			"A=M-1",
			"D=M",
			"@SP",
			"M=M-1",
			fmt.Sprintf("@%s", c.Arg1),
			"D;JNE",
		}
		return res, nil
	}

	return nil, nil
}
