package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	parser := NewParser(reader)
	commands, err := parser.Parse()
	if err != nil {
		panic(err)
	}
	for _, c := range commands {
		asm, err := NewAsmCode(c)
		if err != nil {
			panic(err)
		}
		fmt.Println(strings.Join(asm.Line, "\n"))
	}
}
