package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		panic("input .vm file required")
	}
	path := os.Args[1]
	fileName := filepath.Base(path)
	reader, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	parser := NewParser(reader, fileName)
	commands, err := parser.Parse()
	if err != nil {
		panic(err)
	}
	for i, c := range commands {
		asm, err := NewAsmCode(fileName, i, c)
		if err != nil {
			panic(err)
		}
		if asm == nil {
			continue
		}
		fmt.Println(strings.Join(asm.Code(), "\n"))
	}
}
