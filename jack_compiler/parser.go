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
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseClass] %w", err)
	}
	mayClassKeyword, ok := next.(KeywordToken)
	if !ok || mayClassKeyword.String() != "class" {
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
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseClass] %w", err)
	}
	mayOpenBracket, ok := next.(SymbolToken)
	if !ok || mayOpenBracket.String() != "{" {
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
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseClass] %w", err)
	}
	mayCloseBracket, ok := next.(SymbolToken)
	if !ok || mayCloseBracket.String() != "}" {
		return nil, nil, fmt.Errorf("[parseClass] Invalid symbol %v want }, %v", mayCloseBracket, rest)
	}
	res.AppendChild(mayCloseBracket)

	return res, rest, nil
}

func (p *Parser) parseClassVarDec(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewClassVarDecNode()

	// class var type
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseClassVarDec] %w", err)
	}
	mayClassVarType, ok := next.(KeywordToken)
	if !ok || (mayClassVarType.String() != "static" && mayClassVarType.String() != "field") {
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
		next, r, err := rest.PopNext()
		if err != nil {
			break
		}
		mayComma, ok := next.(SymbolToken)
		if !ok || mayComma.String() != "," {
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
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseClassVarDec] %w", err)
	}
	maySemicolon, ok := next.(SymbolToken)
	if !ok || maySemicolon.String() != ";" {
		return nil, nil, fmt.Errorf("[parseClassVarDec] Invalid keyword %v want ;, %v", maySemicolon, rest)
	}
	res.AppendChild(maySemicolon)

	return res, rest, nil
}

func (p *Parser) parseType(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewTypeNode()

	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseType] %w", err)
	}
	kw, ok := next.(KeywordToken)
	if ok {
		switch kw.String() {
		case "int", "char", "boolean":
			res.AppendChild(kw)
			return res, rest, nil
		}
		return nil, nil, fmt.Errorf("[parseType] Invalid keyword %v want (int|char|boolean), %v", kw, rest)
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
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineDec] %w", err)
	}
	mayFuncType, ok := next.(KeywordToken)
	if !ok || (mayFuncType.String() != "constructor" && mayFuncType.String() != "function" && mayFuncType.String() != "method") {
		return nil, nil, fmt.Errorf("[parseSubroutineDec] Invalid keyword %v want (constructor|function|method), %v", mayFuncType, rest)
	}
	res.AppendChild(mayFuncType)

	// return value type
	next, err = rest.LookAt(0)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineDec] %w", err)
	}
	mayVoid, ok := next.(KeywordToken)
	if ok && mayVoid.String() == "void" {
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
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineDec] %w", err)
	}
	mayOpenParen, ok := next.(SymbolToken)
	if !ok || mayOpenParen.String() != "(" {
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
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineDec] %w", err)
	}
	mayCloseParen, ok := next.(SymbolToken)
	if !ok || mayCloseParen.String() != ")" {
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
			next, rr, err := r.PopNext()
			if err != nil {
				break
			}
			mayComma, ok := next.(SymbolToken)
			if !ok || mayComma.String() != "," {
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
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineBody] %w", err)
	}
	mayOpenBracket, ok := next.(SymbolToken)
	if !ok || mayOpenBracket.String() != "{" {
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
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineBody] %w", err)
	}
	mayCloseBracket, ok := next.(SymbolToken)
	if !ok || mayCloseBracket.String() != "}" {
		return nil, nil, fmt.Errorf("[parseSubroutineBody] Invalid symbol %v want }, %v", mayCloseBracket, rest)
	}
	res.AppendChild(mayCloseBracket)

	return res, rest, nil
}

func (p *Parser) parseVarDec(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewVarDecNode()

	// var keyword
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseVarDec] %w", err)
	}
	mayVarKeyword, ok := next.(KeywordToken)
	if !ok || mayVarKeyword.String() != "var" {
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
		next, r, err := rest.PopNext()
		if err != nil {
			break
		}
		mayComma, ok := next.(SymbolToken)
		if !ok || mayComma.String() != "," {
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
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	maySemicolon, ok := next.(SymbolToken)
	if !ok || maySemicolon.String() != ";" {
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
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseLetStatement] %w", err)
	}
	mayLetKeyword, ok := next.(KeywordToken)
	if !ok || mayLetKeyword.String() != "let" {
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
	next, err = rest.LookAt(0)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseLetStatement] %w", err)
	}
	mayOpenSqBracket, ok := next.(SymbolToken)
	if ok && mayOpenSqBracket.String() == "[" {
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

		next, r, err = r.PopNext()
		if err != nil {
			return nil, nil, fmt.Errorf("[parseLetStatement] %w", err)
		}
		mayCloseSqBracket, ok := next.(SymbolToken)
		if !ok || mayCloseSqBracket != "]" {
			return nil, nil, fmt.Errorf("[parseLetStatement] Invalid symbol %v want ], %v", mayCloseSqBracket, rest)
		}
		res.AppendChild(mayCloseSqBracket)

		rest = r
	}

	// equal
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayEqual, ok := next.(SymbolToken)
	if !ok || mayEqual.String() != "=" {
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
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseLetStatement] %w", err)
	}
	maySemicolon, ok := next.(SymbolToken)
	if !ok || maySemicolon.String() != ";" {
		return nil, nil, fmt.Errorf("[parseLetStatement] Invalid symbol %v want ;, %v", maySemicolon, rest)
	}
	res.AppendChild(maySemicolon)

	return res, rest, nil
}

func (p *Parser) parseIfStatement(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewIfStatementNode()

	// if keyword
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseIfStatement] %w", err)
	}
	mayIfKeyword, ok := next.(KeywordToken)
	if !ok || mayIfKeyword.String() != "if" {
		return nil, nil, fmt.Errorf("[parseIfStatement] Invalid keyword %v want if, %v", mayIfKeyword, rest)
	}
	res.AppendChild(mayIfKeyword)

	// open paren
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseIfStatement] %w", err)
	}
	mayOpenParen, ok := next.(SymbolToken)
	if !ok || mayOpenParen.String() != "(" {
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
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseIfStatement] %w", err)
	}
	mayCloseParen, ok := next.(SymbolToken)
	if !ok || mayCloseParen.String() != ")" {
		return nil, nil, fmt.Errorf("[parseIfStatement] Invalid keyword %v want ), %v", mayCloseParen, rest)
	}
	res.AppendChild(mayCloseParen)

	// open bracket
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseIfStatement] %w", err)
	}
	mayOpenBracketIf, ok := next.(SymbolToken)
	if !ok || mayOpenBracketIf.String() != "{" {
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
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseIfStatement] %w", err)
	}
	mayCloseBracketIf, ok := next.(SymbolToken)
	if !ok || mayCloseBracketIf.String() != "}" {
		return nil, nil, fmt.Errorf("[parseIfStatement] Invalid keyword %v want }, %v", mayCloseBracketIf, rest)
	}
	res.AppendChild(mayCloseBracketIf)

	// else keyword
	next, err = rest.LookAt(0)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseIfStatement] %w", err)
	}
	mayElseKeyword, ok := next.(KeywordToken)
	if !ok || mayElseKeyword.String() != "else" {
		return res, rest, nil
	}
	_, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseIfStatement] %w", err)
	}
	res.AppendChild(mayElseKeyword)

	// open bracket
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseIfStatement] %w", err)
	}
	mayOpenBracketElse, ok := next.(SymbolToken)
	if !ok || mayOpenBracketElse.String() != "{" {
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
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseIfStatement] %w", err)
	}
	mayCloseBracketElse, ok := next.(SymbolToken)
	if !ok || mayCloseBracketElse.String() != "}" {
		return nil, nil, fmt.Errorf("[parseIfStatement] Invalid keyword %v want }, %v", mayCloseBracketElse, rest)
	}
	res.AppendChild(mayCloseBracketElse)

	return res, rest, nil
}

func (p *Parser) parseWhileStatement(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewWhileStatementNode()

	// while keyword
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseWhileStatement] %w", err)
	}
	mayWhileKeyword, ok := next.(KeywordToken)
	if !ok || mayWhileKeyword.String() != "while" {
		return nil, nil, fmt.Errorf("[parseWhileStatement] Invalid keyword %v want while, %v", mayWhileKeyword, rest)
	}
	res.AppendChild(mayWhileKeyword)

	// open paren
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseWhileStatement] %w", err)
	}
	mayOpenParen, ok := next.(SymbolToken)
	if !ok || mayOpenParen.String() != "(" {
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
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseWhileStatement] %w", err)
	}
	mayCloseParen, ok := next.(SymbolToken)
	if !ok || mayCloseParen.String() != ")" {
		return nil, nil, fmt.Errorf("[parseWhileStatement] Invalid keyword %v want ), %v", mayCloseParen, rest)
	}
	res.AppendChild(mayCloseParen)

	// open bracket
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseWhileStatement] %w", err)
	}
	mayOpenBracket, ok := next.(SymbolToken)
	if !ok || mayOpenBracket.String() != "{" {
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
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseWhileStatement] %w", err)
	}
	mayCloseBracket, ok := next.(SymbolToken)
	if !ok || mayCloseBracket.String() != "}" {
		return nil, nil, fmt.Errorf("[parseWhileStatement] Invalid keyword %v want }, %v", mayCloseBracket, rest)
	}
	res.AppendChild(mayCloseBracket)

	return res, rest, nil
}

func (p *Parser) parseDoStatement(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewDoStatementNode()

	// do keyword
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseDoStatement] %w", err)
	}
	mayDoKeyword, ok := next.(KeywordToken)
	if !ok || mayDoKeyword.String() != "do" {
		return nil, nil, fmt.Errorf("[parseDoStatement] Invalid keyword %v want do, %v", mayDoKeyword, rest)
	}

	// subroutine call
	sc, rest, err := p.parseSubroutineCall(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseDoStatement] %w", err)
	}

	// semicolon
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseDoStatement] %w", err)
	}
	maySemicolon, ok := next.(SymbolToken)
	if !ok || maySemicolon.String() != ";" {
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
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseReturnStatement] %w", err)
	}
	mayReturn, ok := next.(KeywordToken)
	if !ok || mayReturn.String() != "return" {
		return nil, nil, fmt.Errorf("[parseReturnStatement] Invalid keyword %v want return, %v", mayReturn, rest)
	}

	// expression
	ex, exRest, exErr := p.parseExpression(rest)
	if exErr == nil {
		rest = exRest
	}

	// `;` symbol
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseReturnStatement] %w", err)
	}
	maySemicolon, ok := next.(SymbolToken)
	if !ok || maySemicolon.String() != ";" {
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
	next, err = tokens.LookAt(0)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseTerm] %w", err)
	}
	_, mayVarNameOrSub := next.(IdentifierToken)
	if mayVarNameOrSub {
		next, err := tokens.LookAt(1)
		if err == nil {
			mayFuncCall, ok := next.(SymbolToken)
			if ok && (mayFuncCall.String() == "(" || mayFuncCall.String() == ".") {
				isSubCall = true
			}
		}
	}

	// VarName
	if mayVarNameOrSub && !isSubCall {
		if vn, rest, err := p.parseVarName(tokens); err == nil {
			res.AppendChild(vn)

			// VarName only
			next, err := rest.LookAt(0)
			if err != nil {
				return nil, nil, fmt.Errorf("[parseTerm] %w", err)
			}
			mayOpenSqBracket, ok := next.(SymbolToken)
			if !ok || mayOpenSqBracket.String() != "[" {
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
			next, rest, err = rest.PopNext()
			if err != nil {
				return nil, nil, fmt.Errorf("[parseTerm] %w", err)
			}
			mayCloseSqBracket, ok := next.(SymbolToken)
			if !ok || mayCloseSqBracket.String() != "]" {
				return nil, nil, fmt.Errorf("[parseTerm] Invalid symbol %v want ], %v", mayCloseSqBracket, rest)
			}
			res.AppendChild(mayOpenSqBracket)
			res.AppendChild(ex)
			res.AppendChild(mayCloseSqBracket)
			return res, rest, nil
		}
	}

	// subroutineCall
	if mayVarNameOrSub && isSubCall {
		if sc, rest, err := p.parseSubroutineCall(tokens); err == nil {
			res.AppendChild(sc)
			return res, rest, nil
		}
	}

	// expression enclosed in paren
	if mayOpenParen, ok := next.(SymbolToken); ok && mayOpenParen.String() == "(" {
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

	next, err := tokens.LookAt(1)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineCall] %w", err)
	}
	del, ok := next.(SymbolToken)
	if !ok || (del.String() != "(" && del.String() != ".") {
		return nil, nil, fmt.Errorf("[parseSubroutineCall] Invalid symbol %v want ((|.)), %v", del, tokens)
	}

	// method call
	rest := tokens
	if del.String() == "." {
		cn, innerRest, err := p.parseClassName(rest) // TODO: how to detect className or varName
		if err != nil {
			return nil, nil, fmt.Errorf("[parseSubroutineCall] %w", err)
		}
		next, innerRest2, err := innerRest.PopNext()
		if err != nil {
			return nil, nil, fmt.Errorf("[parseSubroutineCall] %w", err)
		}
		mayDot, ok := next.(SymbolToken)
		if !ok || mayDot != "." {
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
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineCall] %w", err)
	}
	mayOpenParen, ok := next.(SymbolToken)
	if !ok || mayOpenParen != "(" {
		return nil, nil, fmt.Errorf("[parseSubroutineCall] Invalid symbol %v want (, %v", mayOpenParen, rest)
	}
	el, rest, err := p.parseExpressionList(rest)
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineCall] %w", err)
	}
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseSubroutineCall] %w", err)
	}
	mayCloseParen, ok := next.(SymbolToken)
	if !ok || mayCloseParen != ")" {
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
			next, r, err := rest.PopNext()
			if err != nil {
				return nil, nil, fmt.Errorf("[parseExpressionList] %w", err)
			}
			mayComma, ok := next.(SymbolToken)
			if !ok || mayComma.String() != "," {
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
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseOp] %w", err)
	}
	cur, ok := next.(SymbolToken)
	if !ok {
		return nil, nil, fmt.Errorf("[parseOp] Type mismatch %v want SymbolToken, %v", cur, rest)
	}
	ok = false
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
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseUnaryOp] %w", err)
	}
	cur, ok := next.(SymbolToken)
	if !ok || (cur.String() != "-" && cur.String() != "~") {
		return nil, nil, fmt.Errorf("[parseUnaryOp] Invalid symbol %v want (-|~), %v", cur, rest)
	}
	uon := NewUnaryOpNode()
	uon.AppendChild(cur)
	return uon, rest, nil
}

func (p *Parser) parseKeywordConstant(tokens TokenList) (TreeNode, TokenList, error) {
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, fmt.Errorf("[parseKeywordConstant] %w", err)
	}
	cur, ok := next.(KeywordToken)
	if !ok || (cur.String() != "true" && cur.String() != "false" && cur.String() != "null" && cur.String() != "this") {
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
