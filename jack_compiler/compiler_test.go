package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCompiler_compileSubroutineCall(t *testing.T) {
	type fields struct {
		curClassInfo *classInfo
		curFuncInfo  *funcInfo
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
			fields: fields{&classInfo{name: "FooClass"}, &funcInfo{name: "BarFunc"}},
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
		{
			name:   "current instance method call",
			fields: fields{&classInfo{name: "FooClass"}, &funcInfo{name: "BarFunc"}},
			args: args{
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("MyMethod"))}, SubroutineNameType, true),
					AdaptTokenToNode(SymbolToken("(")),
					MockNodes(nil, ExpressionListType, false),
					AdaptTokenToNode(SymbolToken(")")),
				}, SubroutineCallType, true),
			},
			want: []string{
				"push pointer 0",
				"call FooClass.MyMethod 1",
			},
			wantErr: false,
		},
		{
			name:   "method call from var",
			fields: fields{&classInfo{name: "FooClass"}, &funcInfo{name: "BarFunc"}},
			args: args{
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{AdaptTokenToNodeWithMeta(IdentifierToken("myObj"), &IDMeta{Category: IdCatVar, SymbolInfo: &SymbolInfo{Index: 2, Type: "SomeClass"}})}, VarNameType, true),
					AdaptTokenToNode(SymbolToken(".")),
					MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("MyMethod"))}, SubroutineNameType, true),
					AdaptTokenToNode(SymbolToken("(")),
					MockNodes(nil, ExpressionListType, false),
					AdaptTokenToNode(SymbolToken(")")),
				}, SubroutineCallType, true),
			},
			want: []string{
				"push local 2",
				"call SomeClass.MyMethod 1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Compiler{
				curClassInfo: tt.fields.curClassInfo,
				curFuncInfo:  tt.fields.curFuncInfo,
				vmc:          NewVmCode(),
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
		{
			name: "this",
			args: args{
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{AdaptTokenToNode(KeywordToken("this"))}, KeywordConstantType, true),
					}, TermType, false),
				}, ExpressionType, true),
			},
			want: []string{
				"push pointer 0",
			},
			wantErr: false,
		},
		{
			name: "array",
			args: args{
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{AdaptTokenToNodeWithMeta(IdentifierToken("arr"), &IDMeta{Category: IdCatVar, SymbolInfo: &SymbolInfo{Index: 2}})}, VarNameType, true),
						MockNodes([]TreeNode{AdaptTokenToNode(SymbolToken("["))}, SymbolType, true),
						MockNodes([]TreeNode{
							MockNodes([]TreeNode{
								MockNodes([]TreeNode{AdaptTokenToNodeWithMeta(IdentifierToken("i"), &IDMeta{Category: IdCatVar, SymbolInfo: &SymbolInfo{Index: 3}})}, VarNameType, true),
							}, TermType, false),
						}, ExpressionType, true),
						MockNodes([]TreeNode{AdaptTokenToNode(SymbolToken("]"))}, SymbolType, true),
					}, TermType, false),
				}, ExpressionType, true),
			},
			want: []string{
				"push local 2",
				"push local 3",
				"add",
				"pop pointer 1",
				"push that 0",
			},
			wantErr: false,
		},
		{
			name: "string",
			args: args{
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{AdaptTokenToNode(StrConstToken("abc"))}, TermType, false),
				}, ExpressionType, true),
			},
			want: []string{
				"push constant 3",
				"call String.new 1",
				"push constant 97",
				"call String.appendChar 2",
				"push constant 98",
				"call String.appendChar 2",
				"push constant 99",
				"call String.appendChar 2",
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
		{
			name: "array",
			args: args{
				MockNodes([]TreeNode{
					AdaptTokenToNode(KeywordToken("let")),
					MockNodes([]TreeNode{AdaptTokenToNodeWithMeta(IdentifierToken("x"), &IDMeta{Category: IdCatVar, SymbolInfo: &SymbolInfo{Index: 2}})}, VarNameType, true),
					AdaptTokenToNode(SymbolToken("[")),
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{
							MockNodes([]TreeNode{AdaptTokenToNodeWithMeta(IdentifierToken("i"), &IDMeta{Category: IdCatVar, SymbolInfo: &SymbolInfo{Index: 3}})}, VarNameType, true),
						}, TermType, false),
					}, ExpressionType, true),
					AdaptTokenToNode(SymbolToken("]")),
					AdaptTokenToNode(SymbolToken("=")),
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(100))}, TermType, false),
					}, ExpressionType, false),
					AdaptTokenToNode(SymbolToken(";")),
				}, LetStatementType, true),
			},
			want: []string{
				"push local 2",
				"push local 3",
				"add",
				"pop pointer 1",
				"push constant 100",
				"pop that 0",
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
		curClassInfo *classInfo
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
			fields: fields{&classInfo{name: "MyClass"}, &funcInfo{name: "MyFunc"}, 5},
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
			name:   "if and else block",
			fields: fields{&classInfo{name: "MyClass"}, &funcInfo{name: "MyFunc"}, 5},
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
		{
			name:   "nest",
			fields: fields{&classInfo{name: "MyClass"}, &funcInfo{name: "MyFunc"}, 5},
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
								AdaptTokenToNode(KeywordToken("if")),
								AdaptTokenToNode(SymbolToken("(")),
								MockNodes([]TreeNode{
									MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(10))}, TermType, false),
									MockNodes([]TreeNode{AdaptTokenToNode(SymbolToken(">"))}, OpType, true),
									MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(9))}, TermType, false),
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
				// inner
				"push constant 10",
				"push constant 9",
				"gt",
				"not",
				"if-goto MyClass.MyFunc.6.IF.ELSE",
				"push constant 100",
				"pop local 2",
				"goto MyClass.MyFunc.6.IF.END",
				"label MyClass.MyFunc.6.IF.ELSE",
				"label MyClass.MyFunc.6.IF.END",
				// end of inner
				"goto MyClass.MyFunc.5.IF.END",
				"label MyClass.MyFunc.5.IF.ELSE",
				"label MyClass.MyFunc.5.IF.END",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Compiler{
				curClassInfo: tt.fields.curClassInfo,
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
		curClassInfo *classInfo
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
			name:   "normal",
			fields: fields{&classInfo{name: "MyClass"}, &funcInfo{name: "MyFunc"}, 5},
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
				}, WhileStatementType, true),
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
		{
			name:   "nest",
			fields: fields{&classInfo{name: "MyClass"}, &funcInfo{name: "MyFunc"}, 5},
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
								AdaptTokenToNode(KeywordToken("while")),
								AdaptTokenToNode(SymbolToken("(")),
								MockNodes([]TreeNode{
									MockNodes([]TreeNode{
										MockNodes([]TreeNode{AdaptTokenToNodeWithMeta(IdentifierToken("y"), &IDMeta{Category: IdCatVar, SymbolInfo: &SymbolInfo{Index: 3}})}, VarNameType, true),
									}, TermType, false),
									MockNodes([]TreeNode{AdaptTokenToNode(SymbolToken(">"))}, OpType, true),
									MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(100))}, TermType, false),
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
													MockNodes([]TreeNode{AdaptTokenToNodeWithMeta(IdentifierToken("2"), &IDMeta{Category: IdCatVar, SymbolInfo: &SymbolInfo{Index: 2}})}, VarNameType, true),
												}, TermType, false),
												MockNodes([]TreeNode{AdaptTokenToNode(SymbolToken("+"))}, OpType, true),
												MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(1))}, TermType, false),
											}, ExpressionType, false),
											AdaptTokenToNode(SymbolToken(";")),
										}, LetStatementType, true),
									}, StatementType, true),
									MockNodes([]TreeNode{
										MockNodes([]TreeNode{
											AdaptTokenToNode(KeywordToken("let")),
											MockNodes([]TreeNode{AdaptTokenToNodeWithMeta(IdentifierToken("y"), &IDMeta{Category: IdCatVar, SymbolInfo: &SymbolInfo{Index: 3}})}, VarNameType, true),
											AdaptTokenToNode(SymbolToken("=")),
											MockNodes([]TreeNode{
												MockNodes([]TreeNode{
													MockNodes([]TreeNode{AdaptTokenToNodeWithMeta(IdentifierToken("y"), &IDMeta{Category: IdCatVar, SymbolInfo: &SymbolInfo{Index: 3}})}, VarNameType, true),
												}, TermType, false),
												MockNodes([]TreeNode{AdaptTokenToNode(SymbolToken("+"))}, OpType, true),
												MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(1))}, TermType, false),
											}, ExpressionType, false),
											AdaptTokenToNode(SymbolToken(";")),
										}, LetStatementType, true),
									}, StatementType, true),
								}, StatementsType, false),
								AdaptTokenToNode(SymbolToken("}")),
							}, WhileStatementType, true),
						}, StatementType, true),
					}, StatementsType, false),
					AdaptTokenToNode(SymbolToken("}")),
				}, WhileStatementType, true),
			},
			want: []string{
				"label MyClass.MyFunc.5.WHILE.CONT",
				"push local 2",
				"push constant 10",
				"gt",
				"not",
				"if-goto MyClass.MyFunc.5.WHILE.END",
				// inner
				"label MyClass.MyFunc.6.WHILE.CONT",
				"push local 3",
				"push constant 100",
				"gt",
				"not",
				"if-goto MyClass.MyFunc.6.WHILE.END",
				"push local 2",
				"push constant 1",
				"add",
				"pop local 2",
				"push local 3",
				"push constant 1",
				"add",
				"pop local 3",
				"goto MyClass.MyFunc.6.WHILE.CONT",
				"label MyClass.MyFunc.6.WHILE.END",
				// end of inner
				"goto MyClass.MyFunc.5.WHILE.CONT",
				"label MyClass.MyFunc.5.WHILE.END",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Compiler{
				curClassInfo: tt.fields.curClassInfo,
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

func TestCompiler_compileReturnStatement(t *testing.T) {
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
			name: "return with expression",
			args: args{
				MockNodes([]TreeNode{
					AdaptTokenToNode(KeywordToken("return")),
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(10))}, TermType, false),
					}, ExpressionType, false),
					AdaptTokenToNode(SymbolToken(";")),
				}, ReturnStatementType, true),
			},
			want: []string{
				"push constant 10",
				"return",
			},
			wantErr: false,
		},
		{
			name: "void return",
			args: args{
				MockNodes([]TreeNode{
					AdaptTokenToNode(KeywordToken("return")),
					AdaptTokenToNode(SymbolToken(";")),
				}, ReturnStatementType, true),
			},
			want: []string{
				"push constant 0",
				"return",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Compiler{
				vmc: NewVmCode(),
			}
			got, err := c.compileReturnStatement(tt.args.pt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compiler.compileReturnStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Compiler.compileReturnStatement() diff (-got +want)\n%s", diff)
			}
		})
	}
}

func TestCompiler_compileSubroutineDec(t *testing.T) {
	type fields struct {
		curClassInfo *classInfo
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
			name: "constructor",
			fields: fields{
				curClassInfo: &classInfo{name: "MyClass", fieldCount: 3},
			},
			args: args{
				MockNodes([]TreeNode{
					AdaptTokenToNode(KeywordToken("constructor")),
					MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("MyClass"))}, TypeType, true),
					MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("new"))}, SubroutineNameType, true),
					AdaptTokenToNode(SymbolToken("(")),
					MockNodes([]TreeNode{}, ParameterListType, false),
					AdaptTokenToNode(SymbolToken(")")),
					MockNodes([]TreeNode{
						AdaptTokenToNode(SymbolToken("{")),
						MockNodes([]TreeNode{
							MockNodes([]TreeNode{
								MockNodes([]TreeNode{
									AdaptTokenToNode(KeywordToken("return")),
									MockNodes([]TreeNode{
										MockNodes([]TreeNode{
											MockNodes([]TreeNode{AdaptTokenToNode(KeywordToken("this"))}, KeywordConstantType, true),
										}, TermType, false),
									}, ExpressionType, true),
									AdaptTokenToNode(SymbolToken(";")),
								}, ReturnStatementType, false),
							}, StatementType, true),
						}, StatementsType, false),
						AdaptTokenToNode(SymbolToken("}")),
					}, SubroutineBodyType, false),
				}, SubroutineDecType, false),
			},
			want: []string{
				"function MyClass.new 0",
				"push constant 3",
				"call Memory.alloc 1",
				"pop pointer 0",
				"push pointer 0",
				"return",
			},
			wantErr: false,
		},
		{
			name: "function",
			fields: fields{
				curClassInfo: &classInfo{name: "MyClass", fieldCount: 3},
			},
			args: args{
				MockNodes([]TreeNode{
					AdaptTokenToNode(KeywordToken("function")),
					MockNodes([]TreeNode{AdaptTokenToNode(KeywordToken("void"))}, TypeType, true),
					MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("MyFunc"))}, SubroutineNameType, true),
					AdaptTokenToNode(SymbolToken("(")),
					MockNodes([]TreeNode{}, ParameterListType, false),
					AdaptTokenToNode(SymbolToken(")")),
					MockNodes([]TreeNode{
						AdaptTokenToNode(SymbolToken("{")),
						MockNodes([]TreeNode{
							AdaptTokenToNode(KeywordToken("var")),
							MockNodes([]TreeNode{AdaptTokenToNode(KeywordToken("boolean"))}, TypeType, true),
							MockNodes([]TreeNode{AdaptTokenToNodeWithMeta(IdentifierToken("x"), &IDMeta{Category: IdCatVar, SymbolInfo: &SymbolInfo{Index: 2}})}, VarNameType, true),
							AdaptTokenToNode(SymbolToken(";")),
						}, VarDecType, false),
						MockNodes([]TreeNode{
							MockNodes([]TreeNode{
								MockNodes([]TreeNode{
									AdaptTokenToNode(KeywordToken("return")),
									AdaptTokenToNode(SymbolToken(";")),
								}, ReturnStatementType, false),
							}, StatementType, true),
						}, StatementsType, false),
						AdaptTokenToNode(SymbolToken("}")),
					}, SubroutineBodyType, false),
				}, SubroutineDecType, false),
			},
			want: []string{
				"function MyClass.MyFunc 1",
				"push constant 0",
				"return",
			},
			wantErr: false,
		},
		{
			name: "method",
			fields: fields{
				curClassInfo: &classInfo{name: "MyClass", fieldCount: 3},
			},
			args: args{
				MockNodes([]TreeNode{
					AdaptTokenToNode(KeywordToken("method")),
					MockNodes([]TreeNode{AdaptTokenToNode(KeywordToken("void"))}, TypeType, true),
					MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("myFunc"))}, SubroutineNameType, true),
					AdaptTokenToNode(SymbolToken("(")),
					MockNodes([]TreeNode{}, ParameterListType, false),
					AdaptTokenToNode(SymbolToken(")")),
					MockNodes([]TreeNode{
						AdaptTokenToNode(SymbolToken("{")),
						MockNodes([]TreeNode{
							MockNodes([]TreeNode{
								MockNodes([]TreeNode{
									AdaptTokenToNode(KeywordToken("return")),
									AdaptTokenToNode(SymbolToken(";")),
								}, ReturnStatementType, false),
							}, StatementType, true),
						}, StatementsType, false),
						AdaptTokenToNode(SymbolToken("}")),
					}, SubroutineBodyType, false),
				}, SubroutineDecType, false),
			},
			want: []string{
				"function MyClass.myFunc 0",
				"push argument 0",
				"pop pointer 0",
				"push constant 0",
				"return",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Compiler{
				curClassInfo: tt.fields.curClassInfo,
				vmc:          NewVmCode(),
			}
			got, err := c.compileSubroutineDec(tt.args.pt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compiler.compileSubroutineDec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Compiler.compileSubroutineDec() diff (-got +want)\n%s", diff)
			}
		})
	}
}
