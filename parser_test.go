package gtl

import "testing"

func TestParse(t *testing.T) {
	var tokens []Token
	var ast *AST
	var err error
	var program *Node

	assertValidAST := func(ast *AST) {
		if err != nil {
			t.Fatal(err)
		}
		if ast == nil {
			t.Fatal("ast should not be nil")
		}
		if ast.child == nil {
			t.Fatal("ast should have child")
		}
		if want, got := PROGRAM, ast.child.nodeType; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
	}

	tokens = []Token{
		{Word, "true"},
		{EOF, ""},
	}
	ast, err = Parse(tokens)

	assertValidAST(ast)
	program = ast.child
	if want, got := 1, len(program.children); got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := TRUE, program.children[0].nodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}

	tokens = []Token{
		{Word, "false"},
		{EOF, ""},
	}
	ast, err = Parse(tokens)

	assertValidAST(ast)
	program = ast.child
	if want, got := 1, len(program.children); got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := FALSE, program.children[0].nodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
}
