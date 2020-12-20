package main

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTokenizer_parseLine(t *testing.T) {
	type fields struct {
		state tokenizeState
	}
	type args struct {
		l string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Token
		wantErr bool
	}{
		{
			name:   "class declaration",
			fields: fields{state: ordinal},
			args:   args{l: `class Main {`},
			want: []Token{
				KeywordToken("class"),
				IdentifierToken("Main"),
				SymbolToken("{"),
			},
			wantErr: false,
		},
		{
			name:   "function declaration",
			fields: fields{state: ordinal},
			args:   args{l: `function void main() {`},
			want: []Token{
				KeywordToken("function"),
				KeywordToken("void"),
				IdentifierToken("main"),
				SymbolToken("("),
				SymbolToken(")"),
				SymbolToken("{"),
			},
			wantErr: false,
		},
		{
			name:   "statements",
			fields: fields{state: ordinal},
			args:   args{l: `let a = Array.new(length);`},
			want: []Token{
				KeywordToken("let"),
				IdentifierToken("a"),
				SymbolToken("="),
				IdentifierToken("Array"),
				SymbolToken("."),
				IdentifierToken("new"),
				SymbolToken("("),
				IdentifierToken("length"),
				SymbolToken(")"),
				SymbolToken(";"),
			},
			wantErr: false,
		},
		{
			name:   "should ignore spaces",
			fields: fields{state: ordinal},
			args:   args{l: `   let     a=123   ; `},
			want: []Token{
				KeywordToken("let"),
				IdentifierToken("a"),
				SymbolToken("="),
				IntConstToken(123),
				SymbolToken(";"),
			},
			wantErr: false,
		},
		{
			name:   "divisor",
			fields: fields{state: ordinal},
			args:   args{l: `let div = 10/5;`},
			want: []Token{
				KeywordToken("let"),
				IdentifierToken("div"),
				SymbolToken("="),
				IntConstToken(10),
				SymbolToken("/"),
				IntConstToken(5),
				SymbolToken(";"),
			},
			wantErr: false,
		},
		{
			name:   "single line comment",
			fields: fields{state: ordinal},
			args:   args{l: `let a = 123; // let b = 100;`},
			want: []Token{
				KeywordToken("let"),
				IdentifierToken("a"),
				SymbolToken("="),
				IntConstToken(123),
				SymbolToken(";"),
			},
			wantErr: false,
		},
		{
			name:    "single line comment only",
			fields:  fields{state: ordinal},
			args:    args{l: `   // let b = 100;`},
			want:    nil,
			wantErr: false,
		},
		{
			name:   "single line comment just after identifier",
			fields: fields{state: ordinal},
			args:   args{l: `let b = sum// new line`},
			want: []Token{
				KeywordToken("let"),
				IdentifierToken("b"),
				SymbolToken("="),
				IdentifierToken("sum"),
			},
			wantErr: false,
		},
		{
			name:   "multi line comment opened",
			fields: fields{state: ordinal},
			args:   args{l: `let a = 123; /* let b = 100;`},
			want: []Token{
				KeywordToken("let"),
				IdentifierToken("a"),
				SymbolToken("="),
				IntConstToken(123),
				SymbolToken(";"),
			},
			wantErr: false,
		},
		{
			name:   "multi line comment closed",
			fields: fields{state: multiCommentOpened},
			args:   args{l: `let a = 123; */ let b = 100;`},
			want: []Token{
				KeywordToken("let"),
				IdentifierToken("b"),
				SymbolToken("="),
				IntConstToken(100),
				SymbolToken(";"),
			},
			wantErr: false,
		},
		{
			name:    "multi line comment only",
			fields:  fields{state: ordinal},
			args:    args{l: `/** some document`},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "opening multi line comment only",
			fields:  fields{state: multiCommentOpened},
			args:    args{l: `some document */`},
			want:    nil,
			wantErr: false,
		},
		{
			name:   "string constant",
			fields: fields{state: ordinal},
			args:   args{l: `let length = Keyboard.readInt("HOW MANY NUMBERS? ");`},
			want: []Token{
				KeywordToken("let"),
				IdentifierToken("length"),
				SymbolToken("="),
				IdentifierToken("Keyboard"),
				SymbolToken("."),
				IdentifierToken("readInt"),
				SymbolToken("("),
				StrConstToken("HOW MANY NUMBERS? "),
				SymbolToken(")"),
				SymbolToken(";"),
			},
			wantErr: false,
		},
		{
			name:   "digit only string",
			fields: fields{state: ordinal},
			args:   args{l: `let s = "-123";`},
			want: []Token{
				KeywordToken("let"),
				IdentifierToken("s"),
				SymbolToken("="),
				StrConstToken("-123"),
				SymbolToken(";"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := &Tokenizer{
				state: tt.fields.state,
			}
			got, err := tokenizer.parseLine(tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tokenizer.parseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Tokenizer.parseLine() diff (-got +want)\n%s", diff)
			}
		})
	}
}

func TestTokenizer_delim(t *testing.T) {
	type fields struct {
		state tokenizeState
	}
	type args struct {
		c rune
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "not delim",
			fields: fields{state: ordinal},
			args:   args{c: 'a'},
			want:   false,
		},
		{
			name:   "space is delim",
			fields: fields{state: ordinal},
			args:   args{c: ' '},
			want:   true,
		},
		{
			name:   "symbol is delim",
			fields: fields{state: ordinal},
			args:   args{c: '.'},
			want:   true,
		},
		{
			name:   "not delim if in comment",
			fields: fields{state: multiCommentOpened},
			args:   args{c: '.'},
			want:   false,
		},
		{
			name:   "not delim if in string",
			fields: fields{state: stringOpened},
			args:   args{c: ' '},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := &Tokenizer{
				state: tt.fields.state,
			}
			if got := tokenizer.delim(tt.args.c); got != tt.want {
				t.Errorf("Tokenizer.delim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenizer_bufToToken(t *testing.T) {
	type fields struct {
		state tokenizeState
		buf   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    Token
		wantErr bool
	}{
		{
			name:    "keyword",
			fields:  fields{ordinal, "method"},
			want:    KeywordToken("method"),
			wantErr: false,
		},
		{
			name:    "integer",
			fields:  fields{ordinal, "982"},
			want:    IntConstToken(982),
			wantErr: false,
		},
		{
			name:    "identifier",
			fields:  fields{ordinal, "Main"},
			want:    IdentifierToken("Main"),
			wantErr: false,
		},
		{
			name:    "invalid identifier",
			fields:  fields{ordinal, "8var"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := &Tokenizer{
				state: tt.fields.state,
				buf:   tt.fields.buf,
			}
			got, err := tokenizer.bufToToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("Tokenizer.bufToToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokenizer.bufToToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
