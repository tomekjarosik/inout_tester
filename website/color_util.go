package website

import "github.com/tomekjarosik/inout_tester/internal/testcase"

func ScoreColorFormat(score int) string {
	if score == 0 {
		return "red"
	}
	return "blue"
}

func TestCaseStatusColorFormat(status testcase.Status) string {
	if status == testcase.Accepted {
		return "green lighten-3"
	}
	return " red lighten-3"
}
