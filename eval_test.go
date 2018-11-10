package gtl

func ExampleEval() {
	ast := &AST{Child: &Node{NodeType: True}}
	Eval(ast)
	// Output: &{True [] }
}
