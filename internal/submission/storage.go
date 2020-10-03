package submission

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"sync"

	testcase "github.com/tomekjarosik/inout_tester/internal/testcase"
)

// Storage Persistent storage for Submissions
type Storage interface {
	Init() error
	Destroy() error // !! destroys all the data stored !!

	Upload(meta Metadata, solution io.Reader) error
	Download(meta Metadata) (solution io.ReadCloser, err error)

	Save(Metadata) error
	Get(id ID) (Metadata, bool)
	Remove(id ID) error

	List() []Metadata
	LoadAll() error
}

// SubmissionStorage object holding data about submissions
type defaultStorage struct {
	data map[string]Metadata

	dataDirectory string
	m             sync.Mutex
}

// NewDefaultSubmissionStorage constructor of default implementation of SubmissionStorage
func NewDefaultStorage(dataDir string) Storage {
	return &defaultStorage{
		data:          make(map[string]Metadata, 0),
		dataDirectory: dataDir,
	}
}

func (store *defaultStorage) Init() error {
	return ensureDirectoryExists(store.dataDirectory)
}

// Upload new SubmissionMetadata object with unique ID
func (store *defaultStorage) Upload(meta Metadata, solution io.Reader) error {
	solutionsDirectory := path.Join(store.dataDirectory, meta.ProblemName)
	if err := ensureDirectoryExists(solutionsDirectory); err != nil {
		return err
	}

	submittedSolutionFile, err := os.Create(path.Join(solutionsDirectory, meta.SolutionFilename))
	if err != nil {
		return err
	}
	defer submittedSolutionFile.Close()

	_, err = io.Copy(submittedSolutionFile, solution)
	if err != nil {
		return err
	}
	if err = store.Save(meta); err != nil {
		return err
	}
	return nil
}

func (store *defaultStorage) Download(meta Metadata) (solution io.ReadCloser, err error) {
	solutionFile, err := os.Open(path.Join(store.dataDirectory, meta.ProblemName, meta.SolutionFilename))
	if err != nil {
		return nil, err
	}
	return solutionFile, nil
}

func (store *defaultStorage) Save(metadata Metadata) error {
	store.m.Lock()
	defer store.m.Unlock()
	store.data[metadata.ID.String()] = metadata
	metadataFile, err := os.Create(path.Join(store.dataDirectory, metadata.ID.String()+metaFileExtension))
	if err != nil {
		return err
	}
	defer metadataFile.Close()
	enc := json.NewEncoder(metadataFile)
	enc.SetIndent("", "\t")
	return enc.Encode(metadata)
}

// ByTimestamp is a helper type to implement sorting
type ByTimestamp []Metadata

func (a ByTimestamp) Len() int           { return len(a) }
func (a ByTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimestamp) Less(i, j int) bool { return a[i].SubmittedAt.After(a[j].SubmittedAt) }

func (store *defaultStorage) List() []Metadata {
	store.m.Lock()
	defer store.m.Unlock()
	res := make([]Metadata, 0)
	for _, elem := range store.data {
		res = append(res, elem)
	}
	sort.Sort(ByTimestamp(res))
	return res
}

func (store *defaultStorage) Remove(id ID) error {
	store.m.Lock()
	defer store.m.Unlock()
	delete(store.data, id.String())
	return os.Remove(path.Join(store.dataDirectory, id.String()+metaFileExtension))
}

func (store *defaultStorage) LoadAll() error {

	files, err := ioutil.ReadDir(store.dataDirectory)
	if err != nil {
		return err
	}
	store.m.Lock()
	defer store.m.Unlock()

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), metaFileExtension) {
			continue
		}
		metadataFile, err := os.Open(path.Join(store.dataDirectory, f.Name()))
		if err != nil {
			return err
		}
		defer metadataFile.Close()
		var metadata Metadata
		if err = json.NewDecoder(metadataFile).Decode(&metadata); err != nil {
			log.Println("Json deconding failed")
			return err
		}
		store.data[metadata.ID.String()] = metadata
	}
	log.Printf("Loaded %d submissions into memory\n", len(store.data))
	return nil
}

func (store *defaultStorage) Destroy() error {
	return os.RemoveAll(store.dataDirectory)
}

func (store *defaultStorage) Get(id ID) (Metadata, bool) {
	v, found := store.data[id.String()]
	if !found {
		return Metadata{}, false
	}
	return v, true
}

func (meta Metadata) Score() int {
	res := 0
	for _, tc := range meta.CompletedTestCases {
		if tc.Result.Status == testcase.Accepted {
			res++
		}
	}
	return res
}

func (meta Metadata) MaxScore() int {
	return len(meta.CompletedTestCases)
}
