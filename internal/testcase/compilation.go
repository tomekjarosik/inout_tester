package testcase

//go:generate stringer -type=CompilationMode

import (
	"encoding/json"
	"errors"
	"log"
	"os/exec"
)

// CompilationMode compilation mode
type CompilationMode int

const (
	// Release this is release compilation mode with optimization
	ReleaseMode CompilationMode = iota + 1
	// AnalyzeClang this is mode which maximizes possibility of finding bugs
	AnalyzeClangMode
	//AnalyzeGplusplus
	AnalyzeGplusplusMode
)

func CompilationCommand(sourceCodeFile string, mode CompilationMode, executableFile string) (*exec.Cmd, error) {
	var cmd *exec.Cmd
	switch mode {
	case ReleaseMode:
		cmd = exec.Command("g++", "-std=c++17", "-static", "-O3", sourceCodeFile, "-lm", "-o", executableFile)
	case AnalyzeClangMode:
		cmd = exec.Command("clang++", "-std=c++14", "-Wall", "-O1", "-g", "-fsanitize=address",
			"-fno-omit-frame-pointer", sourceCodeFile, "-lm", "-o", executableFile)
	case AnalyzeGplusplusMode:
		cmd = exec.Command("g++", "-std=c++17", "-Wall", "-O1", "-g", "-fsanitize=address",
			"-fno-omit-frame-pointer", sourceCodeFile, "-lm", "-o", executableFile)
	default:
		return nil, errors.New("unknown compilation mode selected")
	}
	return cmd, nil
}

func CompileSolution(sourceCodeFile string, mode CompilationMode, executableFile string) (output []byte, err error) {
	cmd, err := CompilationCommand(sourceCodeFile, mode, executableFile)
	if err != nil {
		return []byte{}, err
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

func FullCompilationCommadFor(cm CompilationMode) string {
	cmd, err := CompilationCommand("a.cpp", cm, "a.out")
	if err != nil {
		return "unable to convert"
	}
	return cmd.String()
}

func (cm *CompilationMode) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	for i := 0; i <= len(_CompilationMode_index); i++ {
		if CompilationMode(i).String() == s {
			*cm = CompilationMode(i)
			return nil
		}
	}
	return errors.New("invalid CompilationMode status value")
}

func (cm CompilationMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(cm.String())
}
