// Code generated by "stringer -type=Status"; DO NOT EDIT.

package testcase

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NotRunYet-1]
	_ = x[InternalError-2]
	_ = x[TimeLimitExceeded-3]
	_ = x[MemoryLimitExceeded-4]
	_ = x[WrongAnswer-5]
	_ = x[Accepted-6]
	_ = x[RuntimeError-7]
}

const _Status_name = "NotRunYetInternalErrorTimeLimitExceededMemoryLimitExceededWrongAnswerAcceptedRuntimeError"

var _Status_index = [...]uint8{0, 9, 22, 39, 58, 69, 77, 89}

func (i Status) String() string {
	i -= 1
	if i < 0 || i >= Status(len(_Status_index)-1) {
		return "Status(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Status_name[_Status_index[i]:_Status_index[i+1]]
}
