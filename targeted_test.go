package main

import (
	"fmt"
	"reflect"
	"testing"
	"unicode"
)

func TestParseNoRule(t *testing.T) {
	g := &grammar{}
	p := newParser("", []byte(""))
	_, err := p.parse(g)
	if err == nil {
		t.Fatal("want error, got nil")
	}
	el, ok := err.(errList)
	if !ok {
		t.Fatalf("want error type %T, got %T", errList{}, err)
	}
	if len(el) != 1 {
		t.Fatalf("want 1 error, got %d", len(el))
	}
	pe, ok := el[0].(*parserError)
	if !ok {
		t.Fatalf("want single error type %T, got %T", &parserError{}, el[0])
	}
	if pe.Inner != errNoRule {
		t.Fatalf("want error %v, got %v", errNoRule, el[0])
	}
}

func TestParseAnyMatcher(t *testing.T) {
	cases := []struct {
		in  string
		out []byte
	}{
		{"", nil},
		{"a", []byte("a")},
		{"\u2190", []byte("\u2190")},
		{"ab", []byte("a")},
		{"\u2190\U00001100", []byte("\u2190")},
		{"\x0d", []byte("\x0d")},
		{"\xfa", nil},
		{"\nab", []byte("\n")},
	}

	for _, tc := range cases {
		p := newParser("", []byte(tc.in))

		// advance to the first rune
		p.read()

		var want interface{}
		var match bool
		if tc.out != nil {
			want = tc.out
			match = true
		}
		got, ok := p.parseAnyMatcher(&anyMatcher{})
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%q: want %v, got %v", tc.in, tc.out, got)
		}
		if ok != match {
			t.Errorf("%q: want match? %t, got %t", tc.in, match, ok)
		}
		if p.pt.offset != len(tc.out) {
			t.Errorf("%q: want offset %d, got %d", tc.in, len(tc.out), p.pt.offset)
		}
	}
}

func TestParseLitMatcher(t *testing.T) {
	cases := []struct {
		in  string
		lit string
		ic  bool
		out []byte
	}{
		{"", "", false, []byte{}}, // empty literal always matches
		{"", "", true, []byte{}},  // empty literal always matches
		{"a", "", false, []byte{}},
		{"a", "", true, []byte{}},
		{"a", "a", false, []byte("a")},
		{"a", "a", true, []byte("a")},
		{"a", "A", false, nil},
		{"a", "a", true, []byte("a")}, // ignored case literal is always generated lowercase
		{"A", "a", true, []byte("A")},
		{"b", "a", false, nil},
		{"b", "a", true, nil},
		{"abc", "ab", false, []byte("ab")},
		{"abc", "ab", true, []byte("ab")},
		{"ab", "abc", false, nil},
		{"ab", "abc", true, nil},
		{"\u2190a", "\u2190", false, []byte("\u2190")},
		{"\u2190a", "\u2190", true, []byte("\u2190")},
		{"\n", "\n", false, []byte("\n")},
		{"\n", "\n", true, []byte("\n")},
		{"\na", "\n", false, []byte("\n")},
		{"\na", "\n", true, []byte("\n")},
	}

	for _, tc := range cases {
		p := newParser("", []byte(tc.in))

		// advance to the first rune
		p.read()

		var want interface{}
		var match bool
		if tc.out != nil {
			match = true
			want = tc.out
		}
		lbl := fmt.Sprintf("%q (%t): %q", tc.lit, tc.ic, tc.in)

		got, ok := p.parseLitMatcher(&litMatcher{val: tc.lit, ignoreCase: tc.ic})
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: want %v, got %v", lbl, tc.out, got)
		}
		if ok != match {
			t.Errorf("%s: want match? %t, got %t", lbl, match, ok)
		}
		if p.pt.offset != len(tc.out) {
			t.Errorf("%s: want offset %d, got %d", lbl, len(tc.out), p.pt.offset)
		}
	}
}

func TestParseCharClassMatcher(t *testing.T) {
	cases := []struct {
		in      string
		val     string
		chars   []rune
		ranges  []rune
		classes []string
		ic      bool
		iv      bool
		out     []byte
	}{
		{in: "", val: "[]", out: nil}, // empty char class means no char matches
		{in: "", val: "[^]", iv: true, out: []byte{}},
		{in: "", val: "[]i", ic: true, out: nil},
		{in: "", val: "[^]i", ic: true, iv: true, out: []byte{}},
		{in: "a", val: "[]", out: nil},
		{in: "a", val: "[^]", iv: true, out: []byte("a")},
		{in: "a", val: "[]i", ic: true, out: nil},
		{in: "a", val: "[^]i", ic: true, iv: true, out: []byte("a")},

		{in: "a", val: "[a]", chars: []rune{'a'}, out: []byte("a")},
		{in: "a", val: "[a]i", ic: true, chars: []rune{'a'}, out: []byte("a")},
		{in: "A", val: "[a]i", ic: true, chars: []rune{'a'}, out: []byte("A")},
		{in: "a", val: "[^a]", chars: []rune{'a'}, iv: true, out: nil},
		{in: "A", val: "[^a]i", iv: true, ic: true, chars: []rune{'a'}, out: nil},

		{in: "b", val: "[a]", chars: []rune{'a'}, out: nil},
		{in: "b", val: "[a]i", ic: true, chars: []rune{'a'}, out: nil},
		{in: "B", val: "[a]i", ic: true, chars: []rune{'a'}, out: nil},
		{in: "b", val: "[^a]", chars: []rune{'a'}, iv: true, out: []byte("b")},
		{in: "b", val: "[^a]i", iv: true, ic: true, chars: []rune{'a'}, out: []byte("b")},
		{in: "B", val: "[^a]i", iv: true, ic: true, chars: []rune{'a'}, out: []byte("B")},

		{in: "←", val: "[a]", chars: []rune{'a'}, out: nil},
		{in: "←", val: "[a]i", ic: true, chars: []rune{'a'}, out: nil},
		{in: "←", val: "[a]i", ic: true, chars: []rune{'a'}, out: nil},
		{in: "←", val: "[^a]", chars: []rune{'a'}, iv: true, out: []byte("←")},
		{in: "←", val: "[^a]i", iv: true, ic: true, chars: []rune{'a'}, out: []byte("←")},
		{in: "←", val: "[^a]i", iv: true, ic: true, chars: []rune{'a'}, out: []byte("←")},

		{in: "b", val: "[a-c]", ranges: []rune{'a', 'c'}, out: []byte("b")},
		{in: "B", val: "[a-c]", ranges: []rune{'a', 'c'}, out: nil},
		{in: "b", val: "[a-c]i", ic: true, ranges: []rune{'a', 'c'}, out: []byte("b")},
		{in: "B", val: "[a-c]i", ic: true, ranges: []rune{'a', 'c'}, out: []byte("B")},
		{in: "b", val: "[^a-c]", ranges: []rune{'a', 'c'}, iv: true, out: nil},
		{in: "B", val: "[^a-c]", ranges: []rune{'a', 'c'}, iv: true, out: []byte("B")},
		{in: "b", val: "[^a-c]i", iv: true, ic: true, ranges: []rune{'a', 'c'}, out: nil},
		{in: "B", val: "[^a-c]i", iv: true, ic: true, ranges: []rune{'a', 'c'}, out: nil},
		{in: "z", val: "[^a-c]i", iv: true, ic: true, chars: []rune{'a', 'c'}, out: []byte("z")},

		{in: "∝", val: "[a-c]", ranges: []rune{'a', 'c'}, out: nil},
		{in: "∝", val: "[a-c]", ranges: []rune{'a', 'c'}, out: nil},
		{in: "∝", val: "[a-c]i", ic: true, ranges: []rune{'a', 'c'}, out: nil},
		{in: "∝", val: "[a-c]i", ic: true, ranges: []rune{'a', 'c'}, out: nil},
		{in: "∝", val: "[^a-c]", ranges: []rune{'a', 'c'}, iv: true, out: []byte("∝")},
		{in: "∝", val: "[^a-c]", ranges: []rune{'a', 'c'}, iv: true, out: []byte("∝")},
		{in: "∝", val: "[^a-c]i", iv: true, ic: true, ranges: []rune{'a', 'c'}, out: []byte("∝")},
		{in: "∝", val: "[^a-c]i", iv: true, ic: true, ranges: []rune{'a', 'c'}, out: []byte("∝")},
		{in: "∝", val: "[^a-c]i", iv: true, ic: true, chars: []rune{'a', 'c'}, out: []byte("∝")},

		{in: "b", val: "[c-a]", ranges: []rune{'c', 'a'}, out: nil},
		{in: "B", val: "[c-a]i", ic: true, ranges: []rune{'c', 'a'}, out: nil},
		{in: "B", val: "[^c-a]", iv: true, ranges: []rune{'c', 'a'}, out: []byte("B")},
		{in: "B", val: "[^c-a]i", ic: true, iv: true, ranges: []rune{'c', 'a'}, out: []byte("B")},

		{in: "b", val: "[\\pL]", classes: []string{"L"}, out: []byte("b")},
		{in: "b", val: "[\\pL]i", ic: true, classes: []string{"L"}, out: []byte("b")},
		{in: "B", val: "[\\pL]i", ic: true, classes: []string{"L"}, out: []byte("B")},
		{in: "b", val: "[^\\pL]", iv: true, classes: []string{"L"}, out: nil},
		{in: "b", val: "[^\\pL]i", iv: true, ic: true, classes: []string{"L"}, out: nil},
		{in: "B", val: "[^\\pL]i", iv: true, ic: true, classes: []string{"L"}, out: nil},

		{in: "1", val: "[\\pL]", classes: []string{"L"}, out: nil},
		{in: "1", val: "[\\pL]i", ic: true, classes: []string{"L"}, out: nil},
		{in: "1", val: "[\\pL]i", ic: true, classes: []string{"L"}, out: nil},
		{in: "1", val: "[^\\pL]", iv: true, classes: []string{"L"}, out: []byte("1")},
		{in: "1", val: "[^\\pL]i", iv: true, ic: true, classes: []string{"L"}, out: []byte("1")},
		{in: "1", val: "[^\\pL]i", iv: true, ic: true, classes: []string{"L"}, out: []byte("1")},

		{in: "ƛ", val: "[\\pL]", classes: []string{"L"}, out: []byte("ƛ")},
		{in: "ƛ", val: "[\\pL]i", ic: true, classes: []string{"L"}, out: []byte("ƛ")},
		{in: "ƛ", val: "[\\pL]i", ic: true, classes: []string{"L"}, out: []byte("ƛ")},
		{in: "ƛ", val: "[^\\pL]", iv: true, classes: []string{"L"}, out: nil},
		{in: "ƛ", val: "[^\\pL]i", iv: true, ic: true, classes: []string{"L"}, out: nil},
		{in: "ƛ", val: "[^\\pL]i", iv: true, ic: true, classes: []string{"L"}, out: nil},

		{in: "←a", val: "[\\pL]", classes: []string{"L"}, out: nil},
		{in: "←a", val: "[\\pL]i", ic: true, classes: []string{"L"}, out: nil},
		{in: "←a", val: "[\\pL]i", ic: true, classes: []string{"L"}, out: nil},
		{in: "←a", val: "[^\\pL]", iv: true, classes: []string{"L"}, out: []byte("←")},
		{in: "←a", val: "[^\\pL]i", iv: true, ic: true, classes: []string{"L"}, out: []byte("←")},
		{in: "←a", val: "[^\\pL]i", iv: true, ic: true, classes: []string{"L"}, out: []byte("←")},

		{in: "b", val: "[\\p{Latin}]", classes: []string{"Latin"}, out: []byte("b")},
		{in: "b", val: "[\\p{Latin}]i", ic: true, classes: []string{"Latin"}, out: []byte("b")},
		{in: "B", val: "[\\p{Latin}]i", ic: true, classes: []string{"Latin"}, out: []byte("B")},
		{in: "b", val: "[^\\p{Latin}]", iv: true, classes: []string{"Latin"}, out: nil},
		{in: "b", val: "[^\\p{Latin}]i", ic: true, iv: true, classes: []string{"Latin"}, out: nil},
		{in: "B", val: "[^\\p{Latin}]i", iv: true, ic: true, classes: []string{"Latin"}, out: nil},
	}

	for _, tc := range cases {
		p := newParser("", []byte(tc.in))

		// advance to the first rune
		p.read()

		var want interface{}
		var match bool
		if tc.out != nil {
			want = tc.out
			match = true
		}
		lbl := fmt.Sprintf("%q (%t-%t): %q", tc.val, tc.ic, tc.iv, tc.in)

		classes := make([]*unicode.RangeTable, len(tc.classes))
		for i, c := range tc.classes {
			classes[i] = rangeTable(c)
		}

		got, ok := p.parseCharClassMatcher(&charClassMatcher{
			val:        tc.val,
			chars:      tc.chars,
			ranges:     tc.ranges,
			classes:    classes,
			ignoreCase: tc.ic,
			inverted:   tc.iv,
		})
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: want %v, got %v", lbl, tc.out, got)
		}
		if ok != match {
			t.Errorf("%s: want match? %t, got %t", lbl, match, ok)
		}
		if p.pt.offset != len(tc.out) {
			t.Errorf("%s: want offset %d, got %d", lbl, len(tc.out), p.pt.offset)
		}
	}
}

func TestParseZeroOrOneExpr(t *testing.T) {
	cases := []struct {
		in  string
		lit string
		out []byte
	}{
		{"", "", []byte{}},
		{"", "a", nil},
		{"a", "a", []byte("a")},
		{"a", "b", nil},
		{"abc", "ab", []byte("ab")},
		{"ab", "abc", nil},
	}

	for _, tc := range cases {
		p := newParser("", []byte(tc.in))

		// advance to the first rune
		p.read()

		var want interface{}
		if tc.out != nil {
			want = tc.out
		}
		lbl := fmt.Sprintf("%q: %q", tc.lit, tc.in)

		got, ok := p.parseZeroOrOneExpr(&zeroOrOneExpr{expr: &litMatcher{val: tc.lit}})
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%q: want %v, got %v", lbl, tc.out, got)
		}
		// zero or one always matches
		if !ok {
			t.Errorf("%s: want match, got %t", lbl, ok)
		}
		if p.pt.offset != len(tc.out) {
			t.Errorf("%s: want offset %d, got %d", lbl, len(tc.out), p.pt.offset)
		}
	}
}

func TestParseZeroOrMoreExpr(t *testing.T) {
	cases := []struct {
		in  string
		lit string
		out []string
	}{
		// ""* is a pathological case - the empty string always matches, so this
		// is an infinite loop. Not fixing it, because semantically this seems
		// correct.
		// {"", "", []byte{}},

		{"", "a", nil},
		{"a", "a", []string{"a"}},
		{"a", "b", nil},
		{"abc", "ab", []string{"ab"}},
		{"ab", "abc", nil},

		{"aab", "a", []string{"a", "a"}},
		{"bba", "a", nil},
		{"bba", "b", []string{"b", "b"}},
		{"bba", "bb", []string{"bb"}},
		{"aaaaab", "aa", []string{"aa", "aa"}},
		{"aaaaab", "a", []string{"a", "a", "a", "a", "a"}},
	}

	for _, tc := range cases {
		p := newParser("", []byte(tc.in))

		// advance to the first rune
		p.read()

		want := make([]interface{}, len(tc.out))
		for i, v := range tc.out {
			want[i] = []byte(v)
		}
		if tc.out == nil {
			want = nil
		}
		lbl := fmt.Sprintf("%q: %q", tc.lit, tc.in)

		got, ok := p.parseZeroOrMoreExpr(&zeroOrMoreExpr{expr: &litMatcher{val: tc.lit}})
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: want %#v, got %#v", lbl, want, got)
		}
		// zero or more always matches
		if !ok {
			t.Errorf("%s: want match, got %t", lbl, ok)
		}
		wantOffset := 0
		for _, s := range tc.out {
			wantOffset += len(s)
		}
		if p.pt.offset != wantOffset {
			t.Errorf("%s: want offset %d, got %d", lbl, wantOffset, p.pt.offset)
		}
	}
}

func TestParseOneOrMoreExpr(t *testing.T) {
	cases := []struct {
		in  string
		lit string
		out []string
	}{
		// ""+ is a pathological case - the empty string always matches, so this
		// is an infinite loop. Not fixing it, because semantically this seems
		// correct.
		//{"", "", []string{}},

		{"", "a", nil},
		{"a", "a", []string{"a"}},
		{"a", "b", nil},
		{"abc", "ab", []string{"ab"}},
		{"ab", "abc", nil},

		{"aab", "a", []string{"a", "a"}},
		{"bba", "a", nil},
		{"bba", "b", []string{"b", "b"}},
		{"bba", "bb", []string{"bb"}},
		{"aaaaab", "aa", []string{"aa", "aa"}},
		{"aaaaab", "a", []string{"a", "a", "a", "a", "a"}},
	}

	for _, tc := range cases {
		p := newParser("", []byte(tc.in))

		// advance to the first rune
		p.read()

		var want interface{}
		var match bool
		if tc.out != nil {
			vals := make([]interface{}, len(tc.out))
			for i, v := range tc.out {
				vals[i] = []byte(v)
			}
			want = vals
			match = true
		}
		lbl := fmt.Sprintf("%q: %q", tc.lit, tc.in)

		got, ok := p.parseOneOrMoreExpr(&oneOrMoreExpr{expr: &litMatcher{val: tc.lit}})
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: want %#v, got %#v", lbl, want, got)
		}
		if ok != match {
			t.Errorf("%s: want match? %t, got %t", lbl, match, ok)
		}
		wantOffset := 0
		for _, s := range tc.out {
			wantOffset += len(s)
		}
		if p.pt.offset != wantOffset {
			t.Errorf("%s: want offset %d, got %d", lbl, wantOffset, p.pt.offset)
		}
	}
}
