package main

import (
	"reflect"
	"testing"
)

func TestParser_parseLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    *Command
		wantErr bool
	}{
		{
			name: "A Command with no label",
			args: args{line: "@4321"},
			want: &Command{
				Type:   ACommand,
				Symbol: "4321",
				Dest:   "",
				Comp:   "",
				Jump:   "",
			},
			wantErr: false,
		},
		{
			name: "C Command with dest",
			args: args{line: "D=D-A"},
			want: &Command{
				Type:   CCommand,
				Symbol: "",
				Dest:   "D",
				Comp:   "D-A",
				Jump:   "",
			},
			wantErr: false,
		},
		{
			name: "C Command with jump",
			args: args{line: "M+1;JGT"},
			want: &Command{
				Type:   CCommand,
				Symbol: "",
				Dest:   "",
				Comp:   "M+1",
				Jump:   "JGT",
			},
			wantErr: false,
		},
		{
			name: "C Command with dest and jump",
			args: args{line: "D=!D;JMP"},
			want: &Command{
				Type:   CCommand,
				Symbol: "",
				Dest:   "D",
				Comp:   "!D",
				Jump:   "JMP",
			},
			wantErr: false,
		},
		{
			name: "Label symbol",
			args: args{line: "(LOOP)"},
			want: &Command{
				Type:   LCommand,
				Symbol: "LOOP",
				Dest:   "",
				Comp:   "",
				Jump:   "",
			},
			wantErr: false,
		},
		{
			name: "Variable symbol",
			args: args{line: "@.VAR_:$"},
			want: &Command{
				Type:   ACommand,
				Symbol: ".VAR_:$",
				Dest:   "",
				Comp:   "",
				Jump:   "",
			},
			wantErr: false,
		},
		{
			name: "Comment",
			args: args{line: "D=!D // foobar"},
			want: &Command{
				Type:   CCommand,
				Symbol: "",
				Dest:   "D",
				Comp:   "!D",
				Jump:   "",
			},
			wantErr: false,
		},
		{
			name: "Spaces in C Command",
			args: args{line: "  D   = !  D; JM    P   // foobar"},
			want: &Command{
				Type:   CCommand,
				Symbol: "",
				Dest:   "D",
				Comp:   "!D",
				Jump:   "JMP",
			},
			wantErr: false,
		},
		{
			name: "Spaces in A Command",
			args: args{line: "   @S   ymb   ol"},
			want: &Command{
				Type:   ACommand,
				Symbol: "Symbol",
				Dest:   "",
				Comp:   "",
				Jump:   "",
			},
			wantErr: false,
		},
		{
			name:    "Blank line",
			args:    args{line: "// foobar"},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Unknown mnemonic A Command",
			args:    args{line: "#1000"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Unknown mnemonic at comp",
			args:    args{line: "M=P+1"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Unknown mnemonic at jump",
			args:    args{line: "M=1;XYZ"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Symbol cannot start digits",
			args:    args{line: "(1A)"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid syntax, label and C Command",
			args:    args{line: "(1A) D=1"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid syntax, no C Command content",
			args:    args{line: ";JMP"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Symbol cannot contain other than `[a-zA-Z0-9_.$:]`",
			args:    args{line: "@**THIS**"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, err := p.parseLine(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
