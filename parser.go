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

	for i, l := 0, len(tokens); i < l; i++ {
		switch t := tokens[i]; t.tokenType {
		case EOF:
			break
		case Word:
			if t.text == "true" {
				tmp.children = []*Node{
					&Node{nodeType: True},
				}
				continue
			}
			if t.text == "false" {
				tmp.children = []*Node{
					&Node{nodeType: False},
				}
				continue
			}
		}
	}

	return ret, nil
}
