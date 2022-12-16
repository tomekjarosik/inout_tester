package website

import (
	"fmt"
	"github.com/tomekjarosik/inout_tester/internal/testcase"
	"time"
)

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

func TestCaseDurationFormatFunc(duration time.Duration) string {
	return fmt.Sprintf("%ds %3d ms", int(duration.Seconds()), int(duration.Milliseconds()))
}