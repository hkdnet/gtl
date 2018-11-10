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

	Name string // for Variable, LambdaParam
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

func parse(tokens []*Token, env parseEnvironemnt) (*Node, parseEnvironemnt, error) {
	switch t := tokens[env.idx]; t.TokenType {
	case EOF:
		return parseEOF(tokens, env)
	case LParen:
		return parseLParen(tokens, env)
	case Number:
		return parseNumber(tokens, env)
	case Keyword:
		switch t.Text {
		case "true":
			return parseTrue(tokens, env)
		case "false":
			return parseFalse(tokens, env)
		case "then", "else":
			return nil, env, fmt.Errorf("unexpected token %v at %d", t.Text, env.idx)
		case "if":
			return parseIf(tokens, env)
		}
	case Dot: // start param
		return parseDot(tokens, env)
	case Word:
		return parseWord(tokens, env)
	}

	return nil, env, fmt.Errorf("cannot parse at %d", env.idx)
}

func parseEOF(tokens []*Token, env parseEnvironemnt) (*Node, parseEnvironemnt, error) {
	env.idx++
	return nil, env, nil
}

func parseLParen(tokens []*Token, env parseEnvironemnt) (*Node, parseEnvironemnt, error) {
	env.idx++
	env.parenCount++
	ret, nextEnv, err := parse(tokens, env)
	if err != nil {
		return nil, nextEnv, err
	}
	if tokens[nextEnv.idx].TokenType != RParen {
		return nil, env, fmt.Errorf("mismatch lparen at %d", env.idx)
	}
	nextEnv.parenCount--
	if nextEnv.parenCount < 0 {
		return nil, nextEnv, fmt.Errorf("mismatch rparen at %d", nextEnv.idx)
	}
	nextEnv.idx++
	return ret, nextEnv, nil
}

func parseNumber(tokens []*Token, env parseEnvironemnt) (*Node, parseEnvironemnt, error) {
	t := tokens[env.idx]
	if t.Text == "0" {
		env.idx++
		ret := &Node{NodeType: Zero}
		return ret, env, nil
	}
	err := fmt.Errorf("unknown number %v", t.Text)
	return nil, env, err
}

func parseTrue(tokens []*Token, env parseEnvironemnt) (*Node, parseEnvironemnt, error) {
	env.idx++
	return &Node{NodeType: True}, env, nil
}

func parseFalse(tokens []*Token, env parseEnvironemnt) (*Node, parseEnvironemnt, error) {
	env.idx++
	return &Node{NodeType: False}, env, nil
}

func parseIf(tokens []*Token, env parseEnvironemnt) (*Node, parseEnvironemnt, error) {
	env.idx++ // if
	ret := &Node{NodeType: IF, Children: make([]*Node, 3)}
	cond, env, err := parse(tokens, env)
	if err != nil {
		return nil, env, err
	}
	ret.Children[0] = cond
	if thenToken := tokens[env.idx]; thenToken.TokenType != Keyword || thenToken.Text != "then" {
		err := fmt.Errorf("token at %d should be then but %v", env.idx, thenToken)
		return nil, env, err
	}
	env.idx++ // then
	truePart, env, err := parse(tokens, env)
	if err != nil {
		return nil, env, err
	}
	ret.Children[1] = truePart
	if elseToken := tokens[env.idx]; elseToken.TokenType != Keyword || elseToken.Text != "else" {
		err := fmt.Errorf("token at %d should be else but %v", env.idx, elseToken)
		return nil, env, err
	}
	env.idx++ // else
	falsePart, env, err := parse(tokens, env)
	if err != nil {
		return nil, env, err
	}
	ret.Children[2] = falsePart
	return ret, env, err
}

// .x .y -> x y
func parseDot(tokens []*Token, env parseEnvironemnt) (*Node, parseEnvironemnt, error) {
	def := &Node{NodeType: LambdaDef}
	body := &Node{NodeType: LambdaBody}
	ret := &Node{NodeType: Lambda, Children: []*Node{def, body}}
paramLoop:
	for i := env.idx; ; {
		if len(tokens) <= i+1 {
			return nil, env, fmt.Errorf("after dot, there should be a variable but nothing at %d", i+1)
		}
		afterDot := tokens[i+1]
		if afterDot.TokenType != Word {
			return nil, env, fmt.Errorf("after dot, there should be a variable but got %v at %d", afterDot, i+1)
		}
		def.Children = append(def.Children, &Node{NodeType: LambdaParam, Name: afterDot.Text})
		if len(tokens) <= i+1 {
			return nil, env, fmt.Errorf("after a parameter, there should be a dot or arrow but nothing at %d", i+1)
		}
		i++ // skip parameter token
		dotOrArrow := tokens[i+1]
		switch dotOrArrow.TokenType {
		case Arrow:
			env.idx = i + 2 // skip arrow
			break paramLoop
		case Dot:
			i++
		default:
			return nil, env, fmt.Errorf("after a parameter, there should be a dot or arrow but got %v at %d", dotOrArrow, i+1)
		}
	}

	bc, env, err := parse(tokens, env)
	if err != nil {
		return nil, env, err
	}
	body.Children = []*Node{bc}

	return ret, env, nil
}

// x y z -> (x y) z
// a b c d -> ((a b) c) d
func parseWord(tokens []*Token, env parseEnvironemnt) (*Node, parseEnvironemnt, error) {
	if nt := tokens[env.idx+1]; nt.TokenType == EOF || nt.TokenType == RParen {
		ret := &Node{NodeType: Variable, Name: tokens[env.idx].Text}
		env.idx++
		return ret, env, nil
	}
	words := []*Token{tokens[env.idx]}
	env.idx++
applyLoop:
	for i := env.idx; ; {
		switch t := tokens[i]; t.TokenType {
		case RParen, EOF:
			env.idx = i
			break applyLoop
		case Word:
			words = append(words, tokens[i])
			i++
		default:
			return nil, env, fmt.Errorf("unexpected token %v at %d", t, i)
		}
	}
	var app *Node
	for i, l := 1, len(words); i < l; i++ {
		v := &Node{NodeType: Variable, Name: words[i].Text}
		if app == nil { // 1st
			first := &Node{NodeType: Variable, Name: words[0].Text}
			app = &Node{
				NodeType: Apply,
				Children: []*Node{first, v},
			}
		} else {
			app = &Node{
				NodeType: Apply,
				Children: []*Node{app, v},
			}
		}
	}
	return app, env, nil
}
