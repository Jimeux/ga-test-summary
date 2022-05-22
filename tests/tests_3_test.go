package tests

import "testing"

func TestTests3_First(t *testing.T) {
	t.Fatal("failed first")
}

func TestTests3_Second(t *testing.T) {
}

func TestTests3_Third(t *testing.T) {
	t.Fatal("failed third")
}
