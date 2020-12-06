package main

import "fmt"

type Parser struct{}

func NewParser() *Parser {
	return nil
}

func (p *Parser) Parse(tokens []Token) (TreeNode, error) {
	res, rest, err := p.parseClass(TokenList(tokens))
	if err != nil {
		return nil, err
	}
	if len(rest) > 0 {
		return nil, fmt.Errorf("Invalid tokens remaining %v", rest)
	}
	return res, nil
}

func (p *Parser) parseClass(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewClassNode()

	// class keyword
	mayClassKeyword, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseClass] %w", err)
	}
	if mayClassKeyword.Type() != KeywordType || mayClassKeyword.String() != "class" {
		return nil, nil, fmt.Errorf("[parseClass] Invalid keyword %v want class, %v", mayClassKeyword, tokens)
	}
	res.AppendChild(mayClassKeyword)

	// class name
	cn, rest, err := p.parseClassName(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseClass] %w", err)
	}
	res.AppendChild(cn)

	// open bracket
	mayOpenBracket, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseClass] %w", err)
	}
	if mayOpenBracket.Type() != SymbolType || mayOpenBracket.String() != "{" {
		return nil, nil, fmt.Errorf("[parseClass] Invalid symbol %v want {, %v", mayOpenBracket, tokens)
	}
	res.AppendChild(mayOpenBracket)

	// class var declaration
	for true {
		d, r, err := p.parseClassVarDec(rest)
		if err != nil {
			break
		}
		res.AppendChild(d)
		rest = r
	}

	// subroutine declaration
	for true {
		d, r, err := p.parseSubroutineDec(rest)
		if err != nil {
			break
		}
		res.AppendChild(d)
		rest = r
	}

	// close bracket
	mayCloseBracket, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseClass] %w", err)
	}
	if mayCloseBracket.Type() != SymbolType || mayCloseBracket.String() != "}" {
		return nil, nil, fmt.Errorf("[parseClass] Invalid symbol %v want }, %v", mayCloseBracket, rest)
	}
	res.AppendChild(mayCloseBracket)

	return res, rest, nil
}

func (p *Parser) parseClassVarDec(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewClassVarDecNode()

	// class var type
	mayClassVarType, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseClassVarDec] %w", err)
	}
	if mayClassVarType.Type() != KeywordType || (mayClassVarType.String() != "static" && mayClassVarType.String() != "field") {
		return nil, nil, fmt.Errorf("[parseClassVarDec] Invalid keyword %v want (static|field), %v", mayClassVarType, rest)
	}
	res.AppendChild(mayClassVarType)

	// type
	typ, rest, err := p.parseType(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseClassVarDec] %w", err)
	}
	res.AppendChild(typ)

	// var name
	v, rest, err := p.parseVarName(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseClassVarDec] %w", err)
	}
	res.AppendChild(v)

	// following vars
	for true {
		// comma
		mayComma, r, err := rest.PopNext()
		if err != nil {
			break
		}
		if mayComma.Type() != SymbolType || mayComma.String() != "," {
			break
		}

		// var
		v, r, err := p.parseVarName(r)
		if err != nil {
			return nil, nil, fmt.Errorf("[parseClassVarDec] %w", err)
		}

		res.AppendChild(mayComma)
		res.AppendChild(v)
		rest = r
	}

	// semicolon
	maySemicolon, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseClassVarDec] %w", err)
	}
	if maySemicolon.Type() != SymbolType || maySemicolon.String() != ";" {
		return nil, nil, fmt.Errorf("[parseClassVarDec] Invalid keyword %v want ;, %v", maySemicolon, rest)
	}
	res.AppendChild(maySemicolon)

	return res, rest, nil
}

func (p *Parser) parseType(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewTypeNode()

	mayKw, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseType] %w", err)
	}
	if mayKw.Type() == KeywordType {
		switch mayKw.String() {
		case "int", "char", "boolean":
			res.AppendChild(mayKw)
			return res, rest, nil
		}
		return nil, nil, fmt.Errorf("[parseType] Invalid keyword %v want (int|char|boolean), %v", mayKw, rest)
	}

	n, rest, err := p.parseClassName(tokens)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseType] %w", err)
	}
	res.AppendChild(n)
	return res, rest, nil
}

func (p *Parser) parseSubroutineDec(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewSubroutineDecNode()

	// func type
	mayFuncType, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineDec] %w", err)
	}
	if mayFuncType.Type() != KeywordType || (mayFuncType.String() != "constructor" && mayFuncType.String() != "function" && mayFuncType.String() != "method") {
		return nil, nil, fmt.Errorf("[parseSubroutineDec] Invalid keyword %v want (constructor|function|method), %v", mayFuncType, rest)
	}
	res.AppendChild(mayFuncType)

	// return value type
	mayVoid, err := rest.LookAt(0)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineDec] %w", err)
	}
	if mayVoid.Type() == KeywordType && mayVoid.String() == "void" {
		res.AppendChild(mayVoid)
		_, r, err := rest.PopNext()
		if err != nil {
			return nil, nil, fmt.Errorf("[parseSubroutineDec] %w", err)
		}
		rest = r
	} else {
		typ, r, err := p.parseType(rest)
		if err != nil {
			return nil, nil, fmt.Errorf("[parseSubroutineDec] %w", err)
		}
		res.AppendChild(typ)
		rest = r
	}

	// subroutine name
	sn, rest, err := p.parseSubroutineName(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineDec] %w", err)
	}
	res.AppendChild(sn)

	// open paren
	mayOpenParen, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineDec] %w", err)
	}
	if mayOpenParen.Type() != SymbolType || mayOpenParen.String() != "(" {
		return nil, nil, fmt.Errorf("[parseSubroutineDec] Invalid symbol %v want (, %v", mayOpenParen, rest)
	}
	res.AppendChild(mayOpenParen)

	// parameter list
	pl, rest, err := p.parseParameterList(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineDec] %w", err)
	}
	res.AppendChild(pl)

	// close paren
	mayCloseParen, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineDec] %w", err)
	}
	if mayCloseParen.Type() != SymbolType || mayCloseParen.String() != ")" {
		return nil, nil, fmt.Errorf("[parseSubroutineDec] Invalid symbol %v want ), %v", mayCloseParen, rest)
	}
	res.AppendChild(mayCloseParen)

	// subroutine body
	sb, rest, err := p.parseSubroutineBody(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineDec] %w", err)
	}
	res.AppendChild(sb)

	return res, rest, nil
}

func (p *Parser) parseParameterList(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewParameterListNode()
	rest := tokens

	for true {
		r := rest
		// comma
		if len(res.ChildNodes()) > 0 {
			mayComma, rr, err := r.PopNext()
			if err != nil {
				break
			}
			if mayComma.Type() != SymbolType || mayComma.String() != "," {
				break
			}
			res.AppendChild(mayComma)
			r = rr
		}

		// var
		t, r, err := p.parseType(r)
		if err != nil {
			break
		}
		n, r, err := p.parseVarName(r)
		if err != nil {
			return nil, nil, fmt.Errorf("[parseParameterList] %w", err)
		}
		res.AppendChild(t)
		res.AppendChild(n)
		rest = r
	}

	return res, rest, nil
}

func (p *Parser) parseSubroutineBody(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewSubroutineBodyNode()

	// open bracket
	mayOpenBracket, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineBody] %w", err)
	}
	if mayOpenBracket.Type() != SymbolType || mayOpenBracket.String() != "{" {
		return nil, nil, fmt.Errorf("[parseSubroutineBody] Invalid symbol %v want {, %v", mayOpenBracket, rest)
	}
	res.AppendChild(mayOpenBracket)

	// var declarations
	for true {
		v, r, err := p.parseVarDec(rest)
		if err != nil {
			break
		}
		res.AppendChild(v)
		rest = r
	}

	// statements
	stms, rest, err := p.parseStatements(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineBody] %w", err)
	}
	res.AppendChild(stms)

	// close bracket
	mayCloseBracket, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineBody] %w", err)
	}
	if mayCloseBracket.Type() != SymbolType || mayCloseBracket.String() != "}" {
		return nil, nil, fmt.Errorf("[parseSubroutineBody] Invalid symbol %v want }, %v", mayCloseBracket, rest)
	}
	res.AppendChild(mayCloseBracket)

	return res, rest, nil
}

func (p *Parser) parseVarDec(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewVarDecNode()

	// var keyword
	mayVarKeyword, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseVarDec] %w", err)
	}
	if mayVarKeyword.Type() != KeywordType || mayVarKeyword.String() != "var" {
		return nil, nil, fmt.Errorf("[parseVarDec] Invalid keyword %v want var, %v", mayVarKeyword, rest)
	}
	res.AppendChild(mayVarKeyword)

	// type
	typ, rest, err := p.parseType(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseVarDec] %w", err)
	}
	res.AppendChild(typ)

	// varName
	v, rest, err := p.parseVarName(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseVarDec] %w", err)
	}
	res.AppendChild(v)

	// following vars
	for true {
		// comma
		mayComma, r, err := rest.PopNext()
		if err != nil {
			break
		}
		if mayComma.Type() != SymbolType || mayComma.String() != "," {
			break
		}

		// var
		v, r, err := p.parseVarName(r)
		if err != nil {
			return nil, nil, fmt.Errorf("[parseVarDec] %w", err)
		}

		res.AppendChild(mayComma)
		res.AppendChild(v)
		rest = r
	}

	// semicolon
	maySemicolon, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	if maySemicolon.Type() != SymbolType || maySemicolon.String() != ";" {
		return nil, nil, fmt.Errorf("[parseVarDec] Invalid Symbol %v want ;, %v", maySemicolon, rest)
	}
	res.AppendChild(maySemicolon)

	return res, rest, nil
}

func (p *Parser) parseClassName(tokens TokenList) (TreeNode, TokenList, error) {
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseClassName] %w", err)
	}
	if next.Type() != IdentifierType {
		return nil, nil, fmt.Errorf("[parseClassName] Type mismatch %v want IdentifierType, %v, %v", next.Type(), next, rest)
	}
	cn := NewClassNameNode()
	cn.AppendChild(next)
	return cn, rest, nil
}

func (p *Parser) parseSubroutineName(tokens TokenList) (TreeNode, TokenList, error) {
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineName] %w", err)
	}
	if next.Type() != IdentifierType {
		return nil, nil, fmt.Errorf("[parseSubroutineName] Type mismatch %v want IdentifierType, %v, %v", next.Type(), next, rest)
	}
	sn := NewSubroutineNameNode()
	sn.AppendChild(next)
	return sn, rest, nil
}

func (p *Parser) parseVarName(tokens TokenList) (TreeNode, TokenList, error) {
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseVarName] %w", err)
	}
	if next.Type() != IdentifierType {
		return nil, nil, fmt.Errorf("[parseVarName] Type mismatch %v want IdentifierType, %v, %v", next.Type(), next, rest)
	}
	vn := NewVarNameNode()
	vn.AppendChild(next)
	return vn, rest, nil
}

func (p *Parser) parseStatements(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewStatementsNode()
	rest := tokens

	for true {
		n, r, err := p.parseStatement(rest)
		if err != nil {
			break
		}
		res.AppendChild(n)
		rest = r
	}

	return res, rest, nil
}

func (p *Parser) parseStatement(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewStatementNode()

	if n, rest, err := p.parseLetStatement(tokens); err == nil {
		res.AppendChild(n)
		return res, rest, nil
	}

	if n, rest, err := p.parseIfStatement(tokens); err == nil {
		res.AppendChild(n)
		return res, rest, nil
	}

	if n, rest, err := p.parseWhileStatement(tokens); err == nil {
		res.AppendChild(n)
		return res, rest, nil
	}

	if n, rest, err := p.parseDoStatement(tokens); err == nil {
		res.AppendChild(n)
		return res, rest, nil
	}

	if n, rest, err := p.parseReturnStatement(tokens); err == nil {
		res.AppendChild(n)
		return res, rest, nil
	}

	return nil, nil, fmt.Errorf("[parseStatement] Invalid syntax %v", tokens)
}

func (p *Parser) parseLetStatement(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewLetStatementNode()

	// let keyword
	mayLetKeyword, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseLetStatement] %w", err)
	}
	if mayLetKeyword.Type() != KeywordType || mayLetKeyword.String() != "let" {
		return nil, nil, fmt.Errorf("[parseLetStatement] Invalid keyword %v want let, %v", mayLetKeyword, rest)
	}
	res.AppendChild(mayLetKeyword)

	// varName
	vn, rest, err := p.parseVarName(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseLetStatement] %w", err)
	}
	res.AppendChild(vn)

	// array index
	mayOpenSqBracket, err := rest.LookAt(0)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseLetStatement] %w", err)
	}
	if mayOpenSqBracket.Type() == SymbolType && mayOpenSqBracket.String() == "[" {
		openSqBracket, r, err := rest.PopNext()
		if err != nil {
			return nil, nil, fmt.Errorf("[parseLetStatement] %w", err)
		}
		res.AppendChild(openSqBracket)

		ex, r, err := p.parseExpression(r)
		if err != nil {
			return nil, nil, fmt.Errorf("[parseLetStatement] %w", err)
		}
		res.AppendChild(ex)

		mayCloseSqBracket, r, err := r.PopNext()
		if err != nil {
			return nil, nil, fmt.Errorf("[parseLetStatement] %w", err)
		}
		if mayCloseSqBracket.Type() != SymbolType || mayCloseSqBracket.String() != "]" {
			return nil, nil, fmt.Errorf("[parseLetStatement] Invalid symbol %v want ], %v", mayCloseSqBracket, rest)
		}
		res.AppendChild(mayCloseSqBracket)

		rest = r
	}

	// equal
	mayEqual, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	if mayEqual.Type() != SymbolType || mayEqual.String() != "=" {
		return nil, nil, fmt.Errorf("[parseLetStatement] Invalid symbol %v want =, %v", mayEqual, rest)
	}
	res.AppendChild(mayEqual)

	// expression
	ex, rest, err := p.parseExpression(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseLetStatement] %w", err)
	}
	res.AppendChild(ex)

	// semicolon
	maySemicolon, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseLetStatement] %w", err)
	}
	if maySemicolon.Type() != SymbolType || maySemicolon.String() != ";" {
		return nil, nil, fmt.Errorf("[parseLetStatement] Invalid symbol %v want ;, %v", maySemicolon, rest)
	}
	res.AppendChild(maySemicolon)

	return res, rest, nil
}

func (p *Parser) parseIfStatement(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewIfStatementNode()

	// if keyword
	mayIfKeyword, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseIfStatement] %w", err)
	}
	if mayIfKeyword.Type() != KeywordType || mayIfKeyword.String() != "if" {
		return nil, nil, fmt.Errorf("[parseIfStatement] Invalid keyword %v want if, %v", mayIfKeyword, rest)
	}
	res.AppendChild(mayIfKeyword)

	// open paren
	mayOpenParen, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseIfStatement] %w", err)
	}
	if mayOpenParen.Type() != SymbolType || mayOpenParen.String() != "(" {
		return nil, nil, fmt.Errorf("[parseIfStatement] Invalid keyword %v want (, %v", mayOpenParen, rest)
	}
	res.AppendChild(mayOpenParen)

	// expression
	ex, rest, err := p.parseExpression(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseIfStatement] %w", err)
	}
	res.AppendChild(ex)

	// close paren
	mayCloseParen, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseIfStatement] %w", err)
	}
	if mayCloseParen.Type() != SymbolType || mayCloseParen.String() != ")" {
		return nil, nil, fmt.Errorf("[parseIfStatement] Invalid keyword %v want ), %v", mayCloseParen, rest)
	}
	res.AppendChild(mayCloseParen)

	// open bracket
	mayOpenBracketIf, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseIfStatement] %w", err)
	}
	if mayOpenBracketIf.Type() != SymbolType || mayOpenBracketIf.String() != "{" {
		return nil, nil, fmt.Errorf("[parseIfStatement] Invalid keyword %v want {, %v", mayOpenBracketIf, rest)
	}
	res.AppendChild(mayOpenBracketIf)

	// statement
	stIf, rest, err := p.parseStatements(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseIfStatement] %w", err)
	}
	res.AppendChild(stIf)

	// close bracket
	mayCloseBracketIf, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseIfStatement] %w", err)
	}
	if mayCloseBracketIf.Type() != SymbolType || mayCloseBracketIf.String() != "}" {
		return nil, nil, fmt.Errorf("[parseIfStatement] Invalid keyword %v want }, %v", mayCloseBracketIf, rest)
	}
	res.AppendChild(mayCloseBracketIf)

	// else keyword
	mayElseKeyword, err := rest.LookAt(0)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseIfStatement] %w", err)
	}
	if mayElseKeyword.Type() != KeywordType || mayElseKeyword.String() != "else" {
		return res, rest, nil
	}
	_, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseIfStatement] %w", err)
	}
	res.AppendChild(mayElseKeyword)

	// open bracket
	mayOpenBracketElse, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseIfStatement] %w", err)
	}
	if mayOpenBracketElse.Type() != SymbolType || mayOpenBracketElse.String() != "{" {
		return nil, nil, fmt.Errorf("[parseIfStatement] Invalid keyword %v want {, %v", mayOpenBracketElse, rest)
	}
	res.AppendChild(mayOpenBracketElse)

	// statement
	stElse, rest, err := p.parseStatements(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseIfStatement] %w", err)
	}
	res.AppendChild(stElse)

	// close bracket
	mayCloseBracketElse, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseIfStatement] %w", err)
	}
	if mayCloseBracketElse.Type() != SymbolType || mayCloseBracketElse.String() != "}" {
		return nil, nil, fmt.Errorf("[parseIfStatement] Invalid keyword %v want }, %v", mayCloseBracketElse, rest)
	}
	res.AppendChild(mayCloseBracketElse)

	return res, rest, nil
}

func (p *Parser) parseWhileStatement(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewWhileStatementNode()

	// while keyword
	mayWhileKeyword, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseWhileStatement] %w", err)
	}
	if mayWhileKeyword.Type() != KeywordType || mayWhileKeyword.String() != "while" {
		return nil, nil, fmt.Errorf("[parseWhileStatement] Invalid keyword %v want while, %v", mayWhileKeyword, rest)
	}
	res.AppendChild(mayWhileKeyword)

	// open paren
	mayOpenParen, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseWhileStatement] %w", err)
	}
	if mayOpenParen.Type() != SymbolType || mayOpenParen.String() != "(" {
		return nil, nil, fmt.Errorf("[parseWhileStatement] Invalid keyword %v want (, %v", mayOpenParen, rest)
	}
	res.AppendChild(mayOpenParen)

	// expression
	ex, rest, err := p.parseExpression(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseWhileStatement] %w", err)
	}
	res.AppendChild(ex)

	// close paren
	mayCloseParen, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseWhileStatement] %w", err)
	}
	if mayCloseParen.Type() != SymbolType || mayCloseParen.String() != ")" {
		return nil, nil, fmt.Errorf("[parseWhileStatement] Invalid keyword %v want ), %v", mayCloseParen, rest)
	}
	res.AppendChild(mayCloseParen)

	// open bracket
	mayOpenBracket, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseWhileStatement] %w", err)
	}
	if mayOpenBracket.Type() != SymbolType || mayOpenBracket.String() != "{" {
		return nil, nil, fmt.Errorf("[parseWhileStatement] Invalid keyword %v want {, %v", mayOpenBracket, rest)
	}
	res.AppendChild(mayOpenBracket)

	// statements
	st, rest, err := p.parseStatements(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseWhileStatement] %w", err)
	}
	res.AppendChild(st)

	// close bracket
	mayCloseBracket, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseWhileStatement] %w", err)
	}
	if mayCloseBracket.Type() != SymbolType || mayCloseBracket.String() != "}" {
		return nil, nil, fmt.Errorf("[parseWhileStatement] Invalid keyword %v want }, %v", mayCloseBracket, rest)
	}
	res.AppendChild(mayCloseBracket)

	return res, rest, nil
}

func (p *Parser) parseDoStatement(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewDoStatementNode()

	// do keyword
	mayDoKeyword, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseDoStatement] %w", err)
	}
	if mayDoKeyword.Type() != KeywordType || mayDoKeyword.String() != "do" {
		return nil, nil, fmt.Errorf("[parseDoStatement] Invalid keyword %v want do, %v", mayDoKeyword, rest)
	}

	// subroutine call
	sc, rest, err := p.parseSubroutineCall(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseDoStatement] %w", err)
	}

	// semicolon
	maySemicolon, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseDoStatement] %w", err)
	}
	if maySemicolon.Type() != SymbolType || maySemicolon.String() != ";" {
		return nil, nil, fmt.Errorf("[parseDoStatement] Invalid symbol %v want ;, %v", maySemicolon, rest)
	}

	res.AppendChild(mayDoKeyword)
	res.AppendChild(sc)
	res.AppendChild(maySemicolon)

	return res, rest, nil
}

func (p *Parser) parseReturnStatement(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewReturnStatementNode()

	// return keyword
	mayReturn, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseReturnStatement] %w", err)
	}
	if mayReturn.Type() != KeywordType || mayReturn.String() != "return" {
		return nil, nil, fmt.Errorf("[parseReturnStatement] Invalid keyword %v want return, %v", mayReturn, rest)
	}

	// expression
	ex, exRest, exErr := p.parseExpression(rest)
	if exErr == nil {
		rest = exRest
	}

	// `;` symbol
	maySemicolon, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseReturnStatement] %w", err)
	}
	if maySemicolon.Type() != SymbolType || maySemicolon.String() != ";" {
		return nil, nil, fmt.Errorf("[parseReturnStatement] Invalid symbol %v want ;, %v", maySemicolon, rest)
	}

	res.AppendChild(mayReturn)
	if exErr == nil {
		res.AppendChild(ex)
	}
	res.AppendChild(maySemicolon)

	return res, rest, nil
}

func (p *Parser) parseExpression(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewExpressionNode()

	// first term
	t1, rest, err := p.parseTerm(tokens)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseExpression] %w", err)
	}
	res.AppendChild(t1)

	// op (optional)
	mayOp, rest2, err := p.parseOp(rest)
	if err != nil {
		return res, rest, nil
	}
	res.AppendChild(mayOp)

	// second term if op exists
	t2, rest, err := p.parseTerm(rest2)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseExpression] %w", err)
	}
	res.AppendChild(t2)

	return res, rest, nil
}

func (p *Parser) parseTerm(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewTermNode()

	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseTerm] %w", err)
	}

	// Int and string constant
	switch next.Type() {
	case IntConstType, StrConstType:
		res.AppendChild(next)
		return res, rest, nil
	}

	// Keyword Constant
	if kw, rest, err := p.parseKeywordConstant(tokens); err == nil {
		res.AppendChild(kw)
		return res, rest, nil
	}

	isSubCall := false
	mayVarNameOrSub, err := tokens.LookAt(0)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseTerm] %w", err)
	}
	if mayVarNameOrSub.Type() == IdentifierType {
		next, err := tokens.LookAt(1)
		if err == nil {
			mayFuncCall, ok := next.(SymbolToken)
			if ok && (mayFuncCall.String() == "(" || mayFuncCall.String() == ".") {
				isSubCall = true
			}
		}
	}

	// VarName
	if mayVarNameOrSub.Type() == IdentifierType && !isSubCall {
		if vn, rest, err := p.parseVarName(tokens); err == nil {
			res.AppendChild(vn)

			// VarName only
			mayOpenSqBracket, err := rest.LookAt(0)
			if err != nil {
				return nil, nil, fmt.Errorf("[parseTerm] %w", err)
			}
			if mayOpenSqBracket.Type() != SymbolType || mayOpenSqBracket.String() != "[" {
				return res, rest, nil
			}

			// Array index
			_, rest, err := rest.PopNext()
			if err != nil {
				return nil, nil, fmt.Errorf("[parseTerm] %w", err)
			}
			ex, rest, err := p.parseExpression(rest)
			if err != nil {
				return nil, nil, fmt.Errorf("[parseTerm] %w", err)
			}
			mayCloseSqBracket, rest, err := rest.PopNext()
			if err != nil {
				return nil, nil, fmt.Errorf("[parseTerm] %w", err)
			}
			if mayCloseSqBracket.Type() != SymbolType || mayCloseSqBracket.String() != "]" {
				return nil, nil, fmt.Errorf("[parseTerm] Invalid symbol %v want ], %v", mayCloseSqBracket, rest)
			}
			res.AppendChild(mayOpenSqBracket)
			res.AppendChild(ex)
			res.AppendChild(mayCloseSqBracket)
			return res, rest, nil
		}
	}

	// subroutineCall
	if mayVarNameOrSub.Type() == IdentifierType && isSubCall {
		if sc, rest, err := p.parseSubroutineCall(tokens); err == nil {
			res.AppendChild(sc)
			return res, rest, nil
		}
	}

	// expression enclosed in paren
	mayOpenParen := next
	if mayOpenParen.Type() == SymbolType && mayOpenParen.String() == "(" {
		ex, rest, err := p.parseExpression(rest)
		if err != nil {
			return nil, nil, fmt.Errorf("[parseTerm] %w", err)
		}
		mayCloseParen, ok := rest[0].(SymbolToken)
		if !ok || mayCloseParen.String() != ")" {
			return nil, nil, fmt.Errorf("[[arseTerm]] Invalid symbol %v want ), %v", mayCloseParen, rest)
		}
		res.AppendChild(mayOpenParen)
		res.AppendChild(ex)
		res.AppendChild(mayCloseParen)
		rest = rest[1:len(rest)]
		return res, rest, nil
	}

	// unaryOp
	if uo, rest, err := p.parseUnaryOp(tokens); err == nil {
		tm, rest, err := p.parseTerm(rest)
		if err != nil {
			return nil, nil, fmt.Errorf("[parseTerm] %w", err)
		}
		res.AppendChild(uo)
		res.AppendChild(tm)
		return res, rest, nil
	}

	return nil, nil, fmt.Errorf("Invalid syntax %v", tokens)
}

func (p *Parser) parseSubroutineCall(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewSubroutineCallNode()

	del, err := tokens.LookAt(1)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineCall] %w", err)
	}
	if del.Type() != SymbolType || (del.String() != "(" && del.String() != ".") {
		return nil, nil, fmt.Errorf("[parseSubroutineCall] Invalid symbol %v want ((|.)), %v", del, tokens)
	}

	// method call
	rest := tokens
	if del.String() == "." {
		cn, innerRest, err := p.parseClassName(rest) // TODO: how to detect className or varName
		if err != nil {
			return nil, nil, fmt.Errorf("[parseSubroutineCall] %w", err)
		}
		mayDot, innerRest2, err := innerRest.PopNext()
		if err != nil {
			return nil, nil, fmt.Errorf("[parseSubroutineCall] %w", err)
		}
		if mayDot.Type() != SymbolType || mayDot.String() != "." {
			return nil, nil, fmt.Errorf("[parseSubroutineCall] Invalid symbol %v want ., %v", mayDot, innerRest)
		}
		rest = innerRest2
		res.AppendChild(cn)
		res.AppendChild(mayDot)
	}

	// function call
	sn, rest, err := p.parseSubroutineName(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineCall] %w", err)
	}
	mayOpenParen, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineCall] %w", err)
	}
	if mayOpenParen.Type() != SymbolType || mayOpenParen.String() != "(" {
		return nil, nil, fmt.Errorf("[parseSubroutineCall] Invalid symbol %v want (, %v", mayOpenParen, rest)
	}
	el, rest, err := p.parseExpressionList(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineCall] %w", err)
	}
	mayCloseParen, rest, err := rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineCall] %w", err)
	}
	if mayCloseParen.Type() != SymbolType || mayCloseParen.String() != ")" {
		return nil, nil, fmt.Errorf("[parseSubroutineCall] Invalid symbol %v want ), %v", mayCloseParen, rest)
	}
	res.AppendChild(sn)
	res.AppendChild(mayOpenParen)
	res.AppendChild(el)
	res.AppendChild(mayCloseParen)
	return res, rest, nil
}

func (p *Parser) parseExpressionList(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewExpressionListNode()
	rest := tokens
	for true {
		if len(res.ChildNodes()) > 0 {
			mayComma, r, err := rest.PopNext()
			if err != nil {
				return nil, nil, fmt.Errorf("[parseExpressionList] %w", err)
			}
			if mayComma.Type() != SymbolType || mayComma.String() != "," {
				break
			}
			res.AppendChild(mayComma)
			rest = r
		}
		ex, r, err := p.parseExpression(rest)
		if err != nil {
			break
		}
		res.AppendChild(ex)
		rest = r
	}
	return res, rest, nil
}

func (p *Parser) parseOp(tokens TokenList) (TreeNode, TokenList, error) {
	cur, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseOp] %w", err)
	}
	if cur.Type() != SymbolType {
		return nil, nil, fmt.Errorf("[parseOp] Type mismatch %v want SymbolToken, %v", cur, rest)
	}
	ok := false
	switch cur.String() {
	case "+", "-", "*", "/", "&", "|", "<", ">", "=":
		ok = true
	}
	if !ok {
		return nil, nil, fmt.Errorf("[parseOp] Invalid keyword %v want (+|-|*|/|&|||<|>|=), %v", cur, rest)
	}
	on := NewOpNode()
	on.AppendChild(cur)
	return on, rest, nil
}

func (p *Parser) parseUnaryOp(tokens TokenList) (TreeNode, TokenList, error) {
	cur, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseUnaryOp] %w", err)
	}
	if cur.Type() != SymbolType || (cur.String() != "-" && cur.String() != "~") {
		return nil, nil, fmt.Errorf("[parseUnaryOp] Invalid symbol %v want (-|~), %v", cur, rest)
	}
	uon := NewUnaryOpNode()
	uon.AppendChild(cur)
	return uon, rest, nil
}

func (p *Parser) parseKeywordConstant(tokens TokenList) (TreeNode, TokenList, error) {
	cur, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseKeywordConstant] %w", err)
	}
	if cur.Type() != KeywordType || (cur.String() != "true" && cur.String() != "false" && cur.String() != "null" && cur.String() != "this") {
		return nil, nil, fmt.Errorf("[parseKeywordConstant] Invalid keyword %v want (true|false|null|this), %v", cur, rest)
	}
	kc := NewKeyConstNode()
	kc.AppendChild(cur)
	return kc, rest, nil
}

type TokenList []Token

func NewTokenList(t []Token) TokenList {
	return TokenList(t)
}

func (l TokenList) PopAt(at int) (Token, TokenList, error) {
	if len(l) <= at {
		return nil, nil, fmt.Errorf("Invalid index %d", at)
	}
	tkn := l[at]
	rest := NewTokenList(l[at+1 : len(l)])
	return tkn, rest, nil
}

func (l TokenList) PopNext() (Token, TokenList, error) {
	return l.PopAt(0)
}

func (l TokenList) LookAt(at int) (Token, error) {
	if len(l) <= at {
		return nil, fmt.Errorf("Invalid index %d", at)
	}
	tkn := l[at]
	return tkn, nil
}
