package check

import (
	"bufio"
	"flag"
	"fmt"
	"os"
<<<<<<< HEAD
=======
	"runtime"
>>>>>>> save
	"sync"
	"testing"
	"time"
)

// -----------------------------------------------------------------------
// Test suite registry.

var (
	allParallelSuites []interface{}
	allSerialSuites   []interface{}
	)


// Suite registers the given value as a test suite to be run. Any methods
// starting with the Test prefix in the given value will be considered as
// a test method.
func Suite(suite interface{}) interface{} {
	allParallelSuites = append(allParallelSuites, suite)
	return suite
}

// SerialSuites registers the given value as a test suite to be run serially. Any methods
// starting with the Test prefix in the given value will be considered as
// a test method.
func SerialSuites(suite interface{}) interface{} {
	allSerialSuites = append(allSerialSuites, suite)
	return suite
}

// -----------------------------------------------------------------------
// Public running interface.

var (
	oldFilterFlag  = flag.String("gocheck.f", "", "Regular expression selecting which tests and/or suites to run")
	oldVerboseFlag = flag.Bool("gocheck.v", false, "Verbose mode")
	oldStreamFlag  = flag.Bool("gocheck.vv", false, "Super verbose mode (disables output caching)")
	oldBenchFlag   = flag.Bool("gocheck.b", false, "Run benchmarks")
	oldBenchTime   = flag.Duration("gocheck.btime", 1*time.Second, "approximate run time for each benchmark")
	oldListFlag    = flag.Bool("gocheck.list", false, "List the names of all tests that will be run")
	oldWorkFlag    = flag.Bool("gocheck.work", false, "Display and do not remove the test working directory")

	newFilterFlag  = flag.String("check.f", "", "Regular expression selecting which tests and/or suites to run")
	newVerboseFlag = flag.Bool("check.v", false, "Verbose mode")
	newStreamFlag  = flag.Bool("check.vv", false, "Super verbose mode (disables output caching)")
	newBenchFlag   = flag.Bool("check.b", false, "Run benchmarks")
	newBenchTime   = flag.Duration("check.btime", 1*time.Second, "approximate run time for each benchmark")
	newBenchMem    = flag.Bool("check.bmem", false, "Report memory benchmarks")
	newListFlag    = flag.Bool("check.list", false, "List the names of all tests that will be run")
	newWorkFlag    = flag.Bool("check.work", false, "Display and do not remove the test working directory")
	newExcludeFlag = flag.String("check.exclude", "", "Regular expression to exclude tests to run")

	CustomParallelSuiteFlag = flag.Bool("check.p", false, "Run suites in parallel")
)

var CustomVerboseFlag bool

// TestingT runs all test suites registered with the Suite function,
// printing results to stdout, and reporting any failures back to
// the "testing" package.
func TestingT(testingT *testing.T) {
	benchTime := *newBenchTime
	if benchTime == 1*time.Second {
		benchTime = *oldBenchTime
	}
	conf := &RunConf{
		Filter:        *oldFilterFlag + *newFilterFlag,
		Verbose:       *oldVerboseFlag || *newVerboseFlag || CustomVerboseFlag,
		Stream:        *oldStreamFlag || *newStreamFlag,
		Benchmark:     *oldBenchFlag || *newBenchFlag,
		BenchmarkTime: benchTime,
		BenchmarkMem:  *newBenchMem,
		KeepWorkDir:   *oldWorkFlag || *newWorkFlag,
		Exclude:       *newExcludeFlag,
	}
	if *oldListFlag || *newListFlag {
		w := bufio.NewWriter(os.Stdout)
		for _, name := range ListAll(conf) {
			fmt.Fprintln(w, name)
		}
		w.Flush()
		return
	}
	result := RunAll(conf)
	println(result.String())
	if !result.Passed() {
		testingT.Fail()
	}
}

// RunAll runs all test suites registered with the Suite function, using the
// provided run configuration.
func RunAll(runConf *RunConf) *Result {
	result := Result{}
	if !*CustomParallelSuiteFlag {
		// run all suites serially.
		for _, suite := range allParallelSuites {
			result.Add(Run(suite, runConf))
		}

		for _, suite := range allSerialSuites {
			result.Add(Run(suite, runConf))
		}
		return &result
	}
<<<<<<< HEAD
=======
	r := &runner{concurrency: runtime.NumCPU()*4}
	r.start()
	defer r.stop()
>>>>>>> save

	wg := sync.WaitGroup{}
	notifyRunningSuitesCh := make(chan struct{})
	suiteRunners := make([]*suiteRunner, 0, len(allParallelSuites))
	for _, suite := range allParallelSuites {
<<<<<<< HEAD
		suiteRunners = append(suiteRunners, parallelRun(suite, runConf, &wg, notifyRunningSuitesCh))
=======
		suiteRunners = append(suiteRunners, parallelRun(suite, runConf, r, &wg, notifyRunningSuitesCh))
>>>>>>> save
	}
	close(notifyRunningSuitesCh)
	wg.Wait()
	for _, runner := range suiteRunners {
		ret := &runner.tracker.result
		result.Add(ret)
	}

	for _, suite := range allSerialSuites {
		result.Add(Run(suite, runConf))
	}
	return &result
}

<<<<<<< HEAD
func parallelRun(suite interface{}, runConf *RunConf, wg *sync.WaitGroup, notifyRunningSuitesCh chan struct{}) *suiteRunner {
	runner := newSuiteRunner(suite, runConf)
	runner.asyncRun(wg, notifyRunningSuitesCh)
=======
func parallelRun(suite interface{}, runConf *RunConf, r *runner, wg *sync.WaitGroup, notifyRunningSuitesCh chan struct{}) *suiteRunner {
	runner := newSuiteRunner(suite, runConf)
	runner.asyncRun(r, wg, notifyRunningSuitesCh)
>>>>>>> save
	return runner
}

// Run runs the provided test suite using the provided run configuration.
func Run(suite interface{}, runConf *RunConf) *Result {
	runner := newSuiteRunner(suite, runConf)
	return runner.run()
}

// ListAll returns the names of all the test functions registered with the
// Suite function that will be run with the provided run configuration.
func ListAll(runConf *RunConf) []string {
	var names []string
	for _, suite := range allParallelSuites {
		names = append(names, List(suite, runConf)...)
	}

	for _, suite := range allSerialSuites {
		names = append(names, List(suite, runConf)...)
	}
	return names
}

// List returns the names of the test functions in the given
// suite that will be run with the provided run configuration.
func List(suite interface{}, runConf *RunConf) []string {
	var names []string
	runner := newSuiteRunner(suite, runConf)
	for _, t := range runner.tests {
		names = append(names, t.String())
	}
	return names
}

// -----------------------------------------------------------------------
// Result methods.

func (r *Result) Add(other *Result) {
	r.Succeeded += other.Succeeded
	r.Skipped += other.Skipped
	r.Failed += other.Failed
	r.Panicked += other.Panicked
	r.FixturePanicked += other.FixturePanicked
	r.ExpectedFailures += other.ExpectedFailures
	r.Missed += other.Missed
	if r.WorkDir != "" && other.WorkDir != "" {
		r.WorkDir += ":" + other.WorkDir
	} else if other.WorkDir != "" {
		r.WorkDir = other.WorkDir
	}
}

func (r *Result) Passed() bool {
	return (r.Failed == 0 && r.Panicked == 0 &&
		r.FixturePanicked == 0 && r.Missed == 0 &&
		r.RunError == nil)
}

func (r *Result) String() string {
	if r.RunError != nil {
		return "ERROR: " + r.RunError.Error()
	}

	var value string
	if r.Failed == 0 && r.Panicked == 0 && r.FixturePanicked == 0 &&
		r.Missed == 0 {
		value = "OK: "
	} else {
		value = "OOPS: "
	}
	value += fmt.Sprintf("%d passed", r.Succeeded)
	if r.Skipped != 0 {
		value += fmt.Sprintf(", %d skipped", r.Skipped)
	}
	if r.ExpectedFailures != 0 {
		value += fmt.Sprintf(", %d expected failures", r.ExpectedFailures)
	}
	if r.Failed != 0 {
		value += fmt.Sprintf(", %d FAILED", r.Failed)
	}
	if r.Panicked != 0 {
		value += fmt.Sprintf(", %d PANICKED", r.Panicked)
	}
	if r.FixturePanicked != 0 {
		value += fmt.Sprintf(", %d FIXTURE-PANICKED", r.FixturePanicked)
	}
	if r.Missed != 0 {
		value += fmt.Sprintf(", %d MISSED", r.Missed)
	}
	if r.WorkDir != "" {
		value += "\nWORK=" + r.WorkDir
	}
	return value
}
