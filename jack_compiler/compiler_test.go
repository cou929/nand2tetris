package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Compiler.compileSubroutineCall() diff (-got +want)\n%s", diff)
			}
		})
	}
}

func TestCompiler_compileExpression(t *testing.T) {
	type args struct {
		pt TreeNode
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "int calculation",
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
		{
			name: "unaryOp",
			args: args{
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{AdaptTokenToNode(SymbolToken("-"))}, UnaryOpType, true),
						MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(123))}, TermType, true),
					}, TermType, false),
				}, ExpressionType, true),
			},
			want: []string{
				"push constant 123",
				"neg",
			},
			wantErr: false,
		},
		{
			name: "varName",
			args: args{
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{AdaptTokenToNodeWithMeta(IdentifierToken("x"), &IDMeta{Category: IdCatVar, SymbolInfo: &SymbolInfo{Index: 2}})}, VarNameType, true),
					}, TermType, false),
				}, ExpressionType, true),
			},
			want: []string{
				"push local 2",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Compiler{
				vmc: NewVmCode(),
			}
			got, err := c.compileExpression(tt.args.pt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compiler.compileExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Compiler.compileExpression() diff (-got +want)\n%s", diff)
			}
		})
	}
}

func TestCompiler_compileLetStatement(t *testing.T) {
	type args struct {
		pt TreeNode
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "assign to local variable",
			args: args{
				MockNodes([]TreeNode{
					AdaptTokenToNode(KeywordToken("let")),
					MockNodes([]TreeNode{AdaptTokenToNodeWithMeta(IdentifierToken("x"), &IDMeta{Category: IdCatVar, SymbolInfo: &SymbolInfo{Index: 2}})}, VarNameType, true),
					AdaptTokenToNode(SymbolToken("=")),
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(100))}, TermType, false),
					}, ExpressionType, false),
					AdaptTokenToNode(SymbolToken(";")),
				}, LetStatementType, true),
			},
			want: []string{
				"push constant 100",
				"pop local 2",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Compiler{
				vmc: NewVmCode(),
			}
			got, err := c.compileLetStatement(tt.args.pt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compiler.compileLetStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Compiler.compileLetStatement() diff (-got +want)\n%s", diff)
			}
		})
	}
}

func TestCompiler_compileIfStatement(t *testing.T) {
	type fields struct {
		curClassName string
		curFuncInfo  *funcInfo
		ifCounter    int
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
			name:   "if only",
			fields: fields{"MyClass", &funcInfo{name: "MyFunc"}, 5},
			args: args{
				MockNodes([]TreeNode{
					AdaptTokenToNode(KeywordToken("if")),
					AdaptTokenToNode(SymbolToken("(")),
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(100))}, TermType, false),
						MockNodes([]TreeNode{AdaptTokenToNode(SymbolToken(">"))}, OpType, true),
						MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(99))}, TermType, false),
					}, ExpressionType, false),
					AdaptTokenToNode(SymbolToken(")")),
					AdaptTokenToNode(SymbolToken("{")),
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{
							MockNodes([]TreeNode{
								AdaptTokenToNode(KeywordToken("let")),
								MockNodes([]TreeNode{AdaptTokenToNodeWithMeta(IdentifierToken("x"), &IDMeta{Category: IdCatVar, SymbolInfo: &SymbolInfo{Index: 2}})}, VarNameType, true),
								AdaptTokenToNode(SymbolToken("=")),
								MockNodes([]TreeNode{
									MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(100))}, TermType, false),
								}, ExpressionType, false),
								AdaptTokenToNode(SymbolToken(";")),
							}, LetStatementType, true),
						}, StatementType, true),
					}, StatementsType, false),
					AdaptTokenToNode(SymbolToken("}")),
				}, IfStatementType, true),
			},
			want: []string{
				"push constant 100",
				"push constant 99",
				"gt",
				"not",
				"if-goto MyClass.MyFunc.5.IF.ELSE",
				"push constant 100",
				"pop local 2",
				"goto MyClass.MyFunc.5.IF.END",
				"label MyClass.MyFunc.5.IF.ELSE",
				"label MyClass.MyFunc.5.IF.END",
			},
			wantErr: false,
		},
		{
			name:   "if only",
			fields: fields{"MyClass", &funcInfo{name: "MyFunc"}, 5},
			args: args{
				MockNodes([]TreeNode{
					AdaptTokenToNode(KeywordToken("if")),
					AdaptTokenToNode(SymbolToken("(")),
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(100))}, TermType, false),
						MockNodes([]TreeNode{AdaptTokenToNode(SymbolToken(">"))}, OpType, true),
						MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(99))}, TermType, false),
					}, ExpressionType, false),
					AdaptTokenToNode(SymbolToken(")")),
					AdaptTokenToNode(SymbolToken("{")),
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{
							MockNodes([]TreeNode{
								AdaptTokenToNode(KeywordToken("let")),
								MockNodes([]TreeNode{AdaptTokenToNodeWithMeta(IdentifierToken("x"), &IDMeta{Category: IdCatVar, SymbolInfo: &SymbolInfo{Index: 2}})}, VarNameType, true),
								AdaptTokenToNode(SymbolToken("=")),
								MockNodes([]TreeNode{
									MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(100))}, TermType, false),
								}, ExpressionType, false),
								AdaptTokenToNode(SymbolToken(";")),
							}, LetStatementType, true),
						}, StatementType, true),
					}, StatementsType, false),
					AdaptTokenToNode(SymbolToken("}")),
					AdaptTokenToNode(KeywordToken("else")),
					AdaptTokenToNode(SymbolToken("{")),
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{
							MockNodes([]TreeNode{
								AdaptTokenToNode(KeywordToken("let")),
								MockNodes([]TreeNode{AdaptTokenToNodeWithMeta(IdentifierToken("x"), &IDMeta{Category: IdCatVar, SymbolInfo: &SymbolInfo{Index: 2}})}, VarNameType, true),
								AdaptTokenToNode(SymbolToken("=")),
								MockNodes([]TreeNode{
									MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(200))}, TermType, false),
								}, ExpressionType, false),
								AdaptTokenToNode(SymbolToken(";")),
							}, LetStatementType, true),
						}, StatementType, true),
					}, StatementsType, false),
					AdaptTokenToNode(SymbolToken("}")),
				}, IfStatementType, true),
			},
			want: []string{
				"push constant 100",
				"push constant 99",
				"gt",
				"not",
				"if-goto MyClass.MyFunc.5.IF.ELSE",
				"push constant 100",
				"pop local 2",
				"goto MyClass.MyFunc.5.IF.END",
				"label MyClass.MyFunc.5.IF.ELSE",
				"push constant 200",
				"pop local 2",
				"label MyClass.MyFunc.5.IF.END",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Compiler{
				curClassName: tt.fields.curClassName,
				curFuncInfo:  tt.fields.curFuncInfo,
				ifCounter:    tt.fields.ifCounter,
				vmc:          NewVmCode(),
			}
			got, err := c.compileIfStatement(tt.args.pt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compiler.compileIfStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Compiler.compileIfStatement() diff (-got +want)\n%s", diff)
			}
		})
	}
}

func TestCompiler_compileWhileStatement(t *testing.T) {
	type fields struct {
		curClassName string
		curFuncInfo  *funcInfo
		whileCounter int
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
			name:   "if only",
			fields: fields{"MyClass", &funcInfo{name: "MyFunc"}, 5},
			args: args{
				MockNodes([]TreeNode{
					AdaptTokenToNode(KeywordToken("while")),
					AdaptTokenToNode(SymbolToken("(")),
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{
							MockNodes([]TreeNode{AdaptTokenToNodeWithMeta(IdentifierToken("x"), &IDMeta{Category: IdCatVar, SymbolInfo: &SymbolInfo{Index: 2}})}, VarNameType, true),
						}, TermType, false),
						MockNodes([]TreeNode{AdaptTokenToNode(SymbolToken(">"))}, OpType, true),
						MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(10))}, TermType, false),
					}, ExpressionType, false),
					AdaptTokenToNode(SymbolToken(")")),
					AdaptTokenToNode(SymbolToken("{")),
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{
							MockNodes([]TreeNode{
								AdaptTokenToNode(KeywordToken("let")),
								MockNodes([]TreeNode{AdaptTokenToNodeWithMeta(IdentifierToken("x"), &IDMeta{Category: IdCatVar, SymbolInfo: &SymbolInfo{Index: 2}})}, VarNameType, true),
								AdaptTokenToNode(SymbolToken("=")),
								MockNodes([]TreeNode{
									MockNodes([]TreeNode{
										MockNodes([]TreeNode{AdaptTokenToNodeWithMeta(IdentifierToken("x"), &IDMeta{Category: IdCatVar, SymbolInfo: &SymbolInfo{Index: 2}})}, VarNameType, true),
									}, TermType, false),
									MockNodes([]TreeNode{AdaptTokenToNode(SymbolToken("+"))}, OpType, true),
									MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(1))}, TermType, false),
								}, ExpressionType, false),
								AdaptTokenToNode(SymbolToken(";")),
							}, LetStatementType, true),
						}, StatementType, true),
					}, StatementsType, false),
					AdaptTokenToNode(SymbolToken("}")),
				}, IfStatementType, true),
			},
			want: []string{
				"label MyClass.MyFunc.5.WHILE.CONT",
				"push local 2",
				"push constant 10",
				"gt",
				"not",
				"if-goto MyClass.MyFunc.5.WHILE.END",
				"push local 2",
				"push constant 1",
				"add",
				"pop local 2",
				"goto MyClass.MyFunc.5.WHILE.CONT",
				"label MyClass.MyFunc.5.WHILE.END",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Compiler{
				curClassName: tt.fields.curClassName,
				curFuncInfo:  tt.fields.curFuncInfo,
				whileCounter: tt.fields.whileCounter,
				vmc:          NewVmCode(),
			}
			got, err := c.compileWhileStatement(tt.args.pt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compiler.compileWhileStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Compiler.compileWhileStatement() diff (-got +want)\n%s", diff)
			}
		})
	}
}
