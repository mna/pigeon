package vm

import (
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
	"unicode"

	"github.com/PuerkitoBio/pigeon/ast"
	"github.com/PuerkitoBio/pigeon/bootstrap"
)

func TestRun(t *testing.T) {
	cases := []struct {
		grammar string
		input   string
		want    interface{}
		err     error
	}{
		{`A = 'a'`, "a", []byte("a"), nil},
		{`A = 'a'`, "b", nil, errNoMatch},
		{`A = "ab"`, "a", nil, errNoMatch},
		{`A = "ab"`, "b", nil, errNoMatch},
		{`A = "ab"`, "ab", []byte("ab"), nil},
		{`A = "ab"`, "abb", []byte("ab"), nil},

		{`A = 'a'*`, "", []interface{}(nil), nil},
		{`A = 'a'*`, "a", []interface{}{[]byte("a")}, nil},
		{`A = 'a'*`, "aa", []interface{}{[]byte("a"), []byte("a")}, nil},
		{`A = 'a'*`, "aab", []interface{}{[]byte("a"), []byte("a")}, nil},
		{`A = 'a'*`, "baa", []interface{}(nil), nil},

		{`A = 'a'+`, "", nil, errNoMatch},
		{`A = 'a'+`, "a", []interface{}{[]byte("a")}, nil},
		{`A = 'a'+`, "aa", []interface{}{[]byte("a"), []byte("a")}, nil},
		{`A = 'a'+`, "aab", []interface{}{[]byte("a"), []byte("a")}, nil},
		{`A = 'a'+`, "baa", nil, errNoMatch},

		{`A = 'a'?`, "", nil, nil},
		{`A = 'a'?`, "a", []byte("a"), nil},
		{`A = 'a'?`, "aa", []byte("a"), nil},
		{`A = 'a'?`, "aab", []byte("a"), nil},
		{`A = 'a'?`, "baa", nil, nil},

		{`A = 'a' 'b'`, "", nil, errNoMatch},
		{`A = 'a' 'b'`, "a", nil, errNoMatch},
		{`A = 'a' 'b'`, "ab", []interface{}{[]byte("a"), []byte("b")}, nil},
		{`A = 'a' 'b'`, "aab", nil, errNoMatch},
		{`A = 'a' 'b'`, "baa", nil, errNoMatch},

		{`A = 'a' / 'b'`, "", nil, errNoMatch},
		{`A = 'a' / 'b'`, "a", []byte("a"), nil},
		{`A = 'a' / 'b'`, "ab", []byte("a"), nil},
		{`A = 'a' / 'b'`, "aab", []byte("a"), nil},
		{`A = 'a' / 'b'`, "baa", []byte("b"), nil},

		{"A = B\nB= 'a'", "", nil, errNoMatch},
		{"A = B\nB= 'a'", "a", []byte("a"), nil},
		{"A = B\nB = 'a'", "ab", []byte("a"), nil},
		{"A = B\nB = 'a'", "aab", []byte("a"), nil},
		{"A = B\nB = 'a'", "baa", nil, errNoMatch},

		{`A = 'a' &'b'`, "", nil, errNoMatch},
		{`A = 'a' &'b'`, "a", nil, errNoMatch},
		{`A = 'a' &'b'`, "ab", []interface{}{[]byte("a"), nil}, nil},
		{`A = 'a' &'b'`, "aab", nil, errNoMatch},
		{`A = 'a' &'b'`, "baa", nil, errNoMatch},

		{`A = 'a' !'b'`, "", nil, errNoMatch},
		{`A = 'a' !'b'`, "a", []interface{}{[]byte("a"), nil}, nil},
		{`A = 'a' !'b'`, "ab", nil, errNoMatch},
		{`A = 'a' !'b'`, "aab", []interface{}{[]byte("a"), nil}, nil},
		{`A = 'a' !'b'`, "baa", nil, errNoMatch},
	}
	for i, tc := range cases {
		gr, err := bootstrap.NewParser().Parse("", strings.NewReader(tc.grammar))
		if err != nil {
			t.Errorf("%d: parse error: %v", i, err)
			continue
		}

		pg, err := NewGenerator(ioutil.Discard).toProgram(gr)
		if err != nil {
			t.Errorf("%d: generator error: %v", i, err)
			continue
		}

		ϡtheProgram = toϡprogram(t, pg, amockRetCode, bmockRetTrueIfT)
		got, err := Parse("", []byte(tc.input), Debug(true), Recover(false))
		if (err != nil) != (tc.err != nil) {
			t.Errorf("%d: want error? %t, got %v", i, tc.err != nil, err)
			continue
		} else if tc.err != nil {
			pe := err.(errList)[0].(*parserError)
			if tc.err != pe.Inner {
				t.Errorf("%d: want error %v, got %v", i, tc.err, err)
				continue
			}
		}

		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("%d: want %#v, got %#v", i, tc.want, got)
		}
	}
}

func amockRetCode(ti *thunkInfo) func(*ϡvm) (interface{}, error) {
	return func(v *ϡvm) (interface{}, error) {
		return ti.Code, nil
	}
}

func bmockRetTrueIfT(ti *thunkInfo) func(*ϡvm) (bool, error) {
	return func(v *ϡvm) (bool, error) {
		return ti.Code == "T", nil
	}
}

func toϡprogram(t *testing.T, pg *program,
	amock func(*thunkInfo) func(*ϡvm) (interface{}, error),
	bmock func(*thunkInfo) func(*ϡvm) (bool, error)) *ϡprogram {

	vmpg := ϡprogram{
		instrs:      pg.Instrs,
		ss:          pg.Ss,
		instrToRule: pg.InstrToRule,
	}

	// convert matchers
	vmpg.ms = make([]ϡmatcher, len(pg.Ms))
	for i, m := range pg.Ms {
		switch m := m.(type) {
		case *ast.AnyMatcher:
			vmpg.ms[i] = ϡanyMatcher{}
		case *ast.LitMatcher:
			if m.IgnoreCase {
				m.Val = strings.ToLower(m.Val)
			}
			vmpg.ms[i] = ϡstringMatcher{
				ignoreCase: m.IgnoreCase,
				value:      m.Val,
			}
		case *ast.CharClassMatcher:
			if m.IgnoreCase {
				for j, rn := range m.Chars {
					m.Chars[j] = unicode.ToLower(rn)
				}
				for j, rn := range m.Ranges {
					m.Ranges[j] = unicode.ToLower(rn)
				}
			}
			classes := make([]*unicode.RangeTable, len(m.UnicodeClasses))
			for j, cl := range m.UnicodeClasses {
				classes[j] = ϡrangeTable(cl)
			}
			vmpg.ms[i] = ϡcharClassMatcher{
				ignoreCase: m.IgnoreCase,
				inverted:   m.Inverted,
				chars:      m.Chars,
				ranges:     m.Ranges,
				classes:    classes,
			}
		}
	}

	// convert As
	vmpg.as = make([]func(*ϡvm) (interface{}, error), len(pg.As))
	for j, a := range pg.As {
		vmpg.as[j] = amock(a)
	}
	// convert Bs
	vmpg.bs = make([]func(*ϡvm) (bool, error), len(pg.Bs))
	for j, b := range pg.Bs {
		vmpg.bs[j] = bmock(b)
	}
	return &vmpg
}
