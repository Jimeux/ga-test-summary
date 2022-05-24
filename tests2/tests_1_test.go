package tests2

import (
	"log"
	"testing"
)

func TestTests1_First(t *testing.T) {
	log.Println("Example log")
	t.Fatal("failed first")
}

func TestTests1_Second(t *testing.T) {
	t.Log("pass filename regexp-catcher")
}

func TestTests1_Third(t *testing.T) {
	log.Println("Example log third 1")
	log.Println("Example log third 2")
	t.Fatal("failed third")
}

func TestTests1_FourthTable(t *testing.T) {
	t.Log("fail filename regexp-catcher")

	tests := []struct {
		name string
		fail bool
	}{
		{"subtest_1", true},
		{"subtest_2", true},
		{"subtest_3", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.fail {
				t.Fatal("failed sub-test")
			}
		})
	}
}
