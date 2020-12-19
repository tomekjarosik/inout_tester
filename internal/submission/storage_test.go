package submission

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	guuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	testcase "github.com/tomekjarosik/inout_tester/internal/testcase"
)

func readFile(t *testing.T, filename string) string {
	res, err := ioutil.ReadFile(filename)
	assert.NoError(t, err)
	return string(res)
}

func TestDefaultStorage_Upload(t *testing.T) {
	tmpstoragedir := "tmpstoragedir"
	defer os.RemoveAll(tmpstoragedir)

	sp := NewDefaultStorage(tmpstoragedir)
	assert.NoError(t, sp.Init())

	content := "this is a solution"
	m := NewMetadata("testproblem", testcase.ReleaseMode)
	err := sp.Upload(m, strings.NewReader(content))

	assert.NoError(t, err)
	assert.Equal(t, "testproblem", m.ProblemName)
	assert.Equal(t, content, readFile(t, path.Join(tmpstoragedir, "testproblem", m.SolutionFilename)))
	assert.NoError(t, err)
	assert.Equal(t, 0, m.AcceptedCount)
	assert.Equal(t, 0, m.TestCasesCount)
	assert.Equal(t, Queued, m.Status)
}

func TestDefaultStorage_SaveGet(t *testing.T) {
	tmpstoragedir := "tmpstoragedir"
	defer os.RemoveAll(tmpstoragedir)

	sp := NewDefaultStorage(tmpstoragedir)
	assert.NoError(t, sp.Init())

	m := Metadata{
		ID:                 ID(guuid.New()),
		SubmittedAt:        time.Now(),
		SolutionFilename:   "sol.cpp",
		Status:             AllTestsCompleted,
		ProblemName:        "aProblem",
		ExecutableFilename: "a.tsk",
		CompilationMode:    testcase.ReleaseMode,
	}
	err := sp.Save(m)
	assert.NoError(t, err)

	metaRetrieved, exists := sp.Get(m.ID)
	assert.True(t, exists)
	assert.Equal(t, m.ID, metaRetrieved.ID)
	assert.Equal(t, "aProblem", metaRetrieved.ProblemName)
	assert.Equal(t, "a.tsk", metaRetrieved.ExecutableFilename)
	assert.Equal(t, AllTestsCompleted, metaRetrieved.Status)
	assert.Equal(t, "sol.cpp", metaRetrieved.SolutionFilename)
	assert.Equal(t, testcase.ReleaseMode, metaRetrieved.CompilationMode)
}

func TestDefaultStorage_List(t *testing.T) {
	tmpstoragedir := "tmpstoragedir"
	defer os.RemoveAll(tmpstoragedir)

	sp := NewDefaultStorage(tmpstoragedir)
	assert.NoError(t, sp.Init())

	for i := 0; i < 5; i++ {
		m := Metadata{
			ID:                 ID(guuid.New()),
			SubmittedAt:        time.Now(),
			SolutionFilename:   fmt.Sprintf("sol%d.cpp", i),
			Status:             AllTestsCompleted,
			ProblemName:        "aProblem",
			ExecutableFilename: fmt.Sprintf("a%d.tsk", i),
		}
		err := sp.Save(m)
		assert.NoError(t, err)
	}
	list := sp.List()
	assert.Equal(t, 5, len(list))

	assert.Equal(t, "sol4.cpp", list[0].SolutionFilename)
	assert.Equal(t, "sol3.cpp", list[1].SolutionFilename)
	assert.Equal(t, "sol2.cpp", list[2].SolutionFilename)
	assert.Equal(t, "sol1.cpp", list[3].SolutionFilename)
	assert.Equal(t, "sol0.cpp", list[4].SolutionFilename)
}

func TestDefaultStorage_LoadAll(t *testing.T) {
	tmpstoragedir := "tmpstoragedir"
	defer os.RemoveAll(tmpstoragedir)

	sp := NewDefaultStorage(tmpstoragedir)
	assert.NoError(t, sp.Init())

	for i := 0; i < 5; i++ {
		m := Metadata{
			ID:                 ID(guuid.New()),
			SubmittedAt:        time.Now(),
			SolutionFilename:   fmt.Sprintf("sol%d.cpp", i),
			Status:             AllTestsCompleted,
			ProblemName:        "aProblem",
			ExecutableFilename: fmt.Sprintf("a%d.tsk", i),
		}
		err := sp.Save(m)
		assert.NoError(t, err)
	}

	sp2 := NewDefaultStorage(tmpstoragedir)
	assert.NoError(t, sp2.Init())

	err := sp2.LoadAll()
	assert.NoError(t, err)

	list2 := sp2.List()
	assert.Equal(t, 5, len(list2))
	assert.Equal(t, "sol4.cpp", list2[0].SolutionFilename)
	assert.Equal(t, "sol3.cpp", list2[1].SolutionFilename)
	assert.Equal(t, "sol2.cpp", list2[2].SolutionFilename)
	assert.Equal(t, "sol1.cpp", list2[3].SolutionFilename)
	assert.Equal(t, "sol0.cpp", list2[4].SolutionFilename)
	assert.Equal(t, AllTestsCompleted, list2[4].Status)
}

func TestDefaultStorage_Remove(t *testing.T) {
	tmpstoragedir := "tmpstoragedir"
	defer os.RemoveAll(tmpstoragedir)

	sp := NewDefaultStorage(tmpstoragedir)
	assert.NoError(t, sp.Init())

	meta := make([]Metadata, 3)
	for i := 0; i < 3; i++ {
		meta[i] = Metadata{
			ID:                 ID(guuid.New()),
			SubmittedAt:        time.Now(),
			SolutionFilename:   fmt.Sprintf("sol%d.cpp", i),
			Status:             AllTestsCompleted,
			ProblemName:        "aProblem",
			ExecutableFilename: fmt.Sprintf("a%d.tsk", i),
		}
		err := sp.Save(meta[i])
		assert.NoError(t, err)
	}

	sp.Remove(meta[1].ID)

	list := sp.List()
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "sol2.cpp", list[0].SolutionFilename)
	assert.Equal(t, "sol0.cpp", list[1].SolutionFilename)
}
