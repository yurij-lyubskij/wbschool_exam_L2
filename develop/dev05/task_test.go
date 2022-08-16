package main

import (
	"bytes"
	"io"
	"testing"
)

var testName string

func defStart() {
	after = 0
	before = 0
	context = 0
	count = false
	ignore = false
	invert = false
	fixed = false
	line = false
	patterns = false
	testName = "test.txt"
}

func TestCaseSimple(t *testing.T) {
	defStart()
	regexp := "apt"
	args := []string{regexp, testName}
	buf := make([]byte, 1)
	i := bytes.NewBuffer(buf)
	a := io.Writer(i)
	expected := "laptop\n"
	grep(args, a)
	res := i.String()
	if expected != res[1:] {
		t.Fatalf("wrong result, expected:\n%sgot:\n%s", expected, res)
	}
}

func TestCaseRegExp(t *testing.T) {
	defStart()
	regexp := "^d"
	args := []string{regexp, testName}
	buf := make([]byte, 1)
	i := bytes.NewBuffer(buf)
	a := io.Writer(i)
	expected := "data\ndebian\n"
	grep(args, a)
	res := i.String()
	if expected != res[1:] {
		t.Fatalf("wrong result, expected:\n%sgot:\n%s", expected, res)
	}
}

func TestCaseAfter(t *testing.T) {
	defStart()
	after = 1
	regexp := "^d"
	args := []string{regexp, testName}
	buf := make([]byte, 1)
	i := bytes.NewBuffer(buf)
	a := io.Writer(i)
	expected := "data\ndebian\nlaptop\n"
	grep(args, a)
	res := i.String()
	if expected != res[1:] {
		t.Fatalf("wrong result, expected:\n%sgot:\n%s", expected, res)
	}
}

func TestCaseBefore(t *testing.T) {
	defStart()
	before = 1
	regexp := "^d"
	args := []string{regexp, testName}
	buf := make([]byte, 1)
	i := bytes.NewBuffer(buf)
	a := io.Writer(i)
	expected := "computer\ndata\ndebian\n"
	grep(args, a)
	res := i.String()
	if expected != res[1:] {
		t.Fatalf("wrong result, expected:\n%sgot:\n%s", expected, res)
	}
}

func TestCaseContext(t *testing.T) {
	defStart()
	context = 1
	regexp := "^d"
	args := []string{regexp, testName}
	buf := make([]byte, 1)
	i := bytes.NewBuffer(buf)
	a := io.Writer(i)
	expected := "computer\ndata\ndebian\nlaptop\n"
	grep(args, a)
	res := i.String()
	if expected != res[1:] {
		t.Fatalf("wrong result, expected:\n%sgot:\n%s", expected, res)
	}
}

func TestCaseInsensitive(t *testing.T) {
	defStart()
	ignore = true
	regexp := "^lap"
	args := []string{regexp, testName}
	buf := make([]byte, 1)
	i := bytes.NewBuffer(buf)
	a := io.Writer(i)
	expected := "laptop\nLAPTOP\n"
	grep(args, a)
	res := i.String()
	if expected != res[1:] {
		t.Fatalf("wrong result, expected:\n%sgot:\n%s", expected, res)
	}
}

func TestCaseCount(t *testing.T) {
	defStart()
	ignore = true
	count = true
	regexp := "^lap"
	args := []string{regexp, testName}
	buf := make([]byte, 1)
	i := bytes.NewBuffer(buf)
	a := io.Writer(i)
	expected := "2\n"
	grep(args, a)
	res := i.String()
	if expected != res[1:] {
		t.Fatalf("wrong result, expected:\n%sgot:\n%s", expected, res)
	}
}

func TestCaseInvert(t *testing.T) {
	defStart()
	ignore = true
	invert = true
	regexp := "^lap"
	args := []string{regexp, testName}
	buf := make([]byte, 1)
	i := bytes.NewBuffer(buf)
	a := io.Writer(i)
	expected := "computer\ndata\ndebian\nmouse\nRedHat\n"
	grep(args, a)
	res := i.String()
	if expected != res[1:] {
		t.Fatalf("wrong result, expected:\n%sgot:\n%s", expected, res)
	}
}

func TestCaseFixed(t *testing.T) {
	defStart()
	fixed = true
	regexp := "d"
	args := []string{regexp, testName}
	buf := make([]byte, 1)
	i := bytes.NewBuffer(buf)
	a := io.Writer(i)
	expected := "data\ndebian\nRedHat\n"
	grep(args, a)
	res := i.String()
	if expected != res[1:] {
		t.Fatalf("wrong result, expected:\n%sgot:\n%s", expected, res)
	}
}

func TestCaseLine(t *testing.T) {
	defStart()
	fixed = true
	line = true
	regexp := "d"
	args := []string{regexp, testName}
	buf := make([]byte, 1)
	i := bytes.NewBuffer(buf)
	a := io.Writer(i)
	expected := "2 data\n3 debian\n7 RedHat\n"
	grep(args, a)
	res := i.String()
	if expected != res[1:] {
		t.Fatalf("wrong result, expected:\n%sgot:\n%s", expected, res)
	}
}
