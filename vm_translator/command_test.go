package main

import (
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

func TestCommand_SetArg1(t *testing.T) {
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
			if err := c.SetArg1(tt.args.arg); (err != nil) != tt.wantErr {
				t.Errorf("Command.SetArg1() error = %v, wantErr %v", err, tt.wantErr)
			}
			if c.Arg1 != tt.want {
				t.Errorf("Command.SetArg1() = %v, want %v", c.Arg1, tt.want)
			}
		})
	}
}

func TestCommand_SetArg2(t *testing.T) {
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
			if err := c.SetArg2(tt.args.arg); (err != nil) != tt.wantErr {
				t.Errorf("Command.SetArg2() error = %v, wantErr %v", err, tt.wantErr)
			}
			if c.Arg2 != tt.want {
				t.Errorf("Command.SetArg2() = %v, want %v", c.Arg1, tt.want)
			}
		})
	}
}
