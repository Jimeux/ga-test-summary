package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"flag"
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

	moduleName = "github.com/Jimeux/ga-summary"
)

var (
	fileNameExp, _           = regexp.Compile("^\\s*(\\w+_test.go):\\d+:")
	fileTests                = make(map[string][]string, 10)
	testEvents               = make(map[string][]event, 10)
	failedTests              = make(map[string]struct{}, 10)
	passedTests              = make(map[string]struct{}, 100)
	skippedCount             = 0
	parseErrCount            = 0
	in             io.Reader = os.Stdin
	out            io.Writer = os.Stdout
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
	filename := flag.String("file", "", "-file output.json")
	flag.Parse()

	if *filename != "" {
		open, err := os.Open(*filename)
		if err != nil {
			log.Fatalf("file not found for filename %s: %v", *filename, err)
		}
		in = open
	}

	scan := bufio.NewScanner(in)
	for scan.Scan() {
		handleEvent(scan.Bytes())
	}
	writeMarkdown(out)
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
	// store all events in testEvents, and delete when we see event.Action==pass
	testEvents[event.Test] = append(testEvents[event.Test], event)

	switch event.Action {
	case actionSkip:
		skippedCount++
	case actionPass:
		passedTests[event.Test] = struct{}{}
		delete(testEvents, event.Test) // remove events when we see event.Action==pass
	case actionFail:
		failedTests[event.Test] = struct{}{}
	case actionOutput:
		// file names are not available, so search for them in output
		matches := fileNameExp.FindStringSubmatch(event.Output)
		if len(matches) > 1 {
			fullName := strings.TrimPrefix(event.Package+"/"+matches[1], moduleName+"/")
			fileTests[fullName] = append(fileTests[fullName], event.Test)
		}
	}
}

func writeMarkdown(w io.Writer) {
	runLocalExpr := strings.Join(sortedKeys(failedTests), "|")
	_, _ = fmt.Fprintf(w, `# Test Summary

|     Status      | Count |
|-----------------|-------|
| ✅ Passed       | %d   |
| ❌ Failed       | %d   |
| ⏩ Skipped      | %d   |
| 💥 Parse Errors | %d   |

## Run Failed Tests Locally

`+"```"+`bash
go test ./... -v -run '%s'
`+"```"+`

## Failure Details

`, len(passedTests), len(failedTests), skippedCount, parseErrCount, runLocalExpr)

	for _, filename := range sortedKeys(fileTests) {
		_, _ = fmt.Fprint(w, "---\n\n#### `"+filename+"`\n\n")
		for _, testName := range fileTests[filename] {
			_, _ = fmt.Fprint(w, "<details>\n<summary>"+testName+"</summary>\n\n```bash\n")
			for _, event := range testEvents[testName] {
				output := strings.Trim(event.Output, " \n")
				if output == "" {
					continue
				}
				_, _ = fmt.Fprintf(w, output+"\n")
			}
			_, _ = fmt.Fprintf(w, "```\n\n</details>\n\n")
		}
	}
}

func sortedKeys[T any](m map[string]T) []string {
	ord := make([]string, 0, len(m))
	for key := range m {
		ord = append(ord, key)
	}
	sort.Slice(ord, func(i, j int) bool { return ord[i] < ord[j] })
	return ord
}
