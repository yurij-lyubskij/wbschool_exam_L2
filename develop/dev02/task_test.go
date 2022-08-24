package main

import "testing"

func TestEmpty(t *testing.T) {
	s := ""
	exp := ""
	res, err := Unpack(s)
	if res != exp || err != nil {
		t.Fatalf("unpacking string %s \n, got res = %s, err= %v,\n expected %s", s, res, err, exp)
	}
}

func TestIncorrect(t *testing.T) {
	s := "45"
	exp := ""
	res, err := Unpack(s)
	if err == nil || err.Error() != "некорректная строка" {
		t.Fatalf("unpacking string %s \n, got res = %s, err= %v,\n expected %s", s, res, err, exp)
	}
}

func TestSimple(t *testing.T) {
	s := "abcd"
	exp := "abcd"
	res, err := Unpack(s)
	if res != exp || err != nil {
		t.Fatalf("unpacking string %s \n, got res = %s, err= %v,\n expected %s", s, res, err, exp)
	}
}

func TestRepeats(t *testing.T) {
	s := "a4bc2d5e"
	exp := "aaaabccddddde"
	res, err := Unpack(s)
	if res != exp || err != nil {
		t.Fatalf("unpacking string %s \n, got res = %s, err= %v,\n expected %s", s, res, err, exp)
	}
}

func TestEscapeAll(t *testing.T) {
	s := "qwe\\4\\5"
	exp := "qwe45"
	res, err := Unpack(s)
	if res != exp || err != nil {
		t.Fatalf("unpacking string %s \n, got res = %s, err= %v,\n expected %s", s, res, err, exp)
	}
}

func TestEscapeRepeatNumber(t *testing.T) {
	s := "qwe\\45"
	exp := "qwe44444"
	res, err := Unpack(s)
	if res != exp || err != nil {
		t.Fatalf("unpacking string %s \n, got res = %s, err= %v,\n expected %s", s, res, err, exp)
	}
}

func TestEscapeRepeatEscape(t *testing.T) {
	s := "qwe\\\\5"
	exp := "qwe\\\\\\\\\\"
	res, err := Unpack(s)
	if res != exp || err != nil {
		t.Fatalf("unpacking string %s \n, got res = %s, err= %v,\n expected %s", s, res, err, exp)
	}
}

func TestEscapeError(t *testing.T) {
	s := "qwe\\"
	exp := ""
	res, err := Unpack(s)
	if res != exp || err == nil {
		t.Fatalf("unpacking string %s \n, got res = %s, err= %v,\n expected %s", s, res, err, exp)
	}
}
