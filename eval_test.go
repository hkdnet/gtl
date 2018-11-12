package gtl

import "testing"

func evalFromString(str string) (*AST, error) {
	l := NewLexer(str)
	var tokens []*Token
	for l.HasNext() {
		t, err := l.NextToken()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
	}
	return Parse(tokens)
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
}

func Test_evalIf(t *testing.T) {
	ast, err := evalFromString("if true then a else b")
	if err != nil {
		t.Fatal(err)
	}
	ast.show()
	ifNode := ast.Child
	if want, got := IF, ifNode.NodeType; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
	var env evalEnvironment
	n, err := evalIf(ifNode, &env)
	if err != nil {
		t.Fatal(err)
	}
	if want, got := "a", n.Name; got != want {
		t.Errorf("want %v but got %v\n", want, got)
	}
}
