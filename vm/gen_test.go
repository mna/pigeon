package vm

import (
	"io/ioutil"
	"strconv"
	"strings"
	"testing"

	"github.com/PuerkitoBio/pigeon/ast"
	"github.com/PuerkitoBio/pigeon/bootstrap"
)

type testProgram struct {
	Init        string
	Instrs      []ϡinstr
	Ms          []string
	Ss          []string
	As          []*thunkInfo
	Bs          []*thunkInfo
	InstrToRule []int
}

func TestGenerateProgram(t *testing.T) {
	cases := []struct {
		in  string
		out *testProgram
		err error
	}{
		{"", nil, errNoRule},

		// matcher expression
		{"A = 'm'", &testProgram{
			Instrs: combineInstrs(
				mustEncodeInstr(t, ϡopPush, ϡistackID, 3),
				mustEncodeInstr(t, ϡopCall),
				mustEncodeInstr(t, ϡopExit),
				mustEncodeInstr(t, ϡopPush, ϡpstackID),
				mustEncodeInstr(t, ϡopMatch, 0),
				mustEncodeInstr(t, ϡopRestoreIfF),
				mustEncodeInstr(t, ϡopReturn),
			),
			Ms:          []string{`"m"`},
			Ss:          []string{"A"},
			InstrToRule: combineInts(rpt(-1, 3), rpt(0, 4)),
		}, nil},

		// matcher with rule display name
		{`A "Z" = 'm'`, &testProgram{
			Instrs: combineInstrs(
				mustEncodeInstr(t, ϡopPush, ϡistackID, 3),
				mustEncodeInstr(t, ϡopCall),
				mustEncodeInstr(t, ϡopExit),
				mustEncodeInstr(t, ϡopPush, ϡpstackID),
				mustEncodeInstr(t, ϡopMatch, 0),
				mustEncodeInstr(t, ϡopRestoreIfF),
				mustEncodeInstr(t, ϡopReturn),
			),
			Ms:          []string{`"m"`},
			Ss:          []string{"A", "Z"},
			InstrToRule: combineInts(rpt(-1, 3), rpt(1, 4)),
		}, nil},

		// sequence expression
		{`A  = 'm' 'n'`, &testProgram{
			Instrs: combineInstrs(
				mustEncodeInstr(t, ϡopPush, ϡistackID, 11),
				mustEncodeInstr(t, ϡopCall),
				mustEncodeInstr(t, ϡopExit),
				// 'm'
				mustEncodeInstr(t, ϡopPush, ϡpstackID),
				mustEncodeInstr(t, ϡopMatch, 0),
				mustEncodeInstr(t, ϡopRestoreIfF),
				mustEncodeInstr(t, ϡopReturn),
				// 'n'
				mustEncodeInstr(t, ϡopPush, ϡpstackID),
				mustEncodeInstr(t, ϡopMatch, 1),
				mustEncodeInstr(t, ϡopRestoreIfF),
				mustEncodeInstr(t, ϡopReturn),
				// seq (11)
				mustEncodeInstr(t, ϡopPush, ϡpstackID),
				mustEncodeInstr(t, ϡopPush, ϡvstackID, ϡvValFailed),
				mustEncodeInstr(t, ϡopPush, ϡlstackID, 3, 7),
				mustEncodeInstr(t, ϡopTakeLOrJump, 19),
				mustEncodeInstr(t, ϡopCall),
				mustEncodeInstr(t, ϡopCumulOrF),
				mustEncodeInstr(t, ϡopJumpIfF, 19),
				mustEncodeInstr(t, ϡopJump, 14),
				mustEncodeInstr(t, ϡopPop, ϡlstackID),
				mustEncodeInstr(t, ϡopRestoreIfF),
				mustEncodeInstr(t, ϡopReturn),
			),
			Ms:          []string{`"m"`, `"n"`},
			Ss:          []string{"A"},
			InstrToRule: combineInts(rpt(-1, 3), rpt(0, 19)),
		}, nil},

		// choice expression
		{`A  = 'm' / 'n'`, &testProgram{
			Instrs: combineInstrs(
				mustEncodeInstr(t, ϡopPush, ϡistackID, 11),
				mustEncodeInstr(t, ϡopCall),
				mustEncodeInstr(t, ϡopExit),
				// 'm'
				mustEncodeInstr(t, ϡopPush, ϡpstackID),
				mustEncodeInstr(t, ϡopMatch, 0),
				mustEncodeInstr(t, ϡopRestoreIfF),
				mustEncodeInstr(t, ϡopReturn),
				// 'n'
				mustEncodeInstr(t, ϡopPush, ϡpstackID),
				mustEncodeInstr(t, ϡopMatch, 1),
				mustEncodeInstr(t, ϡopRestoreIfF),
				mustEncodeInstr(t, ϡopReturn),
				// choice (11)
				mustEncodeInstr(t, ϡopPush, ϡlstackID, 3, 7),
				mustEncodeInstr(t, ϡopTakeLOrJump, 16),
				mustEncodeInstr(t, ϡopCall),
				mustEncodeInstr(t, ϡopJumpIfT, 16),
				mustEncodeInstr(t, ϡopJump, 12),
				mustEncodeInstr(t, ϡopPop, ϡlstackID),
				mustEncodeInstr(t, ϡopReturn),
			),
			Ms:          []string{`"m"`, `"n"`},
			Ss:          []string{"A"},
			InstrToRule: combineInts(rpt(-1, 3), rpt(0, 15)),
		}, nil},

		// zero or more expression
		{`A  = 'm'*`, &testProgram{
			Instrs: combineInstrs(
				mustEncodeInstr(t, ϡopPush, ϡistackID, 7),
				mustEncodeInstr(t, ϡopCall),
				mustEncodeInstr(t, ϡopExit),
				// 'm'
				mustEncodeInstr(t, ϡopPush, ϡpstackID),
				mustEncodeInstr(t, ϡopMatch, 0),
				mustEncodeInstr(t, ϡopRestoreIfF),
				mustEncodeInstr(t, ϡopReturn),
				// zero or more (7)
				mustEncodeInstr(t, ϡopPush, ϡvstackID, ϡvValEmpty),
				mustEncodeInstr(t, ϡopPush, ϡistackID, 3),
				mustEncodeInstr(t, ϡopCall),
				mustEncodeInstr(t, ϡopPopVJumpIfF, 13),
				mustEncodeInstr(t, ϡopCumulOrF),
				mustEncodeInstr(t, ϡopJump, 8),
				mustEncodeInstr(t, ϡopReturn),
			),
			Ms:          []string{`"m"`},
			Ss:          []string{"A"},
			InstrToRule: combineInts(rpt(-1, 3), rpt(0, 11)),
		}, nil},

		// one or more expression
		{`A  = 'm'+`, &testProgram{
			Instrs: combineInstrs(
				mustEncodeInstr(t, ϡopPush, ϡistackID, 7),
				mustEncodeInstr(t, ϡopCall),
				mustEncodeInstr(t, ϡopExit),
				// 'm'
				mustEncodeInstr(t, ϡopPush, ϡpstackID),
				mustEncodeInstr(t, ϡopMatch, 0),
				mustEncodeInstr(t, ϡopRestoreIfF),
				mustEncodeInstr(t, ϡopReturn),
				// one or more (7)
				mustEncodeInstr(t, ϡopPush, ϡvstackID, ϡvValFailed),
				mustEncodeInstr(t, ϡopPush, ϡistackID, 3),
				mustEncodeInstr(t, ϡopCall),
				mustEncodeInstr(t, ϡopPopVJumpIfF, 13),
				mustEncodeInstr(t, ϡopCumulOrF),
				mustEncodeInstr(t, ϡopJump, 8),
				mustEncodeInstr(t, ϡopReturn),
			),
			Ms:          []string{`"m"`},
			Ss:          []string{"A"},
			InstrToRule: combineInts(rpt(-1, 3), rpt(0, 11)),
		}, nil},

		// zero or one expression
		{`A  = 'm'?`, &testProgram{
			Instrs: combineInstrs(
				mustEncodeInstr(t, ϡopPush, ϡistackID, 7),
				mustEncodeInstr(t, ϡopCall),
				mustEncodeInstr(t, ϡopExit),
				// 'm'
				mustEncodeInstr(t, ϡopPush, ϡpstackID),
				mustEncodeInstr(t, ϡopMatch, 0),
				mustEncodeInstr(t, ϡopRestoreIfF),
				mustEncodeInstr(t, ϡopReturn),
				// zero or one (7)
				mustEncodeInstr(t, ϡopPush, ϡistackID, 3),
				mustEncodeInstr(t, ϡopCall),
				mustEncodeInstr(t, ϡopPopVJumpIfF, 11),
				mustEncodeInstr(t, ϡopReturn),
				mustEncodeInstr(t, ϡopPush, ϡvstackID, ϡvValNil),
				mustEncodeInstr(t, ϡopReturn),
			),
			Ms:          []string{`"m"`},
			Ss:          []string{"A"},
			InstrToRule: combineInts(rpt(-1, 3), rpt(0, 10)),
		}, nil},

		// rule ref expression
		{"A = B\nB = 'm'", &testProgram{
			Instrs: combineInstrs(
				mustEncodeInstr(t, ϡopPush, ϡistackID, 3),
				mustEncodeInstr(t, ϡopCall),
				mustEncodeInstr(t, ϡopExit),
				// rule ref B
				mustEncodeInstr(t, ϡopPush, ϡistackID, 6),
				mustEncodeInstr(t, ϡopCall),
				mustEncodeInstr(t, ϡopReturn),
				// 'm'
				mustEncodeInstr(t, ϡopPush, ϡpstackID),
				mustEncodeInstr(t, ϡopMatch, 0),
				mustEncodeInstr(t, ϡopRestoreIfF),
				mustEncodeInstr(t, ϡopReturn),
			),
			Ms:          []string{`"m"`},
			Ss:          []string{"A", "B"},
			InstrToRule: combineInts(rpt(-1, 3), rpt(0, 3), rpt(1, 4)),
		}, nil},
	}

	for _, tc := range cases {
		gr, err := bootstrap.NewParser().Parse("", strings.NewReader(tc.in))
		if err != nil {
			t.Errorf("%q: parse error: %v", tc.in, err)
			continue
		}

		pg, err := NewGenerator(ioutil.Discard).toProgram(gr)
		if (err != nil) != (tc.err != nil) {
			t.Errorf("%q: want error? %t, got %v", tc.in, tc.err != nil, err)
			continue
		} else if tc.err != err {
			t.Errorf("%q: want error %v, got %v", tc.in, tc.err, err)
			continue
		}

		if tc.err == nil {
			comparePrograms(t, tc.in, tc.out, pg)
		}
	}
}

func combineInts(ints ...[]int) []int {
	var ret []int
	for _, ar := range ints {
		ret = append(ret, ar...)
	}
	return ret
}

func rpt(n, x int) []int {
	ret := make([]int, x)
	for i := 0; i < x; i++ {
		ret[i] = n
	}
	return ret
}

func combineInstrs(instrs ...[]ϡinstr) []ϡinstr {
	var ret []ϡinstr
	for _, ar := range instrs {
		ret = append(ret, ar...)
	}
	return ret
}

func mustEncodeInstr(t *testing.T, op ϡop, args ...int) []ϡinstr {
	instrs, err := ϡencodeInstr(op, args...)
	if err != nil {
		t.Fatal(err)
	}
	return instrs
}

func comparePrograms(t *testing.T, label string, want *testProgram, got *program) {
	// compare Init code
	if want.Init != got.Init {
		t.Errorf("%q: want init %q, got %q", label, want.Init, got.Init)
	}

	// compare instructions
	if len(want.Instrs) != len(got.Instrs) {
		t.Errorf("%q: want %d instructions, got %d", label, len(want.Instrs), len(got.Instrs))
	}
	min := len(want.Instrs)
	if l := len(got.Instrs); l < min {
		min = l
	}
	for i := 0; i < min; i++ {
		if want.Instrs[i] != got.Instrs[i] {
			wop, wn, wa0, _, _ := want.Instrs[i].decode()
			gop, gn, ga0, _, _ := got.Instrs[i].decode()
			t.Errorf("%q: instruction %d: want %s (%d: %d), got %s (%d: %d)",
				label, i, wop, wn, wa0, gop, gn, ga0)
		}
	}

	// compare matchers
	if len(want.Ms) != len(got.Ms) {
		t.Errorf("%q: want %d matchers, got %d", label, len(want.Ms), len(got.Ms))
	}
	min = len(want.Ms)
	if l := len(got.Ms); l < min {
		min = l
	}
	for i := 0; i < min; i++ {
		var raw string
		switch m := got.Ms[i].(type) {
		case *ast.LitMatcher:
			raw = strconv.Quote(m.Val)
			if m.IgnoreCase {
				raw += "i"
			}
		case *ast.CharClassMatcher:
			raw = m.Val
		case *ast.AnyMatcher:
			raw = m.Val
		}
		if want.Ms[i] != raw {
			t.Errorf("%q: matcher %d: want %s, got %s", label, i, want.Ms[i], raw)
		}
	}

	// compare strings
	if len(want.Ss) != len(got.Ss) {
		t.Errorf("%q: want %d strings, got %d", label, len(want.Ss), len(got.Ss))
	}
	min = len(want.Ss)
	if l := len(got.Ss); l < min {
		min = l
	}
	for i := 0; i < min; i++ {
		if want.Ss[i] != got.Ss[i] {
			t.Errorf("%q: string %d: want %q, got %q", label, i, want.Ss[i], got.Ss[i])
		}
	}

	// compare instruction-to-rule mapping
	if len(want.InstrToRule) != len(got.InstrToRule) {
		t.Errorf("%q: want %d instr-to-rule, got %d", label, len(want.InstrToRule), len(got.InstrToRule))
	}
	min = len(want.InstrToRule)
	if l := len(got.InstrToRule); l < min {
		min = l
	}
	for i := 0; i < min; i++ {
		if want.InstrToRule[i] != got.InstrToRule[i] {
			t.Errorf("%q: instr-to-rule %d: want %d, got %d", label, i, want.InstrToRule[i], got.InstrToRule[i])
		}
	}

	// compare A and B thunks
	compareThunkInfos(t, label, "action thunks", want.As, got.As)
	compareThunkInfos(t, label, "bool thunks", want.Bs, got.Bs)
}

func compareThunkInfos(t *testing.T, label, thunkType string, want, got []*thunkInfo) {
	if len(want) != len(got) {
		t.Errorf("%q: want %d %s, got %d", label, len(want), thunkType, len(got))
	}
	min := len(want)
	if l := len(got); l < min {
		min = l
	}
	for i := 0; i < min; i++ {
		// TODO ...
	}
}
