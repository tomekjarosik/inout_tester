package testcase

//go:generate stringer -type=Status

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

// Status status of the
type Status int

const (
	// NotRunYet test is waiting to be run
	NotRunYet Status = iota
	// InternalError something unexpected went wrong
	InternalError
	// TimeLimitExceeded the test took too long to process
	TimeLimitExceeded
	// MemoryLimitExceeded the test run used too much RAM
	MemoryLimitExceeded
	// WrongAnswer test run successfully but test outputs differ
	WrongAnswer
	// Accepted all OK
	Accepted
)

type Runner interface {
	Run(executable string, info Info) Result
}

type defaultRunner struct {
	name    string
	rootDir string
}

// Streams provide data for a test case
type Streams struct {
	Input        io.Reader
	Output       io.ReadWriteSeeker
	GoldenOutput io.Reader
	Close        func() error
}

// Info struct describing results of a single test run
type Info struct {
	Name        string        `json:"name"`
	TimeLimit   time.Duration `json:"timeLimit"`
	MemoryLimit int           `json:"memoryLimit"`
}

type Result struct {
	Status            Status        `json:"status"`
	StatusDescription string        `json:"statusDescription"`
	Duration          time.Duration `json:"duration"`
}

type CompletedTestCase struct {
	Info   Info   `json:"info"`
	Result Result `json:"result"`
}

func NewRunner(name string, problemDataDir string) Runner {
	return &defaultRunner{name: name, rootDir: problemDataDir}
}

// NewTestCase construct of TestCase struct
func NewInfo(testName string, timeLimit time.Duration, memoryLimit int) Info {
	return Info{
		Name:        testName,
		TimeLimit:   timeLimit,
		MemoryLimit: memoryLimit,
	}
}

func (r *defaultRunner) Run(executable string, info Info) Result {
	streams, err := r.provideStreamsFor(info)
	if err != nil {
		return Result{Status: InternalError, StatusDescription: fmt.Sprintf("unable to open data streams, %v", err)}
	}
	defer streams.Close()
	return internalRunSingleTestCase(executable, info, streams)
}

func (r *defaultRunner) provideStreamsFor(info Info) (Streams, error) {
	streams := Streams{}
	inFile, err := os.OpenFile(path.Join(r.rootDir, info.Name+".in"), os.O_RDONLY, 0755)
	if err != nil {
		return streams, err
	}
	goldenOutFile, err := os.OpenFile(path.Join(r.rootDir, info.Name+".out"), os.O_RDONLY, 0755)
	if err != nil {
		return streams, err
	}
	tmpOutput, err := ioutil.TempFile(os.TempDir(), "temp-*.out")
	if err != nil {
		return streams, err
	}
	streams.Input = inFile
	streams.Output = tmpOutput
	streams.GoldenOutput = goldenOutFile
	streams.Close = func() error {
		defer inFile.Close()
		defer goldenOutFile.Close()
		defer os.Remove(tmpOutput.Name())
		// TODO: Handle closing better
		return nil
	}
	return streams, nil
}

// TODO(tjarosik): add static analysis and memory sanitizers for clang
func CompileSolution(sourceCodeFile string, executableFile string) (output []byte, err error) {
	cmd := exec.Command("clang++", "-std=c++14", sourceCodeFile, "-o", executableFile)
	log.Println("About to execute command:", cmd.String())
	output, err = cmd.Output()
	if err != nil {
		return output, errors.New("compilation failed with " + err.Error())
	}
	return output, nil
}

// TODO(tjarosik): handle memory limit (-> ulimit -m 100000 && exec ./my-binary)
func internalRunSingleTestCase(executable string, info Info, streams Streams) Result {
	res := Result{}

	ctx, cancel := context.WithTimeout(context.Background(), info.TimeLimit)
	defer cancel()
	cmd := exec.CommandContext(ctx, executable)
	cmd.Stdin = streams.Input
	cmd.Stdout = streams.Output
	start := time.Now()
	err := cmd.Run()
	res.Duration = time.Since(start)
	if ctx.Err() == context.DeadlineExceeded {
		res.Status = TimeLimitExceeded
		res.StatusDescription = fmt.Sprintf("time limit exceeded: test case was aborted after '%v'", info.TimeLimit)
		return res
	}
	if err != nil {
		log.Println(err)
		res.Status = InternalError
		res.StatusDescription = fmt.Sprintf("unable to run executable '%s' on test input file '%s'", executable, info.Name)
		return res
	}
	streams.Output.Seek(0, io.SeekStart)

	err = compare(streams.GoldenOutput, streams.Output)
	if err != nil {
		res.Status = WrongAnswer
		res.StatusDescription = err.Error()
	} else {
		res.Status = Accepted
		res.StatusDescription = "OK"
	}
	return res
}

func compare(expected, actual io.Reader) error {

	GB := 1024 * 1024 * 1024 // max memory 1GB
	scanner1 := bufio.NewScanner(expected)
	scanner1.Buffer(make([]byte, 16*1024), 1*GB)
	scanner2 := bufio.NewScanner(actual)
	scanner2.Buffer(make([]byte, 16*1024), 1*GB)
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
			return fmt.Errorf("outputs differ in line %d: expected: '%s', actual: '%s'", i, t1Msg, t2Msg)
		}
		i++
	}
	if scanner2.Scan() {
		t2 := strings.TrimRight(scanner2.Text(), "\n\r\t ")
		if len(t2) > 0 {
			return fmt.Errorf("contains additional non-empty lines")
		}
	}
	return nil
}
