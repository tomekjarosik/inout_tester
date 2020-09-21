package main

//go:generate stringer -type=TestCaseStatus

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

// TestCaseStatus status of the
type TestCaseStatus int

const (
	// NotRunYet test is waiting to be run
	NotRunYet TestCaseStatus = iota
	// InternalError something unexpected went wrong
	InternalError
	// CompilationError the solution failed to compile
	CompilationError
	// TimeLimitExceeded the test took too long to process
	TimeLimitExceeded
	// MemoryLimitExceeded the test run used too much RAM
	MemoryLimitExceeded
	// WrongAnswer test run successfully but test outputs differ
	WrongAnswer
	// Accepted all OK
	Accepted
)

// TestCase struct describing results of a single test run (input, output, expected output)
type TestCase struct {
	Name                 string        `json:"name"`
	InputFilename        string        `json:"in"`
	GoldenOutputFilename string        `json:"goldenOut"`
	OutputFilename       string        `json:"out"`
	TimeLimit            time.Duration `json:"timeLimit"`
	MemoryLimit          int           `json:"memoryLimit"`

	Status            TestCaseStatus `json:"status"`
	StatusDescription string         `json:"statusDescription"`
	Duration          time.Duration  `json:"duration"`
}

// NewTestCase construct of TestCase struct
func NewTestCase(solutionID string, testName string, problemDataDir string, outputDir string) TestCase {
	return TestCase{
		Name:                 testName,
		Status:               NotRunYet,
		Duration:             0,
		InputFilename:        path.Join(problemDataDir, testName+".in"),
		GoldenOutputFilename: path.Join(problemDataDir, testName+".out"),
		OutputFilename:       path.Join(outputDir, testName+"."+solutionID+".out"),
		TimeLimit:            10 * time.Second,
	}
}

// TODO(tjarosik): add static analysis and memory sanitizers for clang
func compileSolution(sourceCodeFile string, executableFile string) (output []byte, err error) {
	cmd := exec.Command("clang++", "-std=c++14", sourceCodeFile, "-o", executableFile)
	log.Println("About to execute command:", cmd.String())
	output, err = cmd.Output()
	if err != nil {
		return output, errors.New("compilation failed with " + err.Error())
	}
	return output, nil
}

// TODO(tjarosik): handle memory limit (-> ulimit -m 100000 && exec ./my-binary)
func runSingleTestCase(executable string, tc TestCase) TestCase {
	inFile, err := os.OpenFile(tc.InputFilename, os.O_RDONLY, 0755)
	defer inFile.Close()
	if err != nil {
		log.Println(err)
		tc.Status = InternalError
		tc.StatusDescription = fmt.Sprintf("unable to open input file at: '%s'", tc.InputFilename)
		return tc
	}
	outFile, err := os.Create(tc.OutputFilename)
	defer outFile.Close()
	if err != nil {
		log.Println(err)
		tc.Status = InternalError
		tc.StatusDescription = fmt.Sprintf("unable to create output file at: '%s'", tc.OutputFilename)
		return tc
	}

	ctx, cancel := context.WithTimeout(context.Background(), tc.TimeLimit)
	defer cancel()
	cmd := exec.CommandContext(ctx, executable)
	cmd.Stdin = inFile
	cmd.Stdout = outFile
	start := time.Now()
	err = cmd.Run()
	tc.Duration = time.Since(start)
	if ctx.Err() == context.DeadlineExceeded {
		tc.Status = TimeLimitExceeded
		tc.StatusDescription = fmt.Sprintf("time limit exceeded: test case was aborted after '%v'", tc.TimeLimit)
		return tc
	}
	if err != nil {
		log.Println(err)
		tc.Status = InternalError
		tc.StatusDescription = fmt.Sprintf("unable to run executable '%s' on input file '%s'", executable, tc.InputFilename)
		return tc
	}

	err = compareFiles(tc.GoldenOutputFilename, tc.OutputFilename)
	if err != nil {
		tc.Status = WrongAnswer
		tc.StatusDescription = err.Error()
	} else {
		tc.Status = Accepted
		tc.StatusDescription = "OK"
	}
	return tc
}

func compareFiles(expected, actual string) error {
	f1, err := os.OpenFile(expected, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer f1.Close()
	f2, err := os.OpenFile(actual, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer f2.Close()

	GB := 1024 * 1024 * 1024
	scanner1 := bufio.NewScanner(f1)
	scanner1.Buffer(make([]byte, 16*1024), GB)
	scanner2 := bufio.NewScanner(f2)
	scanner2.Buffer(make([]byte, 16*1024), GB)
	i := 0
	for scanner1.Scan() {
		scanner2.Scan()
		t1 := strings.TrimRight(scanner1.Text(), "\n\r\t ")
		t2 := strings.TrimRight(scanner2.Text(), "\n\r\t ")
		if t1 != t2 {
			t1Msg := t1
			if len(t1Msg) > 256 {
				t1Msg = t1[:256] + "..."
			}
			t2Msg := t2
			if len(t2Msg) > 256 {
				t2Msg = t2[:256] + "..."
			}
			return fmt.Errorf("files differ in line: %d: expected: %s, actual: %s", i, t1Msg, t2Msg)
		}
		i++
	}
	if scanner2.Scan() {
		t2 := strings.TrimRight(scanner2.Text(), "\n\r\t ")
		if len(t2) > 0 {
			return fmt.Errorf("file %s contains additional non-empty lines", actual)
		}
	}
	return nil
}
