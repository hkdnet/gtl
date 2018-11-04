package gtl

import (
	"errors"
	"fmt"
)

type AST struct {
	child *Node
}

type NodeType uint8

type Node struct {
	nodeType NodeType
	children []*Node
}

const (
	Program NodeType = iota
	True
	False
	IF
	Zero
	Succ
	Pred
	IsZero
)

func (n *Node) IsNumericalValue() bool {
	if n.nodeType == Zero {
		return true
	}
	if n.nodeType == Succ {
		c := n.children[0]
		return c.IsNumericalValue()
	}
	return false
}

func (n *Node) IsValue() bool {
	if n.nodeType == True || n.nodeType == False {
		return true
	}
	return n.IsNumericalValue()
}

func Parse(tokens []Token) (*AST, error) {
	var tmp *Node
	tmp = &Node{nodeType: Program}
	ret := &AST{child: tmp}

	t, _, err := parse(tokens, 0)
	if err != nil {
		return nil, err
	}
	// TODO: if nextIdx != len(tokens check ?
	tmp.children = []*Node{t}
	return ret, nil
}

func parse(tokens []Token, i int) (*Node, int, error) {
	fmt.Printf("DEBUG: %d\n", i)
	switch t := tokens[i]; t.tokenType {
	case EOF:
		return nil, i + 1, nil
	case Word:
		if t.text == "true" {
			return &Node{nodeType: True}, i + 1, nil
		}
		if t.text == "false" {
			return &Node{nodeType: False}, i + 1, nil
		}
		if t.text == "if" {
			ret := &Node{nodeType: IF, children: make([]*Node, 3)}
			i = i + 1
			cond, nextIdx, err := parse(tokens, i)
			if err != nil {
				return nil, i, err
			}
			ret.children[0] = cond
			i = nextIdx
			if thenToken := tokens[i]; thenToken.tokenType != Word || thenToken.text != "then" {
				return nil, i, fmt.Errorf("token at %d should be then but %v", i, thenToken)
			}
			i++
			t, nextIdx, err := parse(tokens, i)
			if err != nil {
				return nil, i, err
			}
			ret.children[1] = t
			i = nextIdx
			if elseToken := tokens[i]; elseToken.tokenType != Word || elseToken.text != "else" {
				return nil, i, fmt.Errorf("token at %d should be else but %v", i, elseToken)
			}
			i++
			f, nextIdx, err := parse(tokens, i)
			if err != nil {
				return nil, i, err
			}
			ret.children[2] = f
			i = nextIdx
			return ret, i, nil
		}
	}

	return nil, i, errors.New("cannot parse")
}
