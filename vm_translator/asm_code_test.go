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
					"D=A",
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
					"D=A",
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
					"D=A",
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
					"D=A",
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
		},
		{
			name: "push pointer",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandPush,
					Arg1: "pointer",
					Arg2: 99,
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
					"@R3",
					"D=A",
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
		},
		{
			name: "push temp",
			args: args{
				n: "test.vm",
				i: 2,
				c: &Command{
					Type: CommandPush,
					Arg1: "temp",
					Arg2: 99,
				},
			},
			want: &AsmCode{
				fileName: "test.vm",
				lineNum:  2,
				line: []string{
					"@R5",
					"D=A",
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
