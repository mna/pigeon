package main

import (
	"strings"
	"testing"
)

var tc = `
This is a real U+FFFD: "�"
�`

var invalid = []byte{0xff, 0xfe, 0xfd}

func TestRuneError(t *testing.T) {
	if _, err := Parse("", []byte(tc)); err != nil {
		t.Error("Parsing failed:", err)
	}

	if _, err := Parse("", invalid); err == nil {
		t.Error("Did not fail parsing invalid encoding")
	} else if !strings.Contains(err.Error(), "invalid encoding") {
		t.Error("Unexpected error:", err)
	}
}
