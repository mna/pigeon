package vm

import (
	"io"
	"testing"
)

func TestErrList(t *testing.T) {
	var el errList

	if err := el.ϡerr(); err != nil {
		t.Errorf("want nil, got %v", err)
	}
	cnt := 3
	for i := 0; i < cnt; i++ {
		el.ϡadd(io.EOF)
	}
	if len(el) != cnt {
		t.Errorf("want %d error, got %d", cnt, len(el))
	}

	el.ϡdedupe()
	if len(el) != 1 {
		t.Errorf("want 1 error, got %d", len(el))
	}
	if err := el.ϡerr(); err == nil {
		t.Errorf("want not nil, got nil")
	}
	msg := el.Error()
	if msg != io.EOF.Error() {
		t.Errorf("want message %q, got %q", io.EOF, msg)
	}
	el.ϡadd(errNoMatch)
	msg = el.Error()
	want := io.EOF.Error() + "\n" + errNoMatch.Error()
	if msg != want {
		t.Errorf("want message %q, got %q", want, msg)
	}
}

func TestParserError(t *testing.T) {
	pe := parserError{Inner: io.EOF, ϡprefix: "a"}
	msg := pe.Error()
	if want := "a: " + io.EOF.Error(); want != msg {
		t.Errorf("want message %q, got %q", want, msg)
	}
}
