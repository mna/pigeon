package vm

import (
	"testing"
	"unicode/utf8"
)

func TestParserRead(t *testing.T) {
	cases := []struct {
		in  string
		end position
	}{
		{"", position{line: 1}},
		{" ", position{line: 1, col: 1, offset: 1}},
		{"\n", position{line: 2, col: 0, offset: 1}},
		{"\u03e1", position{line: 1, col: 1, offset: 2}},
		{"a\nb\n\u03e1\n", position{line: 4, col: 0, offset: 7}},
	}

	for _, tc := range cases {
		p := ϡparser{data: []byte(tc.in), pt: ϡsvpt{position: position{line: 1}}}
		for {
			p.read()
			if p.pt.rn == utf8.RuneError {
				break
			}
		}

		if tc.end != p.pt.position {
			t.Errorf("%q: want %s, got %s", tc.in, tc.end, p.pt.position)
		}
		// on normal exit, savepoint is always on utf8.RuneError
		if p.pt.rn != utf8.RuneError {
			t.Errorf("%q: want RuneError on exit, got %#U", tc.in, p.pt.rn)
		}
	}
}

func TestParserInvalidEncodingPanics(t *testing.T) {
	p := ϡparser{data: []byte("ab\xdf"), pt: ϡsvpt{position: position{line: 1}}}

	ok := false
	func() {
		defer func() {
			if e := recover(); e != nil {
				ok = true
			}
		}()
		for {
			p.read()
			if p.pt.rn == utf8.RuneError {
				break
			}
		}
	}()

	if !ok {
		t.Errorf("want panic, got none")
	}
	end := position{line: 1, col: 2, offset: 2}
	if p.pt.position != end {
		t.Errorf("want position %s, got %s", end, p.pt.position)
	}
}
