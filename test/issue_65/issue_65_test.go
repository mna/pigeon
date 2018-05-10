package issue65

import (
	"strings"
	"testing"

	optimized "github.com/mna/pigeon/test/issue_65/optimized"
	optimizedgrammar "github.com/mna/pigeon/test/issue_65/optimized-grammar"
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

	type parser func(string) (interface{}, error)
	parsers := []parser{parseStd, parseOpt, parseOptGrammar}
	for _, p := range parsers {
		for _, c := range cases {
			_, err := p(c.input)
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
}

func parseStd(input string) (interface{}, error) {
	return Parse("", []byte(input))
}

func parseOpt(input string) (interface{}, error) {
	return optimized.Parse("", []byte(input))
}

func parseOptGrammar(input string) (interface{}, error) {
	return optimizedgrammar.Parse("", []byte(input))
}
