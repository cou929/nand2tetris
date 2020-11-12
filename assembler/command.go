package main

import (
	"fmt"
	"regexp"
	"strconv"
)

type CommandType int

const (
	ACommand CommandType = iota + 1
	CCommand
	LCommand
)

type CommandSymbol string

type CommandAddressValue int

type CommandDest string

const (
	DestM   CommandDest = "M"
	DestD               = "D"
	DestMD              = "MD"
	DestA               = "A"
	DestAM              = "AM"
	DestAD              = "AD"
	DestAMD             = "AMD"
)

type CommandComp string

const (
	CompZero    CommandComp = "0"
	Comp1                   = "1"
	CompNeg1                = "-1"
	CompD                   = "D"
	CompA                   = "A"
	CompNotD                = "!D"
	CompNotA                = "!A"
	CompNegD                = "-D"
	CompNegA                = "-A"
	CompDPlus1              = "D+1"
	CompAPlus1              = "A+1"
	CompDMinus1             = "D-1"
	CompAMinus1             = "A-1"
	CompDPlusA              = "D+A"
	CompDMinusA             = "D-A"
	CompAMinusD             = "A-D"
	CompDAndA               = "D&A"
	CompDOrA                = "D|A"
	CompM                   = "M"
	CompNotM                = "!M"
	CompNegM                = "-M"
	CompMPlus1              = "M+1"
	CompMMinus1             = "M-1"
	CompDPlusM              = "D+M"
	CompDMinusM             = "D-M"
	CompMMinusD             = "M-D"
	CompDAndM               = "D&M"
	CompDorM                = "D|M"
)

type CommandJump string

const (
	JGT CommandJump = "JGT"
	JEQ             = "JEQ"
	JGE             = "JGE"
	JLT             = "JLT"
	JNE             = "JNE"
	JLE             = "JLE"
	JMP             = "JMP"
)

type Command struct {
	Type         CommandType
	Symbol       CommandSymbol
	AddressValue CommandAddressValue
	Dest         CommandDest
	Comp         CommandComp
	Jump         CommandJump
}

func NewCommand(typ CommandType) *Command {
	return &Command{Type: typ}
}

func (c *Command) SetSymbolOrValue(in string) error {
	digitsOnly := regexp.MustCompile(`^[0-9]+$`)
	label := regexp.MustCompile(`^[a-zA-Z_.$:][a-zA-Z0-9_.$:]+$`)
	if digitsOnly.MatchString(in) {
		i, err := strconv.Atoi(in)
		if err != nil {
			return fmt.Errorf("Invalid value %s", in)
		}
		c.SetAddressValue(i)
		c.Symbol = ""

		return nil
	}

	if !label.MatchString(in) {
		return fmt.Errorf("Invalid format symbol %s", in)
	}

	c.Symbol = CommandSymbol(in)
	return nil
}

func (c *Command) SetAddressValue(in int) error {
	c.AddressValue = CommandAddressValue(in)
	return nil
}

func (c *Command) SetDest(in string) error {
	switch CommandDest(in) {
	case DestM, DestD, DestMD, DestA, DestAM, DestAD, DestAMD:
		c.Dest = CommandDest(in)
		return nil
	}
	return fmt.Errorf("Invalid format dest %s", in)
}

func (c *Command) SetComp(in string) error {
	switch CommandComp(in) {
	case CompZero, Comp1, CompNeg1, CompD,
		CompA, CompNotD, CompNotA, CompNegD,
		CompNegA, CompDPlus1, CompAPlus1, CompDMinus1,
		CompAMinus1, CompDPlusA, CompDMinusA, CompAMinusD,
		CompDAndA, CompDOrA, CompM, CompNotM, CompNegM,
		CompMPlus1, CompMMinus1, CompDPlusM, CompDMinusM,
		CompMMinusD, CompDAndM, CompDorM:
		c.Comp = CommandComp(in)
		return nil
	}
	return fmt.Errorf("Invalid format comp %s", in)
}

func (c *Command) SetJump(in string) error {
	switch CommandJump(in) {
	case JGT, JEQ, JGE, JLT, JNE, JLE, JMP:
		c.Jump = CommandJump(in)
		return nil
	}
	return fmt.Errorf("Invalid format jump %s", in)
}
