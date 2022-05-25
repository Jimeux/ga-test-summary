package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

	moduleName    = "github.com/Jimeux/ga-summary"
	maxOutputSize = 1_000_000 // $GITHUB_STEP_SUMMARY has a 1024K limit.
)

type (
	// Event holds data from a single test event.
	// See cmd/test2json.
	Event struct {
		Time    time.Time `json:"time"`
		Action  string    `json:"action"`
		Package string    `json:"package"`
		Test    string    `json:"test"`
		Elapsed float64   `json:"elapsed"`
		Output  string    `json:"output"`
	}
	// Test is used to collect and store data for a specific Go test.
	Test struct {
		Package string
		File    string
		Name    string
		Output  []string
	}
)

var (
	source io.Reader = os.Stdin
	dest   io.Writer = os.Stdout
	// fileNameExp is a regex to match a test filename in an output event.
	fileNameExp, _ = regexp.Compile(`^\s*(\w+_test.go):\d+:`)
	// testOutput holds all data during processing, and is updated/deleted from as necessary.
	testOutput map[string]*Test
	// counts for different test statuses.
	failCount, passCount, skipCount int
)

func main() {
	// init data here for testing purposes.
	testOutput = make(map[string]*Test, 200)
	failCount, passCount, skipCount = 0, 0, 0

	scan := bufio.NewScanner(source)
	for scan.Scan() {
		var event Event
		if err := json.Unmarshal(scan.Bytes(), &event); err != nil {
			continue
		}
		handleEvent(event)
	}

	markdown := markdownAll()
	if _, err := fmt.Fprint(dest, markdown); err != nil {
		log.Fatal(err)
	}
}

func handleEvent(event Event) {
	if event.Package == "" || event.Test == "" {
		return
	}

	// get or initialize data for event.Test.
	testKey := event.Package + "/" + event.Test
	test, ok := testOutput[testKey]
	if !ok {
		test = &Test{
			Package: event.Package,
			Name:    event.Test,
		}
		testOutput[testKey] = test
	}

	// fail/pass/skip events come LAST, so initially save all output and update/clear later.
	switch event.Action {
	case actionFail:
		failCount++
	case actionSkip:
		skipCount++
		delete(testOutput, testKey)
	case actionPass:
		passCount++
		delete(testOutput, testKey)
	case actionOutput:
		test.Output = append(test.Output, event.Output)
		// file names are not readily available, so search for them in output.
		matches := fileNameExp.FindStringSubmatch(event.Output)
		if len(matches) > 1 {
			test.File = strings.TrimPrefix(event.Package+"/"+matches[1], moduleName+"/")
		}
	}
}

func markdownAll() string {
	summary := strings.Builder{}
	summary.WriteString(markdownSummaryTable())

	if failCount == 0 {
		return summary.String()
	}

	summary.WriteString("## Failed Tests\n\n")

	packageFileTests := testsByPackageAndFile()
	// ① packages
	for _, pkg := range sortedKeys(packageFileTests) {
		shortPkgName := strings.TrimPrefix(pkg, moduleName+"/")
		summary.WriteString("---\n\n### `" + shortPkgName + "`\n\n")

		fileTests := packageFileTests[pkg]
		summary.WriteString(markdownRunLocally(pkg, fileTests))
		summary.WriteString("\n\n### Details\n\n")
		// ② files
		for _, filename := range sortedKeys(fileTests) {
			tests := fileTests[filename]
			shortFileName := strings.TrimPrefix(filename, shortPkgName+"/")
			summary.WriteString("\n#### `" + shortFileName + "`\n\n")
			// ③ tests
			sort.Slice(tests, func(i, j int) bool { return tests[i].Name < tests[j].Name })
			for _, test := range tests {
				details := markdownTestDetails(test)
				if summary.Len()+len(details) > maxOutputSize {
					summary.WriteString("---\n\n## Test output exceeded the 1024K limit")
					return summary.String()
				}
				summary.WriteString(details + "\n\n")
			}
		}
	}
	return summary.String()
}

func markdownSummaryTable() string {
	return fmt.Sprintf(`## Test Summary

| ✅ Passed | ❌ Failed | ⏩ Skipped |
|-----------|----------|------------|
|    %d     |     %d   |     %d     |

`, passCount, failCount, skipCount)
}

// markdownTestDetails creates a string with all test.Output values inside a <details> tag.
func markdownTestDetails(test *Test) string {
	details := strings.Builder{}
	details.WriteString("<details>\n<summary>" + test.Name + "</summary>\n\n```diff\n")

	for _, output := range test.Output {
		output = strings.Trim(output, " \n")
		if !validateOutput(output) {
			continue
		}
		details.WriteString(output + "\n")
	}
	details.WriteString("```\n\n</details>")
	return details.String()
}

// markdownRunLocally returns a go test command for the given package and test names.
func markdownRunLocally(pkg string, file map[string][]*Test) string {
	names := make([]string, 0, len(file))
	for _, tests := range file {
		for _, test := range tests {
			if !strings.Contains(test.Name, "/") { // exclude sub-tests.
				names = append(names, test.Name)
			}
		}
	}
	sort.Slice(names, func(i, j int) bool { return names[i] < names[j] })

	return `### Run Locally

` + "```" + `bash
go test ` + pkg + ` -run '` + strings.Join(names, "|") + `'
` + "```"
}

// testsByPackageAndFile groups testOutput data in the structure package -> file -> tests.
func testsByPackageAndFile() map[string]map[string][]*Test {
	grouped := make(map[string]map[string][]*Test)
	for _, test := range testOutput {
		// Note: without any logging, parent tests only have output "=== RUN" and "--- FAIL",
		// so the filename will not be recovered.
		if test.File == "" {
			continue
		}
		if _, ok := grouped[test.Package]; !ok {
			grouped[test.Package] = make(map[string][]*Test, 10)
		}
		if _, ok := grouped[test.Package][test.File]; !ok {
			grouped[test.Package][test.File] = make([]*Test, 0, 10)
		}
		grouped[test.Package][test.File] = append(grouped[test.Package][test.File], test)
	}
	return grouped
}

// sortedKeys returns the keys of m in ascending order.
func sortedKeys[T any](m map[string]T) []string {
	ord := make([]string, 0, len(m))
	for key := range m {
		ord = append(ord, key)
	}
	sort.Slice(ord, func(i, j int) bool { return ord[i] < ord[j] })
	return ord
}

// validateOutput determines if a test output string should be written to markdown.
func validateOutput(s string) bool {
	return !(s == "" ||
		strings.HasPrefix(s, "=== RUN") ||
		strings.HasPrefix(s, "=== CONT") ||
		strings.HasPrefix(s, "--- FAIL"))
}
