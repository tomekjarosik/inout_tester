package main

//go:generate stringer -type=SubmissionState

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"time"

	guuid "github.com/google/uuid"
	testcase "github.com/tomekjarosik/inout_tester/internal/testcase"
)

// SubmissionState submission state
type SubmissionState int

const (
	// Queued the solution is queued for processing
	Queued SubmissionState = iota
	// Compiling the solution is being processed
	Compiling
	// CompilationError the solution failed to compile
	CompilationError
	// RunningTests currently running the provided test cases
	RunningTests
	// RunAllTests all done
	RunAllTests
)

// SubmissionMetadata metadata of the submission
type SubmissionMetadata struct {
	ID                 guuid.UUID                   `json:"id"`
	SubmittedAt        time.Time                    `JSON:"submittedAt"`
	ProblemName        string                       `json:"problemName"`
	SolutionFilename   string                       `json:"solutionFilename"`
	State              SubmissionState              `json:"state"`
	ExecutableFilename string                       `json:"executableFilename"`
	CompilationOutput  []byte                       `json:"compilationOutput"`
	CompilationMode    testcase.CompilationMode     `json:"compilationMode"`
	CompletedTestCases []testcase.CompletedTestCase `json:"testCases"`
}

type SubmissionStorage interface {
	Upload(problemName string, solution io.Reader) (SubmissionMetadata, error)
	Save(SubmissionMetadata) error
	Remove(id guuid.UUID) error
	Get(id guuid.UUID) (SubmissionMetadata, error)
	List() ([]SubmissionMetadata, error)
	LoadAll() error
	RootDir() string
}

// ByTimestampt is a helper type to implement sorting
type ByTimestamp []SubmissionMetadata

func (a ByTimestamp) Len() int           { return len(a) }
func (a ByTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimestamp) Less(i, j int) bool { return a[i].SubmittedAt.After(a[j].SubmittedAt) }

// Upload new SubmissionMetadata object with unique ID
func (store *defaultSubmissionStorage) Upload(problemName string, solution io.Reader) (SubmissionMetadata, error) {
	id := guuid.New()
	metadata := SubmissionMetadata{
		ID:                 id,
		SubmittedAt:        time.Now(),
		SolutionFilename:   id.String() + ".cpp",
		State:              Queued,
		ProblemName:        problemName,
		ExecutableFilename: id.String() + ".tsk",
	}

	solutionsDirectory := path.Join(store.RootDir(), problemName)

	submissionFile, err := os.Create(path.Join(solutionsDirectory, metadata.SolutionFilename))
	if err != nil {
		return metadata, err
	}
	defer submissionFile.Close()

	_, err = io.Copy(submissionFile, solution)
	if err != nil {
		return metadata, err
	}
	if err = store.Save(metadata); err != nil {
		return metadata, err
	}
	return metadata, nil
}

// SubmissionStorage object holding data about submissions
type defaultSubmissionStorage struct {
	data map[string]SubmissionMetadata

	submissionsDirectory string
	m                    sync.Mutex
}

// NewDefaultSubmissionStorage constructor of default implementation of SubmissionStorage
func NewDefaultSubmissionStorage() SubmissionStorage {
	return &defaultSubmissionStorage{
		data:                 make(map[string]SubmissionMetadata, 0),
		submissionsDirectory: "submissions",
	}
}

func ensureDirectoryExists(dir string) error {
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func (store *defaultSubmissionStorage) ensureDirectoryStructure() error {
	if err := ensureDirectoryExists(store.submissionsDirectory); err != nil {
		return err
	}
	return nil
}

func (store *defaultSubmissionStorage) Save(metadata SubmissionMetadata) error {
	store.m.Lock()
	defer store.m.Unlock()
	store.data[metadata.ID.String()] = metadata
	metadataFile, err := os.Create(path.Join(store.submissionsDirectory, metadata.ID.String()+".meta"))
	if err != nil {
		return err
	}
	defer metadataFile.Close()
	enc := json.NewEncoder(metadataFile)
	enc.SetIndent("", "\t")
	return enc.Encode(metadata)
}

func (store *defaultSubmissionStorage) List() ([]SubmissionMetadata, error) {
	store.m.Lock()
	defer store.m.Unlock()
	res := make([]SubmissionMetadata, 0)
	for _, elem := range store.data {
		res = append(res, elem)
	}
	sort.Sort(ByTimestamp(res))
	return res, nil
}
func (store *defaultSubmissionStorage) Remove(id guuid.UUID) error {
	//return os.Remove(path.Join(directory, metadata.ID))
	return errors.New("Not implemented")
}

func (store *defaultSubmissionStorage) LoadAll() error {
	files, err := ioutil.ReadDir(store.submissionsDirectory)
	if err != nil {
		return err
	}
	store.m.Lock()
	defer store.m.Unlock()

	const ext = ".meta"
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ext) {
			continue
		}
		log.Printf("loading %s\n", f.Name())
		metadataFile, err := os.Open(path.Join(store.submissionsDirectory, f.Name()))
		if err != nil {
			return err
		}
		var metadata SubmissionMetadata
		if err = json.NewDecoder(metadataFile).Decode(&metadata); err != nil {
			log.Println("Json deconding failed")
			return err
		}
		store.data[metadata.ID.String()] = metadata
	}
	log.Printf("Loaded %d submissions into memory\n", len(store.data))
	return nil
}

func (store *defaultSubmissionStorage) RootDir() string {
	return store.submissionsDirectory
}

func (store *defaultSubmissionStorage) Get(id guuid.UUID) (SubmissionMetadata, error) {
	v, found := store.data[id.String()]
	if !found {
		return SubmissionMetadata{}, fmt.Errorf("item %v not found", id)
	}
	return v, nil
}

func (meta SubmissionMetadata) Score() int {
	res := 0
	for _, tc := range meta.CompletedTestCases {
		if tc.Result.Status == testcase.Accepted {
			res++
		}
	}
	return res
}
