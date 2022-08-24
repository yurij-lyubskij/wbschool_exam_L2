package main

import (
	"bytes"
	"io"
	"testing"
)

var testName string

func defStart() {
	fields = ""
	delimiter = "\t"
	separated = false
	testName = "test.txt"
}

func TestCaseGeneral(t *testing.T) {
	defStart()
	fields = "1-3,4-6"
	delimiter = " "
	args := []string{testName}
	buf := make([]byte, 1)
	i := bytes.NewBuffer(buf)
	a := io.Writer(i)
	expected := "drwx------ 5 user user 12288 янв\n" +
		"drwxr-xr-x 6 user user 4096 дек\n" +
		"drwxr-xr-x 7 user user 4096 июн\n" +
		"drwxr-xr-x 7 user user 4096 окт\n" +
		"drwxr-xr-x 7 user user 4096 янв\n" +
		"drwxr-xr-x 8 user user 12288 янв\n"
	cut(args, a)
	res := i.String()
	if expected != res[1:] {
		t.Fatalf("wrong result, expected:\n%sgot:\n%s", expected, res)
	}
}

func TestCaseSimple(t *testing.T) {
	defStart()
	fields = "1,3,4-6"
	delimiter = " "
	args := []string{testName}
	buf := make([]byte, 1)
	i := bytes.NewBuffer(buf)
	a := io.Writer(i)
	expected := "drwx------ user user 12288 янв\n" +
		"drwxr-xr-x user user 4096 дек\n" +
		"drwxr-xr-x user user 4096 июн\n" +
		"drwxr-xr-x user user 4096 окт\n" +
		"drwxr-xr-x user user 4096 янв\n" +
		"drwxr-xr-x user user 12288 янв\n"
	cut(args, a)
	res := i.String()
	if expected != res[1:] {
		t.Fatalf("wrong result, expected:\n%sgot:\n%s", expected, res)
	}
}

func TestCaseUpper(t *testing.T) {
	defStart()
	fields = "-6"
	delimiter = " "
	args := []string{testName}
	buf := make([]byte, 1)
	i := bytes.NewBuffer(buf)
	a := io.Writer(i)
	expected := "drwx------ 5 user user 12288 янв\n" +
		"drwxr-xr-x 6 user user 4096 дек\n" +
		"drwxr-xr-x 7 user user 4096 июн\n" +
		"drwxr-xr-x 7 user user 4096 окт\n" +
		"drwxr-xr-x 7 user user 4096 янв\n" +
		"drwxr-xr-x 8 user user 12288 янв\n"
	cut(args, a)
	res := i.String()
	if expected != res[1:] {
		t.Fatalf("wrong result, expected:\n%sgot:\n%s", expected, res)
	}
}

func TestCaseLower(t *testing.T) {
	defStart()
	fields = "4-"
	delimiter = " "
	args := []string{testName}
	buf := make([]byte, 1)
	i := bytes.NewBuffer(buf)
	a := io.Writer(i)
	expected := "user 12288 янв 15 14:59 Downloads\n" +
		"user 4096 дек 6 14:29 Android\n" +
		"user 4096 июн 10 2015 Sources\n" +
		"user 4096 окт 31 15:08 VirtualBox\n" +
		"user 4096 янв 13 11:42 Lightworks\n" +
		"user 12288 янв 11 12:33 Pictures\n"
	cut(args, a)
	res := i.String()
	if expected != res[1:] {
		t.Fatalf("wrong result, expected:\n%sgot:\n%s", expected, res)
	}
}

func TestCaseDelimiter(t *testing.T) {
	defStart()
	fields = "4-"
	args := []string{testName}
	buf := make([]byte, 1)
	i := bytes.NewBuffer(buf)
	a := io.Writer(i)
	expected := "drwx------ 5 user user 12288 янв 15 14:59 Downloads\n" +
		"drwxr-xr-x 6 user user 4096 дек 6 14:29 Android\n" +
		"drwxr-xr-x 7 user user 4096 июн 10 2015 Sources\n" +
		"drwxr-xr-x 7 user user 4096 окт 31 15:08 VirtualBox\n" +
		"drwxr-xr-x 7 user user 4096 янв 13 11:42 Lightworks\n" +
		"drwxr-xr-x 8 user user 12288 янв 11 12:33 Pictures\n"
	cut(args, a)
	res := i.String()
	if expected != res[1:] {
		t.Fatalf("wrong result, expected:\n%sgot:\n%s", expected, res)
	}
}

func TestCaseSeparated(t *testing.T) {
	defStart()
	separated = true
	fields = "4-"
	args := []string{testName}
	buf := make([]byte, 1)
	i := bytes.NewBuffer(buf)
	a := io.Writer(i)
	expected := ""
	cut(args, a)
	res := i.String()
	if expected != res[1:] {
		t.Fatalf("wrong result, expected:\n%sgot:\n%s", expected, res)
	}
}
