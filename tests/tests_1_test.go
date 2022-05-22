package tests

import (
	"log"
	"testing"
)

func TestTests1_First(t *testing.T) {
	log.Println("Example log")
	t.Fatal("failed first")
}

func TestTests1_Second(t *testing.T) {
}

func TestTests1_Third(t *testing.T) {
	log.Println("Example log third 1")
	log.Println("Example log third 2")
	t.Fatal("failed third")
}
