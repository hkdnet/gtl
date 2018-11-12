package gtl

import (
	"fmt"
)

type assginment struct {
	name  string
	value *Node
}

type evalEnvironment struct {
	assginments []assginment
}

func (ee *evalEnvironment) Assign(name string, val *Node) {
	ee.assginments = append(ee.assginments, assginment{name, val})
}

func (ee *evalEnvironment) Lookup(name string) *Node {
	for i := len(ee.assginments) - 1; i >= 0; i-- {
		if ee.assginments[i].name == name {
			return ee.assginments[i].value
		}
	}
	return nil
}

func (ee *evalEnvironment) Unassign(name string) error {
	for i := len(ee.assginments) - 1; i >= 0; i-- {
		if ee.assginments[i].name == name {
			ee.assginments = append(ee.assginments[:i], ee.assginments[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("missing unassignment target %s", name)
}

// Eval returns evaluated node
func Eval(ast *AST) (*Node, error) {
	var env evalEnvironment
	n, err := eval(ast.Child, &env)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func eval(n *Node, env *evalEnvironment) (*Node, error) {
	switch n.NodeType {
	case True, False, Zero, FreeVariable, Lambda:
		return n, nil
	case IF:
		return evalIf(n, env)
	case Apply:
		return evalApply(n, env)
	case Variable:
		return evalVariable(n, env)
	default:
		return nil, fmt.Errorf("cannot eval: %s", n.NodeType)
	}
}

func evalIf(n *Node, env *evalEnvironment) (*Node, error) {
	cond, err := eval(n.Children[0], env)
	if err != nil {
		return nil, err
	}
	if cond.NodeType == True {
		return eval(n.Children[1], env)
	}
	if cond.NodeType == False {
		return eval(n.Children[2], env)
	}
	// TODO: type check
	truePart, err := eval(n.Children[1], env)
	if err != nil {
		return nil, err
	}
	falsePart, err := eval(n.Children[2], env)
	if err != nil {
		return nil, err
	}
	return &Node{NodeType: IF, Children: []*Node{cond, truePart, falsePart}}, nil
}

func evalApply(n *Node, env *evalEnvironment) (*Node, error) {
	var err error
	l := n.Children[0]
	r := n.Children[1]
	l, err = eval(l, env)
	if err != nil {
		return nil, err
	}
	r, err = eval(r, env)
	if err != nil {
		return nil, err
	}
	if l.NodeType != Lambda { // cannot eval apply
		return &Node{NodeType: Apply, Children: []*Node{l, r}}, nil
	}
	// l.NodeType == Lambda
	def := l.Children[0]
	body := l.Children[1]
	assigned := def.Children[0]
	env.Assign(assigned.Name, r)
	body, err = eval(body.Children[0], env)
	if err != nil {
		return nil, err
	}
	env.Unassign(assigned.Name)
	if len(def.Children) == 1 {
		return body, nil
	}
	def.Children = def.Children[1:] // tail
	ret := &Node{
		NodeType: Lambda,
		Children: []*Node{def, body},
	}
	return ret, nil
}

func evalVariable(n *Node, env *evalEnvironment) (*Node, error) {
	val := env.Lookup(n.Name)
	if val == nil {
		// nil means this variable is literally bound but its value is not yet bound...
		// (.x .y -> x y) iszero => iszero y
		return n, nil
	}
	return val, nil
}
