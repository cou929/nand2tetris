package main

import (
	"reflect"
	"testing"
)

func TestCommand_SetSymbolOrValue(t *testing.T) {
	type fields struct {
		Symbol CommandSymbol
	}
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    *Command
		wantErr bool
	}{
		{
			name: "decimal",
			args: args{in: "1234"},
			want: &Command{
				AddressValue: 1234,
			},
			wantErr: false,
		},
		{
			name: "string label",
			args: args{in: "LOOP5"},
			want: &Command{
				Symbol: "LOOP5",
			},
			wantErr: false,
		},
		{
			name: "allowed symbols",
			args: args{in: "_.$:"},
			want: &Command{
				Symbol: "_.$:",
			},
			wantErr: false,
		},
		{
			name:    "in case of label, decimal at first char is invalid",
			args:    args{in: "1A"},
			want:    &Command{},
			wantErr: true,
		},
		{
			name:    "in case of label, using not allowed symbol is invalid",
			args:    args{in: "A-B"},
			want:    &Command{},
			wantErr: true,
		},
		{
			name:    "in case of constants, a floating point is invalid",
			args:    args{in: "12.55"},
			want:    &Command{},
			wantErr: true,
		},
		{
			name:    "in case of constants, negative number is invalid",
			args:    args{in: "-99"},
			want:    &Command{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Command{}
			if err := c.SetSymbolOrValue(tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("Command.SetSymbolOrValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(c, tt.want) {
				t.Errorf("Command.SetSymbolOrValue() = %v, want %v", c.Symbol, tt.want)
			}
		})
	}
}

func TestCommand_SetDest(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    CommandDest
		wantErr bool
	}{
		{
			name:    "normal",
			args:    args{in: "AD"},
			want:    DestAD,
			wantErr: false,
		},
		{
			name:    "invalid",
			args:    args{in: "DA"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Command{}
			if err := c.SetDest(tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("Command.SetDest() error = %v, wantErr %v", err, tt.wantErr)
			}
			if c.Dest != tt.want {
				t.Errorf("Command.SetDest() = %v, want %v", c.Dest, tt.want)
			}
		})
	}
}

func TestCommand_SetComp(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    CommandComp
		wantErr bool
	}{
		{
			name:    "normal",
			args:    args{in: "D&A"},
			want:    CompDAndA,
			wantErr: false,
		},
		{
			name:    "invalid",
			args:    args{in: "A+2"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Command{}
			if err := c.SetComp(tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("Command.SetComp() error = %v, wantErr %v", err, tt.wantErr)
			}
			if c.Comp != tt.want {
				t.Errorf("Command.SetComp() = %v, want %v", c.Comp, tt.want)
			}
		})
	}
}

func TestCommand_SetJump(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    CommandJump
		wantErr bool
	}{
		{
			name:    "normal",
			args:    args{in: "JLE"},
			want:    JLE,
			wantErr: false,
		},
		{
			name:    "invalid",
			args:    args{in: "JM1"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Command{}
			if err := c.SetJump(tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("Command.SetJump() error = %v, wantErr %v", err, tt.wantErr)
			}
			if c.Jump != tt.want {
				t.Errorf("Command.SetJump() = %v, want %v", c.Jump, tt.want)
			}
		})
	}
}
