package main

import (
	"reflect"
	"testing"
)

func TestNewAsmCode(t *testing.T) {
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
			name: "push constant",
			args: args{
				c: &Command{
					Type: CommandPush,
					Arg1: "constant",
					Arg2: 99,
				},
			},
			want: &AsmCode{
				Line: []string{
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
			name: "add",
			args: args{
				c: &Command{
					Type: CommandArithmetic,
					Arg1: "add",
				},
			},
			want: &AsmCode{
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
			},
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
