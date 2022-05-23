package main

import (
	"bytes"
	"embed"
	_ "embed"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//go:embed files
var files embed.FS

func Test(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"fail"},
		{"pass"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture, err := files.Open("files/fixture_" + test.name + ".json")
			if err != nil {
				t.Fatal(err)
			}
			source = fixture
			golden, err := files.ReadFile("files/golden_" + test.name + ".md")
			if err != nil {
				t.Fatal(err)
			}
			buf := new(bytes.Buffer)
			dest = buf

			main()

			if diff := cmp.Diff(golden, buf.Bytes()); diff != "" {
				t.Fatal("output did not match golden file:\n" + diff)
			}
		})
	}
}
