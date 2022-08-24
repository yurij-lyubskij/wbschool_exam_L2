package main

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

func TestCasePs(t *testing.T) {
	buf := make([]byte, 1)
	w := bytes.NewBuffer(buf)
	input := "ps"
	expected := "PID= \\d+, PPID=  \\d+, CMD= dev08.test"
	exp := regexp.MustCompile(expected)
	r := strings.NewReader(input)
	shell(r, w)
	res := w.String()
	if !exp.MatchString(res[1:]) {
		t.Fatalf("wrong result, expected:\n%sgot:\n%s", expected, res)
	}
}

func TestCasePwd(t *testing.T) {
	buf := make([]byte, 1)
	w := bytes.NewBuffer(buf)
	input := "pwd"
	expected := "/home/yura11011/wbschool_exam_L2/develop/dev08\n"
	r := strings.NewReader(input)
	shell(r, w)
	res := w.String()
	if res[1:] != expected {
		t.Fatalf("wrong result, expected:\n%sgot:\n%s", expected, res)
	}
}

func TestCasePipeCD(t *testing.T) {
	buf := make([]byte, 1)
	w := bytes.NewBuffer(buf)
	input := "cd ..|pwd"
	expected := "/home/yura11011/wbschool_exam_L2/develop\n"
	r := strings.NewReader(input)
	shell(r, w)
	res := w.String()
	if res[1:] != expected {
		t.Fatalf("wrong result, expected:\n%sgot:\n%s", expected, res)
	}
}

func TestCaseEcho(t *testing.T) {
	buf := make([]byte, 1)
	w := bytes.NewBuffer(buf)
	input := "echo green"
	expected := "green\n"
	r := strings.NewReader(input)
	shell(r, w)
	res := w.String()
	if res[1:] != expected {
		t.Fatalf("wrong result, expected:\n%sgot:\n%s", expected, res)
	}
}
