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
		if want, got := Program, ast.child.nodeType; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
	}

	// case: values
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
	if want, got := True, program.children[0].nodeType; got != want {
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
	if want, got := False, program.children[0].nodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}

	tokens = []Token{
		{Number, "0"},
		{EOF, ""},
	}
	ast, err = Parse(tokens)

	assertValidAST(ast)
	program = ast.child
	if want, got := 1, len(program.children); got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := Zero, program.children[0].nodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}

	// case if
	tokens = []Token{
		{Word, "if"},
		{Word, "true"},
		{Word, "then"},
		{Word, "true"},
		{Word, "else"},
		{Word, "false"},
		{EOF, ""},
	}
	ast, err = Parse(tokens)

	assertValidAST(ast)
	program = ast.child
	if want, got := 1, len(program.children); got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := IF, program.children[0].nodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := True, program.children[0].children[0].nodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}

	// case lambda
	tokens = []Token{
		{Word, "a"},
		{Dot, "."},
		{Arrow, "->"},
		{Word, "a"},
		{EOF, ""},
	}
	ast, err = Parse(tokens)

	assertValidAST(ast)
	program = ast.child
	if want, got := 1, len(program.children); got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := Lambda, program.children[0].nodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	{
		def := program.children[0].children[0]
		body := program.children[0].children[1]
		if want, got := LambdaDef, def.nodeType; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
		if want, got := LambdaBody, body.nodeType; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}

		if want, got := 1, len(def.children); got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
		if want, got := LambdaParam, def.children[0].nodeType; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}

		if want, got := 1, len(body.children); got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
		if want, got := Variable, body.children[0].nodeType; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
	}
}
