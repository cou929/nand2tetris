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
		return nil, nil, err
	}
	kw, ok := next.(KeywordToken)
	if ok {
		switch kw.String() {
		case "int", "char", "boolean":
			res.AppendChild(kw)
			return res, rest, nil
		}
		return nil, nil, fmt.Errorf("Invalid keyword %v", tokens)
	}

	n, rest, err := p.parseClassName(tokens)
	if err != nil {
		return nil, nil, err
	}
	res.AppendChild(n)
	return res, rest, nil
}

func (p *Parser) parseSubroutineDec(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewSubroutineDecNode()

	// func type
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayFuncType, ok := next.(KeywordToken)
	if !ok || (mayFuncType.String() != "constructor" && mayFuncType.String() != "function" && mayFuncType.String() != "method") {
		return nil, nil, fmt.Errorf("Invalid keyword %v", tokens)
	}
	res.AppendChild(mayFuncType)

	// return value type
	next, err = rest.LookAt(0)
	if err != nil {
		return nil, nil, err
	}
	mayVoid, ok := next.(KeywordToken)
	if ok && mayVoid.String() == "void" {
		res.AppendChild(mayVoid)
		_, r, err := rest.PopNext()
		if err != nil {
			return nil, nil, err
		}
		rest = r
	} else {
		typ, r, err := p.parseType(rest)
		if err != nil {
			return nil, nil, err
		}
		res.AppendChild(typ)
		rest = r
	}

	// subroutine name
	sn, rest, err := p.parseSubroutineName(rest)
	if err != nil {
		return nil, nil, err
	}
	res.AppendChild(sn)

	// open paren
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayOpenParen, ok := next.(SymbolToken)
	if !ok || mayOpenParen.String() != "(" {
		return nil, nil, fmt.Errorf("Invalid symbol %v", tokens)
	}
	res.AppendChild(mayOpenParen)

	// parameter list
	pl, rest, err := p.parseParameterList(rest)
	if err != nil {
		return nil, nil, err
	}
	res.AppendChild(pl)

	// close paren
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayCloseParen, ok := next.(SymbolToken)
	if !ok || mayCloseParen.String() != ")" {
		return nil, nil, fmt.Errorf("Invalid symbol %v", tokens)
	}
	res.AppendChild(mayCloseParen)

	// subroutine body
	sb, rest, err := p.parseSubroutineBody(rest)
	if err != nil {
		return nil, nil, err
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
			return nil, nil, err
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
		return nil, nil, err
	}
	mayOpenBracket, ok := next.(SymbolToken)
	if !ok || mayOpenBracket.String() != "{" {
		return nil, nil, fmt.Errorf("Invalid symbol %v", tokens)
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
		return nil, nil, err
	}
	res.AppendChild(stms)

	// close bracket
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayCloseBracket, ok := next.(SymbolToken)
	if !ok || mayCloseBracket.String() != "}" {
		return nil, nil, fmt.Errorf("Invalid symbol %v", tokens)
	}
	res.AppendChild(mayCloseBracket)

	return res, rest, nil
}

func (p *Parser) parseVarDec(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewVarDecNode()

	// var keyword
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayVarKeyword, ok := next.(KeywordToken)
	if !ok || mayVarKeyword.String() != "var" {
		return nil, nil, fmt.Errorf("Invalid keyword %v", tokens)
	}
	res.AppendChild(mayVarKeyword)

	// type
	typ, rest, err := p.parseType(rest)
	if err != nil {
		return nil, nil, err
	}
	res.AppendChild(typ)

	// varName
	v, rest, err := p.parseVarName(rest)
	if err != nil {
		return nil, nil, err
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
			break
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
		return nil, nil, fmt.Errorf("Invalid Symbol %v", tokens)
	}
	res.AppendChild(maySemicolon)

	return res, rest, nil
}

func (p *Parser) parseClassName(tokens TokenList) (TreeNode, TokenList, error) {
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, err
	}
	if next.Type() != IdentifierType {
		return nil, nil, fmt.Errorf("Type mismatch %v", tokens)
	}
	cn := NewClassNameNode()
	cn.AppendChild(next)
	return cn, rest, nil
}

func (p *Parser) parseSubroutineName(tokens TokenList) (TreeNode, TokenList, error) {
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, err
	}
	if next.Type() != IdentifierType {
		return nil, nil, fmt.Errorf("Type mismatch %v", tokens)
	}
	sn := NewSubroutineNameNode()
	sn.AppendChild(next)
	return sn, rest, nil
}

func (p *Parser) parseVarName(tokens TokenList) (TreeNode, TokenList, error) {
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, err
	}
	if next.Type() != IdentifierType {
		return nil, nil, fmt.Errorf("Type mismatch %v", tokens)
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

	return nil, nil, fmt.Errorf("Invalid syntax %v", tokens)
}

func (p *Parser) parseLetStatement(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewLetStatementNode()

	// let keyword
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayLetKeyword, ok := next.(KeywordToken)
	if !ok || mayLetKeyword.String() != "let" {
		return nil, nil, fmt.Errorf("Invalid keyword %v", tokens)
	}
	res.AppendChild(mayLetKeyword)

	// varName
	vn, rest, err := p.parseVarName(rest)
	if err != nil {
		return nil, nil, err
	}
	res.AppendChild(vn)

	// array index
	next, err = rest.LookAt(0)
	if err != nil {
		return nil, nil, err
	}
	mayOpenSqBracket, ok := next.(SymbolToken)
	if ok && mayOpenSqBracket.String() == "[" {
		openSqBracket, r, err := rest.PopNext()
		if err != nil {
			return nil, nil, err
		}
		res.AppendChild(openSqBracket)

		ex, r, err := p.parseExpression(r)
		if err != nil {
			return nil, nil, err
		}
		res.AppendChild(ex)

		next, r, err = r.PopNext()
		if err != nil {
			return nil, nil, err
		}
		mayCloseSqBracket, ok := next.(SymbolToken)
		if !ok || mayCloseSqBracket != "]" {
			return nil, nil, fmt.Errorf("Invalid symbol %v", r)
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
		return nil, nil, fmt.Errorf("Invalid symbol %v", tokens)
	}
	res.AppendChild(mayEqual)

	// expression
	ex, rest, err := p.parseExpression(rest)
	if err != nil {
		return nil, nil, err
	}
	res.AppendChild(ex)

	// semicolon
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	maySemicolon, ok := next.(SymbolToken)
	if !ok || maySemicolon.String() != ";" {
		return nil, nil, fmt.Errorf("Invalid symbol %v", tokens)
	}
	res.AppendChild(maySemicolon)

	return res, rest, nil
}

func (p *Parser) parseIfStatement(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewIfStatementNode()

	// if keyword
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayIfKeyword, ok := next.(KeywordToken)
	if !ok || mayIfKeyword.String() != "if" {
		return nil, nil, fmt.Errorf("Invalid keyword %v", tokens)
	}
	res.AppendChild(mayIfKeyword)

	// open paren
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayOpenParen, ok := next.(SymbolToken)
	if !ok || mayOpenParen.String() != "(" {
		return nil, nil, fmt.Errorf("Invalid keyword %v", tokens)
	}
	res.AppendChild(mayOpenParen)

	// expression
	ex, rest, err := p.parseExpression(rest)
	if err != nil {
		return nil, nil, err
	}
	res.AppendChild(ex)

	// close paren
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayCloseParen, ok := next.(SymbolToken)
	if !ok || mayCloseParen.String() != ")" {
		return nil, nil, fmt.Errorf("Invalid keyword %v", tokens)
	}
	res.AppendChild(mayCloseParen)

	// open bracket
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayOpenBracketIf, ok := next.(SymbolToken)
	if !ok || mayOpenBracketIf.String() != "{" {
		return nil, nil, fmt.Errorf("Invalid keyword %v", tokens)
	}
	res.AppendChild(mayOpenBracketIf)

	// statement
	stIf, rest, err := p.parseStatement(rest)
	if err != nil {
		return nil, nil, err
	}
	res.AppendChild(stIf)

	// close bracket
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayCloseBracketIf, ok := next.(SymbolToken)
	if !ok || mayCloseBracketIf.String() != "}" {
		return nil, nil, fmt.Errorf("Invalid keyword %v", tokens)
	}
	res.AppendChild(mayCloseBracketIf)

	// else keyword
	next, err = rest.LookAt(0)
	if err != nil {
		return nil, nil, err
	}
	mayElseKeyword, ok := next.(KeywordToken)
	if !ok || mayElseKeyword.String() != "else" {
		return res, rest, nil
	}
	_, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	res.AppendChild(mayElseKeyword)

	// open bracket
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayOpenBracketElse, ok := next.(SymbolToken)
	if !ok || mayOpenBracketElse.String() != "{" {
		return nil, nil, fmt.Errorf("Invalid keyword %v", tokens)
	}
	res.AppendChild(mayOpenBracketElse)

	// statement
	stElse, rest, err := p.parseStatement(rest)
	if err != nil {
		return nil, nil, err
	}
	res.AppendChild(stElse)

	// close bracket
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayCloseBracketElse, ok := next.(SymbolToken)
	if !ok || mayCloseBracketElse.String() != "}" {
		return nil, nil, fmt.Errorf("Invalid keyword %v", tokens)
	}
	res.AppendChild(mayCloseBracketElse)

	return res, rest, nil
}

func (p *Parser) parseWhileStatement(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewWhileStatementNode()

	// while keyword
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayWhileKeyword, ok := next.(KeywordToken)
	if !ok || mayWhileKeyword.String() != "while" {
		return nil, nil, fmt.Errorf("Invalid keyword %v", tokens)
	}
	res.AppendChild(mayWhileKeyword)

	// open paren
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayOpenParen, ok := next.(SymbolToken)
	if !ok || mayOpenParen.String() != "(" {
		return nil, nil, fmt.Errorf("Invalid keyword %v", tokens)
	}
	res.AppendChild(mayOpenParen)

	// expression
	ex, rest, err := p.parseExpression(rest)
	if err != nil {
		return nil, nil, err
	}
	res.AppendChild(ex)

	// close paren
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayCloseParen, ok := next.(SymbolToken)
	if !ok || mayCloseParen.String() != ")" {
		return nil, nil, fmt.Errorf("Invalid keyword %v", tokens)
	}
	res.AppendChild(mayCloseParen)

	// open bracket
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayOpenBracket, ok := next.(SymbolToken)
	if !ok || mayOpenBracket.String() != "{" {
		return nil, nil, fmt.Errorf("Invalid keyword %v", tokens)
	}
	res.AppendChild(mayOpenBracket)

	// statement
	st, rest, err := p.parseStatement(rest)
	if err != nil {
		return nil, nil, err
	}
	res.AppendChild(st)

	// close bracket
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayCloseBracket, ok := next.(SymbolToken)
	if !ok || mayCloseBracket.String() != "}" {
		return nil, nil, fmt.Errorf("Invalid keyword %v", tokens)
	}
	res.AppendChild(mayCloseBracket)

	return res, rest, nil
}

func (p *Parser) parseDoStatement(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewDoStatementNode()

	// do keyword
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayDoKeyword, ok := next.(KeywordToken)
	if !ok || mayDoKeyword.String() != "do" {
		return nil, nil, fmt.Errorf("Invalid keyword %v", tokens)
	}

	// subroutine call
	sc, rest, err := p.parseSubroutineCall(rest)
	if err != nil {
		return nil, nil, err
	}

	// semicolon
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	maySemicolon, ok := next.(SymbolToken)
	if !ok || maySemicolon.String() != ";" {
		return nil, nil, fmt.Errorf("Invalid symbol %v", next)
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
		return nil, nil, err
	}
	mayReturn, ok := next.(KeywordToken)
	if !ok || mayReturn.String() != "return" {
		return nil, nil, fmt.Errorf("Invalid keyword %v", tokens)
	}

	// expression
	ex, exRest, exErr := p.parseExpression(rest)
	if exErr == nil {
		rest = exRest
	}

	// `;` symbol
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	maySemicolon, ok := next.(SymbolToken)
	if !ok || maySemicolon.String() != ";" {
		return nil, nil, fmt.Errorf("Invalid symbol %v", next)
	}

	res.AppendChild(mayReturn)
	if exErr == nil {
		res.AppendChild(ex)
	}
	res.AppendChild(maySemicolon)

	return res, rest, nil
}

func (p *Parser) parseExpression(tokens TokenList) (TreeNode, TokenList, error) {
	tm, rest, err := p.parseTerm(tokens)
	if err != nil {
		return nil, nil, err
	}
	res := NewExpressionNode()
	res.AppendChild(tm)
	return res, rest, nil
}

func (p *Parser) parseTerm(tokens TokenList) (TreeNode, TokenList, error) {
	res := NewTermNode()

	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, err
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

	// VarName
	if vn, rest, err := p.parseVarName(tokens); err == nil {
		res.AppendChild(vn)

		// VarName only
		next, err := rest.LookAt(0)
		if err != nil {
			return nil, nil, err
		}
		mayOpenSqBracket, ok := next.(SymbolToken)
		if !ok || mayOpenSqBracket.String() != "[" {
			return res, rest, nil
		}

		// Array index
		_, rest, err := rest.PopNext()
		if err != nil {
			return nil, nil, err
		}
		ex, rest, err := p.parseExpression(rest)
		if err != nil {
			return nil, nil, err
		}
		next, rest, err = rest.PopNext()
		if err != nil {
			return nil, nil, err
		}
		mayCloseSqBracket, ok := next.(SymbolToken)
		if !ok || mayCloseSqBracket.String() != "]" {
			return nil, nil, fmt.Errorf("Invalid symbol %v", rest)
		}
		res.AppendChild(mayOpenSqBracket)
		res.AppendChild(ex)
		res.AppendChild(mayCloseSqBracket)
		return res, rest, nil
	}

	// todo: subroutineCall
	if sc, rest, err := p.parseSubroutineCall(tokens); err == nil {
		res.AppendChild(sc)
		return res, rest, nil
	}

	// expression enclosed in paren
	if mayOpenParen, ok := next.(SymbolToken); ok && mayOpenParen.String() == "(" {
		ex, rest, err := p.parseExpression(rest)
		if err != nil {
			return nil, nil, err
		}
		mayCloseParen, ok := rest[0].(SymbolToken)
		if !ok || mayCloseParen.String() != ")" {
			return nil, nil, fmt.Errorf("Invalid symbol %v", rest)
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
			return nil, nil, err
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
		return nil, nil, err
	}
	del, ok := next.(SymbolToken)
	if !ok || (del.String() != "(" && del.String() != ".") {
		return nil, nil, fmt.Errorf("Invalid symbol %v", tokens)
	}

	// method call
	rest := tokens
	if del.String() == "." {
		cn, innerRest, err := p.parseClassName(rest) // TODO: how to detect className or varName
		if err != nil {
			return nil, nil, err
		}
		next, innerRest2, err := tokens.PopNext()
		if err != nil {
			return nil, nil, err
		}
		mayDot, ok := next.(SymbolToken)
		if !ok || mayDot != "." {
			return nil, nil, fmt.Errorf("Invalid symbol %v", innerRest)
		}
		rest = innerRest2
		res.AppendChild(cn)
		res.AppendChild(mayDot)
	}

	// function call
	sn, rest, err := p.parseSubroutineName(rest)
	if err != nil {
		return nil, nil, err
	}
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayOpenParen, ok := next.(SymbolToken)
	if !ok || mayOpenParen != "(" {
		return nil, nil, fmt.Errorf("Invalid symbol %v", rest)
	}
	el, rest, err := p.parseExpressionList(rest)
	if err != nil {
		return nil, nil, err
	}
	next, rest, err = rest.PopNext()
	if err != nil {
		return nil, nil, err
	}
	mayCloseParen, ok := next.(SymbolToken)
	if !ok || mayCloseParen != ")" {
		return nil, nil, fmt.Errorf("Invalid symbol %v", rest)
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
				return nil, nil, err
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
		return nil, nil, err
	}
	cur, ok := next.(SymbolToken)
	if !ok {
		return nil, nil, fmt.Errorf("Type mismatch %v", tokens)
	}
	ok = false
	switch cur.String() {
	case "+", "-", "*", "/", "&", "|", "<", ">", "=":
		ok = true
	}
	if !ok {
		return nil, nil, fmt.Errorf("Invalid keyword %v", tokens)
	}
	on := NewOpNode()
	on.AppendChild(cur)
	return on, rest, nil
}

func (p *Parser) parseUnaryOp(tokens TokenList) (TreeNode, TokenList, error) {
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, err
	}
	cur, ok := next.(SymbolToken)
	if !ok {
		return nil, nil, fmt.Errorf("Type mismatch %v", tokens)
	}
	if cur.String() != "-" && cur.String() != "~" {
		return nil, nil, fmt.Errorf("Invalid keyword %v", tokens)
	}
	uon := NewUnaryOpNode()
	uon.AppendChild(cur)
	return uon, rest, nil
}

func (p *Parser) parseKeywordConstant(tokens TokenList) (TreeNode, TokenList, error) {
	next, rest, err := tokens.PopNext()
	if err != nil {
		return nil, nil, err
	}
	cur, ok := next.(KeywordToken)
	if !ok {
		return nil, nil, fmt.Errorf("Type mismatch %v", tokens)
	}
	if cur.String() != "true" && cur.String() != "false" && cur.String() != "null" && cur.String() != "this" {
		return nil, nil, fmt.Errorf("Invalid keyword %v", tokens)
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
