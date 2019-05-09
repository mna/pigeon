package issue70

import (
	"errors"
	"regexp"
	"testing"

	optimized "github.com/mna/pigeon/test/issue_70/optimized"
	optimizedgrammar "github.com/mna/pigeon/test/issue_70/optimized-grammar"
)

func TestOptimizeGrammar(t *testing.T) {
	cases := []struct {
		input string
		out   string
		err   error
	}{
		{"Z", "Z", nil},
		{"Y", "Z", errors.New("no match found")},
	}

	type parser func(string) (interface{}, error)
	parsers := map[string]parser{
		"standard":          parseStd,
		"optimized":         parseOpt,
		"optimized-grammar": parseOptGrammar,
	}
	for name, parser := range parsers {
		for _, c := range cases {
			out, err := parser(c.input)
			if (err == nil) != (c.err == nil) || err != nil && !regexp.MustCompile(c.err.Error()).MatchString(err.Error()) {
				t.Errorf("%s: %q: error to be %v, got %v", name, c.input, c.err, err)
				continue
			}
			if err == nil {
				var outStr string
				var ok bool
				if outStr, ok = out.(string); !ok {
					t.Errorf("%s: %q: expect out to be of type string, got %v with type %T", name, c.input, out, out)
					continue
				}
				if outStr != c.out {
					t.Errorf("%s: %q: expect out to be %s, got %s", name, c.input, c.out, outStr)
					continue
				}
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
