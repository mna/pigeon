package alternate_entrypoint

import (
	"strings"
	"testing"
)

func TestAlternateEntrypoint(t *testing.T) {
	src := "bbbcc"

	v, err := Parse("", []byte(src), Entrypoint("Entry2"))
	if err != nil {
		t.Fatal(err)
	}

	got := string(v.([]byte))
	if got != src {
		t.Fatalf("want %s, got %s", src, got)
	}
}

func TestInvalidAlternateEntrypoint(t *testing.T) {
	src := "bbbcc"

	_, err := Parse("", []byte(src), Entrypoint("Z"))
	if err == nil {
		t.Fatal("want error, got none")
	}
	if !strings.Contains(err.Error(), errInvalidEntrypoint.Error()) {
		t.Fatalf("want %s, got %s", errInvalidEntrypoint, err)
	}
}

func TestAlternateInputWithDefaultEntrypoint(t *testing.T) {
	src := "bbbcc"

	_, err := Parse("", []byte(src))
	if err == nil {
		t.Fatal("want error, got none")
	}
	if !strings.Contains(err.Error(), "no match found") {
		t.Fatalf("want 'no match found', got %s", err)
	}
}

func TestDefaultInputWithAlternateEntrypoint(t *testing.T) {
	src := "aacc"

	_, err := Parse("", []byte(src), Entrypoint("Entry2"))
	if err == nil {
		t.Fatal("want error, got none")
	}
	if !strings.Contains(err.Error(), "no match found") {
		t.Fatalf("want 'no match found', got %s", err)
	}
}
