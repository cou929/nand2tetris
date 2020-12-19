package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

var (
	tokenize  = false
	toStdout  = false
	idAttr    = false
	parseTree = false
)

func main() {
	flag.BoolVar(&tokenize, "tokenize", false, "output tokenized result as xml")
	flag.BoolVar(&toStdout, "toStdout", false, "output result to stdout instead of file")
	flag.BoolVar(&idAttr, "idAttr", false, "output attributes of identifier node")
	flag.BoolVar(&parseTree, "parseTree", false, "output parse tree as xml format")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("input dir required")
	}
	dirPath := flag.Args()[0]
	files, err := findJackFiles(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		reader, err := os.Open(f)
		if err != nil {
			log.Fatal(err, f)
		}
		defer reader.Close()

		tokenizer := NewTokenizer(reader)
		tokens, err := tokenizer.Tokenize()
		if err != nil {
			log.Fatal(err, f)
		}

		if tokenize {
			xml := tokens.Xml()
			if toStdout {
				fmt.Println(xml)
				continue
			}
			if err := write(f, xml, "T.out.xml"); err != nil {
				log.Fatal(err)
			}
			continue
		}

		parser := NewParser()
		tree, err := parser.Parse(tokens)
		if err != nil {
			log.Fatal(err, f)
		}

		if parseTree {
			xml := tree.Xml()
			if toStdout {
				fmt.Println(xml)
				continue
			}
			if err := write(f, xml, ".out.xml"); err != nil {
				log.Fatal(err)
			}
			continue
		}

		compiler := NewCompiler()
		vmCode, err := compiler.Compile(tree)
		if err != nil {
			log.Fatal(err, f)
		}

		if toStdout {
			fmt.Println(vmCode)
			continue
		}
		if err := write(f, vmCode, ".vm"); err != nil {
			log.Fatal(err)
		}
	}
}

func findJackFiles(dirPath string) ([]string, error) {
	const suf = ".jack"
	var res []string

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if path.Ext(f.Name()) != suf {
			continue
		}

		p := path.Join(dirPath, f.Name())
		res = append(res, p)
	}

	return res, nil
}

func write(path string, c string, suffix string) error {
	e := filepath.Ext(path)
	n := fmt.Sprintf("%s%s", path[0:len(path)-len(e)], suffix)
	w, err := os.Create(n)
	if err != nil {
		return fmt.Errorf("Failed to create file %s %w", n, err)
	}
	_, err = w.WriteString(c)
	if err != nil {
		return fmt.Errorf("Failed to write to file %s %w", n, err)
	}
	return nil
}
