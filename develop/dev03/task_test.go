package main

import (
	"strings"
	"testing"
)

var defaultIn = "drwxr-xr-x 7 user user 4096 июн 10 2015 Sources\n" +
	"drwxr-xr-x 8 user user 12288 янв 11 12:33 Pictures\n" +
	"drwxr-xr-x 7 user user 4096 янв 13 11:42 Lightworks\n" +
	"drwx------ 5 user user 12288 янв 15 14:59 Downloads\n" +
	"drwxr-xr-x 7 user user 4096 окт 31 15:08 VirtualBox\n" +
	"drwxr-xr-x 6 user user 4096 дек 6 14:29 Android"

func TestReadFile(t *testing.T) {
	expected := defaultIn
	file := "test.txt"
	fSlice, err := ReadFile(file)
	if err != nil {
		t.Fatalf("can't open file")
	}
	fileString := strings.Join(fSlice, "\n")
	if fileString != expected {
		t.Error("file wasn't read correctly")
		t.Fatalf("expected %s\n, got %s", expected, fileString)
	}
}

func TestSortNormal(t *testing.T) {
	input := defaultIn
	var keys Keys
	inSlice := strings.Split(input, "\n")
	sorted := Sort(inSlice, keys)
	sortedString := strings.Join(sorted, "\n")
	var expected = "drwx------ 5 user user 12288 янв 15 14:59 Downloads\n" +
		"drwxr-xr-x 6 user user 4096 дек 6 14:29 Android\n" +
		"drwxr-xr-x 7 user user 4096 июн 10 2015 Sources\n" +
		"drwxr-xr-x 7 user user 4096 окт 31 15:08 VirtualBox\n" +
		"drwxr-xr-x 7 user user 4096 янв 13 11:42 Lightworks\n" +
		"drwxr-xr-x 8 user user 12288 янв 11 12:33 Pictures"
	if sortedString != expected {
		t.Error("a mistake in sorting")
		t.Fatalf("expected %s\n, got %s", expected, sortedString)
	}
}

func TestSortByColumn(t *testing.T) {
	input := defaultIn
	var keys Keys
	keys.Column = 8
	inSlice := strings.Split(input, "\n")
	sorted := Sort(inSlice, keys)
	sortedString := strings.Join(sorted, "\n")
	var expected = "drwxr-xr-x 6 user user 4096 дек 6 14:29 Android\n" +
		"drwx------ 5 user user 12288 янв 15 14:59 Downloads\n" +
		"drwxr-xr-x 7 user user 4096 янв 13 11:42 Lightworks\n" +
		"drwxr-xr-x 8 user user 12288 янв 11 12:33 Pictures\n" +
		"drwxr-xr-x 7 user user 4096 июн 10 2015 Sources\n" +
		"drwxr-xr-x 7 user user 4096 окт 31 15:08 VirtualBox"

	if sortedString != expected {
		t.Error("a mistake in sorting")
		t.Fatalf("expected %s\n, got %s", expected, sortedString)
	}
}

func TestSortNumeric(t *testing.T) {
	input := defaultIn
	var keys Keys
	keys.Column = 6
	keys.Numeric = true
	inSlice := strings.Split(input, "\n")
	sorted := Sort(inSlice, keys)
	sortedString := strings.Join(sorted, "\n")
	var expected = "drwxr-xr-x 6 user user 4096 дек 6 14:29 Android\n" +
		"drwxr-xr-x 7 user user 4096 июн 10 2015 Sources\n" +
		"drwxr-xr-x 8 user user 12288 янв 11 12:33 Pictures\n" +
		"drwxr-xr-x 7 user user 4096 янв 13 11:42 Lightworks\n" +
		"drwx------ 5 user user 12288 янв 15 14:59 Downloads\n" +
		"drwxr-xr-x 7 user user 4096 окт 31 15:08 VirtualBox"

	if sortedString != expected {
		t.Error("a mistake in sorting")
		t.Fatalf("expected %s\n, got %s", expected, sortedString)
	}
}

func TestSortUnique(t *testing.T) {
	input := defaultIn
	var keys Keys
	keys.Column = 8
	keys.Unique = true
	inSlice := strings.Split(input, "\n")
	inSlice = append(inSlice, "drwx------ 5 user user 12288 янв 15 14:59 Downloads")
	sorted := Sort(inSlice, keys)
	sortedString := strings.Join(sorted, "\n")
	var expected = "drwxr-xr-x 6 user user 4096 дек 6 14:29 Android\n" +
		"drwx------ 5 user user 12288 янв 15 14:59 Downloads\n" +
		"drwxr-xr-x 7 user user 4096 янв 13 11:42 Lightworks\n" +
		"drwxr-xr-x 8 user user 12288 янв 11 12:33 Pictures\n" +
		"drwxr-xr-x 7 user user 4096 июн 10 2015 Sources\n" +
		"drwxr-xr-x 7 user user 4096 окт 31 15:08 VirtualBox"

	if sortedString != expected {
		t.Error("a mistake in sorting")
		t.Fatalf("expected %s\n, got %s", expected, sortedString)
	}
}
