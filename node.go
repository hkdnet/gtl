package gtl

import (
	"fmt"
	"strings"
)

// Node has nodeType and children
type Node struct {
	NodeType NodeType
	Children []*Node

	Name string // for Variable, LambdaParam
}

func (n *Node) String() string {
	switch n.NodeType {
	case True:
		return "true"
	case False:
		return "false"
	case IF:
		return fmt.Sprintf("if (%s) then (%s) else (%s)", n.Children[0], n.Children[1], n.Children[2])
	case Zero:
		return "0"
	case Succ:
		return "succ"
	case Pred:
		return "pred"
	case IsZero:
		return "iszero"
	case Variable, FreeVariable:
		return n.Name
	case Lambda:
		return fmt.Sprintf("%s -> (%s)", n.Children[0], n.Children[1])
	case LambdaDef:
		var tmp []string
		for _, p := range n.Children {
			tmp = append(tmp, fmt.Sprintf("%s.", p.Name))
		}
		return strings.Join(tmp, " ")
	case LambdaBody:
		return n.Children[0].String()
	case Apply:
		return fmt.Sprintf("%s %s", n.Children[0], n.Children[1])
	default:
		panic("unknown type?")
	}
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
