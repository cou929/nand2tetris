package main

import "testing"

func TestSymbolTable_AddVariable(t *testing.T) {
	type args struct {
		symbol CommandSymbol
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "first insertion",
			args:    args{symbol: "foo"},
			want:    16,
			wantErr: false,
		},
		{
			name:    "second insertion",
			args:    args{symbol: "bar"},
			want:    17,
			wantErr: false,
		},
		{
			name:    "duplicated symbol",
			args:    args{symbol: "foo"},
			want:    0,
			wantErr: true,
		},
	}
	s := NewSymbolTable()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.AddVariable(tt.args.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("SymbolTable.AddVariable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SymbolTable.AddVariable() = %v, want %v", got, tt.want)
			}
			got2 := s.GetAddress(tt.args.symbol)
			if !tt.wantErr && got2 != tt.want {
				t.Errorf("SymbolTable.AddVariable() => GetAddress() = %v, want %v", got2, tt.want)
			}
		})
	}
}
