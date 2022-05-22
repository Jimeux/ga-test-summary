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
	ActionPass   = "pass"
	ActionFail   = "fail"
	ActionOutput = "output"
	ActionSkip   = "skip"

	moduleName = "github.com/Jimeux/ga-summary"
)

var (
	fileNameExp, _           = regexp.Compile("^\\s*(\\w+_test.go):\\d+:")
	fileTests                = make(map[string][]string, 10)
	testEvents               = make(map[string][]Event, 10)
	failed                   = make(map[string]struct{}, 10)
	passed                   = make(map[string]struct{}, 100)
	skippedCount             = 0
	parseErrCount            = 0
	in             io.Reader = os.Stdin
	out            io.Writer = os.Stdout
)

// Event holds data from a single test event.
// See cmd/test2json.
type Event struct {
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
	var event Event
	if err := json.Unmarshal(b, &event); err != nil {
		parseErrCount++
		return
	}
	if event.Test == "" {
		return
	}
	if _, ok := passed[event.Test]; ok {
		return
	}
	// store all events in testEvents, and delete when we see event.Action==pass
	testEvents[event.Test] = append(testEvents[event.Test], event)

	switch event.Action {
	case ActionSkip:
		skippedCount++
	case ActionPass:
		passed[event.Test] = struct{}{}
		delete(testEvents, event.Test) // remove events when we see pass
	case ActionFail:
		failed[event.Test] = struct{}{}
	case ActionOutput:
		// file names are not available, so search for them in output
		matches := fileNameExp.FindStringSubmatch(event.Output)
		if len(matches) > 1 {
			fullName := strings.TrimPrefix(event.Package+"/"+matches[1], moduleName+"/")
			fileTests[fullName] = append(fileTests[fullName], event.Test)
		}
	}
}

func writeMarkdown(w io.Writer) {
	runExpr := strings.Join(sortedKeys(failed), "|")
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

`, len(passed), len(failed), skippedCount, parseErrCount, runExpr)

	for _, name := range sortedKeys(fileTests) {
		_, _ = fmt.Fprint(w, "---\n\n### `"+name+"`\n\n")
		for _, test := range fileTests[name] {
			_, _ = fmt.Fprint(w, "<details>\n<summary>"+test+"</summary>\n\n```bash\n")
			for _, event := range testEvents[test] {
				line := strings.Trim(event.Output, " \n")
				if line == "" {
					continue
				}
				_, _ = fmt.Fprintf(w, line+"\n")
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
