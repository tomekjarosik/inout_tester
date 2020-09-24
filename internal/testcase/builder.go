package testcase

import (
	"errors"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"
)

// CompilationMode compilation mode
type CompilationMode int

const (
	// Release this is release compilation mode with optimization
	ReleaseMode CompilationMode = iota + 1
	// AnalyzeClang this is mode which maximizes possibility of finding bugs
	AnalyzeModeClang
	//AnalyzeGplusplus
	AnalyzeModeGplusplus
)

func CompileSolution(sourceCodeFile string, mode CompilationMode, executableFile string) (output []byte, err error) {
	var cmd *exec.Cmd
	switch mode {
	case ReleaseMode:
		cmd = exec.Command("g++", "-std=c++17", "-static", "-O3", sourceCodeFile, "-lm", "-o", executableFile)
	case AnalyzeModeClang:
		cmd = exec.Command("clang++", "-std=c++14", "-Wall", "-O1", "-g", "-fsanitize=address",
			"-fno-omit-frame-pointer", sourceCodeFile, "-lm", "-o", executableFile)
	case AnalyzeModeGplusplus:
		cmd = exec.Command("g++", "-std=c++17", "-Wall", "-O1", "-g", "-fsanitize=address",
			"-fno-omit-frame-pointer", sourceCodeFile, "-lm", "-o", executableFile)
	default:
		return []byte{}, errors.New("unknown compilation mode selected")
	}

	log.Println("About to execute command:", cmd.String())
	output, err = cmd.CombinedOutput()
	if err != nil {
		if output != nil {
			log.Println(string(output))
		}
		return output, errors.New("compilation failed with " + err.Error())
	}
	return output, nil
}

// PopulateTestCases searches directory for test case descriptions (.in / .out files, maybe others in the future)
// TODO: implement testcase metadata, which knows about timelimits and memory limits
func Populate(problemDataDir string) (testcases []Info, err error) {

	files, err := ioutil.ReadDir(problemDataDir)
	if err != nil {
		return
	}
	const ext = ".in"
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ext) {
			continue
		}
		tc := NewInfo(strings.TrimSuffix(f.Name(), ext), 10*time.Second, 0)
		testcases = append(testcases, tc)
	}
	return
}
