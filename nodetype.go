package gtl

// NodeType is an enum for node.
type NodeType uint8

const (
	// True is literal true
	True NodeType = iota
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
	// Variable is a variable, which is bound by outer LambdaDef
	Variable
	// FreeVariable is a free variable
	FreeVariable
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
	// NodeNumber is a numerical value
	NodeNumber
)
