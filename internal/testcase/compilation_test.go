package testcase

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TempFileName generates a temporary filename for use in testing or whatever
func TempFileName(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
}

func TestCompilation_ReleaseMode_CorrectFile(t *testing.T) {
	solution := strings.NewReader(`#include <cstdio>
	int main() { printf("OK!"); return 0; }`)
	tmpFileName := TempFileName("testcase", ".tsk")
	out, err := CompileSolution(solution, ReleaseMode, tmpFileName)
	assert.NoError(t, err)
	assert.Equal(t, "", string(out))
	os.Remove(tmpFileName)
}

func TestCompilation_ReleaseMode_SyntaxError(t *testing.T) {
	solution := strings.NewReader(`#include <cstdio>
	int main() { xxx return 0; }`)
	tmpFileName := TempFileName("testcase", ".tsk")
	out, err := CompileSolution(solution, ReleaseMode, tmpFileName)
	assert.Error(t, err)
	assert.EqualError(t, err, "compilation failed with exit status 1")
	assert.Contains(t, string(out), "<stdin>:2:15: error: 'xxx' was not declared in this scope")
	os.Remove(tmpFileName)
}

func TestCompilationMode_FullCommandFor(t *testing.T) {
	s := FullCompilationCommadFor(ReleaseMode)
	// NOTE: full command will include extension and full path to the dir, so don't check if they are exactly the same
	assert.Contains(t, s, "-std=c++17 -static -O3 -x c++ - -lm -o a.out")
	assert.Contains(t, s, "g++")

	s = FullCompilationCommadFor(AnalyzeClangMode)
	assert.Contains(t, s, "-std=c++14 -Wall -Werror -O1 -g -fsanitize=address -fno-omit-frame-pointer -x c++ - -lm -o a.out")
	assert.Contains(t, s, "clang++")

	s = FullCompilationCommadFor(AnalyzeGplusplusMode)
	assert.Contains(t, s, "-std=c++17 -Wall -Werror -O1 -g -fsanitize=address -fno-omit-frame-pointer -x c++ - -lm -o a.out")
	assert.Contains(t, s, "g++")
}

func TestCompilatonMode_UnmarshallJSON(t *testing.T) {
	var cm CompilationMode
	err := cm.UnmarshalJSON([]byte("\"ReleaseMode\""))
	assert.NoError(t, err)
	assert.Equal(t, ReleaseMode, cm)

	err = cm.UnmarshalJSON([]byte("\"AnalyzeClangMode\""))
	assert.NoError(t, err)
	assert.Equal(t, AnalyzeClangMode, cm)

	err = cm.UnmarshalJSON([]byte("\"AnalyzeGplusplusMode\""))
	assert.NoError(t, err)
	assert.Equal(t, AnalyzeGplusplusMode, cm)
}

func TestCompilatonMode_MarshallJSON(t *testing.T) {
	out, err := ReleaseMode.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte("\"ReleaseMode\""), out)

	out, err = AnalyzeClangMode.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte("\"AnalyzeClangMode\""), out)

	out, err = AnalyzeGplusplusMode.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte("\"AnalyzeGplusplusMode\""), out)
}
