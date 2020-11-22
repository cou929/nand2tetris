package main

import (
	"reflect"
	"testing"
)

func TestNewAsmCode_Push(t *testing.T) {
	type args struct {
		n string
		i int
		c *Command
	}
	tests := []struct {
		name    string
		args    args
		want    *AsmCode
		wantErr bool
	}{
		{
			name: "push local",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandPush,
					Arg1: "local",
					Arg2: 99,
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
					"@LCL",
					"D=M",
					"@99",
					"A=D+A",
					"D=M",
					"@SP",
					"A=M",
					"M=D",
					"@SP",
					"M=M+1",
				},
			},
			wantErr: false,
		},
		{
			name: "push argument",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandPush,
					Arg1: "argument",
					Arg2: 99,
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
					"@ARG",
					"D=M",
					"@99",
					"A=D+A",
					"D=M",
					"@SP",
					"A=M",
					"M=D",
					"@SP",
					"M=M+1",
				},
			},
			wantErr: false,
		},
		{
			name: "push this",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandPush,
					Arg1: "this",
					Arg2: 99,
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
					"@THIS",
					"D=M",
					"@99",
					"A=D+A",
					"D=M",
					"@SP",
					"A=M",
					"M=D",
					"@SP",
					"M=M+1",
				},
			},
			wantErr: false,
		},
		{
			name: "push that",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandPush,
					Arg1: "that",
					Arg2: 99,
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
					"@THAT",
					"D=M",
					"@99",
					"A=D+A",
					"D=M",
					"@SP",
					"A=M",
					"M=D",
					"@SP",
					"M=M+1",
				},
			},
			wantErr: false,
		},
		{
			name: "push pointer",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandPush,
					Arg1: "pointer",
					Arg2: 1,
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
					"@R3",
					"D=A",
					"@1",
					"A=D+A",
					"D=M",
					"@SP",
					"A=M",
					"M=D",
					"@SP",
					"M=M+1",
				},
			},
			wantErr: false,
		},
		{
			name: "push temp",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandPush,
					Arg1: "temp",
					Arg2: 6,
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
					"@R5",
					"D=A",
					"@6",
					"A=D+A",
					"D=M",
					"@SP",
					"A=M",
					"M=D",
					"@SP",
					"M=M+1",
				},
			},
			wantErr: false,
		},
		{
			name: "push constant",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandPush,
					Arg1: "constant",
					Arg2: 99,
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
					"@99",
					"D=A", // D=99
					"@SP",
					"A=M",
					"M=D", // M[SP]=99
					"@SP",
					"M=M+1", // increment stack pointer
				},
			},
			wantErr: false,
		},
		{
			name: "push static",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandPush,
					Arg1: "static",
					Arg2: 99,
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
					"@test.vm.99",
					"D=M",
					"@SP",
					"A=M",
					"M=D",
					"@SP",
					"M=M+1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAsmCode(tt.args.n, tt.args.i, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAsmCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAsmCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAsmCode_Pop(t *testing.T) {
	type args struct {
		n string
		i int
		c *Command
	}
	tests := []struct {
		name    string
		args    args
		want    *AsmCode
		wantErr bool
	}{
		{
			name: "pop local",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandPop,
					Arg1: "local",
					Arg2: 99,
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
					// address to set
					"@LCL",
					"D=M",
					"@99",
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
				},
			},
			wantErr: false,
		},
		{
			name: "pop argument",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandPop,
					Arg1: "argument",
					Arg2: 99,
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
					// address to set
					"@ARG",
					"D=M",
					"@99",
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
				},
			},
			wantErr: false,
		},
		{
			name: "pop this",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandPop,
					Arg1: "this",
					Arg2: 99,
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
					// address to set
					"@THIS",
					"D=M",
					"@99",
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
				},
			},
			wantErr: false,
		},
		{
			name: "pop that",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandPop,
					Arg1: "that",
					Arg2: 99,
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
					// address to set
					"@THAT",
					"D=M",
					"@99",
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
				},
			},
			wantErr: false,
		},
		{
			name: "pop pointer",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandPop,
					Arg1: "pointer",
					Arg2: 1,
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
					// address to set
					"@R3",
					"D=A",
					"@1",
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
				},
			},
			wantErr: false,
		},
		{
			name: "pop temp",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandPop,
					Arg1: "temp",
					Arg2: 6,
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
					// address to set
					"@R5",
					"D=A",
					"@6",
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
				},
			},
			wantErr: false,
		},
		{
			name: "pop constant",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandPop,
					Arg1: "constant",
					Arg2: 99,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "pop static",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandPop,
					Arg1: "static",
					Arg2: 99,
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
					// address to set
					"@test.vm.99",
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
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAsmCode(tt.args.n, tt.args.i, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAsmCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAsmCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAsmCode_Arithmetic(t *testing.T) {
	type args struct {
		n string
		i int
		c *Command
	}
	tests := []struct {
		name    string
		args    args
		want    *AsmCode
		wantErr bool
	}{
		{
			name: "add",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandArithmetic,
					Arg1: "add",
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
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
				},
			},
			wantErr: false,
		},
		{
			name: "sub",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandArithmetic,
					Arg1: "sub",
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
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
				},
			},
			wantErr: false,
		},
		{
			name: "neg",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandArithmetic,
					Arg1: "neg",
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
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
				},
			},
			wantErr: false,
		},
		{
			name: "eq",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandArithmetic,
					Arg1: "eq",
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
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
					"@IS_ZERO.test.vm.2",
					"D;JEQ",
					"@IS_NOT_ZERO.test.vm.2",
					"0;JMP",
					// comparison result to D
					"(IS_ZERO.test.vm.2)",
					"@0",
					"D=!A",
					"@END.test.vm.2",
					"0;JMP",
					"(IS_NOT_ZERO.test.vm.2)",
					"@0",
					"D=A",
					"@END.test.vm.2",
					"0;JMP",
					// push added result
					"(END.test.vm.2)",
					"@SP",
					"A=M",
					"M=D",
					"@SP",
					"M=M+1",
				},
			},
			wantErr: false,
		},
		{
			name: "gt",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandArithmetic,
					Arg1: "gt",
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
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
					"@IS_ZERO.test.vm.2",
					"D;JGT",
					"@IS_NOT_ZERO.test.vm.2",
					"0;JMP",
					// comparison result to D
					"(IS_ZERO.test.vm.2)",
					"@0",
					"D=!A",
					"@END.test.vm.2",
					"0;JMP",
					"(IS_NOT_ZERO.test.vm.2)",
					"@0",
					"D=A",
					"@END.test.vm.2",
					"0;JMP",
					// push added result
					"(END.test.vm.2)",
					"@SP",
					"A=M",
					"M=D",
					"@SP",
					"M=M+1",
				},
			},
			wantErr: false,
		},
		{
			name: "lt",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandArithmetic,
					Arg1: "lt",
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
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
					"@IS_ZERO.test.vm.2",
					"D;JLT",
					"@IS_NOT_ZERO.test.vm.2",
					"0;JMP",
					// comparison result to D
					"(IS_ZERO.test.vm.2)",
					"@0",
					"D=!A",
					"@END.test.vm.2",
					"0;JMP",
					"(IS_NOT_ZERO.test.vm.2)",
					"@0",
					"D=A",
					"@END.test.vm.2",
					"0;JMP",
					// push added result
					"(END.test.vm.2)",
					"@SP",
					"A=M",
					"M=D",
					"@SP",
					"M=M+1",
				},
			},
			wantErr: false,
		},
		{
			name: "and",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandArithmetic,
					Arg1: "and",
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
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
				},
			},
			wantErr: false,
		},
		{
			name: "or",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandArithmetic,
					Arg1: "or",
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
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
				},
			},
			wantErr: false,
		},
		{
			name: "not",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandArithmetic,
					Arg1: "not",
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
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
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAsmCode(tt.args.n, tt.args.i, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAsmCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAsmCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAsmCode_Label(t *testing.T) {
	type args struct {
		n string
		i int
		c *Command
	}
	tests := []struct {
		name    string
		args    args
		want    *AsmCode
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandLabel,
					Arg1: "MY_LABEL",
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
					"(MY_LABEL)",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAsmCode(tt.args.n, tt.args.i, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAsmCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAsmCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAsmCode_Goto(t *testing.T) {
	type args struct {
		n string
		i int
		c *Command
	}
	tests := []struct {
		name    string
		args    args
		want    *AsmCode
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandGoto,
					Arg1: "MY_LABEL",
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
					"@MY_LABEL",
					"0;JMP",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAsmCode(tt.args.n, tt.args.i, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAsmCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAsmCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAsmCode_If(t *testing.T) {
	type args struct {
		n string
		i int
		c *Command
	}
	tests := []struct {
		name    string
		args    args
		want    *AsmCode
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandIf,
					Arg1: "MY_LABEL",
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
					"@SP",
					"A=M-1",
					"D=M",
					"@SP",
					"M=M-1",
					"@MY_LABEL",
					"D;JNE",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAsmCode(tt.args.n, tt.args.i, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAsmCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAsmCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
