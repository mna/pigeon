package vm

import (
	"testing"
	"unicode/utf8"
)

func TestAnyMatcher(t *testing.T) {
	cases := []struct {
		in  string
		out bool
	}{
		{"", false},
		{"a", true},
		{"\n", true},
		{" ", true},
		{"ϡ", true},
		{"ab", true},
		{"ϡ進", true},
		{"\n\n", true},
	}

	var m ϡanyMatcher
	for _, tc := range cases {
		pr := testPeekReader{rns: []rune(tc.in)}
		got := m.match(&pr)

		if got != tc.out {
			t.Errorf("%q: want %t, got %t", tc.in, tc.out, got)
		}
	}
}

type testPeekReader struct {
	rns []rune
	ix  int
}

func (pr *testPeekReader) peek() ϡsvpt {
	if pr.ix >= len(pr.rns) {
		return ϡsvpt{rn: utf8.RuneError}
	}
	return ϡsvpt{rn: pr.rns[pr.ix]}
}

func (pr *testPeekReader) read() {
	pr.ix++
}
