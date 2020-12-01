package main

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func MockNodes(tokens []Token, typ NodeType, terminal bool) TreeNode {
	xml := ""
	if !terminal {
		s := typ.String()
		s = strings.Replace(s, "Type", "", -1)
		r := []rune(s)
		head := strings.ToLower(string(r[0]))
		rest := string(r[1:len(r)])
		xml = head + rest
	}
	return &GeneralNode{
		Children:  tokens,
		Typ:       typ,
		XMLHeader: xml,
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
			want: MockNodes([]Token{KeywordToken("true")}, KeyConstType, true),
			want1: []Token{
				SymbolToken(";"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, got1, err := p.parseKeywordConstant(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseKeywordConstant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
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
			want: MockNodes([]Token{SymbolToken("-")}, UnaryOpType, true),
			want1: []Token{
				IdentifierToken("i"),
				SymbolToken(";"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, got1, err := p.parseUnaryOp(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseUnaryOp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
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
			want: MockNodes([]Token{SymbolToken("+")}, OpType, true),
			want1: []Token{
				IdentifierToken("123"),
				SymbolToken(";"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, got1, err := p.parseOp(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseOp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
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
			want: MockNodes([]Token{IntConstToken(123)}, TermType, false),
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
			want: MockNodes([]Token{StrConstToken("string")}, TermType, false),
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
			want: MockNodes([]Token{MockNodes([]Token{KeywordToken("true")}, KeyConstType, true)}, TermType, false),
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
			want: MockNodes([]Token{MockNodes([]Token{IdentifierToken("args")}, VarNameType, true)}, TermType, false),
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
				[]Token{
					MockNodes([]Token{IdentifierToken("args")}, VarNameType, true),
					SymbolToken("["),
					MockNodes([]Token{
						MockNodes([]Token{
							MockNodes([]Token{IdentifierToken("i")}, VarNameType, true),
						}, TermType, false),
					}, ExpressionType, false),
					SymbolToken("]"),
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
			want: MockNodes([]Token{
				MockNodes([]Token{
					MockNodes([]Token{IdentifierToken("fooFunc")}, SubroutineNameType, true),
					SymbolToken("("),
					MockNodes([]Token{
						MockNodes([]Token{
							MockNodes([]Token{
								MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
							}, TermType, false),
						}, ExpressionType, false),
					}, ExpressionListType, false),
					SymbolToken(")"),
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
			want: MockNodes([]Token{
				SymbolToken("("),
				MockNodes([]Token{
					MockNodes([]Token{
						MockNodes([]Token{IdentifierToken("ok")}, VarNameType, true),
					}, TermType, false),
				}, ExpressionType, false),
				SymbolToken(")"),
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
			want: MockNodes([]Token{
				MockNodes([]Token{SymbolToken("~")}, UnaryOpType, true),
				MockNodes([]Token{
					MockNodes([]Token{IdentifierToken("some")}, VarNameType, true),
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
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, got1, err := p.parseTerm(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseTerm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
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
			want: MockNodes([]Token{IdentifierToken("args")}, VarNameType, true),
			want1: []Token{
				SymbolToken(";"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, got1, err := p.parseVarName(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseVarName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
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
			want: MockNodes([]Token{
				MockNodes([]Token{
					MockNodes([]Token{
						IdentifierToken("x"),
					}, VarNameType, true),
				}, TermType, false),
			}, ExpressionType, false),
			want1: []Token{
				SymbolToken("]"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, got1, err := p.parseExpression(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
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
			want: MockNodes([]Token(nil), ExpressionListType, false),
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
			want: MockNodes([]Token{
				MockNodes([]Token{
					MockNodes([]Token{
						MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
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
			want: MockNodes([]Token{
				MockNodes([]Token{
					MockNodes([]Token{
						MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
					}, TermType, false),
				}, ExpressionType, false),
				SymbolToken(","),
				MockNodes([]Token{
					MockNodes([]Token{
						MockNodes([]Token{IdentifierToken("y")}, VarNameType, true),
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
			p := &Parser{}
			got, got1, err := p.parseExpressionList(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseExpressionList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
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
			want: MockNodes([]Token{
				IdentifierToken("myFunc"),
			}, SubroutineNameType, true),
			want1: []Token{
				SymbolToken("("),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, got1, err := p.parseSubroutineName(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseSubroutineName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
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
			want: MockNodes([]Token{
				IdentifierToken("MyClass"),
			}, ClassNameType, true),
			want1: []Token{
				SymbolToken("."),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, got1, err := p.parseClassName(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseClassName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Parser.parseClassName() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseClassName() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}

func TestParser_parseSubroutineCall(t *testing.T) {
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
			name: "function call",
			args: args{[]Token{
				IdentifierToken("MyFunc"),
				SymbolToken("("),
				IdentifierToken("x"),
				SymbolToken(","),
				IdentifierToken("y"),
				SymbolToken(")"),
				SymbolToken(";"),
			}},
			want: MockNodes([]Token{
				MockNodes([]Token{IdentifierToken("MyFunc")}, SubroutineNameType, true),
				SymbolToken("("),
				MockNodes([]Token{
					MockNodes([]Token{
						MockNodes([]Token{
							MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
						}, TermType, false),
					}, ExpressionType, false),
					SymbolToken(","),
					MockNodes([]Token{
						MockNodes([]Token{
							MockNodes([]Token{IdentifierToken("y")}, VarNameType, true),
						}, TermType, false),
					}, ExpressionType, false),
				}, ExpressionListType, false),
				SymbolToken(")"),
			}, SubroutineCallType, true),
			want1: []Token{
				SymbolToken(";"),
			},
			wantErr: false,
		},
		{
			name: "function call with no argument",
			args: args{[]Token{
				IdentifierToken("MyFunc"),
				SymbolToken("("),
				SymbolToken(")"),
				SymbolToken(";"),
			}},
			want: MockNodes([]Token{
				MockNodes([]Token{IdentifierToken("MyFunc")}, SubroutineNameType, true),
				SymbolToken("("),
				MockNodes(nil, ExpressionListType, false),
				SymbolToken(")"),
			}, SubroutineCallType, true),
			want1: []Token{
				SymbolToken(";"),
			},
			wantErr: false,
		},
		{
			name: "method call",
			args: args{[]Token{
				IdentifierToken("MyClass"),
				SymbolToken("."),
				IdentifierToken("MyMethod"),
				SymbolToken("("),
				SymbolToken(")"),
				SymbolToken(";"),
			}},
			want: MockNodes([]Token{
				MockNodes([]Token{IdentifierToken("MyClass")}, ClassNameType, true),
				SymbolToken("."),
				MockNodes([]Token{IdentifierToken("MyMethod")}, SubroutineNameType, true),
				SymbolToken("("),
				MockNodes(nil, ExpressionListType, false),
				SymbolToken(")"),
			}, SubroutineCallType, true),
			want1: []Token{
				SymbolToken(";"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, got1, err := p.parseSubroutineCall(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseSubroutineCall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
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
			want: MockNodes([]Token{
				KeywordToken("return"),
				SymbolToken(";"),
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
			want: MockNodes([]Token{
				KeywordToken("return"),
				MockNodes([]Token{
					MockNodes([]Token{
						MockNodes([]Token{
							IdentifierToken("res"),
						}, VarNameType, true),
					}, TermType, false),
				}, ExpressionType, false),
				SymbolToken(";"),
			}, ReturnStatementType, false),
			want1: []Token{
				SymbolToken("}"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, got1, err := p.parseReturnStatement(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseReturnStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
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
			want: MockNodes([]Token{
				KeywordToken("do"),
				MockNodes([]Token{
					MockNodes([]Token{IdentifierToken("MyFunc")}, SubroutineNameType, true),
					SymbolToken("("),
					MockNodes(nil, ExpressionListType, false),
					SymbolToken(")"),
				}, SubroutineCallType, true),
				SymbolToken(";"),
			}, DoStatementType, false),
			want1: []Token{
				SymbolToken("}"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, got1, err := p.parseDoStatement(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseDoStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
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
				SymbolToken("}"),
				KeywordToken("return"),
			}},
			want: MockNodes([]Token{
				KeywordToken("while"),
				SymbolToken("("),
				MockNodes([]Token{
					MockNodes([]Token{
						MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
					}, TermType, false),
				}, ExpressionType, false),
				SymbolToken(")"),
				SymbolToken("{"),
				MockNodes([]Token{
					MockNodes([]Token{
						KeywordToken("do"),
						MockNodes([]Token{
							MockNodes([]Token{IdentifierToken("MyFunc")}, SubroutineNameType, true),
							SymbolToken("("),
							MockNodes(nil, ExpressionListType, false),
							SymbolToken(")"),
						}, SubroutineCallType, true),
						SymbolToken(";"),
					}, DoStatementType, false),
				}, StatementType, true),
				SymbolToken("}"),
			}, WhileStatementType, false),
			want1: []Token{
				KeywordToken("return"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, got1, err := p.parseWhileStatement(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseWhileStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
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
				KeywordToken("do"),
				IdentifierToken("MyFunc"),
				SymbolToken("("),
				SymbolToken(")"),
				SymbolToken(";"),
				SymbolToken("}"),
				KeywordToken("return"),
			}},
			want: MockNodes([]Token{
				KeywordToken("if"),
				SymbolToken("("),
				MockNodes([]Token{
					MockNodes([]Token{
						MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
					}, TermType, false),
				}, ExpressionType, false),
				SymbolToken(")"),
				SymbolToken("{"),
				MockNodes([]Token{
					MockNodes([]Token{
						KeywordToken("do"),
						MockNodes([]Token{
							MockNodes([]Token{IdentifierToken("MyFunc")}, SubroutineNameType, true),
							SymbolToken("("),
							MockNodes(nil, ExpressionListType, false),
							SymbolToken(")"),
						}, SubroutineCallType, true),
						SymbolToken(";"),
					}, DoStatementType, false),
				}, StatementType, true),
				SymbolToken("}"),
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
				KeywordToken("do"),
				IdentifierToken("MyFuncElse"),
				SymbolToken("("),
				SymbolToken(")"),
				SymbolToken(";"),
				SymbolToken("}"),
				KeywordToken("return"),
			}},
			want: MockNodes([]Token{
				KeywordToken("if"),
				SymbolToken("("),
				MockNodes([]Token{
					MockNodes([]Token{
						MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
					}, TermType, false),
				}, ExpressionType, false),
				SymbolToken(")"),
				SymbolToken("{"),
				MockNodes([]Token{
					MockNodes([]Token{
						KeywordToken("do"),
						MockNodes([]Token{
							MockNodes([]Token{IdentifierToken("MyFunc")}, SubroutineNameType, true),
							SymbolToken("("),
							MockNodes(nil, ExpressionListType, false),
							SymbolToken(")"),
						}, SubroutineCallType, true),
						SymbolToken(";"),
					}, DoStatementType, false),
				}, StatementType, true),
				SymbolToken("}"),
				KeywordToken("else"),
				SymbolToken("{"),
				MockNodes([]Token{
					MockNodes([]Token{
						KeywordToken("do"),
						MockNodes([]Token{
							MockNodes([]Token{IdentifierToken("MyFuncElse")}, SubroutineNameType, true),
							SymbolToken("("),
							MockNodes(nil, ExpressionListType, false),
							SymbolToken(")"),
						}, SubroutineCallType, true),
						SymbolToken(";"),
					}, DoStatementType, false),
				}, StatementType, true),
				SymbolToken("}"),
			}, IfStatementType, false),
			want1: []Token{
				KeywordToken("return"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, got1, err := p.parseIfStatement(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseIfStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
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
			want: MockNodes([]Token{
				KeywordToken("let"),
				MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
				SymbolToken("="),
				MockNodes([]Token{
					MockNodes([]Token{IntConstToken(100)}, TermType, false),
				}, ExpressionType, false),
				SymbolToken(";"),
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
			want: MockNodes([]Token{
				KeywordToken("let"),
				MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
				SymbolToken("["),
				MockNodes([]Token{
					MockNodes([]Token{IntConstToken(3)}, TermType, false),
				}, ExpressionType, false),
				SymbolToken("]"),
				SymbolToken("="),
				MockNodes([]Token{
					MockNodes([]Token{IntConstToken(100)}, TermType, false),
				}, ExpressionType, false),
				SymbolToken(";"),
			}, LetStatementType, false),
			want1: []Token{
				KeywordToken("return"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, got1, err := p.parseLetStatement(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseLetStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
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
			want: MockNodes([]Token{
				MockNodes([]Token{
					MockNodes([]Token{
						KeywordToken("let"),
						MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
						SymbolToken("="),
						MockNodes([]Token{
							MockNodes([]Token{IntConstToken(100)}, TermType, false),
						}, ExpressionType, false),
						SymbolToken(";"),
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
			}},
			want: MockNodes([]Token{
				MockNodes([]Token{
					MockNodes([]Token{
						KeywordToken("let"),
						MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
						SymbolToken("="),
						MockNodes([]Token{
							MockNodes([]Token{IntConstToken(100)}, TermType, false),
						}, ExpressionType, false),
						SymbolToken(";"),
					}, LetStatementType, false),
				}, StatementType, true),
				MockNodes([]Token{
					MockNodes([]Token{
						KeywordToken("let"),
						MockNodes([]Token{IdentifierToken("y")}, VarNameType, true),
						SymbolToken("="),
						MockNodes([]Token{
							MockNodes([]Token{IntConstToken(200)}, TermType, false),
						}, ExpressionType, false),
						SymbolToken(";"),
					}, LetStatementType, false),
				}, StatementType, true),
			}, StatementsType, false),
			want1: []Token{
				KeywordToken("return"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, got1, err := p.parseStatements(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseStatements() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
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
			want: MockNodes([]Token{
				KeywordToken("var"),
				MockNodes([]Token{KeywordToken("boolean")}, TypeType, true),
				MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
				SymbolToken(";"),
			}, VarDecType, false),
			want1: []Token{
				KeywordToken("return"),
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
			want: MockNodes([]Token{
				KeywordToken("var"),
				MockNodes([]Token{KeywordToken("int")}, TypeType, true),
				MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
				SymbolToken(","),
				MockNodes([]Token{IdentifierToken("y")}, VarNameType, true),
				SymbolToken(";"),
			}, VarDecType, false),
			want1: []Token{
				KeywordToken("return"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, got1, err := p.parseVarDec(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseVarDec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Parser.parseVarDec() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseVarDec() diff (-got1 +want1)\n%s", diff)
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
			want: MockNodes([]Token{
				KeywordToken("boolean"),
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
			want: MockNodes([]Token{
				MockNodes([]Token{IdentifierToken("MyClass")}, ClassNameType, true),
			}, TypeType, true),
			want1: []Token{
				IdentifierToken("x"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, got1, err := p.parseType(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
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
			want: MockNodes([]Token{
				SymbolToken("{"),
				MockNodes([]Token{
					MockNodes([]Token{
						MockNodes([]Token{
							KeywordToken("let"),
							MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
							SymbolToken("="),
							MockNodes([]Token{
								MockNodes([]Token{IntConstToken(100)}, TermType, false),
							}, ExpressionType, false),
							SymbolToken(";"),
						}, LetStatementType, false),
					}, StatementType, true),
				}, StatementsType, false),
				SymbolToken("}"),
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
			want: MockNodes([]Token{
				SymbolToken("{"),
				MockNodes([]Token{
					KeywordToken("var"),
					MockNodes([]Token{KeywordToken("int")}, TypeType, true),
					MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
					SymbolToken(";"),
				}, VarDecType, false),
				MockNodes([]Token{
					KeywordToken("var"),
					MockNodes([]Token{KeywordToken("char")}, TypeType, true),
					MockNodes([]Token{IdentifierToken("y")}, VarNameType, true),
					SymbolToken(";"),
				}, VarDecType, false),
				MockNodes([]Token{
					MockNodes([]Token{
						MockNodes([]Token{
							KeywordToken("let"),
							MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
							SymbolToken("="),
							MockNodes([]Token{
								MockNodes([]Token{IntConstToken(100)}, TermType, false),
							}, ExpressionType, false),
							SymbolToken(";"),
						}, LetStatementType, false),
					}, StatementType, true),
				}, StatementsType, false),
				SymbolToken("}"),
			}, SubroutineBodyType, false),
			want1: []Token{
				KeywordToken("class"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, got1, err := p.parseSubroutineBody(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseSubroutineBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
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
			want: MockNodes([]Token{
				MockNodes([]Token{KeywordToken("int")}, TypeType, true),
				MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
				SymbolToken(","),
				MockNodes([]Token{KeywordToken("int")}, TypeType, true),
				MockNodes([]Token{IdentifierToken("y")}, VarNameType, true),
			}, ParameterListType, false),
			want1: []Token{
				SymbolToken(")"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, got1, err := p.parseParameterList(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseParameterList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Parser.parseParameterList() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseParameterList() diff (-got1 +want1)\n%s", diff)
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
			want: MockNodes([]Token{
				KeywordToken("constructor"),
				MockNodes([]Token{KeywordToken("int")}, TypeType, true),
				MockNodes([]Token{IdentifierToken("myFunc")}, SubroutineNameType, true),
				SymbolToken("("),
				MockNodes([]Token{
					MockNodes([]Token{KeywordToken("int")}, TypeType, true),
					MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
				}, ParameterListType, false),
				SymbolToken(")"),
				MockNodes([]Token{
					SymbolToken("{"),
					MockNodes([]Token{
						MockNodes([]Token{
							MockNodes([]Token{
								KeywordToken("return"),
								SymbolToken(";"),
							}, ReturnStatementType, false),
						}, StatementType, true),
					}, StatementsType, false),
					SymbolToken("}"),
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
			p := &Parser{}
			got, got1, err := p.parseSubroutineDec(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseSubroutineDec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
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
			want: MockNodes([]Token{
				KeywordToken("static"),
				MockNodes([]Token{KeywordToken("boolean")}, TypeType, true),
				MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
				SymbolToken(";"),
			}, ClassVarDecType, false),
			want1: []Token{
				KeywordToken("return"),
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
			want: MockNodes([]Token{
				KeywordToken("field"),
				MockNodes([]Token{KeywordToken("int")}, TypeType, true),
				MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
				SymbolToken(","),
				MockNodes([]Token{IdentifierToken("y")}, VarNameType, true),
				SymbolToken(";"),
			}, ClassVarDecType, false),
			want1: []Token{
				KeywordToken("return"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, got1, err := p.parseClassVarDec(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseClassVarDec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Parser.parseClassVarDec() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseClassVarDec() diff (-got1 +want1)\n%s", diff)
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
			want: MockNodes([]Token{
				KeywordToken("class"),
				MockNodes([]Token{
					IdentifierToken("MyClass"),
				}, ClassNameType, true),
				SymbolToken("{"),
				MockNodes([]Token{
					KeywordToken("static"),
					MockNodes([]Token{KeywordToken("boolean")}, TypeType, true),
					MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
					SymbolToken(";"),
				}, ClassVarDecType, false),
				MockNodes([]Token{
					KeywordToken("constructor"),
					MockNodes([]Token{KeywordToken("int")}, TypeType, true),
					MockNodes([]Token{IdentifierToken("myFunc")}, SubroutineNameType, true),
					SymbolToken("("),
					MockNodes([]Token{
						MockNodes([]Token{KeywordToken("int")}, TypeType, true),
						MockNodes([]Token{IdentifierToken("x")}, VarNameType, true),
					}, ParameterListType, false),
					SymbolToken(")"),
					MockNodes([]Token{
						SymbolToken("{"),
						MockNodes([]Token{
							MockNodes([]Token{
								MockNodes([]Token{
									KeywordToken("return"),
									SymbolToken(";"),
								}, ReturnStatementType, false),
							}, StatementType, true),
						}, StatementsType, false),
						SymbolToken("}"),
					}, SubroutineBodyType, false),
				}, SubroutineDecType, false),
				SymbolToken("}"),
			}, ClassType, false),
			want1: []Token{
				KeywordToken("return"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{}
			got, got1, err := p.parseClass(tt.args.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseClass() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Parser.parseClass() diff (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("Parser.parseClass() diff (-got1 +want1)\n%s", diff)
			}
		})
	}
}
