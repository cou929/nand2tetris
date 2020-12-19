package main

import "fmt"

type SymbolTable struct {
	classScopeTable ScopedTable
	funcScopeTable  ScopedTable
	index           map[VarKind]int
}

type ScopedTable map[string]*SymbolTableEntry

type SymbolTableEntry struct {
	Name  string
	Typ   string
	Kind  VarKind
	Index int
}

type VarKind int

const (
	Static VarKind = iota + 1
	Field
	Argument
	Var
)

func (k VarKind) String() string {
	switch k {
	case Static:
		return "Static"
	case Field:
		return "Field"
	case Argument:
		return "Argument"
	case Var:
		return "Var"
	}
	return "Invalid"
}

func NewVarKind(in string) (VarKind, error) {
	switch in {
	case "static":
		return Static, nil
	case "field":
		return Field, nil
	case "argument":
		return Argument, nil
	case "var":
		return Var, nil
	}
	return 0, fmt.Errorf("Invalid kind name of variable %v", in)
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		classScopeTable: make(ScopedTable),
		funcScopeTable:  make(ScopedTable),
		index:           make(map[VarKind]int),
	}
}

func (s *SymbolTable) Clear() {
	s.classScopeTable = make(ScopedTable)
	s.funcScopeTable = make(ScopedTable)
	s.index = make(map[VarKind]int)
}

func (s *SymbolTable) ClearFuncTable() {
	s.funcScopeTable = make(ScopedTable)
	s.index[Argument] = 0
	s.index[Var] = 0
}

func (s *SymbolTable) Define(name string, typ string, kind VarKind) error {
	ste := &SymbolTableEntry{
		Name:  name,
		Typ:   typ,
		Kind:  kind,
		Index: s.index[kind],
	}

	tt := s.classScopeTable
	if kind == Argument || kind == Var {
		tt = s.funcScopeTable
	}

	if _, ok := tt[name]; ok {
		return fmt.Errorf("Duplicate symbol %v %v %v", name, typ, kind)
	}

	tt[name] = ste
	s.index[kind]++
	return nil
}

func (s *SymbolTable) LookUp(name string) *SymbolTableEntry {
	ste, ok := s.funcScopeTable[name]
	if ok {
		return ste
	}
	ste, ok = s.classScopeTable[name]
	if ok {
		return ste
	}
	return nil
}
