package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLeafNode_SetMeta(t *testing.T) {
	type fields struct {
		Typ NodeType
		V   string
	}
	type args struct {
		parent      NodeType
		grandParent NodeType
		s           *SymbolTableEntry
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    IDMeta
		wantErr bool
	}{
		{
			name: "ClassName",
			fields: fields{
				Typ: IdentifierType,
				V:   "id",
			},
			args: args{
				parent:      ClassNameType,
				grandParent: ClassType,
				s:           nil,
			},
			want: IDMeta{
				Category:    IdCatClass,
				Declaration: true,
				SymbolInfo:  nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &LeafNode{
				Typ: tt.fields.Typ,
				V:   tt.fields.V,
			}
			if err := n.SetMeta(tt.args.parent, tt.args.grandParent, tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("LeafNode.SetMeta() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(n.IDMeta, tt.want); diff != "" {
				t.Errorf("LeafNode.SetMeta() diff (-got +want)\n%s", diff)
			}
		})
	}
}
