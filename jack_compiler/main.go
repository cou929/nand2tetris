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

func main() {
	tokenize := false
	toStdout := false
	flag.BoolVar(&tokenize, "tokenize", false, "output tokenized result as xml")
	flag.BoolVar(&toStdout, "toStdout", false, "output result to stdout instead of file")
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

			e := filepath.Ext(f)
			n := fmt.Sprintf("%sT.out.xml", f[0:len(f)-len(e)])
			w, err := os.Create(n)
			if err != nil {
				log.Fatal(err, f)
			}
			_, err = w.WriteString(xml)
			if err != nil {
				log.Fatal(err, f)
			}

			continue
		}

		parser := NewParser()
		tree, err := parser.Parse(tokens)
		if err != nil {
			log.Fatal(err, f)
		}

		xml := tree.Xml()
		if toStdout {
			fmt.Println(xml)
			continue
		}
		e := filepath.Ext(f)
		n := fmt.Sprintf("%s.out.xml", f[0:len(f)-len(e)])
		w, err := os.Create(n)
		if err != nil {
			log.Fatal(err, f)
		}
		_, err = w.WriteString(xml)
		if err != nil {
			log.Fatal(err, f)
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
