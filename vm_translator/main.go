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
		fmt.Println(c)
	}
}
