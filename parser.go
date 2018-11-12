package gtl

import (
	"errors"
	"fmt"
)

// AST is a abstract syntax tree. It contains only one Program Node.
type AST struct {
	Child *Node
}

func (ast *AST) show() {
	ast.Child.show("")
}

// Node has nodeType and children
type Node struct {
	NodeType NodeType
	Children []*Node

	Name string // for Variable, LambdaParam
}

func (n *Node) show(indent string) {
	fmt.Printf("%s%s\n", indent, n.NodeType)
	nextIndent := indent + "  "
	for _, c := range n.Children {
		c.show(nextIndent)
	}
}

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

	knownWords []string
}

func (e *parseEnvironemnt) AddKnownWord(name string) {
	e.knownWords = append(e.knownWords, name)
}

func (e *parseEnvironemnt) RemoveKnownWord(name string) error {
	for i := len(e.knownWords) - 1; i >= 0; i-- {
		if e.knownWords[i] == name {
			e.knownWords = append(e.knownWords[:i], e.knownWords[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("unknown word %s", name)
}

func (e *parseEnvironemnt) IsBound(name string) bool {
	for i := len(e.knownWords) - 1; i >= 0; i-- {
		if e.knownWords[i] == name {
			return true
		}
	}
	return false
}

// Parse returns an AST for tokens.
func Parse(tokens []*Token) (*AST, error) {
	var nodes []*Node
	var env parseEnvironemnt
parseLoop:
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
				return nil, fmt.Errorf("parse returns nil pointer at %d", env.idx)
			}
			break parseLoop
		}
		nodes = append(nodes, t)
	}
	if l := len(nodes); l == 0 {
		return nil, errors.New("no nodes")
	} else if l == 1 {
		ret := &AST{Child: nodes[0]}
		return ret, nil
	}
	app := nodesToApply(nodes)
	ret := &AST{Child: app}
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
	case KeywordTrue:
		return parseTrue(tokens, env)
	case KeywordFalse:
		return parseFalse(tokens, env)
	case KeywordThen, KeywordElse:
		return nil, env, fmt.Errorf("unexpected token %v at %d", t.Text, env.idx)
	case KeywordIf:
		return parseIf(tokens, env)
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
	if thenToken := tokens[env.idx]; thenToken.TokenType != KeywordThen {
		err := fmt.Errorf("token at %d should be then but %v", env.idx, thenToken)
		return nil, env, err
	}
	env.idx++ // then
	truePart, env, err := parse(tokens, env)
	if err != nil {
		return nil, env, err
	}
	ret.Children[1] = truePart
	if elseToken := tokens[env.idx]; elseToken.TokenType != KeywordElse {
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

	for _, p := range def.Children {
		env.AddKnownWord(p.Name)
	}
	bc, env, err := parse(tokens, env)
	if err != nil {
		return nil, env, err
	}
	for _, p := range def.Children {
		env.RemoveKnownWord(p.Name)
	}
	body.Children = []*Node{bc}

	return ret, env, nil
}

func buildVariableNode(env parseEnvironemnt, name string) *Node {
	var nt NodeType
	if env.IsBound(name) {
		nt = Variable
	} else {
		nt = FreeVariable
	}
	return &Node{NodeType: nt, Name: name}
}

// len(nodes) must be more than 1
func nodesToApply(nodes []*Node) *Node {
	var app *Node
	for i, l := 1, len(nodes); i < l; i++ {
		if app == nil { // 1st
			app = &Node{
				NodeType: Apply,
				Children: []*Node{nodes[0], nodes[i]},
			}
		} else {
			app = &Node{
				NodeType: Apply,
				Children: []*Node{app, nodes[i]},
			}
		}
	}
	return app
}

// x y z -> (x y) z
// a b c d -> ((a b) c) d
func parseWord(tokens []*Token, env parseEnvironemnt) (*Node, parseEnvironemnt, error) {
	if nt := tokens[env.idx+1]; nt.TokenType == EOF || nt.TokenType == RParen || nt.TokenType == KeywordThen || nt.TokenType == KeywordElse {
		ret := buildVariableNode(env, tokens[env.idx].Text)
		env.idx++
		return ret, env, nil
	}
	nodes := []*Node{buildVariableNode(env, tokens[env.idx].Text)}
	env.idx++
applyLoop:
	for i := env.idx; ; {
		switch t := tokens[i]; t.TokenType {
		case RParen, EOF:
			env.idx = i
			break applyLoop
		case Word:
			v := buildVariableNode(env, tokens[i].Text)
			nodes = append(nodes, v)
			i++
		case KeywordTrue:
			v := &Node{NodeType: True}
			nodes = append(nodes, v)
			i++
		case KeywordFalse:
			v := &Node{NodeType: False}
			nodes = append(nodes, v)
			i++
		case KeywordThen, KeywordElse:
			env.idx = i
			break applyLoop
		default:
			return nil, env, fmt.Errorf("unexpected token %v at %d", t, i)
		}
	}
	app := nodesToApply(nodes)
	return app, env, nil
}
