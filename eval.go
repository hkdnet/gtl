package gtl

import (
	"errors"
	"fmt"
)

type evalEnvironment struct {
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
	default:
		return nil, errors.New("error")
	}
}
