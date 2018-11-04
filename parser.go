package gtl

import (
	"errors"
	"fmt"
)

// AST is a abstract syntax tree. It contains only one Program Node.
type AST struct {
	Child *Node
}

// NodeType is an enum for node.
type NodeType uint8

// Node has nodeType and children
type Node struct {
	NodeType NodeType
	Children []*Node
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
	// Variable is a variable
	// TODO: better comment ...
	Variable
	// Lambda is a function. a lambda's children are always [LambdaDef, LambdaBody]
	Lambda
	// LambdaDef has some LambdaParams
	LambdaDef
	// LambdaParam represents a parameter of Lambda
	LambdaParam
	// LambdaBody has single child
	LambdaBody
)

// IsNumericalValue returns whether a node is a numerical value or not.
func (n *Node) IsNumericalValue() bool {
	if n.NodeType == Zero {
		return true
	}
	if n.NodeType == Succ {
		c := n.Children[0]
		return c.IsNumericalValue()
	}
	return false
}

// IsValue returns whether a node is a value or not.
func (n *Node) IsValue() bool {
	if n.NodeType == True || n.NodeType == False {
		return true
	}
	return n.IsNumericalValue()
}

// Parse returns an AST for tokens.
func Parse(tokens []*Token) (*AST, error) {
	var tmp *Node
	tmp = &Node{NodeType: Program}
	ret := &AST{Child: tmp}

	t, _, err := parse(tokens, 0)
	if err != nil {
		return nil, err
	}
	// TODO: if nextIdx != len(tokens check ?
	tmp.Children = []*Node{t}
	return ret, nil
}

func parse(tokens []*Token, i int) (*Node, int, error) {
	switch t := tokens[i]; t.tokenType {
	case EOF:
		return nil, i + 1, nil
	case Number:
		if t.text == "0" {
			return &Node{NodeType: Zero}, i + 1, nil
		}
		return nil, i, fmt.Errorf("unknown number %v", t.text)
	case Word:
		if t.text == "true" {
			return &Node{NodeType: True}, i + 1, nil
		}
		if t.text == "false" {
			return &Node{NodeType: False}, i + 1, nil
		}
		if t.text == "then" || t.text == "else" {
			return nil, i, fmt.Errorf("unexpected token %v at %d", t.text, i)
		}
		if t.text == "if" {
			ret := &Node{NodeType: IF, Children: make([]*Node, 3)}
			i = i + 1
			cond, nextIdx, err := parse(tokens, i)
			if err != nil {
				return nil, i, err
			}
			ret.Children[0] = cond
			i = nextIdx
			if thenToken := tokens[i]; thenToken.tokenType != Word || thenToken.text != "then" {
				return nil, i, fmt.Errorf("token at %d should be then but %v", i, thenToken)
			}
			i++
			t, nextIdx, err := parse(tokens, i)
			if err != nil {
				return nil, i, err
			}
			ret.Children[1] = t
			i = nextIdx
			if elseToken := tokens[i]; elseToken.tokenType != Word || elseToken.text != "else" {
				return nil, i, fmt.Errorf("token at %d should be else but %v", i, elseToken)
			}
			i++
			f, nextIdx, err := parse(tokens, i)
			if err != nil {
				return nil, i, err
			}
			ret.Children[2] = f
			i = nextIdx
			return ret, i, nil
		}

		if i+1 >= len(tokens) {
			return nil, i, fmt.Errorf("the last token should be eof but got word %v", t.text)
		}

		// variable name
		switch nextToken := tokens[i+1]; nextToken.tokenType {
		case EOF:
			return &Node{NodeType: Variable}, i + 1, nil
		case Dot:
			i += 2

			def := &Node{NodeType: LambdaDef}
			body := &Node{NodeType: LambdaBody}
			ret := &Node{NodeType: Lambda, Children: []*Node{def, body}}

			def.Children = append(def.Children, &Node{NodeType: LambdaParam}) // TODO: parameter name?
			for i+1 < len(tokens) {
				if tokens[i].tokenType == Arrow {
					break
				}
				if tokens[i].tokenType == Word && tokens[i+1].tokenType == Dot {
					def.Children = append(def.Children, &Node{NodeType: LambdaParam}) // TODO: parameter name?
					i += 2
					continue
				}
				return nil, i, fmt.Errorf("invalid lambda definition at %d-%d", i, i+1)
			}

			i++ // skip arrow

			bc, nextIdx, err := parse(tokens, i)
			if err != nil {
				return nil, nextIdx, err
			}
			body.Children = []*Node{bc}

			return ret, nextIdx, nil
		}

	}

	return nil, i, errors.New("cannot parse")
}
