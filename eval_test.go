package gtl

import "testing"

func ExampleEval() {
	ast := &AST{Child: &Node{NodeType: True}}
	Eval(ast)
	// Output: &{True [] }
}

func Test_evalEnvironment(t *testing.T) {
	var ee evalEnvironment

	node1 := &Node{}
	node2 := &Node{}
	node3 := &Node{}
	ee.Assign("a", node1)
	ee.Assign("aa", node2)
	v, err := ee.Lookup("a")
	if err != nil {
		t.Fatal(err)
	}
	if v != node1 {
		t.Errorf("Lookup should return exact match")
	}
	ee.Assign("a", node3) // overwrite
	v, err = ee.Lookup("a")
	if err != nil {
		t.Fatal(err)
	}
	if v != node3 {
		t.Errorf("should return latest assignment")
	}
	err = ee.Unassign("aa")
	if err != nil {
		t.Fatal(err)
	}
	_, err = ee.Lookup("aa")
	if err == nil {
		t.Error("Lookup for unassigned name should return error")
	}
}
