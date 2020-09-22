package main

//go:generate go build -o testdata/multiply2.exe testdata/multiply2.go
//go:generate go build -o testdata/multiply3.exe testdata/multiply3.go
//go:generate go build -o testdata/infinite_loop.exe testdata/infinite_loop.go

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunTestCase_Success(t *testing.T) {
	tc := TestCase{
		Name:                 "t1",
		InputFilename:        "testdata/problems/volvo/t1.in",
		GoldenOutputFilename: "testdata/problems/volvo/t1.out",
		OutputFilename:       "testdata/submissions/volvo/t1.unittests2.out",
		TimeLimit:            3 * time.Second,
	}
	outTc := runSingleTestCase("testdata/multiply2.exe", tc)
	assert.Equal(t, Accepted, outTc.Status)
}

func TestRunTestCase_WrongAnswer(t *testing.T) {
	tc := TestCase{
		Name:                 "t1",
		InputFilename:        "testdata/problems/volvo/t1.in",
		GoldenOutputFilename: "testdata/problems/volvo/t1.out",
		OutputFilename:       "testdata/submissions/volvo/t1.unittests3.out",
		TimeLimit:            3 * time.Second,
	}
	outTc := runSingleTestCase("testdata/multiply3.exe", tc)
	assert.Equal(t, Accepted, outTc.Status)
}

func TestRunTestCase_MissingInputFile(t *testing.T) {
	tc := TestCase{
		Name:                 "t1",
		InputFilename:        "testdata/problems/volvo/t1Invalid.in",
		GoldenOutputFilename: "testdata/problems/volvo/t1.out",
		OutputFilename:       "testdata/submissions/volvo/t1.unittests.out",
		TimeLimit:            3 * time.Second,
	}
	outTc := runSingleTestCase("testdata/multiply2.exe", tc)
	assert.Equal(t, InternalError, outTc.Status)
	assert.Equal(t, "unable to open input file at: 'testdata/problems/volvo/t1Invalid.in'", outTc.StatusDescription)
}

func TestRunTestCase_InvalidOutputPath(t *testing.T) {
	tc := TestCase{
		Name:                 "t1",
		InputFilename:        "testdata/problems/volvo/t1.in",
		GoldenOutputFilename: "testdata/problems/volvo/t1.out",
		OutputFilename:       "testdata/submissions/volvo11111/t1.unittests.out",
		TimeLimit:            3 * time.Second,
	}
	outTc := runSingleTestCase("testdata/multiply2.exe", tc)
	assert.Equal(t, InternalError, outTc.Status)
	assert.Equal(t, "unable to create output file at: 'testdata/submissions/volvo11111/t1.unittests.out'", outTc.StatusDescription)
}

func TestRunTestCase_TimeLimitExceeded(t *testing.T) {
	tc := TestCase{
		Name:                 "t1",
		InputFilename:        "testdata/problems/volvo/t1.in",
		GoldenOutputFilename: "testdata/problems/volvo/t1.out",
		OutputFilename:       "testdata/submissions/volvo/t1.unittests.out",
		TimeLimit:            1000 * time.Millisecond,
	}
	outTc := runSingleTestCase("testdata/infinite_loop.exe", tc)
	assert.Equal(t, TimeLimitExceeded, outTc.Status)
	assert.Equal(t, "time limit exceeded: test case was aborted after '1s'", outTc.StatusDescription)
}

const diffDataDir = "testdata/filediffdata/"

func TestCompareFiles_Identical(t *testing.T) {
	assert.NoError(t, compareFiles(diffDataDir+"a.txt", diffDataDir+"a.txt"))
}
func TestCompareFiles_Different1(t *testing.T) {
	err := compareFiles(diffDataDir+"a.txt", diffDataDir+"b.txt")
	assert.Equal(t, err, errors.New("files differ in line 1: expected: 456, actual: gggggg  gggggg"))
}
func TestCompareFiles_AdditionalLines(t *testing.T) {
	err := compareFiles(diffDataDir+"a.txt", diffDataDir+"a_extended.txt")
	assert.Equal(t, err, errors.New("file testdata/filediffdata/a_extended.txt contains additional non-empty lines"))
	err = compareFiles(diffDataDir+"a_extended.txt", diffDataDir+"a.txt")
	assert.Equal(t, err, errors.New("files differ in line 3: expected: 10 12 13, actual: "))
}

func TestCompareFiles_DifferentLongLines(t *testing.T) {
	err := compareFiles(diffDataDir+"long.txt", diffDataDir+"a_extended.txt")
	assert.Equal(t, err, errors.New("files differ in line 0: expected: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa..., actual: 123"))
	err = compareFiles(diffDataDir+"a_extended.txt", diffDataDir+"long.txt")
	assert.Equal(t, err, errors.New("files differ in line 0: expected: 123, actual: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa..."))
}
