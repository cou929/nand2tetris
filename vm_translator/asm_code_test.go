package main

import (
	"reflect"
	"testing"
)

func TestNewAsmCode_Push(t *testing.T) {
	type args struct {
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
				c: &Command{
					Type: CommandPush,
					Arg1: "local",
					Arg2: 99,
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
				c: &Command{
					Type: CommandPush,
					Arg1: "argument",
					Arg2: 99,
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
				c: &Command{
					Type: CommandPush,
					Arg1: "this",
					Arg2: 99,
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
				c: &Command{
					Type: CommandPush,
					Arg1: "that",
					Arg2: 99,
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
				c: &Command{
					Type: CommandPush,
					Arg1: "pointer",
					Arg2: 1,
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
				c: &Command{
					Type: CommandPush,
					Arg1: "temp",
					Arg2: 6,
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
				c: &Command{
					Type: CommandPush,
					Arg1: "constant",
					Arg2: 99,
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
				c: &Command{
					Type: CommandPush,
					Arg1: "static",
					Arg2: 99,
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
				line: []string{
					"@TestClass.vm.99",
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
			got, err := NewAsmCode(tt.args.c)
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
				c: &Command{
					Type: CommandPop,
					Arg1: "local",
					Arg2: 99,
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
				c: &Command{
					Type: CommandPop,
					Arg1: "argument",
					Arg2: 99,
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
				c: &Command{
					Type: CommandPop,
					Arg1: "this",
					Arg2: 99,
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
				c: &Command{
					Type: CommandPop,
					Arg1: "that",
					Arg2: 99,
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
				c: &Command{
					Type: CommandPop,
					Arg1: "pointer",
					Arg2: 1,
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
				c: &Command{
					Type: CommandPop,
					Arg1: "temp",
					Arg2: 6,
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
				c: &Command{
					Type: CommandPop,
					Arg1: "constant",
					Arg2: 99,
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "pop static",
			args: args{
				c: &Command{
					Type: CommandPop,
					Arg1: "static",
					Arg2: 99,
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
				line: []string{
					// address to set
					"@TestClass.vm.99",
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
			got, err := NewAsmCode(tt.args.c)
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
				c: &Command{
					Type: CommandArithmetic,
					Arg1: "add",
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
				c: &Command{
					Type: CommandArithmetic,
					Arg1: "sub",
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
				c: &Command{
					Type: CommandArithmetic,
					Arg1: "neg",
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
				c: &Command{
					Type: CommandArithmetic,
					Arg1: "eq",
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
					"@IS_ZERO.TestClass.vm.2",
					"D;JEQ",
					"@IS_NOT_ZERO.TestClass.vm.2",
					"0;JMP",
					// comparison result to D
					"(IS_ZERO.TestClass.vm.2)",
					"@0",
					"D=!A",
					"@END.TestClass.vm.2",
					"0;JMP",
					"(IS_NOT_ZERO.TestClass.vm.2)",
					"@0",
					"D=A",
					"@END.TestClass.vm.2",
					"0;JMP",
					// push added result
					"(END.TestClass.vm.2)",
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
				c: &Command{
					Type: CommandArithmetic,
					Arg1: "gt",
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
					"@IS_ZERO.TestClass.vm.2",
					"D;JGT",
					"@IS_NOT_ZERO.TestClass.vm.2",
					"0;JMP",
					// comparison result to D
					"(IS_ZERO.TestClass.vm.2)",
					"@0",
					"D=!A",
					"@END.TestClass.vm.2",
					"0;JMP",
					"(IS_NOT_ZERO.TestClass.vm.2)",
					"@0",
					"D=A",
					"@END.TestClass.vm.2",
					"0;JMP",
					// push added result
					"(END.TestClass.vm.2)",
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
				c: &Command{
					Type: CommandArithmetic,
					Arg1: "lt",
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
					"@IS_ZERO.TestClass.vm.2",
					"D;JLT",
					"@IS_NOT_ZERO.TestClass.vm.2",
					"0;JMP",
					// comparison result to D
					"(IS_ZERO.TestClass.vm.2)",
					"@0",
					"D=!A",
					"@END.TestClass.vm.2",
					"0;JMP",
					"(IS_NOT_ZERO.TestClass.vm.2)",
					"@0",
					"D=A",
					"@END.TestClass.vm.2",
					"0;JMP",
					// push added result
					"(END.TestClass.vm.2)",
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
				c: &Command{
					Type: CommandArithmetic,
					Arg1: "and",
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
				c: &Command{
					Type: CommandArithmetic,
					Arg1: "or",
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
				c: &Command{
					Type: CommandArithmetic,
					Arg1: "not",
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
			got, err := NewAsmCode(tt.args.c)
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
				c: &Command{
					Type: CommandLabel,
					Arg1: "MY_LABEL",
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
				line: []string{
					"(MY_LABEL)",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAsmCode(tt.args.c)
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
				c: &Command{
					Type: CommandGoto,
					Arg1: "MY_LABEL",
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
			got, err := NewAsmCode(tt.args.c)
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
				c: &Command{
					Type: CommandIf,
					Arg1: "MY_LABEL",
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
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
			got, err := NewAsmCode(tt.args.c)
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

func TestNewAsmCode_Function(t *testing.T) {
	type args struct {
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
				c: &Command{
					Type: CommandFunction,
					Arg1: "myFunc",
					Arg2: 3,
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
				line: []string{
					"(myFunc)",
					"@LCL",
					"D=M",
					"@0",
					"A=D+A",
					"M=0",
					"@SP",
					"M=M+1",
					"@LCL",
					"D=M",
					"@1",
					"A=D+A",
					"M=0",
					"@SP",
					"M=M+1",
					"@LCL",
					"D=M",
					"@2",
					"A=D+A",
					"M=0",
					"@SP",
					"M=M+1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAsmCode(tt.args.c)
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

func TestNewAsmCode_Return(t *testing.T) {
	type args struct {
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
				c: &Command{
					Type: CommandReturn,
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
				line: []string{
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
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAsmCode(tt.args.c)
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

func TestNewAsmCode_Call(t *testing.T) {
	type args struct {
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
				c: &Command{
					Type: CommandCall,
					Arg1: "myFunc",
					Arg2: 3,
					Meta: &CommandMeta{
						"TestClass.vm",
						"TestClass.fooFn",
						2,
					},
				},
			},
			want: &AsmCode{
				line: []string{
					// hold return address
					"@Return:TestClass.vm.TestClass.fooFn.2",
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
					"@3",
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
					"@myFunc",
					"0;JMP",
					// mark return address
					"(Return:TestClass.vm.TestClass.fooFn.2)",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAsmCode(tt.args.c)
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

func TestBootstrapLine(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "normal",
			want: []string{
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BootstrapLine(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BootstrapLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
