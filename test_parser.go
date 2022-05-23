package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	actionPass   = "pass"
	actionFail   = "fail"
	actionOutput = "output"
	actionSkip   = "skip"

	moduleName = "github.com/Jimeux/ga-summary"
)

var (
	source io.Reader = os.Stdin
	dest   io.Writer = os.Stdout
	// regex to match a filename in an output event
	fileNameExp, _ = regexp.Compile(`^\s*(\w+_test.go):\d+:`)

	// store test name -> file path mapping. Add file path to failedTests or delete when result is known
	testFilePath map[string]string
	// store all events, and delete key when non-fail result (skip/pass) is known
	testEvents map[string][]event
	// store fail tests names together with their file paths
	failedTests map[string]string
	// store pass test names for filtering/count purposes
	passedTests map[string]struct{}
	// store skip test names for filtering/count purposes
	skippedTests map[string]struct{}
	// count event parse failures
	// Deprecated: remove when functionality is properly validated
	parseErrCount int
)

// event holds data from a single test event.
// See cmd/test2json.
type event struct {
	Time    time.Time `json:"time"`
	Action  string    `json:"action"`
	Package string    `json:"package"`
	Test    string    `json:"test"`
	Elapsed float64   `json:"elapsed"`
	Output  string    `json:"output"`
}

func main() {
	initData()
	scan := bufio.NewScanner(source)
	for scan.Scan() {
		handleEvent(scan.Bytes())
	}
	writeMarkdown(dest)
}

// initData initializes data maps (exists for testing purposes).
func initData() {
	testFilePath = make(map[string]string, 10)
	testEvents = make(map[string][]event, 10)
	failedTests = make(map[string]string, 10)
	passedTests = make(map[string]struct{}, 100)
	skippedTests = make(map[string]struct{}, 100)
	parseErrCount = 0
}

func handleEvent(b []byte) {
	var event event
	if err := json.Unmarshal(b, &event); err != nil {
		parseErrCount++
		return
	}
	if event.Test == "" {
		return
	}
	if _, ok := passedTests[event.Test]; ok {
		return
	}
	if _, ok := skippedTests[event.Test]; ok {
		return
	}
	// output events come before pass/fail events, so insert all events into testEvents
	testEvents[event.Test] = append(testEvents[event.Test], event)

	switch event.Action {
	case actionFail:
		// fail action comes last, so it's safe to get/set file name
		failedTests[event.Test] = testFilePath[event.Test]
	case actionSkip:
		skippedTests[event.Test] = struct{}{}
		delete(testEvents, event.Test)
		delete(testFilePath, event.Test)
	case actionPass:
		passedTests[event.Test] = struct{}{}
		delete(testEvents, event.Test)
		delete(testFilePath, event.Test)
	case actionOutput:
		// file names are not readily available, so search for them in output
		matches := fileNameExp.FindStringSubmatch(event.Output)
		if len(matches) > 1 {
			testFilePath[event.Test] = strings.TrimPrefix(event.Package+"/"+matches[1], moduleName+"/")
		}
	}
}

func writeMarkdown(w io.Writer) {
	_, _ = fmt.Fprint(w, summaryTable())
	if len(failedTests) == 0 {
		return
	}
	_, _ = fmt.Fprintf(w, `## Run Failed Tests Locally

`+"```"+`bash
go test ./... -run '%s'
`+"```"+`

## Failure Details
`, runLocalCommand())

	fileTests := testsByFile()
	for _, filename := range sortedKeys(fileTests) {
		tests := fileTests[filename]
		sort.Slice(tests, func(i, j int) bool { return tests[i] < tests[j] })
		_, _ = fmt.Fprint(w, "---\n\n#### `"+filename+"`\n\n")

		for _, testName := range tests {
			_, _ = fmt.Fprint(w, "<details>\n<summary>"+testName+"</summary>\n\n```diff\n")
			for _, event := range testEvents[testName] {
				output := strings.Trim(event.Output, " \n")
				if output == "" ||
					strings.HasPrefix(output, "=== RUN") ||
					strings.HasPrefix(output, "--- FAIL") {
					continue
				}
				_, _ = fmt.Fprintf(w, output+"\n")
			}
			_, _ = fmt.Fprintf(w, "```\n\n</details>\n\n")
		}
	}
}

func summaryTable() string {
	return fmt.Sprintf(`# Test Summary

|     Status      | Count |
|-----------------|-------|
| ✅ Passed       | %d   |
| ❌ Failed       | %d   |
| ⏩ Skipped      | %d   |
| 💥 Parse Errors | %d   |

`, len(passedTests), len(failedTests), len(skippedTests), parseErrCount)
}

func testsByFile() map[string][]string {
	fileTests := make(map[string][]string, 10)
	for test, file := range testFilePath {
		if _, ok := failedTests[test]; ok {
			fileTests[file] = append(fileTests[file], test)
		}
	}
	return fileTests
}

func runLocalCommand() string {
	testNames := sortedKeys(failedTests)
	parentsOnly := make([]string, 0, len(testNames))
	for _, name := range testNames {
		if !strings.Contains(name, "/") { // exclude sub-tests
			parentsOnly = append(parentsOnly, name)
		}
	}
	return strings.Join(parentsOnly, "|")
}

func sortedKeys[T any](m map[string]T) []string {
	ord := make([]string, 0, len(m))
	for key := range m {
		ord = append(ord, key)
	}
	sort.Slice(ord, func(i, j int) bool { return ord[i] < ord[j] })
	return ord
}
