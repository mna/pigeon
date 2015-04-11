package main

import (
	"fmt"
	"reflect"
	"testing"
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
		in    string
		out   []byte
		match bool
	}{
		{"", nil, false},
		{"a", []byte("a"), true},
		{"\u2190", []byte("\u2190"), true},
		{"ab", []byte("a"), true},
		{"\u2190\U00001100", []byte("\u2190"), true},
		{"\x0d", []byte("\x0d"), true},
		{"\xfa", nil, false},
		{"\nab", []byte("\n"), true},
	}

	for _, tc := range cases {
		p := newParser("", []byte(tc.in))

		// advance to the first rune
		p.read()

		var want interface{}
		if tc.out != nil {
			want = tc.out
		}
		got, ok := p.parseAnyMatcher(&anyMatcher{})
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%q: want %v, got %v", tc.in, tc.out, got)
		}
		if ok != tc.match {
			t.Errorf("%q: want match? %t, got %t", tc.in, tc.match, ok)
		}
		if p.pt.offset != len(tc.out) {
			t.Errorf("%q: want offset %d, got %d", tc.in, len(tc.out), p.pt.offset)
		}
	}
}

func TestParseLitMatcher(t *testing.T) {
	cases := []struct {
		in    string
		lit   string
		ic    bool
		out   []byte
		match bool
	}{
		{"", "", false, []byte{}, true}, // empty literal always matches
		{"", "", true, []byte{}, true},  // empty literal always matches
		{"a", "", false, []byte{}, true},
		{"a", "", true, []byte{}, true},
		{"a", "a", false, []byte("a"), true},
		{"a", "a", true, []byte("a"), true},
		{"a", "A", false, nil, false},
		{"a", "a", true, []byte("a"), true}, // ignored case literal is always generated lowercase
		{"A", "a", true, []byte("A"), true},
		{"b", "a", false, nil, false},
		{"b", "a", true, nil, false},
		{"abc", "ab", false, []byte("ab"), true},
		{"abc", "ab", true, []byte("ab"), true},
		{"ab", "abc", false, nil, false},
		{"ab", "abc", true, nil, false},
		{"\u2190a", "\u2190", false, []byte("\u2190"), true},
		{"\u2190a", "\u2190", true, []byte("\u2190"), true},
	}

	for _, tc := range cases {
		p := newParser("", []byte(tc.in))

		// advance to the first rune
		p.read()

		var want interface{}
		if tc.out != nil {
			want = tc.out
		}
		lbl := fmt.Sprintf("%q (%t): %q", tc.lit, tc.ic, tc.in)

		got, ok := p.parseLitMatcher(&litMatcher{val: tc.lit, ignoreCase: tc.ic})
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: want %v, got %v", lbl, tc.out, got)
		}
		if ok != tc.match {
			t.Errorf("%s: want match? %t, got %t", lbl, tc.match, ok)
		}
		if p.pt.offset != len(tc.out) {
			t.Errorf("%s: want offset %d, got %d", lbl, len(tc.out), p.pt.offset)
		}
	}
}
