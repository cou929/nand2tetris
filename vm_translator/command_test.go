package main

import (
	"reflect"
	"testing"
)

func TestArithmeticOp_Valid(t *testing.T) {
	tests := []struct {
		name string
		a    ArithmeticOp
		want bool
	}{
		{
			name: "valid",
			a:    "lt",
			want: true,
		},
		{
			name: "unknown operator",
			a:    "mult",
			want: false,
		},

		{
			name: "case sensitive",
			a:    "Add",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Valid(); got != tt.want {
				t.Errorf("ArithmeticOp.Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommand_setArg1(t *testing.T) {
	type fields struct {
		Type CommandType
	}
	type args struct {
		arg string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    CommandArg1
		wantErr bool
	}{
		{
			name:    "valid arithmetic command",
			fields:  fields{Type: CommandArithmetic},
			args:    args{arg: "or"},
			want:    "or",
			wantErr: false,
		},
		{
			name:    "invalid command type",
			fields:  fields{Type: CommandReturn},
			args:    args{arg: "blah"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Command{
				Type: tt.fields.Type,
			}
			if err := c.setArg1(tt.args.arg); (err != nil) != tt.wantErr {
				t.Errorf("Command.setArg1() error = %v, wantErr %v", err, tt.wantErr)
			}
			if c.Arg1 != tt.want {
				t.Errorf("Command.setArg1() = %v, want %v", c.Arg1, tt.want)
			}
		})
	}
}

func TestCommand_setArg2(t *testing.T) {
	type fields struct {
		Type CommandType
	}
	type args struct {
		arg int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    CommandArg2
		wantErr bool
	}{
		{
			name:    "normal",
			fields:  fields{Type: CommandPush},
			args:    args{arg: 3},
			want:    3,
			wantErr: false,
		},
		{
			name:    "unsupported type",
			fields:  fields{Type: CommandArithmetic},
			args:    args{arg: 3},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Command{
				Type: tt.fields.Type,
			}
			if err := c.setArg2(tt.args.arg); (err != nil) != tt.wantErr {
				t.Errorf("Command.setArg2() error = %v, wantErr %v", err, tt.wantErr)
			}
			if c.Arg2 != tt.want {
				t.Errorf("Command.setArg2() = %v, want %v", c.Arg1, tt.want)
			}
		})
	}
}

func TestNewCommand(t *testing.T) {
	type args struct {
		line []string
	}
	tests := []struct {
		name    string
		args    args
		want    *Command
		wantErr bool
	}{
		{
			name: "valid arithmetic command",
			args: args{
				[]string{"sub"},
			},
			want: &Command{
				Type: CommandArithmetic,
				Arg1: "sub",
			},
			wantErr: false,
		},
		{
			name: "invalid command",
			args: args{
				[]string{"multi"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid arithmetic arg1 specified",
			args: args{
				[]string{"add", "1"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid push command",
			args: args{
				[]string{"push", "argument", "1"},
			},
			want: &Command{
				Type: CommandPush,
				Arg1: "argument",
				Arg2: 1,
			},
			wantErr: false,
		},
		{
			name: "invalid push arg2 format",
			args: args{
				[]string{"push", "static", "abc"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid push arg2 empty",
			args: args{
				[]string{"push", "static"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid push number of argument",
			args: args{
				[]string{"push", "static", "1", "2"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid pop command",
			args: args{
				[]string{"pop", "local", "11"},
			},
			want: &Command{
				Type: CommandPop,
				Arg1: "local",
				Arg2: 11,
			},
			wantErr: false,
		},
		{
			name: "invalid label arg2 specified",
			args: args{
				[]string{"label", "LOOP", "22"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid label command",
			args: args{
				[]string{"label", "LOOP"},
			},
			want: &Command{
				Type: CommandLabel,
				Arg1: "LOOP",
			},
			wantErr: false,
		},
		{
			name: "invalid label",
			args: args{
				[]string{"label", "LO-OP"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid goto command",
			args: args{
				[]string{"goto", "LOOP"},
			},
			want: &Command{
				Type: CommandGoto,
				Arg1: "LOOP",
			},
			wantErr: false,
		},
		{
			name: "valid if-goto command",
			args: args{
				[]string{"if-goto", "LOOP"},
			},
			want: &Command{
				Type: CommandIf,
				Arg1: "LOOP",
			},
			wantErr: false,
		},
		{
			name: "valid function command",
			args: args{
				[]string{"function", "test", "3"},
			},
			want: &Command{
				Type: CommandFunction,
				Arg1: "test",
				Arg2: 3,
			},
			wantErr: false,
		},
		{
			name: "valid call command",
			args: args{
				[]string{"call", "test", "2"},
			},
			want: &Command{
				Type: CommandCall,
				Arg1: "test",
				Arg2: 2,
			},
			wantErr: false,
		},
		{
			name: "valid return command",
			args: args{
				[]string{"return"},
			},
			want: &Command{
				Type: CommandReturn,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCommand(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemorySegment_Valid(t *testing.T) {
	tests := []struct {
		name string
		m    MemorySegment
		want bool
	}{
		{
			name: "valid",
			m:    "this",
			want: true,
		},
		{
			name: "unknown segment",
			m:    "dynamic",
			want: false,
		},
		{
			name: "case sensitive",
			m:    "Pointer",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Valid(); got != tt.want {
				t.Errorf("MemorySegment.Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
