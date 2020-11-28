package main

import "testing"

func TestNewIdentifierToken(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name  string
		args  args
		want  IdentifierToken
		want1 bool
	}{
		{
			name:  "normal",
			args:  args{in: "var_1"},
			want:  "var_1",
			want1: true,
		},
		{
			name:  "single char",
			args:  args{in: "i"},
			want:  "i",
			want1: true,
		},
		{
			name:  "first letter should not a digit",
			args:  args{in: "1var"},
			want:  "",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := NewIdentifierToken(tt.args.in)
			if got != tt.want {
				t.Errorf("NewIdentifierToken() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("NewIdentifierToken() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
