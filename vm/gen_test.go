package vm

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/PuerkitoBio/pigeon/ast"
	"github.com/PuerkitoBio/pigeon/bootstrap"
)

type testProgram struct {
	Init   string
	Instrs []ϡinstr
	Ms     []string
	Ss     []string
	As     []*thunkInfo
	Bs     []*thunkInfo
}

func TestNoRule(t *testing.T) {
	gr := ast.NewGrammar(ast.Pos{})
	_, err := NewGenerator(ioutil.Discard).toProgram(gr)
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if err != errNoRule {
		t.Fatalf("want error %v, got %v", errNoRule, err)
	}
}

func TestGenerateProgram(t *testing.T) {
	cases := []struct {
		in  string
		out *testProgram
		err error
	}{
		{"", nil, errNoRule},
		{"A = B", nil, errors.New(`undefined rule "B"`)},

		// matcher expression
		{"A = 'm'", &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 4), -1),
				setRuleNmIx(encodeMatcher(t, 0), 0), // 'm'
			),
			Ms: []string{`"m"`},
			Ss: []string{"A"},
		}, nil},

		// matcher expression with an Init
		{"{x}\nA = 'm'", &testProgram{
			Init: "x",
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 4), -1),
				setRuleNmIx(encodeMatcher(t, 0), 0), // 'm'
			),
			Ms: []string{`"m"`},
			Ss: []string{"A"},
		}, nil},

		// matcher with rule display name
		{`A "Z" = 'm'`, &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 4), -1),
				setRuleNmIx(encodeMatcher(t, 0), 1), // 'm'
			),
			Ms: []string{`"m"`},
			Ss: []string{"A", "Z"},
		}, nil},

		// sequence expression
		{`A  = 'm' 'n'`, &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 12), -1),
				encodeMatcher(t, 0), // 4: 'm'
				encodeMatcher(t, 1), // 8: 'n'
				encodeSequence(t, 12, 4, 8),
			),
			Ms: []string{`"m"`, `"n"`},
			Ss: []string{"A"},
		}, nil},

		// choice expression
		{`A  = 'm' / 'n'`, &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 12), -1),
				encodeMatcher(t, 0), // 4: 'm'
				encodeMatcher(t, 1), // 8: 'n'
				encodeChoice(t, 12, 4, 8),
			),
			Ms: []string{`"m"`, `"n"`},
			Ss: []string{"A"},
		}, nil},

		// zero or more expression
		{`A  = 'm'*`, &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 8), -1),
				encodeMatcher(t, 0), // 4: 'm'
				encodeRepetition(t, 8, ϡvValEmpty, 4),
			),
			Ms: []string{`"m"`},
			Ss: []string{"A"},
		}, nil},

		// one or more expression
		{`A  = 'm'+`, &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 8), -1),
				encodeMatcher(t, 0), // 4: 'm'
				encodeRepetition(t, 8, ϡvValFailed, 4),
			),
			Ms: []string{`"m"`},
			Ss: []string{"A"},
		}, nil},

		// zero or one expression
		{`A  = 'm'?`, &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 8), -1),
				encodeMatcher(t, 0), // 4: 'm'
				encodeOption(t, 8, 4),
			),
			Ms: []string{`"m"`},
			Ss: []string{"A"},
		}, nil},

		// rule ref expression
		{"A = B\nB = 'm'", &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 4), -1),
				setRuleNmIx(encodeRuleRef(t, 9), 0),
				setRuleNmIx(encodeMatcher(t, 0), 1), // 9: 'm'
			),
			Ms: []string{`"m"`},
			Ss: []string{"A", "B"},
		}, nil},

		// and expression
		{"A = &'m'", &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 8), -1),
				encodeMatcher(t, 0), // 4: 'm'
				encodePredicate(t, true, 4),
			),
			Ms: []string{`"m"`},
			Ss: []string{"A"},
		}, nil},

		// not expression
		{"A = !'m'", &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 8), -1),
				encodeMatcher(t, 0), // 4: 'm'
				encodePredicate(t, false, 4),
			),
			Ms: []string{`"m"`},
			Ss: []string{"A"},
		}, nil},

		// and code expression
		{"A = &{x}", &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 4), -1),
				encodeCodePredicate(t, true, 0),
			),
			Ss: []string{"A"},
			Bs: []*thunkInfo{&thunkInfo{
				RuleNm: "A",
				Code:   "x",
				ExprIx: 1,
			}},
		}, nil},

		// not code expression
		{"A = !{x}", &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 4), -1),
				encodeCodePredicate(t, false, 0),
			),
			Ss: []string{"A"},
			Bs: []*thunkInfo{&thunkInfo{
				RuleNm: "A",
				Code:   "x",
				ExprIx: 1,
			}},
		}, nil},

		// labeled expression
		{"A = label:'m'", &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 8), -1),
				encodeMatcher(t, 0), // 4: 'm'
				encodeLabel(t, 1, 4),
			),
			Ms: []string{`"m"`},
			Ss: []string{"A", "label"},
		}, nil},

		// action expression
		{"A = 'm' {x}", &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 8), -1),
				encodeMatcher(t, 0), // 4: 'm'
				encodeAction(t, 8, 0, 4),
			),
			Ms: []string{`"m"`},
			Ss: []string{"A"},
			As: []*thunkInfo{&thunkInfo{
				RuleNm: "A",
				ExprIx: 1,
				Code:   "x",
			}},
		}, nil},

		// label+action expression
		{"A = label:'m' {x}", &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 14), -1),
				encodeMatcher(t, 0),  // 4: 'm'
				encodeLabel(t, 1, 4), // 8
				encodeAction(t, 14, 0, 8),
			),
			Ms: []string{`"m"`},
			Ss: []string{"A", "label"},
			As: []*thunkInfo{&thunkInfo{
				Parms:  []string{"label"},
				RuleNm: "A",
				ExprIx: 1,
				Code:   "x",
			}},
		}, nil},

		// multi-label+action expression
		{"A = l1:'m' l2:'n' {x}", &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 35), -1),
				encodeMatcher(t, 0),          // 4: 'm'
				encodeLabel(t, 1, 4),         // 8: l1
				encodeMatcher(t, 1),          // 14: 'n'
				encodeLabel(t, 2, 14),        // 18: l2
				encodeSequence(t, 24, 8, 18), // 24
				encodeAction(t, 35, 0, 24),   // 35
			),
			Ms: []string{`"m"`, `"n"`},
			Ss: []string{"A", "l1", "l2"},
			As: []*thunkInfo{&thunkInfo{
				Parms:  []string{"l1", "l2"},
				RuleNm: "A",
				ExprIx: 1,
				Code:   "x",
			}},
		}, nil},

		// choice resets the params
		{"A = l1:'m' / l2:'n' {x}", &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 32), -1),
				encodeMatcher(t, 0),        // 4: 'm'
				encodeLabel(t, 1, 4),       // 8: l1
				encodeMatcher(t, 1),        // 14: 'n'
				encodeLabel(t, 2, 14),      // 18: l2
				encodeAction(t, 24, 0, 18), // 24
				encodeChoice(t, 32, 8, 24), // 32
			),
			Ms: []string{`"m"`, `"n"`},
			Ss: []string{"A", "l1", "l2"},
			As: []*thunkInfo{&thunkInfo{
				Parms:  []string{"l2"},
				RuleNm: "A",
				ExprIx: 4,
				Code:   "x",
			}},
		}, nil},

		// scope of params
		{"A = l1:'m' l2:(l3:'n' {y}) {x}", &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 49), -1),
				encodeMatcher(t, 0),          // 4: 'm'
				encodeLabel(t, 1, 4),         // 8: l1
				encodeMatcher(t, 1),          // 14: 'n'
				encodeLabel(t, 3, 14),        // 18: l3
				encodeAction(t, 24, 0, 18),   // 24: y
				encodeLabel(t, 2, 24),        // 32: l2
				encodeSequence(t, 38, 8, 32), // 38
				encodeAction(t, 49, 1, 38),   // 49
			),
			Ms: []string{`"m"`, `"n"`},
			Ss: []string{"A", "l1", "l2", "l3"},
			As: []*thunkInfo{&thunkInfo{
				Parms:  []string{"l3"},
				RuleNm: "A",
				ExprIx: 6,
				Code:   "y",
			}, &thunkInfo{
				Parms:  []string{"l1", "l2"},
				RuleNm: "A",
				ExprIx: 1,
				Code:   "x",
			}},
		}, nil},

		// code predicates have access to params too
		{"A = l1:'m' / l2:'n' &{x} l3:'o' {y}", &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 56), -1),
				encodeMatcher(t, 0),               // 4: 'm'
				encodeLabel(t, 1, 4),              // 8: l1
				encodeMatcher(t, 1),               // 14: 'n'
				encodeLabel(t, 2, 14),             // 18: l2
				encodeCodePredicate(t, true, 0),   // 24
				encodeMatcher(t, 2),               // 27: 'o'
				encodeLabel(t, 3, 27),             // 31: l3
				encodeSequence(t, 37, 18, 24, 31), // 37
				encodeAction(t, 48, 0, 37),        // 48
				encodeChoice(t, 56, 8, 48),        // 56
			),
			Ms: []string{`"m"`, `"n"`, `"o"`},
			Ss: []string{"A", "l1", "l2", "l3"},
			As: []*thunkInfo{&thunkInfo{
				Parms:  []string{"l2", "l3"},
				RuleNm: "A",
				ExprIx: 4,
				Code:   "y",
			}},
			Bs: []*thunkInfo{&thunkInfo{
				Parms:  []string{"l2"},
				RuleNm: "A",
				ExprIx: 8,
				Code:   "x",
			}},
		}, nil},

		// normalization of matchers
		{"A = `m` 'm' `m`i", &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 16), -1),
				encodeMatcher(t, 0), // 4: 'm'
				encodeMatcher(t, 0), // 8: 'm'
				encodeMatcher(t, 1), // 12: 'm'i
				encodeSequence(t, 16, 4, 8, 12),
			),
			Ms: []string{`"m"`, `"m"i`},
			Ss: []string{"A"},
		}, nil},

		// test char class and any matchers
		{"A = [a-z] .", &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 12), -1),
				encodeMatcher(t, 0), // 4: [a-z]
				encodeMatcher(t, 1), // 8: .
				encodeSequence(t, 12, 4, 8),
			),
			Ms: []string{`[a-z]`, `.`},
			Ss: []string{"A"},
		}, nil},

		// action with a sequence of 4 (the grammar's Grammar rule bug)
		{"A = B C D E {x}\nB = 'b'\nC = 'c'\nD = 'd'\nE = 'e'", &testProgram{
			Instrs: combineInstrs(
				setRuleNmIx(encodeBootstrap(t, 35), -1),
				encodeRuleRef(t, 43),                // 4
				encodeRuleRef(t, 47),                // 9
				encodeRuleRef(t, 51),                // 14
				encodeRuleRef(t, 55),                // 19
				encodeSequence(t, 24, 4, 9, 14, 19), // 24
				encodeAction(t, 35, 0, 24),          // 35
				setRuleNmIx(encodeMatcher(t, 0), 1), // 43
				setRuleNmIx(encodeMatcher(t, 1), 2), // 47
				setRuleNmIx(encodeMatcher(t, 2), 3), // 51
				setRuleNmIx(encodeMatcher(t, 3), 4), // 55
			),
			Ms: []string{`"b"`, `"c"`, `"d"`, `"e"`},
			Ss: []string{"A", "B", "C", "D", "E"},
			As: []*thunkInfo{&thunkInfo{
				RuleNm: "A",
				ExprIx: 1,
				Code:   "x",
			}},
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
			if tc.err != nil && fmt.Sprintf("%T", tc.err) == fmt.Sprintf("%T", err) {
				if tc.err.Error() != err.Error() {
					t.Errorf("%q: want error %v, got %v", tc.in, tc.err, err)
					continue
				}
			} else {
				t.Errorf("%q: want error %v, got %v", tc.in, tc.err, err)
				continue
			}
		}

		if tc.err == nil {
			comparePrograms(t, tc.in, tc.out, pg)
		}
	}
}

func setRuleNmIx(instrs []ϡinstr, ix int) []ϡinstr {
	for i, ins := range instrs {
		ins.ruleNmIx = ix
		instrs[i] = ins
	}
	return instrs
}

func combineInstrs(instrs ...[]ϡinstr) []ϡinstr {
	var ret []ϡinstr
	for _, ar := range instrs {
		ret = append(ret, ar...)
	}
	return ret
}

func mustEncodeInstr(t *testing.T, op ϡop, args ...uint16) ϡinstr {
	return ϡinstr{op: op, args: args}
}

func encodeBootstrap(t *testing.T, start uint16) []ϡinstr {
	return []ϡinstr{
		mustEncodeInstr(t, ϡopPush, ϡistackID, start),
		mustEncodeInstr(t, ϡopPush, ϡastackID),
		mustEncodeInstr(t, ϡopCall),
		mustEncodeInstr(t, ϡopExit),
	}
}

func encodeMatcher(t *testing.T, mIx uint16) []ϡinstr {
	return []ϡinstr{
		mustEncodeInstr(t, ϡopPush, ϡpstackID),
		mustEncodeInstr(t, ϡopMatch, mIx),
		mustEncodeInstr(t, ϡopRestoreIfF),
		mustEncodeInstr(t, ϡopReturn),
	}
}

func encodeSequence(t *testing.T, start uint16, ls ...uint16) []ϡinstr {
	return []ϡinstr{
		mustEncodeInstr(t, ϡopPush, ϡpstackID),
		mustEncodeInstr(t, ϡopPush, ϡvstackID, ϡvValFailed),
		mustEncodeInstr(t, ϡopPush, append([]uint16{ϡlstackID}, ls...)...),
		mustEncodeInstr(t, ϡopTakeLOrJump, start+8),
		mustEncodeInstr(t, ϡopCall),
		mustEncodeInstr(t, ϡopCumulOrF),
		mustEncodeInstr(t, ϡopJumpIfF, start+8),
		mustEncodeInstr(t, ϡopJump, start+3),
		mustEncodeInstr(t, ϡopPop, ϡlstackID),
		mustEncodeInstr(t, ϡopRestoreIfF),
		mustEncodeInstr(t, ϡopReturn),
	}
}

func encodeChoice(t *testing.T, start uint16, ls ...uint16) []ϡinstr {
	return []ϡinstr{
		mustEncodeInstr(t, ϡopPush, append([]uint16{ϡlstackID}, ls...)...),
		mustEncodeInstr(t, ϡopTakeLOrJump, start+8),
		mustEncodeInstr(t, ϡopPush, ϡastackID),
		mustEncodeInstr(t, ϡopCall),
		mustEncodeInstr(t, ϡopPop, ϡastackID),
		mustEncodeInstr(t, ϡopJumpIfT, start+9),
		mustEncodeInstr(t, ϡopPop, ϡvstackID),
		mustEncodeInstr(t, ϡopJump, start+1),
		mustEncodeInstr(t, ϡopPush, ϡvstackID, ϡvValFailed),
		mustEncodeInstr(t, ϡopPop, ϡlstackID),
		mustEncodeInstr(t, ϡopReturn),
	}
}

func encodeRepetition(t *testing.T, start uint16, vVal uint16, ix uint16) []ϡinstr {
	return []ϡinstr{
		mustEncodeInstr(t, ϡopPush, ϡvstackID, vVal),
		mustEncodeInstr(t, ϡopPush, ϡistackID, ix),
		mustEncodeInstr(t, ϡopPush, ϡastackID),
		mustEncodeInstr(t, ϡopCall),
		mustEncodeInstr(t, ϡopPop, ϡastackID),
		mustEncodeInstr(t, ϡopPopVJumpIfF, start+8),
		mustEncodeInstr(t, ϡopCumulOrF),
		mustEncodeInstr(t, ϡopJump, start+1),
		mustEncodeInstr(t, ϡopReturn),
	}
}

func encodeOption(t *testing.T, start uint16, ix uint16) []ϡinstr {
	return []ϡinstr{
		mustEncodeInstr(t, ϡopPush, ϡistackID, ix),
		mustEncodeInstr(t, ϡopPush, ϡastackID),
		mustEncodeInstr(t, ϡopCall),
		mustEncodeInstr(t, ϡopPop, ϡastackID),
		mustEncodeInstr(t, ϡopPopVJumpIfF, start+6),
		mustEncodeInstr(t, ϡopReturn),
		mustEncodeInstr(t, ϡopPush, ϡvstackID, ϡvValNil),
		mustEncodeInstr(t, ϡopReturn),
	}
}

func encodeRuleRef(t *testing.T, ix uint16) []ϡinstr {
	return []ϡinstr{
		mustEncodeInstr(t, ϡopPush, ϡistackID, ix),
		mustEncodeInstr(t, ϡopPush, ϡastackID),
		mustEncodeInstr(t, ϡopCall),
		mustEncodeInstr(t, ϡopPop, ϡastackID),
		mustEncodeInstr(t, ϡopReturn),
	}
}

func encodePredicate(t *testing.T, and bool, ix uint16) []ϡinstr {
	op := ϡopNilIfF
	if and {
		op = ϡopNilIfT
	}
	return []ϡinstr{
		mustEncodeInstr(t, ϡopPush, ϡpstackID),
		mustEncodeInstr(t, ϡopPush, ϡistackID, ix),
		mustEncodeInstr(t, ϡopPush, ϡastackID),
		mustEncodeInstr(t, ϡopCall),
		mustEncodeInstr(t, ϡopPop, ϡastackID),
		mustEncodeInstr(t, op),
		mustEncodeInstr(t, ϡopRestore),
		mustEncodeInstr(t, ϡopReturn),
	}
}

func encodeCodePredicate(t *testing.T, and bool, bIx uint16) []ϡinstr {
	op := ϡopNilIfF
	if and {
		op = ϡopNilIfT
	}
	return []ϡinstr{
		mustEncodeInstr(t, ϡopCallB, bIx),
		mustEncodeInstr(t, op),
		mustEncodeInstr(t, ϡopReturn),
	}
}

func encodeLabel(t *testing.T, lblIx, ix uint16) []ϡinstr {
	return []ϡinstr{
		mustEncodeInstr(t, ϡopPush, ϡistackID, ix),
		mustEncodeInstr(t, ϡopPush, ϡastackID),
		mustEncodeInstr(t, ϡopCall),
		mustEncodeInstr(t, ϡopPop, ϡastackID),
		mustEncodeInstr(t, ϡopStoreIfT, lblIx),
		mustEncodeInstr(t, ϡopReturn),
	}
}

func encodeAction(t *testing.T, start, actIx, ix uint16) []ϡinstr {
	return []ϡinstr{
		mustEncodeInstr(t, ϡopPush, ϡpstackID),
		mustEncodeInstr(t, ϡopPush, ϡistackID, ix),
		mustEncodeInstr(t, ϡopCall),
		mustEncodeInstr(t, ϡopJumpIfF, start+6),
		mustEncodeInstr(t, ϡopCallA, actIx),
		mustEncodeInstr(t, ϡopReturn),
		mustEncodeInstr(t, ϡopPop, ϡpstackID),
		mustEncodeInstr(t, ϡopReturn),
	}
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
		wi, gi := want.Instrs[i], got.Instrs[i]
		if wi.op != gi.op || !reflect.DeepEqual(wi.args, gi.args) || wi.ruleNmIx != gi.ruleNmIx {
			t.Errorf("%q: instruction %d: want %s %v (%d), got %s %v (%d)",
				label, i, wi.op, wi.args, wi.ruleNmIx, gi.op, gi.args, gi.ruleNmIx)
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
		compareThunkInfo(t, label, thunkType, i, want[i], got[i])
	}
}

func compareThunkInfo(t *testing.T, label, thunkType string, id int, want, got *thunkInfo) {
	// compare parameters
	if len(want.Parms) != len(got.Parms) {
		t.Errorf("%q: %s %d: want %d params, got %d", label, thunkType, id, len(want.Parms), len(got.Parms))
	}
	min := len(want.Parms)
	if l := len(got.Parms); l < min {
		min = l
	}
	for i := 0; i < min; i++ {
		if want.Parms[i] != got.Parms[i] {
			t.Errorf("%q: %s %d: param %d: want %q, got %q", label, thunkType, id, i, want.Parms[i], got.Parms[i])
		}
	}

	if want.RuleNm != got.RuleNm {
		t.Errorf("%q: %s %d: want rule name %q, got %q", label, thunkType, id, want.RuleNm, got.RuleNm)
	}
	if want.ExprIx != got.ExprIx {
		t.Errorf("%q: %s %d: want expression index %d, got %d", label, thunkType, id, want.ExprIx, got.ExprIx)
	}
	if want.Code != got.Code {
		t.Errorf("%q: %s %d: want code %q, got %q", label, thunkType, id, want.Code, got.Code)
	}
}
