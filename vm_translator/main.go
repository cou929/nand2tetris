package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("input dir required")
	}
	dirPath := os.Args[1]

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	var codes []string
	codes = BootstrapLine()

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if path.Ext(f.Name()) != ".vm" {
			continue
		}

		p := path.Join(dirPath, f.Name())
		reader, err := os.Open(p)
		if err != nil {
			log.Fatal(err)
		}
		parser := NewParser(reader, f.Name())
		commands, err := parser.Parse()
		if err != nil {
			log.Fatal(err)
		}
		for _, c := range commands {
			asm, err := NewAsmCode(c)
			if err != nil {
				log.Fatal(err)
			}
			if asm == nil {
				continue
			}
			codes = append(codes, asm.Code()...)
		}
	}

	fmt.Println(strings.Join(codes, "\n"))
}
