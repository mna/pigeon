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

func TestStringMatcher(t *testing.T) {
	cases := []struct {
		in  string
		val string
		ic  bool
		out bool
	}{
		{"", "", false, true},
		{"", "", true, true},
		{"a", "", false, true},
		{"a", "", true, true},
		{"", "a", false, false},
		{"", "a", true, false},
		{"a", "a", false, true},
		{"a", "a", true, true},
		{"A", "a", false, false},
		{"A", "a", true, true},
		{"abc", "a", false, true},
		{"abc", "a", true, true},
		{"ABc", "a", false, false},
		{"ABc", "a", true, true},
		{"a", "ab", false, false},
		{"a", "ab", true, false},
		{"A", "ab", false, false},
		{"A", "ab", true, false},
		{"abc", "ab", false, true},
		{"abc", "ab", true, true},
		{"ABc", "ab", false, false},
		{"ABc", "ab", true, true},
		{"ϡ", "a", false, false},
		{"ϡ", "a", true, false},
		{"ϡ", "ϡ", false, true},
		{"ϡ", "ϡ", true, true},
		{"ϡ\n", "ϡ", false, true},
		{"ϡ\n", "ϡ", true, true},
		{"ϡ", "ϡ\n", false, false},
		{"ϡ", "ϡ\n", true, false},
		{"ϡ\n", "ϡ\n", false, true},
		{"ϡ\n", "ϡ\n", true, true},
	}

	var m ϡstringMatcher
	for _, tc := range cases {
		pr := testPeekReader{rns: []rune(tc.in)}
		m.value = tc.val
		m.ignoreCase = tc.ic
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
