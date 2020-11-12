package main

import "fmt"

type SymbolTable struct {
	t               map[CommandSymbol]int
	variableCounter int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		t: map[CommandSymbol]int{
			"SP":     0,
			"LCL":    1,
			"ARG":    2,
			"THIS":   3,
			"THAT":   4,
			"R0":     0,
			"R1":     1,
			"R2":     2,
			"R3":     3,
			"R4":     4,
			"R5":     5,
			"R6":     6,
			"R7":     7,
			"R8":     8,
			"R9":     9,
			"R10":    10,
			"R11":    11,
			"R12":    12,
			"R13":    13,
			"R14":    14,
			"R15":    15,
			"SCREEN": 16384,
			"KBD":    24576,
		},
		variableCounter: 16,
	}
}

func (s *SymbolTable) AddEntry(symbol CommandSymbol, address int) {
	s.t[symbol] = address
}

func (s *SymbolTable) Contains(symbol CommandSymbol) bool {
	_, ok := s.t[symbol]
	return ok
}

func (s *SymbolTable) GetAddress(symbol CommandSymbol) int {
	return s.t[symbol]
}

func (s *SymbolTable) AddVariable(symbol CommandSymbol) (int, error) {
	if s.Contains(symbol) {
		return 0, fmt.Errorf("Variable symbol already stored %s", symbol)
	}
	s.AddEntry(symbol, s.variableCounter)
	s.variableCounter = s.variableCounter + 1
	return s.GetAddress(symbol), nil
}
