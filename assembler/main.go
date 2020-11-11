package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	parser := NewParser(reader)
	commands, _ := parser.Parse()
	for _, c := range commands {
		b, err := NewBinaryCode(c)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%016b\n", b.Line)
	}
}
