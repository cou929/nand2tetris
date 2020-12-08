package main

// NodeType represents both terminal and non-terminal symbols.
// This is also used for tokenized symbols
type NodeType int

const (
	// terminal symbols
	KeywordType NodeType = iota + 1
	SymbolType
	IntConstType
	StrConstType
	IdentifierType

	// non-terminal symbols
	ClassType
	ClassVarDecType
	SubroutineDecType
	ParameterListType
	SubroutineBodyType
	VarDecType
	StatementsType
	WhileStatementType
	IfStatementType
	ReturnStatementType
	LetStatementType
	DoStatementType
	ExpressionType
	TermType
	ExpressionListType

	// internal use
	TypeType
	ClassNameType
	SubroutineNameType
	VarNameType
	StatementType
	SubroutineCallType
	OpType
	UnaryOpType
	KeywordConstantType
)

func (typ NodeType) String() string {
	switch typ {
	case KeywordType:
		return "KeywordType"
	case SymbolType:
		return "SymbolType"
	case IntConstType:
		return "IntConstType"
	case StrConstType:
		return "StrConstType"
	case IdentifierType:
		return "IdentifierType"
	case ClassType:
		return "ClassType"
	case ClassVarDecType:
		return "ClassVarDecType"
	case SubroutineDecType:
		return "SubroutineDecType"
	case ParameterListType:
		return "ParameterListType"
	case SubroutineBodyType:
		return "SubroutineBodyType"
	case VarDecType:
		return "VarDecType"
	case StatementsType:
		return "StatementsType"
	case WhileStatementType:
		return "WhileStatementType"
	case IfStatementType:
		return "IfStatementType"
	case ReturnStatementType:
		return "ReturnStatementType"
	case LetStatementType:
		return "LetStatementType"
	case DoStatementType:
		return "DoStatementType"
	case ExpressionType:
		return "ExpressionType"
	case TermType:
		return "TermType"
	case ExpressionListType:
		return "ExpressionListType"
	case TypeType:
		return "TypeType"
	case ClassNameType:
		return "ClassNameType"
	case SubroutineNameType:
		return "SubroutineNameType"
	case VarNameType:
		return "VarNameType"
	case StatementType:
		return "StatementType"
	case SubroutineCallType:
		return "SubroutineCallType"
	case OpType:
		return "OpType"
	case UnaryOpType:
		return "UnaryOpType"
	case KeywordConstantType:
		return "KeywordConstantType"
	}
	return "invalid type"
}
