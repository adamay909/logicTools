// Code generated by "stringer -type LogicalConstant"; DO NOT EDIT.

package gentzen

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[None-0]
	_ = x[Neg-1]
	_ = x[Conj-2]
	_ = x[Disj-3]
	_ = x[Cond-4]
	_ = x[Uni-5]
	_ = x[Ex-6]
	_ = x[Ident-7]
	_ = x[Nec-8]
	_ = x[Pos-9]
}

const _LogicalConstant_name = "NoneNegConjDisjCondUniExIdentNecPos"

var _LogicalConstant_index = [...]uint8{0, 4, 7, 11, 15, 19, 22, 24, 29, 32, 35}

func (i LogicalConstant) String() string {
	if i < 0 || i >= LogicalConstant(len(_LogicalConstant_index)-1) {
		return "LogicalConstant(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _LogicalConstant_name[_LogicalConstant_index[i]:_LogicalConstant_index[i+1]]
}
