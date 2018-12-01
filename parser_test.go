package gtl

import (
	"testing"
)

func TestParse(t *testing.T) {
	var tokens []*Token
	var ast *AST
	var err error

	assertValidAST := func(ast *AST) {
		if ast == nil {
			t.Fatal("ast should not be nil")
		}
		if ast.Child == nil {
			t.Fatal("ast should have child")
		}
	}

	// case: values
	tokens = []*Token{
		{KeywordTrue, "true"},
		{EOF, ""},
	}
	ast, err = Parse(tokens)
	if err != nil {
		t.Fatal(err)
	}

	assertValidAST(ast)
	if want, got := True, ast.Child.NodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}

	tokens = []*Token{
		{KeywordFalse, "false"},
		{EOF, ""},
	}
	ast, err = Parse(tokens)
	if err != nil {
		t.Fatal(err)
	}

	assertValidAST(ast)
	if want, got := False, ast.Child.NodeType; got != want {
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
	if want, got := Zero, ast.Child.NodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}

	// case if
	tokens = []*Token{
		{KeywordIf, "if"},
		{KeywordTrue, "true"},
		{KeywordThen, "then"},
		{KeywordTrue, "true"},
		{KeywordElse, "else"},
		{KeywordFalse, "false"},
		{EOF, ""},
	}
	ast, err = Parse(tokens)
	if err != nil {
		t.Fatal(err)
	}

	assertValidAST(ast)
	if want, got := IF, ast.Child.NodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := True, ast.Child.Children[0].NodeType; got != want {
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
	if want, got := Lambda, ast.Child.NodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	{
		def := ast.Child.Children[0]
		body := ast.Child.Children[1]
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

	tokens = []*Token{
		{LParen, "("},
		{Dot, "."},
		{Word, "a"},
		{Arrow, "->"},
		{Word, "a"},
		{RParen, ")"},
		{EOF, ""},
	}
	ast, err = Parse(tokens)
	if err != nil {
		t.Fatal(err)
	}
	assertValidAST(ast)
}

func Test_parseIf(t *testing.T) {
	var env parseEnvironemnt
	tokens := []*Token{
		{KeywordIf, "if"},
		{Word, "a"},
		{KeywordThen, "then"},
		{Word, "b"},
		{KeywordElse, "else"},
		{Word, "c"},
		{EOF, ""},
	}
	node, _, err := parseIf(tokens, env)
	if err != nil {
		t.Fatal(err)
	}
	if want, got := IF, node.NodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := 3, len(node.Children); got != want {
		t.Errorf("want %v but got %v\n", want, got)
		return
	}
	cond := node.Children[0]
	truePart := node.Children[1]
	falsePart := node.Children[2]
	if want, got := FreeVariable, cond.NodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := "a", cond.Name; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := FreeVariable, truePart.NodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := "b", truePart.Name; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := FreeVariable, falsePart.NodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := "c", falsePart.Name; got != want {
		t.Errorf("want %v but got %v\n", want, got)
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
	if want, got := Lambda, node.NodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := Variable, node.Children[1].Children[0].NodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	env = parseEnvironemnt{}
	tokens = []*Token{
		{Dot, "."},
		{Word, "a"},
		{Arrow, "->"},
		{Word, "b"},
		{EOF, ""},
	}
	node, env, err = parseDot(tokens, env)
	if err != nil {
		t.Fatal(err)
	}
	if want, got := FreeVariable, node.Children[1].Children[0].NodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
}

func Test_parseWord(t *testing.T) {
	var env parseEnvironemnt
	tokens := []*Token{
		{Word, "a"},
		{Word, "b"},
		{Word, "c"},
		{Word, "d"},
		{EOF, ""},
	}
	node, env, err := parseWord(tokens, env)
	if err != nil {
		t.Fatal(err)
	}
	if want, got := 4, env.idx; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := Apply, node.NodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := "d", node.Children[1].Name; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	node = node.Children[0]
	if want, got := Apply, node.NodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := "c", node.Children[1].Name; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	node = node.Children[0]
	if want, got := Apply, node.NodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := "b", node.Children[1].Name; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	if want, got := "a", node.Children[0].Name; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
}

func Test_buildVariableNode(t *testing.T) {
	var env parseEnvironemnt
	env.AddKnownWord("a")
	n := buildVariableNode(env, "a")
	if want, got := Variable, n.NodeType; want != got {
		t.Errorf("want %v but got %v\n", want, got)
	}
	n = buildVariableNode(env, "b")
	if want, got := FreeVariable, n.NodeType; want != got {
		t.Errorf("want %v but got %v\n", want, got)
	}
}
