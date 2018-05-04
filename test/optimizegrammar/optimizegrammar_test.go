package optimizegrammar

import (
	"strings"
	"testing"
)

func TestOptimizeGrammar(t *testing.T) {
	cases := []struct {
		input  string
		errMsg string
	}{
		{"X", ""},
		{"Y", "no match found"},
		{"XY", "YY"},
		{"XX", "no match found"},
		{"X,X", ""},
		{"X,XY", "YY"},
		{"X,XY,X", "YY"},
	}
	for _, c := range cases {
		_, err := Parse("", []byte(c.input))
		if c.errMsg == "" && err != nil {
			t.Errorf("%q: want no error, got %v", c.input, err)
			continue
		}
		if c.errMsg != "" && err == nil {
			t.Errorf("%q: want error %q, got none", c.input, c.errMsg)
			continue
		}
		if c.errMsg != "" && !strings.Contains(err.Error(), c.errMsg) {
			t.Errorf("%q: want error to contain %q, got %q", c.input, c.errMsg, err)
			continue
		}
	}
}
