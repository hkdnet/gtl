package gtl

import "testing"

func TestParse(t *testing.T) {
	var tokens []*Token
	var ast *AST
	var err error
	var program *Node

	assertValidAST := func(ast *AST) {
		if ast == nil {
			t.Fatal("ast should not be nil")
		}
		if ast.Child == nil {
			t.Fatal("ast should have child")
		}
		if want, got := Program, ast.Child.NodeType; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
	}

	// case: values
	tokens = []*Token{
		{Keyword, "true"},
		{EOF, ""},
	}
	ast, err = Parse(tokens)
	if err != nil {
		t.Fatal(err)
	}

	assertValidAST(ast)
	program = ast.Child
	if want, got := 1, len(program.Children); got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := True, program.Children[0].NodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}

	tokens = []*Token{
		{Keyword, "false"},
		{EOF, ""},
	}
	ast, err = Parse(tokens)
	if err != nil {
		t.Fatal(err)
	}

	assertValidAST(ast)
	program = ast.Child
	if want, got := 1, len(program.Children); got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := False, program.Children[0].NodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}

	tokens = []*Token{
		{Number, "0"},
		{EOF, ""},
	}
	ast, err = Parse(tokens)
	if err != nil {
		t.Fatal(err)
	}

	assertValidAST(ast)
	program = ast.Child
	if want, got := 1, len(program.Children); got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := Zero, program.Children[0].NodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}

	// case if
	tokens = []*Token{
		{Keyword, "if"},
		{Keyword, "true"},
		{Keyword, "then"},
		{Keyword, "true"},
		{Keyword, "else"},
		{Keyword, "false"},
		{EOF, ""},
	}
	ast, err = Parse(tokens)
	if err != nil {
		t.Fatal(err)
	}

	assertValidAST(ast)
	program = ast.Child
	if want, got := 1, len(program.Children); got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := IF, program.Children[0].NodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := True, program.Children[0].Children[0].NodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}

	// case lambda
	tokens = []*Token{
		{Dot, "."},
		{Word, "a"},
		{Arrow, "->"},
		{Word, "a"},
		{EOF, ""},
	}
	ast, err = Parse(tokens)
	if err != nil {
		t.Fatal(err)
	}

	assertValidAST(ast)
	program = ast.Child
	if want, got := 1, len(program.Children); got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := Lambda, program.Children[0].NodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	{
		def := program.Children[0].Children[0]
		body := program.Children[0].Children[1]
		if want, got := LambdaDef, def.NodeType; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
		if want, got := LambdaBody, body.NodeType; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}

		if want, got := 1, len(def.Children); got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
		if want, got := LambdaParam, def.Children[0].NodeType; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}

		if want, got := 1, len(body.Children); got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
		if want, got := Variable, body.Children[0].NodeType; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
	}
}

func Test_parseDot(t *testing.T) {
	var env parseEnvironemnt
	tokens := []*Token{
		{Dot, "."},
		{Word, "a"},
		{Arrow, "->"},
		{Word, "a"},
		{EOF, ""},
	}
	node, env, err := parseDot(tokens, env)
	if err != nil {
		t.Fatal(err)
	}
	if want, got := 4, env.idx; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := Lambda, node.NodeType; got != got {
		t.Errorf("want %v but got %v\n", want, got)
	}
}
