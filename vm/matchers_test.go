package vm

import (
	"testing"
	"unicode"
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
			t.Errorf("%q with %q: want %t, got %t", tc.in, tc.val, tc.out, got)
		}
	}
}

func TestCharClassMatcher(t *testing.T) {
	cases := []struct {
		in      string
		val     string
		chars   []rune
		ranges  []rune
		classes []string
		ic      bool
		iv      bool
		out     bool
	}{
		{in: "", val: "[]", out: false},            // empty char class means no char matches
		{in: "", val: "[^]", iv: true, out: false}, // can't match EOF
		{in: "", val: "[]i", ic: true, out: false},
		{in: "", val: "[^]i", ic: true, iv: true, out: false}, // can't match EOF
		{in: "a", val: "[]", out: false},
		{in: "a", val: "[^]", iv: true, out: true},
		{in: "a", val: "[]i", ic: true, out: false},
		{in: "a", val: "[^]i", ic: true, iv: true, out: true},

		{in: "a", val: "[a]", chars: []rune{'a'}, out: true},
		{in: "a", val: "[a]i", ic: true, chars: []rune{'a'}, out: true},
		{in: "A", val: "[a]i", ic: true, chars: []rune{'a'}, out: true},
		{in: "a", val: "[^a]", chars: []rune{'a'}, iv: true, out: false},
		{in: "A", val: "[^a]i", iv: true, ic: true, chars: []rune{'a'}, out: false},

		{in: "b", val: "[a]", chars: []rune{'a'}, out: false},
		{in: "b", val: "[a]i", ic: true, chars: []rune{'a'}, out: false},
		{in: "B", val: "[a]i", ic: true, chars: []rune{'a'}, out: false},
		{in: "b", val: "[^a]", chars: []rune{'a'}, iv: true, out: true},
		{in: "b", val: "[^a]i", iv: true, ic: true, chars: []rune{'a'}, out: true},
		{in: "B", val: "[^a]i", iv: true, ic: true, chars: []rune{'a'}, out: true},

		{in: "←", val: "[a]", chars: []rune{'a'}, out: false},
		{in: "←", val: "[a]i", ic: true, chars: []rune{'a'}, out: false},
		{in: "←", val: "[a]i", ic: true, chars: []rune{'a'}, out: false},
		{in: "←", val: "[^a]", chars: []rune{'a'}, iv: true, out: true},
		{in: "←", val: "[^a]i", iv: true, ic: true, chars: []rune{'a'}, out: true},
		{in: "←", val: "[^a]i", iv: true, ic: true, chars: []rune{'a'}, out: true},

		{in: "b", val: "[a-c]", ranges: []rune{'a', 'c'}, out: true},
		{in: "B", val: "[a-c]", ranges: []rune{'a', 'c'}, out: false},
		{in: "b", val: "[a-c]i", ic: true, ranges: []rune{'a', 'c'}, out: true},
		{in: "B", val: "[a-c]i", ic: true, ranges: []rune{'a', 'c'}, out: true},
		{in: "b", val: "[^a-c]", ranges: []rune{'a', 'c'}, iv: true, out: false},
		{in: "B", val: "[^a-c]", ranges: []rune{'a', 'c'}, iv: true, out: true},
		{in: "b", val: "[^a-c]i", iv: true, ic: true, ranges: []rune{'a', 'c'}, out: false},
		{in: "B", val: "[^a-c]i", iv: true, ic: true, ranges: []rune{'a', 'c'}, out: false},
		{in: "z", val: "[^a-c]i", iv: true, ic: true, chars: []rune{'a', 'c'}, out: true},

		{in: "∝", val: "[a-c]", ranges: []rune{'a', 'c'}, out: false},
		{in: "∝", val: "[a-c]", ranges: []rune{'a', 'c'}, out: false},
		{in: "∝", val: "[a-c]i", ic: true, ranges: []rune{'a', 'c'}, out: false},
		{in: "∝", val: "[a-c]i", ic: true, ranges: []rune{'a', 'c'}, out: false},
		{in: "∝", val: "[^a-c]", ranges: []rune{'a', 'c'}, iv: true, out: true},
		{in: "∝", val: "[^a-c]", ranges: []rune{'a', 'c'}, iv: true, out: true},
		{in: "∝", val: "[^a-c]i", iv: true, ic: true, ranges: []rune{'a', 'c'}, out: true},
		{in: "∝", val: "[^a-c]i", iv: true, ic: true, ranges: []rune{'a', 'c'}, out: true},
		{in: "∝", val: "[^a-c]i", iv: true, ic: true, chars: []rune{'a', 'c'}, out: true},

		{in: "b", val: "[c-a]", ranges: []rune{'c', 'a'}, out: false},
		{in: "B", val: "[c-a]i", ic: true, ranges: []rune{'c', 'a'}, out: false},
		{in: "B", val: "[^c-a]", iv: true, ranges: []rune{'c', 'a'}, out: true},
		{in: "B", val: "[^c-a]i", ic: true, iv: true, ranges: []rune{'c', 'a'}, out: true},

		{in: "b", val: "[\\pL]", classes: []string{"L"}, out: true},
		{in: "b", val: "[\\pL]i", ic: true, classes: []string{"L"}, out: true},
		{in: "B", val: "[\\pL]i", ic: true, classes: []string{"L"}, out: true},
		{in: "b", val: "[^\\pL]", iv: true, classes: []string{"L"}, out: false},
		{in: "b", val: "[^\\pL]i", iv: true, ic: true, classes: []string{"L"}, out: false},
		{in: "B", val: "[^\\pL]i", iv: true, ic: true, classes: []string{"L"}, out: false},

		{in: "1", val: "[\\pL]", classes: []string{"L"}, out: false},
		{in: "1", val: "[\\pL]i", ic: true, classes: []string{"L"}, out: false},
		{in: "1", val: "[\\pL]i", ic: true, classes: []string{"L"}, out: false},
		{in: "1", val: "[^\\pL]", iv: true, classes: []string{"L"}, out: true},
		{in: "1", val: "[^\\pL]i", iv: true, ic: true, classes: []string{"L"}, out: true},
		{in: "1", val: "[^\\pL]i", iv: true, ic: true, classes: []string{"L"}, out: true},

		{in: "ƛ", val: "[\\pL]", classes: []string{"L"}, out: true},
		{in: "ƛ", val: "[\\pL]i", ic: true, classes: []string{"L"}, out: true},
		{in: "ƛ", val: "[\\pL]i", ic: true, classes: []string{"L"}, out: true},
		{in: "ƛ", val: "[^\\pL]", iv: true, classes: []string{"L"}, out: false},
		{in: "ƛ", val: "[^\\pL]i", iv: true, ic: true, classes: []string{"L"}, out: false},
		{in: "ƛ", val: "[^\\pL]i", iv: true, ic: true, classes: []string{"L"}, out: false},

		{in: "←a", val: "[\\pL]", classes: []string{"L"}, out: false},
		{in: "←a", val: "[\\pL]i", ic: true, classes: []string{"L"}, out: false},
		{in: "←a", val: "[\\pL]i", ic: true, classes: []string{"L"}, out: false},
		{in: "←a", val: "[^\\pL]", iv: true, classes: []string{"L"}, out: true},
		{in: "←a", val: "[^\\pL]i", iv: true, ic: true, classes: []string{"L"}, out: true},
		{in: "←a", val: "[^\\pL]i", iv: true, ic: true, classes: []string{"L"}, out: true},

		{in: "b", val: "[\\p{Latin}]", classes: []string{"Latin"}, out: true},
		{in: "b", val: "[\\p{Latin}]i", ic: true, classes: []string{"Latin"}, out: true},
		{in: "B", val: "[\\p{Latin}]i", ic: true, classes: []string{"Latin"}, out: true},
		{in: "b", val: "[^\\p{Latin}]", iv: true, classes: []string{"Latin"}, out: false},
		{in: "b", val: "[^\\p{Latin}]i", ic: true, iv: true, classes: []string{"Latin"}, out: false},
		{in: "B", val: "[^\\p{Latin}]i", iv: true, ic: true, classes: []string{"Latin"}, out: false},

		{in: "", val: "[^<]", iv: true, chars: []rune{'<'}, out: false},
	}

	var m ϡcharClassMatcher
	for _, tc := range cases {
		pr := testPeekReader{rns: []rune(tc.in)}
		m.chars = tc.chars
		m.ranges = tc.ranges
		m.classes = make([]*unicode.RangeTable, len(tc.classes))
		for i, cl := range tc.classes {
			m.classes[i] = ϡrangeTable(cl)
		}
		m.ignoreCase = tc.ic
		m.inverted = tc.iv
		got := m.match(&pr)

		if got != tc.out {
			t.Errorf("%q with %q: want %t, got %t", tc.in, tc.val, tc.out, got)
		}
	}
}

type testPeekReader struct {
	rns []rune
	ix  int
}

func (pr *testPeekReader) peek() rune {
	if pr.ix >= len(pr.rns) {
		return utf8.RuneError
	}
	return pr.rns[pr.ix]
}

func (pr *testPeekReader) read() {
	pr.ix++
}
