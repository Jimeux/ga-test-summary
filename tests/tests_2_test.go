package tests

import "testing"

func TestTests2_First(t *testing.T) {
	t.Fatal("failed first")
}

func TestTests2_Second(t *testing.T) {
}

func TestTests2_Third(t *testing.T) {
	t.Fatal("failed third")
}
