package main

import (
	"fmt"
	"io"
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

func TestParseSeqExpr(t *testing.T) {
	cases := []struct {
		in   string
		lits []string
		out  []string
	}{
		{"", nil, []string{}}, // empty seq (impossible case via the parser) always matches
		{"", []string{"a"}, nil},
		{"a", []string{"a"}, []string{"a"}},
		{"a", []string{"a", "b"}, nil},
		{"abc", []string{"a", "b"}, []string{"a", "b"}},
		{"abc", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"ab", []string{"a", "b", "c"}, nil},
	}

	for _, tc := range cases {
		p := newParser("", []byte(tc.in))

		// advance to the first rune
		p.read()

		var want interface{}
		var match bool
		if tc.out != nil {
			var vals []interface{}
			for _, v := range tc.out {
				vals = append(vals, []byte(v))
			}
			want = vals
			match = true
		}
		lbl := fmt.Sprintf("%v: %q", tc.lits, tc.in)

		lits := make([]interface{}, len(tc.lits))
		for i, l := range tc.lits {
			lits[i] = &litMatcher{val: l}
		}

		got, ok := p.parseSeqExpr(&seqExpr{exprs: lits})
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

func TestParseRuleRefExpr(t *testing.T) {
	p := newParser("", []byte(""))

	func() {
		defer func() {
			if e := recover(); e != nil {
				return
			}
			t.Fatal("want panic, got none")
		}()
		p.parseRuleRefExpr(&ruleRefExpr{})
	}()

	p.parseRuleRefExpr(&ruleRefExpr{name: "a"})
	if p.errs.err() == nil {
		t.Fatal("want error, got none")
	}
}

func TestParseNotExpr(t *testing.T) {
	cases := []struct {
		in    string
		lit   string
		match bool
	}{
		{"", "", false},
		{"", "a", true},
		{"a", "a", false},
		{"b", "a", true},
		{"ab", "a", false},
		{"ab", "ab", false},
		{"ab", "abc", true},
		{"abc", "abc", false},
		{"abc", "ab", false},
		{"abc", "ac", true},
	}
	for _, tc := range cases {
		p := newParser("", []byte(tc.in))

		// advance to the first rune
		p.read()

		lbl := fmt.Sprintf("%q: %q", tc.lit, tc.in)

		_, ok := p.parseNotExpr(&notExpr{expr: &litMatcher{val: tc.lit}})
		if ok != tc.match {
			t.Errorf("%s: want match? %t, got %t", lbl, tc.match, ok)
		}
		if p.pt.offset != 0 {
			t.Errorf("%s: want offset %d, got %d", lbl, 0, p.pt.offset)
		}
	}
}

func TestParseAndExpr(t *testing.T) {
	cases := []struct {
		in    string
		lit   string
		match bool
	}{
		{"", "", true},
		{"", "a", false},
		{"a", "a", true},
		{"b", "a", false},
		{"ab", "a", true},
		{"ab", "ab", true},
		{"ab", "abc", false},
		{"abc", "abc", true},
		{"abc", "ab", true},
		{"abc", "ac", false},
	}
	for _, tc := range cases {
		p := newParser("", []byte(tc.in))

		// advance to the first rune
		p.read()

		lbl := fmt.Sprintf("%q: %q", tc.lit, tc.in)

		_, ok := p.parseAndExpr(&andExpr{expr: &litMatcher{val: tc.lit}})
		if ok != tc.match {
			t.Errorf("%s: want match? %t, got %t", lbl, tc.match, ok)
		}
		if p.pt.offset != 0 {
			t.Errorf("%s: want offset %d, got %d", lbl, 0, p.pt.offset)
		}
	}
}

func TestParseNotCodeExpr(t *testing.T) {
	cases := []struct {
		in  string
		b   bool
		err error
	}{
		{"", true, nil},
		{"", true, io.EOF},
		{"", false, nil},
		{"", false, io.EOF},
		{"a", true, nil},
		{"a", true, io.EOF},
		{"a", false, nil},
		{"a", false, io.EOF},
	}

	for _, tc := range cases {
		fn := func(_ *parser) (bool, error) {
			return tc.b, tc.err
		}
		p := newParser("", []byte(tc.in))

		// advance to the first rune
		p.read()

		lbl := fmt.Sprintf("%q: %t-%t", tc.in, tc.b, tc.err == nil)

		_, ok := p.parseNotCodeExpr(&notCodeExpr{run: fn})
		if ok != !tc.b {
			t.Errorf("%s: want match? %t, got %t", lbl, !tc.b, ok)
		}

		el := *p.errs
		wantn := 0
		if tc.err != nil {
			wantn = 1
		}
		if len(el) != wantn {
			t.Errorf("%s: want %d error, got %d", lbl, wantn, len(el))
		} else if wantn == 1 {
			ie := el[0].(*parserError).Inner
			if ie != tc.err {
				t.Errorf("%s: want error %v, got %v", lbl, tc.err, ie)
			}
		}

		if p.pt.offset != 0 {
			t.Errorf("%s: want offset %d, got %d", lbl, 0, p.pt.offset)
		}
	}
}

func TestParseAndCodeExpr(t *testing.T) {
	cases := []struct {
		in  string
		b   bool
		err error
	}{
		{"", true, nil},
		{"", true, io.EOF},
		{"", false, nil},
		{"", false, io.EOF},
		{"a", true, nil},
		{"a", true, io.EOF},
		{"a", false, nil},
		{"a", false, io.EOF},
	}

	for _, tc := range cases {
		fn := func(_ *parser) (bool, error) {
			return tc.b, tc.err
		}
		p := newParser("", []byte(tc.in))

		// advance to the first rune
		p.read()

		lbl := fmt.Sprintf("%q: %t-%t", tc.in, tc.b, tc.err == nil)

		_, ok := p.parseAndCodeExpr(&andCodeExpr{run: fn})
		if ok != tc.b {
			t.Errorf("%s: want match? %t, got %t", lbl, tc.b, ok)
		}

		el := *p.errs
		wantn := 0
		if tc.err != nil {
			wantn = 1
		}
		if len(el) != wantn {
			t.Errorf("%s: want %d error, got %d", lbl, wantn, len(el))
		} else if wantn == 1 {
			ie := el[0].(*parserError).Inner
			if ie != tc.err {
				t.Errorf("%s: want error %v, got %v", lbl, tc.err, ie)
			}
		}

		if p.pt.offset != 0 {
			t.Errorf("%s: want offset %d, got %d", lbl, 0, p.pt.offset)
		}
	}
}

func TestParseLabeledExpr(t *testing.T) {
	cases := []struct {
		in  string
		lit string
		out []byte
	}{
		{"", "", []byte{}},
		{"", "a", nil},
		{"a", "a", []byte("a")},
		{"a", "ab", nil},
		{"ab", "a", []byte("a")},
		{"ab", "ab", []byte("ab")},
		{"ab", "abc", nil},
		{"abc", "ab", []byte("ab")},
	}

	for _, tc := range cases {
		p := newParser("", []byte(tc.in))

		// advance to the first rune
		p.read()
		p.pushV()

		var want interface{}
		var match bool
		if tc.out != nil {
			match = true
			want = tc.out
		}
		lbl := fmt.Sprintf("%q: %q", tc.lit, tc.in)

		got, ok := p.parseLabeledExpr(&labeledExpr{label: "l", expr: &litMatcher{val: tc.lit}})
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: want %v, got %v", lbl, tc.out, got)
		}
		if ok != match {
			t.Errorf("%s: want match? %t, got %t", lbl, match, ok)
		} else {
			// must be 1 var set on the stack
			if len(p.vstack) != 1 {
				t.Errorf("%s: want %d var sets on the stack, got %d", lbl, 1, len(p.vstack))
			} else {
				vs := p.vstack[0]
				if !reflect.DeepEqual(vs["l"], got) {
					t.Errorf("%s: want %v on the stack for this label, got %v", lbl, got, vs["l"])
				}
			}
		}

		if p.pt.offset != len(tc.out) {
			t.Errorf("%s: want offset %d, got %d", lbl, len(tc.out), p.pt.offset)
		}
	}
}

func TestParseChoiceExpr(t *testing.T) {
	cases := []struct {
		in   string
		lits []string
		out  []byte
	}{
		{"", nil, nil}, // empty choice (impossible case via the parser)

		{"", []string{"a"}, nil},
		{"a", []string{"a"}, []byte("a")},
		{"a", []string{"b"}, nil},
		{"ab", []string{"b"}, nil},
		{"ba", []string{"b"}, []byte("b")},
		{"a", []string{"a", "b"}, []byte("a")},
		{"a", []string{"b", "a"}, []byte("a")},
		{"ab", []string{"a", "b"}, []byte("a")},
		{"ab", []string{"b", "a"}, []byte("a")},
		{"cb", []string{"a", "b"}, nil},
		{"cb", []string{"b", "a"}, nil},
		{"abcd", []string{"abc", "ab", "a"}, []byte("abc")},
		{"abcd", []string{"a", "ab", "abc"}, []byte("a")},
		{"bcd", []string{"a", "ab", "abc"}, nil},
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
		lbl := fmt.Sprintf("%v: %q", tc.lits, tc.in)

		lits := make([]interface{}, len(tc.lits))
		for i, l := range tc.lits {
			lits[i] = &litMatcher{val: l}
		}

		got, ok := p.parseChoiceExpr(&choiceExpr{alternatives: lits})
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: want %#v, got %#v", lbl, want, got)
		}
		if ok != match {
			t.Errorf("%s: want match? %t, got %t", lbl, match, ok)
		}
		if p.pt.offset != len(tc.out) {
			t.Errorf("%s: want offset %d, got %d", lbl, len(tc.out), p.pt.offset)
		}
	}
}

func TestParseActionExpr(t *testing.T) {
	cases := []struct {
		in  string
		lit string
		v   interface{}
		err error
	}{
		{"", "", 1, nil}, // empty string always matches
		{"", "", 1, io.EOF},
		{"", "a", nil, nil},
		{"a", "a", 1, nil},
		{"a", "a", 1, io.EOF},
		{"ab", "a", 1, nil},
		{"ab", "a", 1, io.EOF},
		{"ba", "a", nil, nil},
	}

	for _, tc := range cases {
		called := false
		fn := func(_ *parser) (interface{}, error) {
			called = true
			return tc.v, tc.err
		}
		p := newParser("", []byte(tc.in))

		// advance to the first rune
		p.read()

		lbl := fmt.Sprintf("%q: %q", tc.in, tc.lit)

		match := tc.v != nil

		got, ok := p.parseActionExpr(&actionExpr{run: fn, expr: &litMatcher{val: tc.lit}})
		if ok != match {
			t.Errorf("%s: want match? %t, got %t", lbl, match, ok)
		}
		if !reflect.DeepEqual(got, tc.v) {
			t.Errorf("%s: want %#v, got %#v", lbl, tc.v, got)
		}
		if match != called {
			t.Errorf("%s: want action code to be called? %t, got %t", lbl, match, called)
		}

		el := *p.errs
		wantn := 0
		if tc.err != nil {
			wantn = 1
		}
		if len(el) != wantn {
			t.Errorf("%s: want %d error, got %d", lbl, wantn, len(el))
		} else if wantn == 1 {
			ie := el[0].(*parserError).Inner
			if ie != tc.err {
				t.Errorf("%s: want error %v, got %v", lbl, tc.err, ie)
			}
		}

		wantOffset := 0
		if match {
			wantOffset = len(tc.lit)
		}
		if p.pt.offset != wantOffset {
			t.Errorf("%s: want offset %d, got %d", lbl, wantOffset, p.pt.offset)
		}
	}
}
