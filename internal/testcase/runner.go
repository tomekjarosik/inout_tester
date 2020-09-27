package testcase

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Runner interface {
	Run(executable string, info Info) Result
}

type defaultRunner struct {
	name            string
	streamsProvider StreamsProvider
}

// Info struct describing results of a single test run
type Info struct {
	Name        string        `json:"name"`
	TimeLimit   time.Duration `json:"timeLimit"`
	MemoryLimit int           `json:"memoryLimit"`
}

type Result struct {
	Status      Status        `json:"status"`
	Description string        `json:"description"`
	Duration    time.Duration `json:"duration"`
}

type CompletedTestCase struct {
	Info   Info   `json:"info"`
	Result Result `json:"result"`
}

func NewRunner(name string, streamsProvider StreamsProvider) Runner {
	return &defaultRunner{
		name:            name,
		streamsProvider: streamsProvider}
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
	streams, err := r.streamsProvider(info)
	if err != nil {
		return Result{Status: InternalError, Description: fmt.Sprintf("unable to open data streams, %v", err)}
	}
	defer streams.Close()

	return runTestWithTmpOutput(executable, info, streams)
}

func runTestWithTmpOutput(executable string, info Info, streams Streams) Result {
	tmpStdOutput, err := ioutil.TempFile(os.TempDir(), "tempstd-*.out")
	if err != nil {
		return Result{Status: InternalError, Description: fmt.Sprintf("unable to open temporary output file: %v", err)}
	}
	defer os.Remove(tmpStdOutput.Name())
	defer tmpStdOutput.Close()

	tmpErrorOutput, err := ioutil.TempFile(os.TempDir(), "temperr-*.out")
	if err != nil {
		return Result{Status: InternalError, Description: fmt.Sprintf("unable to open temporary output file: %v", err)}
	}
	defer os.Remove(tmpErrorOutput.Name())
	defer tmpErrorOutput.Close()

	return RunTest(executable, info, streams, tmpStdOutput, tmpErrorOutput)
}

// TODO(tjarosik): handle memory limit (-> ulimit -m 100000 && exec ./my-binary)
func RunTest(executable string, info Info, streams Streams, generatedStdOutput io.ReadWriteSeeker, generatedErrorOutput io.ReadWriteSeeker) Result {
	ctx, cancel := context.WithTimeout(context.Background(), info.TimeLimit)
	defer cancel()
	cmd := exec.CommandContext(ctx, executable)
	cmd.Stdin = streams.Input
	cmd.Stdout = generatedStdOutput
	cmd.Stderr = generatedErrorOutput
	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)
	if ctx.Err() == context.DeadlineExceeded {
		return Result{Status: TimeLimitExceeded,
			Description: fmt.Sprintf("time limit exceeded: test case was aborted after '%v'", info.TimeLimit),
			Duration:    duration}
	}
	if err != nil {
		log.Println(err)
		_, err = generatedErrorOutput.Seek(0, io.SeekStart)
		stderrOutput, err := ioutil.ReadAll(generatedErrorOutput)
		if err != nil {
			stderrOutput = []byte("unable to read stderr")
		}
		return Result{Status: RuntimeError,
			Description: fmt.Sprintf("unable to run executable '%s' on test input file '%s'. Stderr:%s", executable, info.Name, string(stderrOutput)),
			Duration:    duration}
	}
	_, err = generatedStdOutput.Seek(0, io.SeekStart)
	if err != nil {
		return Result{Status: InternalError,
			Description: fmt.Sprintf("unable to rewind generated output for test '%s'", info.Name),
			Duration:    duration}
	}

	err = compare(streams.Output, generatedStdOutput)
	if err != nil {
		return Result{Status: WrongAnswer,
			Description: err.Error(), Duration: duration}
	}

	return Result{Status: Accepted, Description: "OK", Duration: duration}
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
