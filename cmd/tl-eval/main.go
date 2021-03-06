package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hkdnet/gtl"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "USAGE: %s FILENAME", os.Args[0])
	}
	filename := os.Args[1]

	err := run(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	source := string(b)

	l := gtl.NewLexer(source)
	var tokens []*gtl.Token
	for l.HasNext() {
		token, err := l.NextToken()
		if err != nil {
			return err
		}
		tokens = append(tokens, token)
	}
	ast, err := gtl.Parse(tokens)
	if err != nil {
		return err
	}
	result, err := gtl.Eval(ast)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
