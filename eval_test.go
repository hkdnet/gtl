package gtl

import (
	"testing"
)

// NOTE: this function may cause panic
func buildASTFromString(str string) *AST {
	l := NewLexer(str)
	var tokens []*Token
	for l.HasNext() {
		t, err := l.NextToken()
		if err != nil {
			panic(err)
		}
		tokens = append(tokens, t)
	}
	ast, err := Parse(tokens)
	if err != nil {
		panic(err) // for convenience
	}
	return ast
}

// helper function
func assertEval(source string, assert func(*Node)) {
	ast := buildASTFromString(source)
	node := ast.Child
	var env evalEnvironment
	n, err := eval(node, &env)
	if err != nil {
		panic(err)
	}
	assert(n)
}

func Test_evalEnvironment(t *testing.T) {
	var ee evalEnvironment

	node1 := &Node{}
	node2 := &Node{}
	node3 := &Node{}
	ee.Assign("a", node1)
	ee.Assign("aa", node2)
	v := ee.Lookup("a")
	if v != node1 {
		t.Errorf("Lookup should return exact match")
	}
	ee.Assign("a", node3) // overwrite
	v = ee.Lookup("a")
	if v != node3 {
		t.Errorf("should return latest assignment")
	}
	err := ee.Unassign("aa")
	if err != nil {
		t.Fatal(err)
	}
	v = ee.Lookup("aa")
	if v != nil {
		t.Error("Lookup for unassigned name should return nil")
	}
	err = ee.Unassign("no-such-key")
	if err == nil {
		t.Error("Unassign for not-assigned name should return error")
	}
}

func Test_evalIf(t *testing.T) {
	assertEval("if true then a else b", func(n *Node) {
		if want, got := "a", n.Name; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
	})
	assertEval("if false then a else b", func(n *Node) {
		if want, got := "b", n.Name; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
	})
	assertEval("if a then b else c", func(n *Node) {
		if want, got := IF, n.NodeType; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
		if want, got := "if (a) then (b) else (c)", n.String(); got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
	})
}

func Test_evalApply(t *testing.T) {
	assertEval("(.a -> a) true", func(n *Node) {
		if want, got := True, n.NodeType; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
	})
	assertEval("a true", func(n *Node) {
		if want, got := Apply, n.NodeType; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
		if want, got := FreeVariable, n.Children[0].NodeType; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
		if want, got := True, n.Children[1].NodeType; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
	})
	assertEval("iszero 0", func(n *Node) {
		if want, got := True, n.NodeType; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
	})

	assertEval("(.a .b -> a b) iszero", func(n *Node) {
		// .b -> iszero b
		if want, got := Lambda, n.NodeType; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
		if want, got := 2, len(n.Children); got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
		if want, got := LambdaDef, n.Children[0].NodeType; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
		if want, got := LambdaBody, n.Children[1].NodeType; got != want {
			t.Errorf("want %v but got %v\n", want, got)
		}
	})
}
