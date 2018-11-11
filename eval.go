package gtl

import (
	"errors"
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

func (ee *evalEnvironment) Lookup(name string) (*Node, error) {
	for i := len(ee.assginments) - 1; i >= 0; i-- {
		if ee.assginments[i].name == name {
			return ee.assginments[i].value, nil
		}
	}
	return nil, fmt.Errorf("missing assignment for %s", name)
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

// Eval shows evaluated value
func Eval(ast *AST) error {
	var env evalEnvironment
	n, err := eval(ast.Child, env)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", n)
	return nil
}

func eval(n *Node, env evalEnvironment) (*Node, error) {
	switch n.NodeType {
	case True, False, Zero:
		return n, nil
	case IF:
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
		return nil, errors.New("cond should be bool but not")
	case Variable:
		val, err := env.Lookup(n.Name)
		if err != nil {
			return nil, err
		}
		return val, nil
	default:
		return nil, errors.New("error")
	}
}
