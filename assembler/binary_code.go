package main

import (
	"fmt"
	"strconv"
)

type BinaryCode struct {
	Dest int
	Comp int
	Jump int
	Line int
}

var compToBin = map[CommandComp]int{
	"0":   0b0_101010,
	"1":   0b0_111111,
	"-1":  0b0_111010,
	"D":   0b0_001100,
	"A":   0b0_110000,
	"!D":  0b0_001101,
	"!A":  0b0_110001,
	"-D":  0b0_001111,
	"-A":  0b0_110011,
	"D+1": 0b0_011111,
	"A+1": 0b0_110111,
	"D-1": 0b0_001110,
	"A-1": 0b0_110010,
	"D+A": 0b0_000010,
	"D-A": 0b0_010011,
	"A-D": 0b0_000111,
	"D&A": 0b0_000000,
	"D|A": 0b0_010101,
	"M":   0b1_110000,
	"!M":  0b1_110001,
	"-M":  0b1_110011,
	"M+1": 0b1_110111,
	"M-1": 0b1_110010,
	"D+M": 0b1_000010,
	"D-M": 0b1_010011,
	"M-D": 0b1_000111,
	"D&M": 0b1_000000,
	"D|M": 0b1_010101,
}

var destToBin = map[CommandDest]int{
	"":    0b000,
	"M":   0b001,
	"D":   0b010,
	"MD":  0b011,
	"A":   0b100,
	"AM":  0b101,
	"AD":  0b110,
	"AMD": 0b111,
}

var jumpToBin = map[CommandJump]int{
	"":    0b000,
	"JGT": 0b001,
	"JEQ": 0b010,
	"JGE": 0b011,
	"JLT": 0b100,
	"JNE": 0b101,
	"JLE": 0b110,
	"JMP": 0b111,
}

func NewBinaryCode(command *Command) (*BinaryCode, error) {
	if command.Type == ACommand {
		i, err := strconv.Atoi(string(command.Symbol))
		if err != nil {
			return nil, fmt.Errorf("Invalid symbol %s, %w", command.Symbol, err)
		}
		return &BinaryCode{
			Line: i,
		}, nil
	}

	// C Command
	b := &BinaryCode{}
	b.Comp = compToBin[command.Comp]
	b.Dest = destToBin[command.Dest]
	b.Jump = jumpToBin[command.Jump]
	b.Line = (0b111 << 13) | (b.Comp << 6) | (b.Dest << 3) | (b.Jump)

	return b, nil
}
