package main

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSymbolTable_Define(t *testing.T) {
	type fields struct {
		classScopeTable ScopedTable
		funcScopeTable  ScopedTable
		index           map[VarKind]int
	}
	type args struct {
		name string
		typ  string
		kind VarKind
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want1   ScopedTable
		want2   ScopedTable
		want3   map[VarKind]int
		wantErr bool
	}{
		{
			name: "define class scoped symbol",
			fields: fields{
				classScopeTable: make(ScopedTable),
				funcScopeTable:  make(ScopedTable),
				index:           make(map[VarKind]int),
			},
			args: args{
				name: "foo",
				typ:  "int",
				kind: Static,
			},
			want1: ScopedTable{
				"foo": &SymbolTableEntry{
					Name:  "foo",
					Typ:   "int",
					Kind:  Static,
					Index: 0,
				},
			},
			want2: ScopedTable{},
			want3: map[VarKind]int{
				Static: 1,
			},
			wantErr: false,
		},
		{
			name: "define func scoped symbol",
			fields: fields{
				classScopeTable: ScopedTable{
					"foo": &SymbolTableEntry{ // same name ok at other scope
						Name:  "foo",
						Typ:   "int",
						Kind:  Field,
						Index: 1,
					},
				},
				funcScopeTable: make(ScopedTable),
				index: map[VarKind]int{
					Field: 2,
					Var:   5,
				},
			},
			args: args{
				name: "foo",
				typ:  "int",
				kind: Var,
			},
			want1: ScopedTable{
				"foo": &SymbolTableEntry{
					Name:  "foo",
					Typ:   "int",
					Kind:  Field,
					Index: 1,
				},
			},
			want2: ScopedTable{
				"foo": &SymbolTableEntry{
					Name:  "foo",
					Typ:   "int",
					Kind:  Var,
					Index: 5,
				},
			},
			want3: map[VarKind]int{
				Field: 2,
				Var:   6,
			},
			wantErr: false,
		},
		{
			name: "naming conflict",
			fields: fields{
				classScopeTable: make(ScopedTable),
				funcScopeTable: ScopedTable{
					"foo": &SymbolTableEntry{
						Name:  "foo",
						Typ:   "char",
						Kind:  Argument,
						Index: 0,
					},
				},
				index: map[VarKind]int{
					Argument: 1,
					Var:      5,
				},
			},
			args: args{
				name: "foo",
				typ:  "int",
				kind: Var,
			},
			want1: ScopedTable{},
			want2: ScopedTable{
				"foo": &SymbolTableEntry{
					Name:  "foo",
					Typ:   "char",
					Kind:  Argument,
					Index: 0,
				},
			},
			want3: map[VarKind]int{
				Argument: 1,
				Var:      5,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SymbolTable{
				classScopeTable: tt.fields.classScopeTable,
				funcScopeTable:  tt.fields.funcScopeTable,
				index:           tt.fields.index,
			}
			if err := s.Define(tt.args.name, tt.args.typ, tt.args.kind); (err != nil) != tt.wantErr {
				t.Errorf("SymbolTable.Define() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(s.classScopeTable, tt.want1); diff != "" {
				t.Errorf("SymbolTable.Define() classScopeTable diff (-got +want1)\n%s", diff)
			}
			if diff := cmp.Diff(s.funcScopeTable, tt.want2); diff != "" {
				t.Errorf("SymbolTable.Define() funcScopeTable diff (-got +want2)\n%s", diff)
			}
			if diff := cmp.Diff(s.index, tt.want3); diff != "" {
				t.Errorf("SymbolTable.Define() index diff (-got +want3)\n%s", diff)
			}
		})
	}
}

func TestSymbolTable_LookUp(t *testing.T) {
	type fields struct {
		classScopeTable ScopedTable
		funcScopeTable  ScopedTable
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *SymbolTableEntry
	}{
		{
			name: "normal",
			fields: fields{
				classScopeTable: ScopedTable{
					"foo": &SymbolTableEntry{
						Name:  "foo",
						Typ:   "int",
						Kind:  Field,
						Index: 1,
					},
				},
				funcScopeTable: ScopedTable{
					"foo": &SymbolTableEntry{
						Name:  "foo",
						Typ:   "int",
						Kind:  Argument,
						Index: 2,
					},
				},
			},
			args: args{
				name: "foo",
			},
			want: &SymbolTableEntry{
				Name:  "foo",
				Typ:   "int",
				Kind:  Argument,
				Index: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SymbolTable{
				classScopeTable: tt.fields.classScopeTable,
				funcScopeTable:  tt.fields.funcScopeTable,
			}
			got := s.LookUp(tt.args.name)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SymbolTable.LookUp() got = %v, want %v", got, tt.want)
			}
		})
	}
}
