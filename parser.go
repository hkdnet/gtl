package gtl

import (
	"errors"
	"fmt"
)

// AST is a abstract syntax tree. It contains only one Program Node.
type AST struct {
	child *Node
}

// NodeType is an enum for node.
type NodeType uint8

// Node has nodeType and children
type Node struct {
	nodeType NodeType
	children []*Node
}

const (
	// Program is a toplevel node
	Program NodeType = iota
	// True is literal true
	True
	// False is literal false
	False
	// IF is a if expression
	IF
	// Zero is literal 0
	Zero
	// Succ is a builtin function, succ
	Succ
	// Pred is a builtin function, pred
	Pred
	// IsZero is a builtin function, iszero
	IsZero
)

// IsNumericalValue returns whether a node is a numerical value or not.
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

// IsValue returns whether a node is a value or not.
func (n *Node) IsValue() bool {
	if n.nodeType == True || n.nodeType == False {
		return true
	}
	return n.IsNumericalValue()
}

// Parse returns an AST for tokens.
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
	case Number:
		if t.text == "0" {
			return &Node{nodeType: Zero}, i + 1, nil
		}
		return nil, i, fmt.Errorf("unknown number %v", t.text)
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
