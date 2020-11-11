package main

import (
	"reflect"
	"testing"
)

func TestNewBinaryCode(t *testing.T) {
	type args struct {
		command *Command
	}
	tests := []struct {
		name    string
		args    args
		want    *BinaryCode
		wantErr bool
	}{
		{
			name: "A Command",
			args: args{
				command: &Command{
					Type:   ACommand,
					Symbol: "4321",
				},
			},
			want: &BinaryCode{
				Line: 0b0_001_0000_1110_0001,
			},
			wantErr: false,
		},
		{
			name: "C Command `dest=comp`",
			args: args{
				command: &Command{
					Type: CCommand,
					Comp: "M+1",
					Dest: "AM",
				},
			},
			want: &BinaryCode{
				Comp: 0b1_110111,
				Dest: 0b101,
				Line: 0b111_1_110111_101_000,
			},
			wantErr: false,
		},
		{
			name: "C Command `comp;jump`",
			args: args{
				command: &Command{
					Type: CCommand,
					Comp: "D&A",
					Jump: "JNE",
				},
			},
			want: &BinaryCode{
				Comp: 0b0_000000,
				Jump: 0b101,
				Line: 0b111_0_000000_000_101,
			},
			wantErr: false,
		},
		{
			name: "C Command `dest=comp;jump`",
			args: args{
				command: &Command{
					Type: CCommand,
					Comp: "!M",
					Dest: "M",
					Jump: "JMP",
				},
			},
			want: &BinaryCode{
				Comp: 0b1_110001,
				Dest: 0b001,
				Jump: 0b111,
				Line: 0b111_1_110001_001_111,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBinaryCode(tt.args.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBinaryCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBinaryCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
