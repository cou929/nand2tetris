package main

import (
	"reflect"
	"testing"
)

func TestCompiler_compileSubroutineCall(t *testing.T) {
	type fields struct {
		curClassName string
		curFuncInfo  *funcInfo
		vmc          *VmCode
	}
	type args struct {
		pt TreeNode
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:   "function",
			fields: fields{"FooClass", &funcInfo{name: "BarFunc"}, NewVmCode()},
			args: args{
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("MyClass"))}, ClassNameType, true),
					AdaptTokenToNode(SymbolToken(".")),
					MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("MyMethod"))}, SubroutineNameType, true),
					AdaptTokenToNode(SymbolToken("(")),
					MockNodes(nil, ExpressionListType, false),
					AdaptTokenToNode(SymbolToken(")")),
				}, SubroutineCallType, true),
			},
			want: []string{
				"call MyClass.MyMethod 0",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Compiler{
				curClassName: tt.fields.curClassName,
				curFuncInfo:  tt.fields.curFuncInfo,
				vmc:          tt.fields.vmc,
			}
			got, err := c.compileSubroutineCall(tt.args.pt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compiler.compileSubroutineCall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Compiler.compileSubroutineCall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompiler_compileExpression(t *testing.T) {
	type fields struct {
		vmc *VmCode
	}
	type args struct {
		pt TreeNode
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:   "int calculation",
			fields: fields{NewVmCode()},
			args: args{
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(10))}, TermType, false),
					MockNodes([]TreeNode{AdaptTokenToNode(SymbolToken("+"))}, OpType, true),
					MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(11))}, TermType, false),
					MockNodes([]TreeNode{AdaptTokenToNode(SymbolToken("+"))}, OpType, true),
					MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(12))}, TermType, false),
				}, ExpressionType, true),
			},
			want: []string{
				"push constant 10",
				"push constant 11",
				"push constant 12",
				"add",
				"add",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Compiler{
				vmc: tt.fields.vmc,
			}
			got, err := c.compileExpression(tt.args.pt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compiler.compileExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Compiler.compileExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}
