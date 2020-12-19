package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Compiler struct {
	curClassInfo    *classInfo
	curFuncInfo     *funcInfo
	callingArgCount int
	ifCounter       int
	whileCounter    int
	vmc             *VmCode
}

type classInfo struct {
	name       string
	fieldCount int
}

type funcInfo struct {
	name          string
	kind          funcKind
	localVarCount int
}

type funcKind int

const (
	Constructor funcKind = iota + 1
	Function
	Method
)

func (f funcKind) String() string {
	switch f {
	case Constructor:
		return "Constructor"
	case Function:
		return "Function"
	case Method:
		return "Method"
	}
	return "Invalid funcKind"
}

func NewCompiler() *Compiler {
	return &Compiler{
		vmc: NewVmCode(),
	}
}

func (c *Compiler) Compile(pt TreeNode) (string, error) {
	codes, err := c.compile(pt)
	if err != nil {
		return "", err
	}
	return strings.Join(codes, "\n"), nil
}

func (c *Compiler) compile(pt TreeNode) ([]string, error) {
	switch pt.Type() {
	case ClassType:
		return c.compileClass(pt)
	case ClassVarDecType:
		return c.compileClassVarDec(pt)
	case TypeType:
		return c.compileType(pt)
	case SubroutineDecType:
		return c.compileSubroutineDec(pt)
	case ParameterListType:
		return c.compileParameterList(pt)
	case SubroutineBodyType:
		return c.compileSubroutineBody(pt)
	case VarDecType:
		return c.compileVarDec(pt)
	case StatementsType:
		return c.compileStatements(pt)
	case StatementType:
		return c.compileStatement(pt)
	case LetStatementType:
		return c.compileLetStatement(pt)
	case IfStatementType:
		return c.compileIfStatement(pt)
	case WhileStatementType:
		return c.compileWhileStatement(pt)
	case DoStatementType:
		return c.compileDoStatement(pt)
	case ReturnStatementType:
		return c.compileReturnStatement(pt)
	case ExpressionType:
		return c.compileExpression(pt)
	case SubroutineCallType:
		return c.compileSubroutineCall(pt)
	case ExpressionListType:
		return c.compileExpressionList(pt)
	case OpType:
		return c.compileOp(pt)
	case UnaryOpType:
		return c.compileUnaryOp(pt)
	}
	return []string{}, fmt.Errorf("Not supported %v", pt.Type())
}

func (c *Compiler) resetAllState() {
	c.resetClassState()
	c.resetFuncState()
	c.resetCallingArgCount()
	c.resetIfCounter()
	c.resetWhileCounter()
}

func (c *Compiler) resetClassState() {
	c.curClassInfo = nil
}

func (c *Compiler) setClassName(in string) {
	if c.curClassInfo == nil {
		c.curClassInfo = &classInfo{}
	}
	c.curClassInfo.name = in
}

func (c *Compiler) incClassFieldCount() {
	if c.curClassInfo == nil {
		c.curClassInfo = &classInfo{}
	}
	c.curClassInfo.fieldCount++
}

func (c *Compiler) resetFuncState() {
	c.curFuncInfo = nil
}

func (c *Compiler) setFuncName(in string) {
	if c.curFuncInfo == nil {
		c.curFuncInfo = &funcInfo{}
	}
	c.curFuncInfo.name = in
}

func (c *Compiler) setFuncKind(in string) error {
	if c.curFuncInfo == nil {
		c.curFuncInfo = &funcInfo{}
	}
	switch in {
	case "constructor":
		c.curFuncInfo.kind = Constructor
	case "function":
		c.curFuncInfo.kind = Function
	case "method":
		c.curFuncInfo.kind = Method
	}
	return fmt.Errorf("Invalid function kind %s", in)
}

func (c *Compiler) incLocalVarCount() {
	if c.curFuncInfo == nil {
		c.curFuncInfo = &funcInfo{}
	}
	c.curFuncInfo.localVarCount++
}

func (c *Compiler) resetCallingArgCount() {
	c.callingArgCount = 0
}

func (c *Compiler) incCallingArgCount() {
	c.callingArgCount++
}

func (c *Compiler) resetIfCounter() {
	c.ifCounter = 0
}

func (c *Compiler) incIfCounter() {
	c.ifCounter++
}

func (c *Compiler) resetWhileCounter() {
	c.whileCounter = 0
}

func (c *Compiler) incWhileCounter() {
	c.whileCounter++
}

func (c *Compiler) compileClass(pt TreeNode) ([]string, error) {
	c.resetAllState()

	var res []string
	for _, node := range pt.ChildNodes() {
		var codes []string
		var err error

		switch node.Type() {
		case ClassNameType:
			c.setClassName(node.Value())
			continue
		case ClassVarDecType:
			codes, err = c.compile(node)
		case SubroutineDecType:
			codes, err = c.compile(node)
		}

		if err != nil {
			return nil, fmt.Errorf("[compileClass] %w", err)
		}
		res = append(res, codes...)
	}
	return res, nil
}

func (c *Compiler) compileClassVarDec(pt TreeNode) ([]string, error) {
	for _, node := range pt.ChildNodes() {
		if node.Type() == VarNameType && node.Meta().Category == IdCatField {
			c.incClassFieldCount()
		}
	}
	return nil, nil
}

func (c *Compiler) compileType(pt TreeNode) ([]string, error) {
	return nil, nil
}

func (c *Compiler) compileSubroutineDec(pt TreeNode) ([]string, error) {
	c.resetFuncState()

	var res []string
	for i, node := range pt.ChildNodes() {
		var codes []string
		var err error

		if i == 0 {
			c.setFuncKind(node.Value())
			continue
		}

		switch node.Type() {
		case SubroutineNameType:
			c.setFuncName(node.Value())
			continue
		case ParameterListType:
			codes, err = c.compile(node)
		case SubroutineBodyType:
			codes, err = c.compile(node)
		}

		if err != nil {
			return nil, fmt.Errorf("[compileSubroutineDec] %w", err)
		}
		res = append(res, codes...)
	}

	// prepare this segment
	switch c.curFuncInfo.kind {
	case Method:
		// set arg0 (= base address of the instance) to pointer 0 (this)
		codes := []string{
			c.vmc.push("argument", 0),
			c.vmc.pop("pointer", 0),
		}
		res = append(codes, res...) // prepend
	case Constructor:
		// allocate memory of the object
		// object size is number of filed the class has
		codes := []string{
			c.vmc.pushConstant(c.curClassInfo.fieldCount),
			c.vmc.call("Memory", "alloc", 1),
			c.vmc.pop("pointer", 0),
		}
		res = append(codes, res...) // prepend
	}

	// prepend function declaration
	res = append([]string{c.vmc.function(c.curClassInfo.name, c.curFuncInfo.name, c.curFuncInfo.localVarCount)}, res...)

	return res, nil
}

func (c *Compiler) compileParameterList(pt TreeNode) ([]string, error) {
	return []string{}, nil
}

func (c *Compiler) compileSubroutineBody(pt TreeNode) ([]string, error) {
	var res []string
	for _, node := range pt.ChildNodes() {
		var codes []string
		var err error

		switch node.Type() {
		case VarDecType:
			codes, err = c.compile(node)
		case StatementsType:
			codes, err = c.compile(node)
		}

		if err != nil {
			return nil, fmt.Errorf("[compileSubroutineDec] %w", err)
		}
		res = append(res, codes...)
	}
	return res, nil
}

func (c *Compiler) compileVarDec(pt TreeNode) ([]string, error) {
	var res []string
	for _, node := range pt.ChildNodes() {
		if node.Type() == VarNameType && node.Meta() != nil && node.Meta().Category == IdCatVar {
			c.incLocalVarCount()
		}
	}
	return res, nil
}

func (c *Compiler) compileStatements(pt TreeNode) ([]string, error) {
	var res []string
	for _, node := range pt.ChildNodes() {
		codes, err := c.compile(node)
		if err != nil {
			return nil, fmt.Errorf("[compileStatements] %w", err)
		}
		res = append(res, codes...)
	}
	return res, nil
}

func (c *Compiler) compileStatement(pt TreeNode) ([]string, error) {
	var res []string
	for _, node := range pt.ChildNodes() {
		codes, err := c.compile(node)
		if err != nil {
			return nil, fmt.Errorf("[compileStatement] %w", err)
		}
		res = append(res, codes...)
	}
	return res, nil
}

func (c *Compiler) compileLetStatement(pt TreeNode) ([]string, error) {
	var res []string
	varName := pt.ChildNodes()[1]

	// case of array
	if len(pt.ChildNodes()) == 8 {
		switch varName.Meta().Category {
		case IdCatStatic:
			res = append(res, c.vmc.push("static", varName.Meta().SymbolInfo.Index))
		case IdCatField:
			res = append(res, c.vmc.push("this", varName.Meta().SymbolInfo.Index))
		case IdCatArg:
			idx := varName.Meta().SymbolInfo.Index
			if c.curFuncInfo.kind == Method {
				idx++
			}
			res = append(res, c.vmc.push("argument", idx))
		case IdCatVar:
			res = append(res, c.vmc.push("local", varName.Meta().SymbolInfo.Index))
		}

		idxExp := pt.ChildNodes()[3]
		idx, err := c.compile(idxExp)
		if err != nil {
			return nil, fmt.Errorf("[compileStatement] %w", err)
		}
		res = append(res, idx...)
		res = append(res, c.vmc.add())
		res = append(res, c.vmc.pop("pointer", 1))
		exp := pt.ChildNodes()[6]
		codes, err := c.compile(exp)
		if err != nil {
			return nil, fmt.Errorf("[compileLetStatement] %w", err)
		}
		res = append(res, codes...)
		res = append(res, c.vmc.pop("that", 0))
		return res, nil
	}

	expIndex := 3
	exp := pt.ChildNodes()[expIndex]
	codes, err := c.compile(exp)
	if err != nil {
		return nil, fmt.Errorf("[compileLetStatement] %w", err)
	}
	res = append(res, codes...)

	switch varName.Meta().Category {
	case IdCatStatic:
		res = append(res, c.vmc.pop("static", varName.Meta().SymbolInfo.Index))
	case IdCatField:
		res = append(res, c.vmc.pop("this", varName.Meta().SymbolInfo.Index))
	case IdCatArg:
		idx := varName.Meta().SymbolInfo.Index
		if c.curFuncInfo.kind == Method {
			idx++
		}
		res = append(res, c.vmc.pop("argument", idx))
	case IdCatVar:
		res = append(res, c.vmc.pop("local", varName.Meta().SymbolInfo.Index))
	}

	return res, nil
}

func (c *Compiler) compileIfStatement(pt TreeNode) ([]string, error) {
	var res []string
	suffix := c.ifCounter
	c.incIfCounter()
	endLabel := fmt.Sprintf("%s.%s.%d.IF.END", c.curClassInfo.name, c.curFuncInfo.name, suffix)
	elseLabel := fmt.Sprintf("%s.%s.%d.IF.ELSE", c.curClassInfo.name, c.curFuncInfo.name, suffix)

	// ~(cond)
	cond, err := c.compile(pt.ChildNodes()[2])
	if err != nil {
		return nil, fmt.Errorf("[compileIfStatement] %w", err)
	}
	condNot := append(cond, c.vmc.not())
	res = append(res, condNot...)

	// if statement
	res = append(res, c.vmc.ifGoTo(elseLabel))
	ifStatement, err := c.compile(pt.ChildNodes()[5])
	if err != nil {
		return nil, fmt.Errorf("[compileIfStatement] %w", err)
	}
	res = append(res, ifStatement...)
	res = append(res, c.vmc.goTo(endLabel))

	// else statement
	res = append(res, c.vmc.label(elseLabel))
	if len(pt.ChildNodes()) > 10 {
		elseStatement, err := c.compile(pt.ChildNodes()[9])
		if err != nil {
			return nil, fmt.Errorf("[compileIfStatement] %w", err)
		}
		res = append(res, elseStatement...)
	}
	res = append(res, c.vmc.label(endLabel))

	return res, nil
}

func (c *Compiler) compileWhileStatement(pt TreeNode) ([]string, error) {
	var res []string
	suffix := c.whileCounter
	c.incWhileCounter()
	contLabel := fmt.Sprintf("%s.%s.%d.WHILE.CONT", c.curClassInfo.name, c.curFuncInfo.name, suffix)
	endLabel := fmt.Sprintf("%s.%s.%d.WHILE.END", c.curClassInfo.name, c.curFuncInfo.name, suffix)

	res = append(res, c.vmc.label(contLabel))

	// ~(cond)
	cond, err := c.compile(pt.ChildNodes()[2])
	if err != nil {
		return nil, fmt.Errorf("[compileWhileStatement] %w", err)
	}
	condNot := append(cond, c.vmc.not())
	res = append(res, condNot...)
	res = append(res, c.vmc.ifGoTo(endLabel))

	// statement
	stmt, err := c.compile(pt.ChildNodes()[5])
	if err != nil {
		return nil, fmt.Errorf("[compileWhileStatement] %w", err)
	}
	res = append(res, stmt...)
	res = append(res, c.vmc.goTo(contLabel))
	res = append(res, c.vmc.label(endLabel))

	return res, nil
}

func (c *Compiler) compileDoStatement(pt TreeNode) ([]string, error) {
	var res []string
	for _, node := range pt.ChildNodes() {
		if node.Type() != SubroutineCallType {
			continue
		}
		codes, err := c.compile(node)
		if err != nil {
			return nil, fmt.Errorf("[compileDoStatement] %w", err)
		}
		res = append(res, codes...)
	}
	// discard return value of void function
	res = append(res, c.vmc.pop("temp", 7)) // todo: what happen if anyone uses temp 7 ?
	return res, nil
}

func (c *Compiler) compileReturnStatement(pt TreeNode) ([]string, error) {
	var res []string
	if len(pt.ChildNodes()) == 2 {
		// return;
		res = append(res, c.vmc.pushConstant(0))
	} else {
		// return expression;
		exp, err := c.compile(pt.ChildNodes()[1])
		if err != nil {
			return nil, fmt.Errorf("[compileReturnStatement] %w", err)
		}
		res = append(res, exp...)
	}
	res = append(res, c.vmc.ret())
	return res, nil
}

func (c *Compiler) compileExpression(pt TreeNode) ([]string, error) {
	return c.traverseExpression(pt.ChildNodes())
}

func (c *Compiler) traverseExpression(exps []TreeNode) ([]string, error) {
	var res []string

	if len(exps) == 1 {
		term := exps[0]
		child := term.ChildNodes()[0]
		switch child.Type() {
		case IntConstType:
			i, err := strconv.Atoi(child.Value())
			if err != nil {
				return nil, fmt.Errorf("[compileExpression] %w", err)
			}
			res = append(res, c.vmc.pushConstant(i))
		case VarNameType:
			switch child.Meta().Category {
			case IdCatStatic:
				res = append(res, c.vmc.push("static", child.Meta().SymbolInfo.Index))
			case IdCatField:
				res = append(res, c.vmc.push("this", child.Meta().SymbolInfo.Index))
			case IdCatArg:
				idx := child.Meta().SymbolInfo.Index
				if c.curFuncInfo.kind == Method {
					idx++
				}
				res = append(res, c.vmc.push("argument", idx))
			case IdCatVar:
				res = append(res, c.vmc.push("local", child.Meta().SymbolInfo.Index))
			}
			if len(term.ChildNodes()) == 4 {
				// case of array
				idxExp := term.ChildNodes()[2]
				idx, err := c.compile(idxExp)
				if err != nil {
					return nil, fmt.Errorf("[compileExpression] %w", err)
				}
				res = append(res, idx...)
				res = append(res, c.vmc.add())
				res = append(res, c.vmc.pop("pointer", 1))
				res = append(res, c.vmc.push("that", 0))
			}
		case SubroutineCallType:
			codes, err := c.compile(child)
			if err != nil {
				return nil, fmt.Errorf("[compileExpression] %w", err)
			}
			res = append(res, codes...)
		case SymbolType: // (expression)
			exp := term.ChildNodes()[1]
			codes, err := c.compile(exp)
			if err != nil {
				return nil, fmt.Errorf("[compileExpression] %w", err)
			}
			res = append(res, codes...)
		case UnaryOpType:
			op := term.ChildNodes()[0]
			term := term.ChildNodes()[1]
			codes, err := c.traverseExpression([]TreeNode{term})
			if err != nil {
				return nil, fmt.Errorf("[compileExpression] %w", err)
			}
			res = append(res, codes...)
			codes, err = c.compile(op)
			if err != nil {
				return nil, fmt.Errorf("[compileExpression] %w", err)
			}
			res = append(res, codes...)
		case KeywordConstantType:
			switch child.Value() {
			case "true":
				res = append(res, c.vmc.true()...)
			case "false":
				res = append(res, c.vmc.false())
			case "null":
				res = append(res, c.vmc.null())
			case "this":
				res = append(res, c.vmc.push("pointer", 0))
			}
		case StrConstType:
			res = append(res, c.vmc.newStr(child.Value())...)
		default:
			return nil, fmt.Errorf("[compileExpression] Invalid node %v, %v", term, child.Type())
		}
	}

	if len(exps) >= 3 {
		term := exps[0]
		op := exps[1]
		rest := exps[2:len(exps)]

		codes, err := c.traverseExpression([]TreeNode{term})
		if err != nil {
			return nil, fmt.Errorf("[compileExpression] %w", err)
		}
		res = append(res, codes...)

		codes, err = c.traverseExpression(rest)
		if err != nil {
			return nil, fmt.Errorf("[compileExpression] %w", err)
		}
		res = append(res, codes...)

		codes, err = c.compile(op)
		if err != nil {
			return nil, fmt.Errorf("[compileExpression] %w", err)
		}
		res = append(res, codes...)

		return res, nil
	}

	return res, nil
}

func (c *Compiler) compileSubroutineCall(pt TreeNode) ([]string, error) {
	var res []string
	className := ""
	subName := ""
	c.resetCallingArgCount()

	for i, node := range pt.ChildNodes() {
		if i == 0 {
			switch node.Type() {
			case SubroutineNameType:
				// calling method of current instance
				className = c.curClassInfo.name
				subName = node.Value()
				// pass current this as argument[0] of callee
				res = append(res, c.vmc.push("pointer", 0))
				c.incCallingArgCount()
			case ClassNameType:
				className = node.Value()
			case VarNameType:
				className = node.Meta().SymbolInfo.Type
				// push base address of the instance as argument[0] of callee
				switch node.Meta().Category {
				case IdCatStatic:
					res = append(res, c.vmc.push("static", node.Meta().SymbolInfo.Index))
				case IdCatField:
					res = append(res, c.vmc.push("this", node.Meta().SymbolInfo.Index))
				case IdCatArg:
					idx := node.Meta().SymbolInfo.Index
					if c.curFuncInfo.kind == Method {
						idx++
					}
					res = append(res, c.vmc.push("argument", idx))
				case IdCatVar:
					res = append(res, c.vmc.push("local", node.Meta().SymbolInfo.Index))
				}
				c.incCallingArgCount()
			}
		}

		if i == 2 && node.Type() == SubroutineNameType {
			subName = node.Value()
		}

		if node.Type() == ExpressionListType {
			codes, err := c.compile(node)
			if err != nil {
				return nil, fmt.Errorf("[compileSubroutineCall] %w", err)
			}
			res = append(res, codes...)
		}
	}
	res = append(res, c.vmc.call(className, subName, c.callingArgCount))
	return res, nil
}

func (c *Compiler) compileExpressionList(pt TreeNode) ([]string, error) {
	var res []string
	for _, node := range pt.ChildNodes() {
		if node.Type() != ExpressionType {
			continue
		}
		codes, err := c.compile(node)
		if err != nil {
			return nil, fmt.Errorf("[compileExpressionList] %w", err)
		}
		res = append(res, codes...)
		c.incCallingArgCount()
	}
	return res, nil
}

func (c *Compiler) compileOp(pt TreeNode) ([]string, error) {
	switch pt.Value() {
	case "+":
		return []string{c.vmc.add()}, nil
	case "-":
		return []string{c.vmc.sub()}, nil
	case "*":
		return []string{c.vmc.mul()}, nil
	case "/":
		return []string{c.vmc.div()}, nil
	case "&":
		return []string{c.vmc.and()}, nil
	case "|":
		return []string{c.vmc.or()}, nil
	case "<":
		return []string{c.vmc.lt()}, nil
	case ">":
		return []string{c.vmc.gt()}, nil
	case "=":
		return []string{c.vmc.eq()}, nil
	}
	return nil, fmt.Errorf("[compileOp] Invalid op %v", pt)
}

func (c *Compiler) compileUnaryOp(pt TreeNode) ([]string, error) {
	switch pt.Value() {
	case "-":
		return []string{c.vmc.neg()}, nil
	case "~":
		return []string{c.vmc.not()}, nil
	}
	return nil, fmt.Errorf("[compileOp] Invalid op %v", pt)
}
