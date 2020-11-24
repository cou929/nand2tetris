package main

import "fmt"

type AsmCode struct {
	line []string
}

func (a AsmCode) Code() []string {
	return a.line
}

func NewAsmCode(c *Command) (*AsmCode, error) {
	res := &AsmCode{}

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
			base = fmt.Sprintf("@%s.%d", c.Meta.fileName, c.Arg2)
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
			base = fmt.Sprintf("@%s.%d", c.Meta.fileName, c.Arg2)
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
				fmt.Sprintf("@IS_ZERO.%s.%d", c.Meta.fileName, c.Meta.lineNum),
				"D;JEQ",
				fmt.Sprintf("@IS_NOT_ZERO.%s.%d", c.Meta.fileName, c.Meta.lineNum),
				"0;JMP",
				// comparison result to D
				fmt.Sprintf("(IS_ZERO.%s.%d)", c.Meta.fileName, c.Meta.lineNum),
				"@0",
				"D=!A",
				fmt.Sprintf("@END.%s.%d", c.Meta.fileName, c.Meta.lineNum),
				"0;JMP",
				fmt.Sprintf("(IS_NOT_ZERO.%s.%d)", c.Meta.fileName, c.Meta.lineNum),
				"@0",
				"D=A",
				fmt.Sprintf("@END.%s.%d", c.Meta.fileName, c.Meta.lineNum),
				"0;JMP",
				// push added result
				fmt.Sprintf("(END.%s.%d)", c.Meta.fileName, c.Meta.lineNum),
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
				fmt.Sprintf("@IS_ZERO.%s.%d", c.Meta.fileName, c.Meta.lineNum),
				"D;JGT",
				fmt.Sprintf("@IS_NOT_ZERO.%s.%d", c.Meta.fileName, c.Meta.lineNum),
				"0;JMP",
				// comparison result to D
				fmt.Sprintf("(IS_ZERO.%s.%d)", c.Meta.fileName, c.Meta.lineNum),
				"@0",
				"D=!A",
				fmt.Sprintf("@END.%s.%d", c.Meta.fileName, c.Meta.lineNum),
				"0;JMP",
				fmt.Sprintf("(IS_NOT_ZERO.%s.%d)", c.Meta.fileName, c.Meta.lineNum),
				"@0",
				"D=A",
				fmt.Sprintf("@END.%s.%d", c.Meta.fileName, c.Meta.lineNum),
				"0;JMP",
				// push added result
				fmt.Sprintf("(END.%s.%d)", c.Meta.fileName, c.Meta.lineNum),
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
				fmt.Sprintf("@IS_ZERO.%s.%d", c.Meta.fileName, c.Meta.lineNum),
				"D;JLT",
				fmt.Sprintf("@IS_NOT_ZERO.%s.%d", c.Meta.fileName, c.Meta.lineNum),
				"0;JMP",
				// comparison result to D
				fmt.Sprintf("(IS_ZERO.%s.%d)", c.Meta.fileName, c.Meta.lineNum),
				"@0",
				"D=!A",
				fmt.Sprintf("@END.%s.%d", c.Meta.fileName, c.Meta.lineNum),
				"0;JMP",
				fmt.Sprintf("(IS_NOT_ZERO.%s.%d)", c.Meta.fileName, c.Meta.lineNum),
				"@0",
				"D=A",
				fmt.Sprintf("@END.%s.%d", c.Meta.fileName, c.Meta.lineNum),
				"0;JMP",
				// push added result
				fmt.Sprintf("(END.%s.%d)", c.Meta.fileName, c.Meta.lineNum),
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

	if c.Type == CommandFunction {
		res.line = []string{
			fmt.Sprintf("(%s)", c.Arg1),
		}
		for i := 0; i < int(c.Arg2); i++ {
			t := []string{
				"@LCL",
				"D=M",
				fmt.Sprintf("@%d", i),
				"A=D+A",
				"M=0",
				"@SP",
				"M=M+1",
			}
			res.line = append(res.line, t...)
		}
		return res, nil
	}

	if c.Type == CommandReturn {
		res.line = []string{
			// remember return address at R5 (temp segment)
			"@LCL",
			"D=M",
			"@5",
			"A=D-A",
			"D=M",
			"@R5",
			"M=D",
			// pop result and  set to ARG as returned value
			"@SP",
			"A=M-1",
			"D=M",
			"@SP", // unnecessary?
			"M=M-1",
			"@ARG",
			"A=M",
			"M=D",
			// resume caller's frame (SP, THAT, THIS, ARG, LCL)
			"@ARG",
			"D=M+1",
			"@SP",
			"M=D",

			"@LCL",
			"D=M",
			"@1",
			"D=D-A",
			"A=D",
			"D=M",
			"@THAT",
			"M=D",

			"@LCL",
			"D=M",
			"@2",
			"D=D-A",
			"A=D",
			"D=M",
			"@THIS",
			"M=D",

			"@LCL",
			"D=M",
			"@3",
			"D=D-A",
			"A=D",
			"D=M",
			"@ARG",
			"M=D",

			"@LCL",
			"D=M",
			"@4",
			"D=D-A",
			"A=D",
			"D=M",
			"@LCL",
			"M=D",
			// jump to return address
			"@R5",
			"A=M",
			"0;JMP",
		}
		return res, nil
	}

	if c.Type == CommandCall {
		retAddr := fmt.Sprintf("Return:%s.%s.%d", c.Meta.fileName, c.Meta.funcName, c.Meta.lineNum)
		res.line = []string{
			// hold return address
			fmt.Sprintf("@%s", retAddr),
			"D=A",
			"@SP",
			"A=M",
			"M=D",
			"@SP",
			"M=M+1",
			// hold LCL
			"@LCL",
			"D=M",
			"@SP",
			"A=M",
			"M=D",
			"@SP",
			"M=M+1",
			// hold ARG
			"@ARG",
			"D=M",
			"@SP",
			"A=M",
			"M=D",
			"@SP",
			"M=M+1",
			// hold THIS
			"@THIS",
			"D=M",
			"@SP",
			"A=M",
			"M=D",
			"@SP",
			"M=M+1",
			// hold THAT
			"@THAT",
			"D=M",
			"@SP",
			"A=M",
			"M=D",
			"@SP",
			"M=M+1",
			// move ARG
			"@SP",
			"D=M",
			fmt.Sprintf("@%d", c.Arg2),
			"D=D-A",
			"@5",
			"D=D-A",
			"@ARG",
			"M=D",
			// move LCL (SP is same position at first)
			"@SP",
			"D=M",
			"@LCL",
			"M=D",
			// jump to the func
			fmt.Sprintf("@%s", c.Arg1),
			"0;JMP",
			// mark return address
			fmt.Sprintf("(%s)", retAddr),
		}
		return res, nil
	}

	return nil, nil
}

func BootstrapLine() []string {
	return []string{
		// initialize SP
		"@256",
		"D=A",
		"@SP",
		"M=D",
		// call Sys.init
		// hold return address
		"@Return:vm:bootstrap",
		"D=A",
		"@SP",
		"A=M",
		"M=D",
		"@SP",
		"M=M+1",
		// hold LCL
		"@LCL",
		"D=M",
		"@SP",
		"A=M",
		"M=D",
		"@SP",
		"M=M+1",
		// hold ARG
		"@ARG",
		"D=M",
		"@SP",
		"A=M",
		"M=D",
		"@SP",
		"M=M+1",
		// hold THIS
		"@THIS",
		"D=M",
		"@SP",
		"A=M",
		"M=D",
		"@SP",
		"M=M+1",
		// hold THAT
		"@THAT",
		"D=M",
		"@SP",
		"A=M",
		"M=D",
		"@SP",
		"M=M+1",
		// move ARG
		"@SP",
		"D=M",
		"@0",
		"D=D-A",
		"@5",
		"D=D-A",
		"@ARG",
		"M=D",
		// move LCL (SP is same position at first)
		"@SP",
		"D=M",
		"@LCL",
		"M=D",
		// jump to the func
		"@Sys.init",
		"0;JMP",
		// mark return address
		"(Return:vm:bootstrap)",
	}
}
