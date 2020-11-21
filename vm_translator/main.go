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
	for i, c := range commands {
		asm, err := NewAsmCode("stdin", i, c)
		if err != nil {
			panic(err)
		}
		if asm == nil {
			continue
		}
		fmt.Println(strings.Join(asm.Code(), "\n"))
	}
}
