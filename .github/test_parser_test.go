package main

import (
	"bytes"
	_ "embed"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//go:embed fixture.json
var fixture string

//go:embed golden.md
var golden string

func Test(t *testing.T) {
	in = strings.NewReader(fixture)
	buf := new(bytes.Buffer)
	out = buf

	main()

	if diff := cmp.Diff(golden, buf.String()); diff != "" {
		t.Fatal("output did not match golden file:\n" + diff)
	}
}
