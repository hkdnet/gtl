// Code generated by "stringer -type=NodeType"; DO NOT EDIT.

package gtl

import "strconv"

const _NodeType_name = "TrueFalseIFZeroSuccPredIsZeroVariableFreeVariableLambdaLambdaDefLambdaParamLambdaBodyApplyNodeNumber"

var _NodeType_index = [...]uint8{0, 4, 9, 11, 15, 19, 23, 29, 37, 49, 55, 64, 75, 85, 90, 100}

func (i NodeType) String() string {
	if i >= NodeType(len(_NodeType_index)-1) {
		return "NodeType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _NodeType_name[_NodeType_index[i]:_NodeType_index[i+1]]
}
