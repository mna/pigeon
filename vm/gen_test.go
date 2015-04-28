package vm

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/PuerkitoBio/pigeon/bootstrap"
)

func TestGenerateProgram(t *testing.T) {
	t.Skip()
	cases := []struct {
		in  string
		out *program
		err error
	}{
		{"", nil, errNoRule},
		{"A = 'a'", &program{
			Instrs: combineInstrs(
				mustEncodeInstr(t, ϡopPush, 3),
			),
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

func comparePrograms(t *testing.T, label string, want, got *program) {
	if want.Init != got.Init {
		t.Errorf("%q: want init %q, got %q", label, want.Init, got.Init)
	}

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
			t.Errorf("%q: want %s (%d: %d), got %s (%d: %d)", label, wop, wn, wa0,
				gop, gn, ga0)
		}
	}
}
