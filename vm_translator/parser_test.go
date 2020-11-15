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
			name: "valid arithmetic command",
			args: args{line: "add"},
			want: &Command{
				Type: CommandArithmetic,
				Arg1: "add",
			},
			wantErr: false,
		},
		{
			name: "valid push command",
			args: args{line: "push local 99"},
			want: &Command{
				Type: CommandPush,
				Arg1: "local",
				Arg2: 99,
			},
			wantErr: false,
		},
		{
			name: "ignore white spaces",
			args: args{line: "  pop    local 3  "},
			want: &Command{
				Type: CommandPop,
				Arg1: "local",
				Arg2: 3,
			},
			wantErr: false,
		},
		{
			name: "comment",
			args: args{line: "pop local 0  // comment"},
			want: &Command{
				Type: CommandPop,
				Arg1: "local",
				Arg2: 0,
			},
			wantErr: false,
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
