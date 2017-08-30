package main

import (
	"strings"
	"testing"
)

func TestMaxExprCnt(t *testing.T) {
	maxExprMsg := "max expr count hit"
	_, err := Parse("", []byte("infinite parse"), MaxExpressions(5, maxExprMsg))
	if err == nil {
		t.Errorf("expected non nil error message for testing max expr cnt option.")
		t.Fail()
	}
	if !strings.Contains(err.Error(), maxExprMsg) {
		t.Errorf("expected error to contain max expr error message")
	}
}
