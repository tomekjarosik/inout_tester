package testcase

//go:generate go build -o testdata/multiply2.exe testdata/multiply2.go
//go:generate go build -o testdata/multiply3.exe testdata/multiply3.go
//go:generate go build -o testdata/infinite_loop.exe testdata/infinite_loop.go

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunTestCase_Success(t *testing.T) {
	info := Info{Name: "test1", TimeLimit: 3 * time.Second}
	streams := Streams{
		Input:  strings.NewReader("1\n"),
		Output: strings.NewReader("2\n"),
	}
	res := runTestWithTmpOutput("testdata/multiply2.exe", info, streams)
	assert.Equal(t, Accepted, res.Status)
}

func TestRunTestCase_WrongAnswer(t *testing.T) {
	info := Info{Name: "test1", TimeLimit: 3 * time.Second}
	streams := Streams{
		Input:  strings.NewReader("1\n"),
		Output: strings.NewReader("2\n"),
	}
	res := runTestWithTmpOutput("testdata/multiply3.exe", info, streams)
	assert.Equal(t, WrongAnswer, res.Status)
}

func TestRunTestCase_TimeLimitExceeded(t *testing.T) {
	info := Info{Name: "test1", TimeLimit: 1000 * time.Millisecond}
	streams := Streams{
		Input:  strings.NewReader("1\n"),
		Output: strings.NewReader("2\n"),
	}
	res := runTestWithTmpOutput("testdata/infinite_loop.exe", info, streams)
	assert.Equal(t, TimeLimitExceeded, res.Status)
	assert.Equal(t, "time limit exceeded: test case was aborted after '1s'", res.Description)
}

func TestCompare_Identical(t *testing.T) {
	assert.NoError(t, compare(strings.NewReader("1\n2\n3\n"), strings.NewReader("1\n2\n3\n")))
}

func TestCompare_Different1(t *testing.T) {
	err := compare(strings.NewReader("123\n456\n789\n"), strings.NewReader("123\nggggg\n789\n"))
	assert.Equal(t, err, errors.New("outputs differ in line 1: expected: '456', actual: 'ggggg'"))
}

func TestCompare_AdditionalLines(t *testing.T) {
	err := compare(strings.NewReader("123\n456\n789\n"), strings.NewReader("123\n456\n789\nabcde\n"))
	assert.Equal(t, err, errors.New("contains additional non-empty lines"))
	err = compare(strings.NewReader("123\n456\n789\nabcde\n"), strings.NewReader("123\n456\n789\n"))
	assert.Equal(t, err, errors.New("outputs differ in line 3: expected: 'abcde', actual: ''"))
}

func TestCompare_LongLines(t *testing.T) {
	str1 := strings.Repeat("a", 300)
	str2 := strings.Repeat("a", 350)

	err := compare(strings.NewReader(str1), strings.NewReader(str2))
	assert.Equal(t, err, errors.New("outputs differ in line 0: expected: 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa...', actual: 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa...'"))

	err = compare(strings.NewReader("abc"), strings.NewReader(str2))
	assert.Equal(t, err, errors.New("outputs differ in line 0: expected: 'abc', actual: 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa...'"))

}
