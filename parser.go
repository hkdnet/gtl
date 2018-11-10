package gtl

import (
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
	// Apply is "function call"
	Apply
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

type parseEnvironemnt struct {
	idx        int
	parenCount int
}

// Parse returns an AST for tokens.
func Parse(tokens []*Token) (*AST, error) {
	var tmp *Node
	tmp = &Node{NodeType: Program}
	ret := &AST{Child: tmp}

	var env parseEnvironemnt
	for l := len(tokens); env.idx < l; {
		var t *Node
		var err error
		t, env, err = parse(tokens, env)
		if err != nil {
			return nil, err
		}
		// FIXME: too tricky...
		if t == nil {
			if env.idx != l {
				return ret, fmt.Errorf("parse returns nil pointer at %d", env.idx)
			}
			return ret, nil
		}
		tmp.Children = append(tmp.Children, t)
	}
	return ret, nil
}

func parse(tokens []*Token, _env parseEnvironemnt) (ret *Node, env parseEnvironemnt, err error) {
	env = _env
	switch t := tokens[env.idx]; t.TokenType {
	case EOF:
		env.idx++
		return
	case LParen:
		env.idx++
		env.parenCount++
		return parse(tokens, env)
	case Number:
		if t.Text == "0" {
			env.idx++
			ret = &Node{NodeType: Zero}
			return
		}
		err = fmt.Errorf("unknown number %v", t.Text)
		return
	case Word:
		if t.Text == "true" {
			env.idx++
			ret = &Node{NodeType: True}
			return
		}
		if t.Text == "false" {
			env.idx++
			ret = &Node{NodeType: False}
			return
		}
		if t.Text == "then" || t.Text == "else" {
			err = fmt.Errorf("unexpected token %v at %d", t.Text, env.idx)
			return
		}
		if t.Text == "if" {
			ret = &Node{NodeType: IF, Children: make([]*Node, 3)}
			env.idx++
			var cond *Node
			var truePart *Node
			var falsePart *Node
			cond, env, err = parse(tokens, env)
			if err != nil {
				return
			}
			ret.Children[0] = cond
			if thenToken := tokens[env.idx]; thenToken.TokenType != Word || thenToken.Text != "then" {
				err = fmt.Errorf("token at %d should be then but %v", env.idx, thenToken)
				return
			}
			env.idx++
			truePart, env, err = parse(tokens, env)
			if err != nil {
				return
			}
			ret.Children[1] = truePart
			if elseToken := tokens[env.idx]; elseToken.TokenType != Word || elseToken.Text != "else" {
				err = fmt.Errorf("token at %d should be else but %v", env.idx, elseToken)
				return
			}
			env.idx++
			falsePart, env, err = parse(tokens, env)
			if err != nil {
				return
			}
			ret.Children[2] = falsePart
			env.idx++
			return
		}

		if env.idx+1 >= len(tokens) {
			err = fmt.Errorf("the last token should be eof but got word %v", t.Text)
			return
		}

		// variable name
		switch nextToken := tokens[env.idx+1]; nextToken.TokenType {
		case RParen:
			if env.parenCount < 1 {
				err = fmt.Errorf("paren mismatch at %d", env.idx+1)
				return
			}
			env.parenCount--
			env.idx += 2
			ret = &Node{NodeType: Variable}
			return
		case EOF:
			env.idx++
			ret = &Node{NodeType: Variable}
			return
		case Dot:
			env.idx += 2

			def := &Node{NodeType: LambdaDef}
			body := &Node{NodeType: LambdaBody}
			ret = &Node{NodeType: Lambda, Children: []*Node{def, body}}

			def.Children = append(def.Children, &Node{NodeType: LambdaParam}) // TODO: parameter name?
			for env.idx+1 < len(tokens) {
				if tokens[env.idx].TokenType == Arrow {
					break
				}
				if tokens[env.idx].TokenType == Word && tokens[env.idx+1].TokenType == Dot {
					def.Children = append(def.Children, &Node{NodeType: LambdaParam}) // TODO: parameter name?
					env.idx += 2
					continue
				}
				err = fmt.Errorf("invalid lambda definition at %d-%d", env.idx, env.idx+1)
				return
			}

			env.idx++ // skip arrow

			var bc *Node
			bc, env, err = parse(tokens, env)
			if err != nil {
				return
			}
			body.Children = []*Node{bc}

			return
		case Word: // TODO: is this only apply?
			env.idx += 2
			children := []*Node{}
			ret = &Node{NodeType: Apply, Children: children}
			return
		}
	}

	err = fmt.Errorf("cannot parse at %d", env.idx)
	return
}
