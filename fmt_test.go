package emoji

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSprint(t *testing.T) {
	var (
		input    = "I am :man_technologist: from :flag_for_turkey:. Tests are :thumbs_up:"
		expected = fmt.Sprintf("I am %v from %v. Tests are %v", ManTechnologist, FlagForTurkey, ThumbsUp)
	)

	got := Sprint(input)
	if got != expected {
		t.Fatalf("test case fail: got: %v, expected: %v", got, expected)
	}
}

func TestSprintf(t *testing.T) {
	var (
		input    = "I am :man_technologist:. Tests are :thumbs_up:. %v is formatted."
		args     = "this string"
		expected = fmt.Sprintf("I am %v. Tests are %v. %v is formatted.", ManTechnologist, ThumbsUp, args)
	)

	got := Sprintf(input, args)
	if got != expected {
		t.Fatalf("test case fail: got: %v, expected: %v", got, expected)
	}
}

func TestSprintln(t *testing.T) {
	var (
		input    = "I am :man_technologist: from :flag_for_turkey:. Tests are :thumbs_up:"
		expected = fmt.Sprintf("I am %v from %v. Tests are %v\n", ManTechnologist, FlagForTurkey, ThumbsUp)
	)

	got := Sprintln(input)
	if got != expected {
		t.Fatalf("test case fail: got: %v, expected: %v", got, expected)
	}
}

func TestPrint(t *testing.T) {
	var (
		input = "I am :man_technologist: from :flag_for_turkey:. Tests are :thumbs_up:"
	)

	n, err := Print(input)
	if err != nil || n == 0 {
		t.Fatalf("test case fail n: %v: %v", n, err)
	}
}

func TestPrintf(t *testing.T) {
	var (
		input = "I am :man_technologist:. Tests are :thumbs_up:. %v is formatted."
		args  = "this string"
	)

	n, err := Printf(input, args)
	if err != nil || n == 0 {
		t.Fatalf("test case fail n: %v: %v", n, err)
	}
}

func TestPrintln(t *testing.T) {
	var (
		input = "I am :man_technologist: from :flag_for_turkey:. Tests are :thumbs_up:"
	)

	n, err := Println(input)
	if err != nil || n == 0 {
		t.Fatalf("test case fail n: %v: %v", n, err)
	}
}

func TestFprint(t *testing.T) {
	var (
		input    = "I am :man_technologist: from :flag_for_turkey:. Tests are :thumbs_up:"
		expected = fmt.Sprintf("I am %v from %v. Tests are %v", ManTechnologist, FlagForTurkey, ThumbsUp)
	)

	var w bytes.Buffer
	_, err := Fprint(&w, input)
	if err != nil {
		t.Fatalf("test case fail: %v", err)
	}

	got := w.String()
	if got != expected {
		t.Fatalf("test case fail: got: %v, expected: %v", got, expected)
	}
}

func TestFprintf(t *testing.T) {
	var (
		input    = "I am :man_technologist:. Tests are :thumbs_up:. %v is formatted."
		args     = "this string"
		expected = fmt.Sprintf("I am %v. Tests are %v. %v is formatted.", ManTechnologist, ThumbsUp, args)
	)

	var w bytes.Buffer
	_, err := Fprintf(&w, input, args)
	if err != nil {
		t.Fatalf("test case fail: %v", err)
	}

	got := w.String()
	if got != expected {
		t.Fatalf("test case fail: got: %v, expected: %v", got, expected)
	}
}

func TestFprintln(t *testing.T) {
	var (
		input    = "I am :man_technologist: from :flag_for_turkey:. Tests are :thumbs_up:"
		expected = fmt.Sprintf("I am %v from %v. Tests are %v\n", ManTechnologist, FlagForTurkey, ThumbsUp)
	)

	var w bytes.Buffer
	_, err := Fprintln(&w, input)
	if err != nil {
		t.Fatalf("test case fail: %v", err)
	}

	got := w.String()
	if got != expected {
		t.Fatalf("test case fail: got: %v, expected: %v", got, expected)
	}
}

func TestErrorf(t *testing.T) {
	var (
		input    = "I am :man_technologist:. Tests are :thumbs_up:. %v is formatted."
		args     = "this string"
		expected = fmt.Sprintf("I am %v. Tests are %v. %v is formatted.", ManTechnologist, ThumbsUp, args)
	)

	got := Errorf(input, args).Error()
	if got != expected {
		t.Fatalf("test case fail: got: %v, expected: %v", got, expected)
	}
}
