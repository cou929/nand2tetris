package main

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func MockNodes(nodes []TreeNode, typ NodeType, terminal bool) TreeNode {
	s := typ.String()
	s = s[0:strings.LastIndex(s, "Type")]
	r := []rune(s)
	head := strings.ToLower(string(r[0]))
	rest := string(r[1:len(r)])
	name := head + rest
	switch typ {
	case TypeType, ClassNameType, SubroutineNameType, VarNameType, OpType, UnaryOpType, KeywordConstantType:
		return &OneChildNode{
			Children:  nodes,
			Typ:       typ,
			N:         name,
			XMLMarkup: !terminal,
		}
	}
	return &InnerNode{
		Children:  nodes,
		Typ:       typ,
		N:         name,
		XMLMarkup: !terminal,
	}
}

func TestParser_parseKeywordConstant(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		args    args
		want    TreeNode
		want1   TokenList
		wantErr bool
	}{
		{
			name: "normal",
			args: args{[]Token{
				KeywordToken("true"),
				SymbolToken(";"),
			}},
			want: MockNodes([]TreeNode{AdaptTokenToNode(KeywordToken("true"))}, KeywordConstantType, true),
			want1: []Token{
				SymbolToken(";"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseKeywordConstant(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseKeywordConstant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseKeywordConstant() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseKeywordConstant() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}

func TestParser_parseUnaryOp(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    TreeNode
		want1   TokenList
		wantErr bool
	}{
		{
			name: "normal",
			args: args{[]Token{
				SymbolToken("-"),
				IdentifierToken("i"),
				SymbolToken(";"),
			}},
			want: MockNodes([]TreeNode{AdaptTokenToNode(SymbolToken("-"))}, UnaryOpType, true),
			want1: []Token{
				IdentifierToken("i"),
				SymbolToken(";"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseUnaryOp(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseUnaryOp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseUnaryOp() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseUnaryOp() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}

func TestParser_parseOp(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    TreeNode
		want1   TokenList
		wantErr bool
	}{
		{
			name: "normal",
			args: args{[]Token{
				SymbolToken("+"),
				IdentifierToken("123"),
				SymbolToken(";"),
			}},
			want: MockNodes([]TreeNode{AdaptTokenToNode(SymbolToken("+"))}, OpType, true),
			want1: []Token{
				IdentifierToken("123"),
				SymbolToken(";"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseOp(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseOp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseOp() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseOp() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}

func TestParser_parseTerm(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    TreeNode
		want1   TokenList
		wantErr bool
	}{
		{
			name: "int const",
			args: args{[]Token{
				IntConstToken(123),
				SymbolToken("+"),
				IntConstToken(456),
			}},
			want: MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(123))}, TermType, false),
			want1: []Token{
				SymbolToken("+"),
				IntConstToken(456),
			},
			wantErr: false,
		},
		{
			name: "str const",
			args: args{[]Token{
				StrConstToken("string"),
				SymbolToken(";"),
			}},
			want: MockNodes([]TreeNode{AdaptTokenToNode(StrConstToken("string"))}, TermType, false),
			want1: []Token{
				SymbolToken(";"),
			},
			wantErr: false,
		},
		{
			name: "keyword const",
			args: args{[]Token{
				KeywordToken("true"),
				SymbolToken(";"),
			}},
			want: MockNodes([]TreeNode{MockNodes([]TreeNode{AdaptTokenToNode(KeywordToken("true"))}, KeywordConstantType, true)}, TermType, false),
			want1: []Token{
				SymbolToken(";"),
			},
			wantErr: false,
		},
		{
			name: "varName only",
			args: args{[]Token{
				IdentifierToken("args"),
				SymbolToken(";"),
			}},
			want: MockNodes([]TreeNode{MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("args"))}, VarNameType, true)}, TermType, false),
			want1: []Token{
				SymbolToken(";"),
			},
			wantErr: false,
		},
		{
			name: "varName with array index and no expression",
			args: args{[]Token{
				IdentifierToken("args"),
				SymbolToken("["),
				IdentifierToken("i"),
				SymbolToken("]"),
				SymbolToken(";"),
			}},
			want: MockNodes(
				[]TreeNode{
					MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("args"))}, VarNameType, true),
					AdaptTokenToNode(SymbolToken("[")),
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{
							MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("i"))}, VarNameType, true),
						}, TermType, false),
					}, ExpressionType, false),
					AdaptTokenToNode(SymbolToken("]")),
				},
				TermType, false),
			want1: []Token{
				SymbolToken(";"),
			},
			wantErr: false,
		},
		{
			name: "subroutineCall",
			args: args{[]Token{
				IdentifierToken("fooFunc"),
				SymbolToken("("),
				IdentifierToken("x"),
				SymbolToken(")"),
				SymbolToken(";"),
			}},
			want: MockNodes([]TreeNode{
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("fooFunc"))}, SubroutineNameType, true),
					AdaptTokenToNode(SymbolToken("(")),
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{
							MockNodes([]TreeNode{
								MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
							}, TermType, false),
						}, ExpressionType, false),
					}, ExpressionListType, false),
					AdaptTokenToNode(SymbolToken(")")),
				}, SubroutineCallType, true),
			}, TermType, false),
			want1: []Token{
				SymbolToken(";"),
			},
			wantErr: false,
		},
		{
			name: "expression enclosed in paren",
			args: args{[]Token{
				SymbolToken("("),
				IdentifierToken("ok"),
				SymbolToken(")"),
				SymbolToken(";"),
			}},
			want: MockNodes([]TreeNode{
				AdaptTokenToNode(SymbolToken("(")),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("ok"))}, VarNameType, true),
					}, TermType, false),
				}, ExpressionType, false),
				AdaptTokenToNode(SymbolToken(")")),
			}, TermType, false),
			want1: []Token{
				SymbolToken(";"),
			},
			wantErr: false,
		},
		{
			name: "unary operator",
			args: args{[]Token{
				SymbolToken("~"),
				IdentifierToken("some"),
				SymbolToken(";"),
			}},
			want: MockNodes([]TreeNode{
				MockNodes([]TreeNode{AdaptTokenToNode(SymbolToken("~"))}, UnaryOpType, true),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("some"))}, VarNameType, true),
				}, TermType, false),
			}, TermType, false),
			want1: []Token{
				SymbolToken(";"),
			},
			wantErr: false,
		},
		{
			name: "invalid",
			args: args{[]Token{
				SymbolToken(")"),
			}},
			want:    (*InnerNode)(nil),
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseTerm(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseTerm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseTerm() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseTerm() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}

func TestParser_parseVarName(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    TreeNode
		want1   TokenList
		wantErr bool
	}{
		{
			name: "normal",
			args: args{[]Token{
				IdentifierToken("args"),
				SymbolToken(";"),
			}},
			want: MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("args"))}, VarNameType, true),
			want1: []Token{
				SymbolToken(";"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseVarName(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseVarName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseVarName() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseVarName() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}

func TestParser_parseExpression(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    TreeNode
		want1   TokenList
		wantErr bool
	}{
		{
			name: "one term only",
			args: args{[]Token{
				IdentifierToken("x"),
				SymbolToken("]"),
			}},
			want: MockNodes([]TreeNode{
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						AdaptTokenToNode(IdentifierToken("x")),
					}, VarNameType, true),
				}, TermType, false),
			}, ExpressionType, false),
			want1: []Token{
				SymbolToken("]"),
			},
			wantErr: false,
		},
		{
			name: "op",
			args: args{[]Token{
				IdentifierToken("x"),
				SymbolToken("|"),
				IdentifierToken("y"),
				SymbolToken(";"),
			}},
			want: MockNodes([]TreeNode{
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
				}, TermType, false),
				MockNodes([]TreeNode{AdaptTokenToNode(SymbolToken("|"))}, OpType, true),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("y"))}, VarNameType, true),
				}, TermType, false),
			}, ExpressionType, false),
			want1: []Token{
				SymbolToken(";"),
			},
			wantErr: false,
		},
		{
			name: "multiple op",
			args: args{[]Token{
				IdentifierToken("x"),
				SymbolToken("+"),
				IdentifierToken("y"),
				SymbolToken("+"),
				IntConstToken(10),
				SymbolToken(";"),
			}},
			want: MockNodes([]TreeNode{
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
				}, TermType, false),
				MockNodes([]TreeNode{AdaptTokenToNode(SymbolToken("+"))}, OpType, true),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("y"))}, VarNameType, true),
				}, TermType, false),
				MockNodes([]TreeNode{AdaptTokenToNode(SymbolToken("+"))}, OpType, true),
				MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(10))}, TermType, false),
			}, ExpressionType, false),
			want1: []Token{
				SymbolToken(";"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseExpression(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseExpression() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseExpression() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}

func TestParser_parseExpressionList(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    TreeNode
		want1   TokenList
		wantErr bool
	}{
		{
			name: "no expression",
			args: args{[]Token{
				SymbolToken(")"),
			}},
			want: MockNodes([]TreeNode(nil), ExpressionListType, false),
			want1: []Token{
				SymbolToken(")"),
			},
			wantErr: false,
		},
		{
			name: "single expression",
			args: args{[]Token{
				IdentifierToken("x"),
				SymbolToken(")"),
			}},
			want: MockNodes([]TreeNode{
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
					}, TermType, false),
				}, ExpressionType, false),
			}, ExpressionListType, false),
			want1: []Token{
				SymbolToken(")"),
			},
			wantErr: false,
		},
		{
			name: "multiple expression",
			args: args{[]Token{
				IdentifierToken("x"),
				SymbolToken(","),
				IdentifierToken("y"),
				SymbolToken(")"),
			}},
			want: MockNodes([]TreeNode{
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
					}, TermType, false),
				}, ExpressionType, false),
				AdaptTokenToNode(SymbolToken(",")),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("y"))}, VarNameType, true),
					}, TermType, false),
				}, ExpressionType, false),
			}, ExpressionListType, false),
			want1: []Token{
				SymbolToken(")"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseExpressionList(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseExpressionList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseExpressionList() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseExpressionList() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}

func TestParser_parseSubroutineName(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    TreeNode
		want1   TokenList
		wantErr bool
	}{
		{
			name: "normal",
			args: args{[]Token{
				IdentifierToken("myFunc"),
				SymbolToken("("),
			}},
			want: MockNodes([]TreeNode{
				AdaptTokenToNode(IdentifierToken("myFunc")),
			}, SubroutineNameType, true),
			want1: []Token{
				SymbolToken("("),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseSubroutineName(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseSubroutineName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseSubroutineName() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseSubroutineName() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}

func TestParser_parseClassName(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    TreeNode
		want1   TokenList
		wantErr bool
	}{
		{
			name: "normal",
			args: args{[]Token{
				IdentifierToken("MyClass"),
				SymbolToken("."),
			}},
			want: MockNodes([]TreeNode{
				AdaptTokenToNode(IdentifierToken("MyClass")),
			}, ClassNameType, true),
			want1: []Token{
				SymbolToken("."),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseClassName(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseClassName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseClassName() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseClassName() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}

func TestParser_parseSubroutineCall(t *testing.T) {
	type fields struct {
		symbolTable *SymbolTable
	}
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    TreeNode
		want1   TokenList
		wantErr bool
	}{
		{
			name: "function call",
			fields: fields{
				symbolTable: &SymbolTable{
					classScopeTable: ScopedTable{},
					funcScopeTable: ScopedTable{
						"x": &SymbolTableEntry{
							Name:  "x",
							Typ:   "int",
							Kind:  Var,
							Index: 0,
						},
						"y": &SymbolTableEntry{
							Name:  "y",
							Typ:   "int",
							Kind:  Var,
							Index: 0,
						},
					},
					index: map[VarKind]int{
						Var: 1,
					},
				},
			},
			args: args{[]Token{
				IdentifierToken("MyFunc"),
				SymbolToken("("),
				IdentifierToken("x"),
				SymbolToken(","),
				IdentifierToken("y"),
				SymbolToken(")"),
				SymbolToken(";"),
			}},
			want: MockNodes([]TreeNode{
				MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("MyFunc"))}, SubroutineNameType, true),
				AdaptTokenToNode(SymbolToken("(")),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{
							MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
						}, TermType, false),
					}, ExpressionType, false),
					AdaptTokenToNode(SymbolToken(",")),
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{
							MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("y"))}, VarNameType, true),
						}, TermType, false),
					}, ExpressionType, false),
				}, ExpressionListType, false),
				AdaptTokenToNode(SymbolToken(")")),
			}, SubroutineCallType, true),
			want1: []Token{
				SymbolToken(";"),
			},
			wantErr: false,
		},
		{
			name:   "function call with no argument",
			fields: fields{&SymbolTable{}},
			args: args{[]Token{
				IdentifierToken("MyFunc"),
				SymbolToken("("),
				SymbolToken(")"),
				SymbolToken(";"),
			}},
			want: MockNodes([]TreeNode{
				MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("MyFunc"))}, SubroutineNameType, true),
				AdaptTokenToNode(SymbolToken("(")),
				MockNodes(nil, ExpressionListType, false),
				AdaptTokenToNode(SymbolToken(")")),
			}, SubroutineCallType, true),
			want1: []Token{
				SymbolToken(";"),
			},
			wantErr: false,
		},
		{
			name:   "method call",
			fields: fields{&SymbolTable{}},
			args: args{[]Token{
				IdentifierToken("MyClass"),
				SymbolToken("."),
				IdentifierToken("MyMethod"),
				SymbolToken("("),
				SymbolToken(")"),
				SymbolToken(";"),
			}},
			want: MockNodes([]TreeNode{
				MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("MyClass"))}, ClassNameType, true),
				AdaptTokenToNode(SymbolToken(".")),
				MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("MyMethod"))}, SubroutineNameType, true),
				AdaptTokenToNode(SymbolToken("(")),
				MockNodes(nil, ExpressionListType, false),
				AdaptTokenToNode(SymbolToken(")")),
			}, SubroutineCallType, true),
			want1: []Token{
				SymbolToken(";"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				symbolTable: tt.fields.symbolTable,
			}
			got, got1, err := p.parseSubroutineCall(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseSubroutineCall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseSubroutineCall() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseSubroutineCall() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}

func TestParser_parseReturnStatement(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    TreeNode
		want1   TokenList
		wantErr bool
	}{
		{
			name: "return only",
			args: args{[]Token{
				KeywordToken("return"),
				SymbolToken(";"),
				SymbolToken("}"),
			}},
			want: MockNodes([]TreeNode{
				AdaptTokenToNode(KeywordToken("return")),
				AdaptTokenToNode(SymbolToken(";")),
			}, ReturnStatementType, false),
			want1: []Token{
				SymbolToken("}"),
			},
			wantErr: false,
		},
		{
			name: "return with expression",
			args: args{[]Token{
				KeywordToken("return"),
				IdentifierToken("res"),
				SymbolToken(";"),
				SymbolToken("}"),
			}},
			want: MockNodes([]TreeNode{
				AdaptTokenToNode(KeywordToken("return")),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{
							AdaptTokenToNode(IdentifierToken("res")),
						}, VarNameType, true),
					}, TermType, false),
				}, ExpressionType, false),
				AdaptTokenToNode(SymbolToken(";")),
			}, ReturnStatementType, false),
			want1: []Token{
				SymbolToken("}"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseReturnStatement(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseReturnStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseReturnStatement() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseReturnStatement() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}

func TestParser_parseDoStatement(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    TreeNode
		want1   TokenList
		wantErr bool
	}{
		{
			name: "normal",
			args: args{[]Token{
				KeywordToken("do"),
				IdentifierToken("MyFunc"),
				SymbolToken("("),
				SymbolToken(")"),
				SymbolToken(";"),
				SymbolToken("}"),
			}},
			want: MockNodes([]TreeNode{
				AdaptTokenToNode(KeywordToken("do")),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("MyFunc"))}, SubroutineNameType, true),
					AdaptTokenToNode(SymbolToken("(")),
					MockNodes(nil, ExpressionListType, false),
					AdaptTokenToNode(SymbolToken(")")),
				}, SubroutineCallType, true),
				AdaptTokenToNode(SymbolToken(";")),
			}, DoStatementType, false),
			want1: []Token{
				SymbolToken("}"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseDoStatement(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseDoStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseDoStatement() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseDoStatement() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}

func TestParser_parseWhileStatement(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    TreeNode
		want1   TokenList
		wantErr bool
	}{
		{
			name: "normal",
			args: args{[]Token{
				KeywordToken("while"),
				SymbolToken("("),
				IdentifierToken("x"),
				SymbolToken(")"),
				SymbolToken("{"),
				KeywordToken("do"),
				IdentifierToken("MyFunc"),
				SymbolToken("("),
				SymbolToken(")"),
				SymbolToken(";"),
				KeywordToken("do"),
				IdentifierToken("MyFunc"),
				SymbolToken("("),
				SymbolToken(")"),
				SymbolToken(";"),
				SymbolToken("}"),
				KeywordToken("return"),
			}},
			want: MockNodes([]TreeNode{
				AdaptTokenToNode(KeywordToken("while")),
				AdaptTokenToNode(SymbolToken("(")),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
					}, TermType, false),
				}, ExpressionType, false),
				AdaptTokenToNode(SymbolToken(")")),
				AdaptTokenToNode(SymbolToken("{")),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{
							AdaptTokenToNode(KeywordToken("do")),
							MockNodes([]TreeNode{
								MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("MyFunc"))}, SubroutineNameType, true),
								AdaptTokenToNode(SymbolToken("(")),
								MockNodes(nil, ExpressionListType, false),
								AdaptTokenToNode(SymbolToken(")")),
							}, SubroutineCallType, true),
							AdaptTokenToNode(SymbolToken(";")),
						}, DoStatementType, false),
					}, StatementType, true),
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{
							AdaptTokenToNode(KeywordToken("do")),
							MockNodes([]TreeNode{
								MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("MyFunc"))}, SubroutineNameType, true),
								AdaptTokenToNode(SymbolToken("(")),
								MockNodes(nil, ExpressionListType, false),
								AdaptTokenToNode(SymbolToken(")")),
							}, SubroutineCallType, true),
							AdaptTokenToNode(SymbolToken(";")),
						}, DoStatementType, false),
					}, StatementType, true),
				}, StatementsType, false),
				AdaptTokenToNode(SymbolToken("}")),
			}, WhileStatementType, false),
			want1: []Token{
				KeywordToken("return"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseWhileStatement(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseWhileStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseWhileStatement() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseWhileStatement() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}

func TestParser_parseIfStatement(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    TreeNode
		want1   TokenList
		wantErr bool
	}{
		{
			name: "if block only",
			args: args{[]Token{
				KeywordToken("if"),
				SymbolToken("("),
				IdentifierToken("x"),
				SymbolToken(")"),
				SymbolToken("{"),
				KeywordToken("let"),
				IdentifierToken("x"),
				SymbolToken("="),
				IntConstToken(100),
				SymbolToken(";"),
				KeywordToken("do"),
				IdentifierToken("MyFunc"),
				SymbolToken("("),
				SymbolToken(")"),
				SymbolToken(";"),
				SymbolToken("}"),
				KeywordToken("return"),
			}},
			want: MockNodes([]TreeNode{
				AdaptTokenToNode(KeywordToken("if")),
				AdaptTokenToNode(SymbolToken("(")),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
					}, TermType, false),
				}, ExpressionType, false),
				AdaptTokenToNode(SymbolToken(")")),
				AdaptTokenToNode(SymbolToken("{")),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{
							AdaptTokenToNode(KeywordToken("let")),
							MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
							AdaptTokenToNode(SymbolToken("=")),
							MockNodes([]TreeNode{
								MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(100))}, TermType, false),
							}, ExpressionType, false),
							AdaptTokenToNode(SymbolToken(";")),
						}, LetStatementType, false),
					}, StatementType, true),
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{
							AdaptTokenToNode(KeywordToken("do")),
							MockNodes([]TreeNode{
								MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("MyFunc"))}, SubroutineNameType, true),
								AdaptTokenToNode(SymbolToken("(")),
								MockNodes(nil, ExpressionListType, false),
								AdaptTokenToNode(SymbolToken(")")),
							}, SubroutineCallType, true),
							AdaptTokenToNode(SymbolToken(";")),
						}, DoStatementType, false),
					}, StatementType, true),
				}, StatementsType, false),
				AdaptTokenToNode(SymbolToken("}")),
			}, IfStatementType, false),
			want1: []Token{
				KeywordToken("return"),
			},
			wantErr: false,
		},
		{
			name: "if and else block",
			args: args{[]Token{
				KeywordToken("if"),
				SymbolToken("("),
				IdentifierToken("x"),
				SymbolToken(")"),
				SymbolToken("{"),
				KeywordToken("do"),
				IdentifierToken("MyFunc"),
				SymbolToken("("),
				SymbolToken(")"),
				SymbolToken(";"),
				SymbolToken("}"),
				KeywordToken("else"),
				SymbolToken("{"),
				KeywordToken("let"),
				IdentifierToken("x"),
				SymbolToken("="),
				IntConstToken(100),
				SymbolToken(";"),
				KeywordToken("do"),
				IdentifierToken("MyFuncElse"),
				SymbolToken("("),
				SymbolToken(")"),
				SymbolToken(";"),
				SymbolToken("}"),
				KeywordToken("return"),
			}},
			want: MockNodes([]TreeNode{
				AdaptTokenToNode(KeywordToken("if")),
				AdaptTokenToNode(SymbolToken("(")),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
					}, TermType, false),
				}, ExpressionType, false),
				AdaptTokenToNode(SymbolToken(")")),
				AdaptTokenToNode(SymbolToken("{")),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{
							AdaptTokenToNode(KeywordToken("do")),
							MockNodes([]TreeNode{
								MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("MyFunc"))}, SubroutineNameType, true),
								AdaptTokenToNode(SymbolToken("(")),
								MockNodes(nil, ExpressionListType, false),
								AdaptTokenToNode(SymbolToken(")")),
							}, SubroutineCallType, true),
							AdaptTokenToNode(SymbolToken(";")),
						}, DoStatementType, false),
					}, StatementType, true),
				}, StatementsType, false),
				AdaptTokenToNode(SymbolToken("}")),
				AdaptTokenToNode(KeywordToken("else")),
				AdaptTokenToNode(SymbolToken("{")),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{
							AdaptTokenToNode(KeywordToken("let")),
							MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
							AdaptTokenToNode(SymbolToken("=")),
							MockNodes([]TreeNode{
								MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(100))}, TermType, false),
							}, ExpressionType, false),
							AdaptTokenToNode(SymbolToken(";")),
						}, LetStatementType, false),
					}, StatementType, true),
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{
							AdaptTokenToNode(KeywordToken("do")),
							MockNodes([]TreeNode{
								MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("MyFuncElse"))}, SubroutineNameType, true),
								AdaptTokenToNode(SymbolToken("(")),
								MockNodes(nil, ExpressionListType, false),
								AdaptTokenToNode(SymbolToken(")")),
							}, SubroutineCallType, true),
							AdaptTokenToNode(SymbolToken(";")),
						}, DoStatementType, false),
					}, StatementType, true),
				}, StatementsType, false),
				AdaptTokenToNode(SymbolToken("}")),
			}, IfStatementType, false),
			want1: []Token{
				KeywordToken("return"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseIfStatement(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseIfStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseIfStatement() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseIfStatement() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}

func TestParser_parseLetStatement(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    TreeNode
		want1   TokenList
		wantErr bool
	}{
		{
			name: "normal",
			args: args{[]Token{
				KeywordToken("let"),
				IdentifierToken("x"),
				SymbolToken("="),
				IntConstToken(100),
				SymbolToken(";"),
				KeywordToken("return"),
			}},
			want: MockNodes([]TreeNode{
				AdaptTokenToNode(KeywordToken("let")),
				MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
				AdaptTokenToNode(SymbolToken("=")),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(100))}, TermType, false),
				}, ExpressionType, false),
				AdaptTokenToNode(SymbolToken(";")),
			}, LetStatementType, false),
			want1: []Token{
				KeywordToken("return"),
			},
			wantErr: false,
		},
		{
			name: "assign to array index",
			args: args{[]Token{
				KeywordToken("let"),
				IdentifierToken("x"),
				SymbolToken("["),
				IntConstToken(3),
				SymbolToken("]"),
				SymbolToken("="),
				IntConstToken(100),
				SymbolToken(";"),
				KeywordToken("return"),
			}},
			want: MockNodes([]TreeNode{
				AdaptTokenToNode(KeywordToken("let")),
				MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
				AdaptTokenToNode(SymbolToken("[")),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(3))}, TermType, false),
				}, ExpressionType, false),
				AdaptTokenToNode(SymbolToken("]")),
				AdaptTokenToNode(SymbolToken("=")),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(100))}, TermType, false),
				}, ExpressionType, false),
				AdaptTokenToNode(SymbolToken(";")),
			}, LetStatementType, false),
			want1: []Token{
				KeywordToken("return"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseLetStatement(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseLetStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseLetStatement() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseLetStatement() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}

func TestParser_parseStatements(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    TreeNode
		want1   TokenList
		wantErr bool
	}{
		{
			name: "no statement",
			args: args{[]Token{
				KeywordToken("return"),
			}},
			want: MockNodes(nil, StatementsType, false),
			want1: []Token{
				KeywordToken("return"),
			},
			wantErr: false,
		},
		{
			name: "single statement",
			args: args{[]Token{
				KeywordToken("let"),
				IdentifierToken("x"),
				SymbolToken("="),
				IntConstToken(100),
				SymbolToken(";"),
				KeywordToken("return"),
			}},
			want: MockNodes([]TreeNode{
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						AdaptTokenToNode(KeywordToken("let")),
						MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
						AdaptTokenToNode(SymbolToken("=")),
						MockNodes([]TreeNode{
							MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(100))}, TermType, false),
						}, ExpressionType, false),
						AdaptTokenToNode(SymbolToken(";")),
					}, LetStatementType, false),
				}, StatementType, true),
			}, StatementsType, false),
			want1: []Token{
				KeywordToken("return"),
			},
			wantErr: false,
		},
		{
			name: "multi statements",
			args: args{[]Token{
				KeywordToken("let"),
				IdentifierToken("x"),
				SymbolToken("="),
				IntConstToken(100),
				SymbolToken(";"),
				KeywordToken("let"),
				IdentifierToken("y"),
				SymbolToken("="),
				IntConstToken(200),
				SymbolToken(";"),
				KeywordToken("return"),
				SymbolToken(";"),
			}},
			want: MockNodes([]TreeNode{
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						AdaptTokenToNode(KeywordToken("let")),
						MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
						AdaptTokenToNode(SymbolToken("=")),
						MockNodes([]TreeNode{
							MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(100))}, TermType, false),
						}, ExpressionType, false),
						AdaptTokenToNode(SymbolToken(";")),
					}, LetStatementType, false),
				}, StatementType, true),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						AdaptTokenToNode(KeywordToken("let")),
						MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("y"))}, VarNameType, true),
						AdaptTokenToNode(SymbolToken("=")),
						MockNodes([]TreeNode{
							MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(200))}, TermType, false),
						}, ExpressionType, false),
						AdaptTokenToNode(SymbolToken(";")),
					}, LetStatementType, false),
				}, StatementType, true),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						AdaptTokenToNode(KeywordToken("return")),
						AdaptTokenToNode(SymbolToken(";")),
					}, ReturnStatementType, false),
				}, StatementType, true),
			}, StatementsType, false),
			want1:   []Token{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseStatements(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseStatements() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseStatements() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseStatements() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}

func TestParser_parseVarDec(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    TreeNode
		want1   TokenList
		want2   *SymbolTable
		wantErr bool
	}{
		{
			name: "single var",
			args: args{[]Token{
				KeywordToken("var"),
				KeywordToken("boolean"),
				IdentifierToken("x"),
				SymbolToken(";"),
				KeywordToken("return"),
			}},
			want: MockNodes([]TreeNode{
				AdaptTokenToNode(KeywordToken("var")),
				MockNodes([]TreeNode{AdaptTokenToNode(KeywordToken("boolean"))}, TypeType, true),
				MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
				AdaptTokenToNode(SymbolToken(";")),
			}, VarDecType, false),
			want1: []Token{
				KeywordToken("return"),
			},
			want2: &SymbolTable{
				classScopeTable: ScopedTable{},
				funcScopeTable: ScopedTable{
					"x": &SymbolTableEntry{
						Name:  "x",
						Typ:   "boolean",
						Kind:  Var,
						Index: 0,
					},
				},
				index: map[VarKind]int{
					Var: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "multi vars",
			args: args{[]Token{
				KeywordToken("var"),
				KeywordToken("int"),
				IdentifierToken("x"),
				SymbolToken(","),
				IdentifierToken("y"),
				SymbolToken(";"),
				KeywordToken("return"),
			}},
			want: MockNodes([]TreeNode{
				AdaptTokenToNode(KeywordToken("var")),
				MockNodes([]TreeNode{AdaptTokenToNode(KeywordToken("int"))}, TypeType, true),
				MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
				AdaptTokenToNode(SymbolToken(",")),
				MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("y"))}, VarNameType, true),
				AdaptTokenToNode(SymbolToken(";")),
			}, VarDecType, false),
			want1: []Token{
				KeywordToken("return"),
			},
			want2: &SymbolTable{
				classScopeTable: ScopedTable{},
				funcScopeTable: ScopedTable{
					"x": &SymbolTableEntry{
						Name:  "x",
						Typ:   "int",
						Kind:  Var,
						Index: 0,
					},
					"y": &SymbolTableEntry{
						Name:  "y",
						Typ:   "int",
						Kind:  Var,
						Index: 1,
					},
				},
				index: map[VarKind]int{
					Var: 2,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseVarDec(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseVarDec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseVarDec() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseVarDec() diff (-got1 +want1)\n%s", diff)
			}
			opt = cmp.AllowUnexported(*p.symbolTable)
			if diff := cmp.Diff(p.symbolTable, tt.want2, opt); diff != "" {
				t.Errorf("Parser.parseVarDec() SymbolTable (-got +want2)\n%s", diff)
			}
		})
	}
}

func TestParser_parseType(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    TreeNode
		want1   TokenList
		wantErr bool
	}{
		{
			name: "keyword type",
			args: args{[]Token{
				KeywordToken("boolean"),
				IdentifierToken("x"),
			}},
			want: MockNodes([]TreeNode{
				AdaptTokenToNode(KeywordToken("boolean")),
			}, TypeType, true),
			want1: []Token{
				IdentifierToken("x"),
			},
			wantErr: false,
		},
		{
			name: "class name",
			args: args{[]Token{
				IdentifierToken("MyClass"),
				IdentifierToken("x"),
			}},
			want: MockNodes([]TreeNode{
				MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("MyClass"))}, ClassNameType, true),
			}, TypeType, true),
			want1: []Token{
				IdentifierToken("x"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseType(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseType() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseType() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}

func TestParser_parseSubroutineBody(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    TreeNode
		want1   TokenList
		wantErr bool
	}{
		{
			name: "no var dec",
			args: args{[]Token{
				SymbolToken("{"),
				KeywordToken("let"),
				IdentifierToken("x"),
				SymbolToken("="),
				IntConstToken(100),
				SymbolToken(";"),
				SymbolToken("}"),
				KeywordToken("class"),
			}},
			want: MockNodes([]TreeNode{
				AdaptTokenToNode(SymbolToken("{")),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{
							AdaptTokenToNode(KeywordToken("let")),
							MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
							AdaptTokenToNode(SymbolToken("=")),
							MockNodes([]TreeNode{
								MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(100))}, TermType, false),
							}, ExpressionType, false),
							AdaptTokenToNode(SymbolToken(";")),
						}, LetStatementType, false),
					}, StatementType, true),
				}, StatementsType, false),
				AdaptTokenToNode(SymbolToken("}")),
			}, SubroutineBodyType, false),
			want1: []Token{
				KeywordToken("class"),
			},
			wantErr: false,
		},
		{
			name: "with var dec",
			args: args{[]Token{
				SymbolToken("{"),
				KeywordToken("var"),
				KeywordToken("int"),
				IdentifierToken("x"),
				SymbolToken(";"),
				KeywordToken("var"),
				KeywordToken("char"),
				IdentifierToken("y"),
				SymbolToken(";"),
				KeywordToken("let"),
				IdentifierToken("x"),
				SymbolToken("="),
				IntConstToken(100),
				SymbolToken(";"),
				SymbolToken("}"),
				KeywordToken("class"),
			}},
			want: MockNodes([]TreeNode{
				AdaptTokenToNode(SymbolToken("{")),
				MockNodes([]TreeNode{
					AdaptTokenToNode(KeywordToken("var")),
					MockNodes([]TreeNode{AdaptTokenToNode(KeywordToken("int"))}, TypeType, true),
					MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
					AdaptTokenToNode(SymbolToken(";")),
				}, VarDecType, false),
				MockNodes([]TreeNode{
					AdaptTokenToNode(KeywordToken("var")),
					MockNodes([]TreeNode{AdaptTokenToNode(KeywordToken("char"))}, TypeType, true),
					MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("y"))}, VarNameType, true),
					AdaptTokenToNode(SymbolToken(";")),
				}, VarDecType, false),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{
							AdaptTokenToNode(KeywordToken("let")),
							MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
							AdaptTokenToNode(SymbolToken("=")),
							MockNodes([]TreeNode{
								MockNodes([]TreeNode{AdaptTokenToNode(IntConstToken(100))}, TermType, false),
							}, ExpressionType, false),
							AdaptTokenToNode(SymbolToken(";")),
						}, LetStatementType, false),
					}, StatementType, true),
				}, StatementsType, false),
				AdaptTokenToNode(SymbolToken("}")),
			}, SubroutineBodyType, false),
			want1: []Token{
				KeywordToken("class"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseSubroutineBody(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseSubroutineBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseSubroutineBody() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseSubroutineBody() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}

func TestParser_parseParameterList(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    TreeNode
		want1   TokenList
		want2   *SymbolTable
		wantErr bool
	}{
		{
			name: "no var",
			args: args{[]Token{
				SymbolToken(")"),
			}},
			want: MockNodes(nil, ParameterListType, false),
			want1: []Token{
				SymbolToken(")"),
			},
			want2: &SymbolTable{
				classScopeTable: ScopedTable{},
				funcScopeTable:  ScopedTable{},
				index:           map[VarKind]int{},
			},
			wantErr: false,
		},
		{
			name: "multi vars",
			args: args{[]Token{
				KeywordToken("int"),
				IdentifierToken("x"),
				SymbolToken(","),
				KeywordToken("int"),
				IdentifierToken("y"),
				SymbolToken(")"),
			}},
			want: MockNodes([]TreeNode{
				MockNodes([]TreeNode{AdaptTokenToNode(KeywordToken("int"))}, TypeType, true),
				MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
				AdaptTokenToNode(SymbolToken(",")),
				MockNodes([]TreeNode{AdaptTokenToNode(KeywordToken("int"))}, TypeType, true),
				MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("y"))}, VarNameType, true),
			}, ParameterListType, false),
			want1: []Token{
				SymbolToken(")"),
			},
			want2: &SymbolTable{
				classScopeTable: ScopedTable{},
				funcScopeTable: ScopedTable{
					"x": &SymbolTableEntry{
						Name:  "x",
						Typ:   "int",
						Kind:  Argument,
						Index: 0,
					},
					"y": &SymbolTableEntry{
						Name:  "y",
						Typ:   "int",
						Kind:  Argument,
						Index: 1,
					},
				},
				index: map[VarKind]int{
					Argument: 2,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseParameterList(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseParameterList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseParameterList() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseParameterList() diff (-got1 +want1)\n%s", diff)
			}
			opt = cmp.AllowUnexported(*p.symbolTable)
			if diff := cmp.Diff(p.symbolTable, tt.want2, opt); diff != "" {
				t.Errorf("Parser.parseParameterList() SymbolTable (-got +want2)\n%s", diff)
			}
		})
	}
}

func TestParser_parseSubroutineDec(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    TreeNode
		want1   TokenList
		wantErr bool
	}{
		{
			name: "normal",
			args: args{[]Token{
				KeywordToken("constructor"),
				KeywordToken("int"),
				IdentifierToken("myFunc"),
				SymbolToken("("),
				KeywordToken("int"),
				IdentifierToken("x"),
				SymbolToken(")"),
				SymbolToken("{"),
				KeywordToken("return"),
				SymbolToken(";"),
				SymbolToken("}"),
				KeywordToken("class"),
			}},
			want: MockNodes([]TreeNode{
				AdaptTokenToNode(KeywordToken("constructor")),
				MockNodes([]TreeNode{AdaptTokenToNode(KeywordToken("int"))}, TypeType, true),
				MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("myFunc"))}, SubroutineNameType, true),
				AdaptTokenToNode(SymbolToken("(")),
				MockNodes([]TreeNode{
					MockNodes([]TreeNode{AdaptTokenToNode(KeywordToken("int"))}, TypeType, true),
					MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
				}, ParameterListType, false),
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
			want1: []Token{
				KeywordToken("class"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseSubroutineDec(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseSubroutineDec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseSubroutineDec() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseSubroutineDec() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}

func TestParser_parseClassVarDec(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    TreeNode
		want1   TokenList
		want2   *SymbolTable
		wantErr bool
	}{
		{
			name: "single var",
			args: args{[]Token{
				KeywordToken("static"),
				KeywordToken("boolean"),
				IdentifierToken("x"),
				SymbolToken(";"),
				KeywordToken("return"),
			}},
			want: MockNodes([]TreeNode{
				AdaptTokenToNode(KeywordToken("static")),
				MockNodes([]TreeNode{AdaptTokenToNode(KeywordToken("boolean"))}, TypeType, true),
				MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
				AdaptTokenToNode(SymbolToken(";")),
			}, ClassVarDecType, false),
			want1: []Token{
				KeywordToken("return"),
			},
			want2: &SymbolTable{
				classScopeTable: ScopedTable{
					"x": &SymbolTableEntry{
						Name:  "x",
						Typ:   "boolean",
						Kind:  Static,
						Index: 0,
					},
				},
				funcScopeTable: ScopedTable{},
				index: map[VarKind]int{
					Static: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "multi vars",
			args: args{[]Token{
				KeywordToken("field"),
				KeywordToken("int"),
				IdentifierToken("x"),
				SymbolToken(","),
				IdentifierToken("y"),
				SymbolToken(";"),
				KeywordToken("return"),
			}},
			want: MockNodes([]TreeNode{
				AdaptTokenToNode(KeywordToken("field")),
				MockNodes([]TreeNode{AdaptTokenToNode(KeywordToken("int"))}, TypeType, true),
				MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
				AdaptTokenToNode(SymbolToken(",")),
				MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("y"))}, VarNameType, true),
				AdaptTokenToNode(SymbolToken(";")),
			}, ClassVarDecType, false),
			want1: []Token{
				KeywordToken("return"),
			},
			want2: &SymbolTable{
				classScopeTable: ScopedTable{
					"x": &SymbolTableEntry{
						Name:  "x",
						Typ:   "int",
						Kind:  Field,
						Index: 0,
					},
					"y": &SymbolTableEntry{
						Name:  "y",
						Typ:   "int",
						Kind:  Field,
						Index: 1,
					},
				},
				funcScopeTable: ScopedTable{},
				index: map[VarKind]int{
					Field: 2,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseClassVarDec(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseClassVarDec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseClassVarDec() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseClassVarDec() diff (-got1 +want1)\n%s", diff)
			}
			opt = cmp.AllowUnexported(*p.symbolTable)
			if diff := cmp.Diff(p.symbolTable, tt.want2, opt); diff != "" {
				t.Errorf("Parser.parseClassVarDec() SymbolTable (-got +want2)\n%s", diff)
			}
		})
	}
}

func TestParser_parseClass(t *testing.T) {
	type args struct {
		tokens TokenList
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    TreeNode
		want1   TokenList
		wantErr bool
	}{
		{
			name: "normal",
			args: args{[]Token{
				KeywordToken("class"),
				IdentifierToken("MyClass"),
				SymbolToken("{"),
				KeywordToken("static"),
				KeywordToken("boolean"),
				IdentifierToken("x"),
				SymbolToken(";"),
				KeywordToken("constructor"),
				KeywordToken("int"),
				IdentifierToken("myFunc"),
				SymbolToken("("),
				KeywordToken("int"),
				IdentifierToken("x"),
				SymbolToken(")"),
				SymbolToken("{"),
				KeywordToken("return"),
				SymbolToken(";"),
				SymbolToken("}"),
				SymbolToken("}"),
				KeywordToken("return"),
			}},
			want: MockNodes([]TreeNode{
				AdaptTokenToNode(KeywordToken("class")),
				MockNodes([]TreeNode{
					AdaptTokenToNode(IdentifierToken("MyClass")),
				}, ClassNameType, true),
				AdaptTokenToNode(SymbolToken("{")),
				MockNodes([]TreeNode{
					AdaptTokenToNode(KeywordToken("static")),
					MockNodes([]TreeNode{AdaptTokenToNode(KeywordToken("boolean"))}, TypeType, true),
					MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
					AdaptTokenToNode(SymbolToken(";")),
				}, ClassVarDecType, false),
				MockNodes([]TreeNode{
					AdaptTokenToNode(KeywordToken("constructor")),
					MockNodes([]TreeNode{AdaptTokenToNode(KeywordToken("int"))}, TypeType, true),
					MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("myFunc"))}, SubroutineNameType, true),
					AdaptTokenToNode(SymbolToken("(")),
					MockNodes([]TreeNode{
						MockNodes([]TreeNode{AdaptTokenToNode(KeywordToken("int"))}, TypeType, true),
						MockNodes([]TreeNode{AdaptTokenToNode(IdentifierToken("x"))}, VarNameType, true),
					}, ParameterListType, false),
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
				AdaptTokenToNode(SymbolToken("}")),
			}, ClassType, false),
			want1: []Token{
				KeywordToken("return"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, got1, err := p.parseClass(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseClass() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			opt := cmpopts.IgnoreFields(LeafNode{}, "IDMeta")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("Parser.parseClass() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseClass() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}
