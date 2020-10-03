package submission

import (
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"time"

	testcase "github.com/tomekjarosik/inout_tester/internal/testcase"
)

// Processor processes submissions
type Processor interface {
	Submit(meta Metadata)
	Process() error
	Quit()
}

type defaultProcessor struct {
	queue           chan Metadata
	store           Storage
	testcaseArchive testcase.Archive
	workersCount    int
}

// NewProcessor constructor of the Processor
func NewProcessor(store Storage, testcaseArchive testcase.Archive) Processor {
	return &defaultProcessor{
		queue:           make(chan Metadata, 1000),
		store:           store,
		testcaseArchive: testcaseArchive,
	}
}

func (p *defaultProcessor) Submit(meta Metadata) {
	p.queue <- meta
}

func testcaseProcessor(runner testcase.Runner, executable string, jobs <-chan testcase.Info, results chan<- testcase.CompletedTestCase) {
	for tc := range jobs {
		results <- testcase.CompletedTestCase{Info: tc, Result: runner.Run(executable, tc)}
	}
	log.Println("worker exited")
}

func (p *defaultProcessor) processSubmission(submission Metadata) (res Metadata, err error) {
	fmt.Println("Processing submission:", submission)
	start := time.Now()
	submission.Status = Compiling
	p.store.Save(submission)

	solution, err := p.store.Download(submission)
	defer solution.Close()

	executable := path.Join(os.TempDir(), submission.ProblemName+"-"+submission.ID.String()+".out")
	defer os.Remove(executable)

	submission.CompilationOutput, err = testcase.CompileSolution(solution, submission.CompilationMode, executable)

	if err != nil {
		submission.Status = CompilationError
		p.store.Save(submission)
		return submission, err
	}

	submission.Status = RunningTests
	p.store.Save(submission)

	testcases, err := p.testcaseArchive.Testcases(submission.ProblemName)
	if err != nil {
		return submission, err
	}

	runner := p.testcaseArchive.Runner(submission.ProblemName)

	// Put all TestCases into buffered channel
	infoChan := make(chan testcase.Info, len(testcases))
	for _, tc := range testcases {
		infoChan <- tc
	}
	close(infoChan)

	resultChan := make(chan testcase.CompletedTestCase, len(testcases))
	for i := 0; i < submission.WorkerCount; i++ {
		go testcaseProcessor(runner, executable, infoChan, resultChan)
	}

	processedTestCases := make([]testcase.CompletedTestCase, 0)
	for i := 0; i < len(testcases); i++ {
		completedTc := <-resultChan
		processedTestCases = append(processedTestCases, completedTc)
		// TODO: Implement saving submission using channels not mutex
		sort.Sort(testcase.ByTestcaseName(processedTestCases))
		submission.CompletedTestCases = processedTestCases
		p.store.Save(submission)
	}

	submission.Status = AllTestsCompleted
	submission.TotalProcessingTime = time.Since(start)
	err = p.store.Save(submission)
	log.Println("Processed submission", submission)
	return submission, err
}

func (p *defaultProcessor) Process() error {
	if err := p.store.LoadAll(); err != nil {
		log.Panic(err)
	}
	for submission := range p.queue {
		_, err := p.processSubmission(submission)
		if err != nil {
			log.Println("ProcessSubmission returned error: ", err)
		}
	}
	fmt.Println("defaultSubmissionProcessor has exited successfully.")
	return nil
}

func (p *defaultProcessor) Quit() {
	close(p.queue)
}
