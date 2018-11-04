package gtl

type AST struct {
	child *Node
}

type NodeType uint8

type Node struct {
	nodeType NodeType
	children []*Node
}

const (
	PROGRAM NodeType = iota
	TRUE
	FALSE
	IF
	ZERO
	SUCC
	PRED
	ISZERO
)

func (n *Node) IsNumericalValue() bool {
	if n.nodeType == ZERO {
		return true
	}
	if n.nodeType == SUCC {
		c := n.children[0]
		return c.IsNumericalValue()
	}
	return false
}

func (n *Node) IsValue() bool {
	if n.nodeType == TRUE || n.nodeType == FALSE {
		return true
	}
	return n.IsNumericalValue()
}

func Parse(tokens []Token) (*AST, error) {
	var tmp *Node
	tmp = &Node{nodeType: PROGRAM}
	ret := &AST{child: tmp}

	for i, l := 0, len(tokens); i < l; i++ {
		switch t := tokens[i]; t.tokenType {
		case EOF:
			break
		case IDENTIFIER:
			if t.text == "true" {
				tmp.children = []*Node{
					&Node{nodeType: TRUE},
				}
				continue
			}
			if t.text == "false" {
				tmp.children = []*Node{
					&Node{nodeType: FALSE},
				}
				continue
			}
		}
	}

	return ret, nil
}
