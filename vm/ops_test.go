package vm

import "testing"

func TestConsts(t *testing.T) {
	if ϡlPerI != 4 {
		t.Errorf("want 4 lPerI, got %d", ϡlPerI)
	}
	if ϡlMask != 65535 {
		t.Errorf("want lMask to be 65535, got %d", ϡlMask)
	}
	if ϡnMask != 1023 {
		t.Errorf("want nMask to be 1023, got %d", ϡnMask)
	}
	if ϡoMask != 63 {
		t.Errorf("want oMask to be 63, got %d", ϡoMask)
	}
	if ϡiBits-ϡoBits != 58 {
		t.Errorf("want iBits - oBits to be 58, got %d", ϡiBits-ϡoBits)
	}
}

func TestDecodeInstr(t *testing.T) {
	cases := []struct {
		ins  uint64
		op   ϡop
		args []int
	}{
		{0, ϡopExit, nil},
		{1 << 58, ϡopCall, nil},
		{2<<58 | 1<<48 | 10<<32, ϡopCallA, []int{10}},
		{12<<58 | 2<<48 | 10<<32 | 12345<<16, ϡopPush, []int{10, 12345}},
		{12<<58 | 3<<48 | 1<<32 | 2<<16 | 3, ϡopPush, []int{1, 2, 3}},
	}

	for i, tc := range cases {
		op, n, ix0, ix1, ix2 := ϡinstr(tc.ins).decode()
		if op != tc.op {
			t.Errorf("%d: want op %s, got %s", i, tc.op, op)
		}
		if n != len(tc.args) {
			t.Errorf("%d: want %d arguments, got %d", i, len(tc.args), n)
			continue
		}
		for j, arg := range []int{ix0, ix1, ix2}[:n] {
			if tc.args[j] != arg {
				t.Errorf("%d: arg %d: want %d, got %d", i, j, tc.args[j], arg)
			}
		}
	}
}

func TestDecodeLs(t *testing.T) {
	cases := []struct {
		ins  uint64
		args []int
	}{
		{0, nil},
		{1 << 48, []int{1}},
		{1<<48 | 2<<32, []int{1, 2}},
		{1<<48 | 2<<32 | 3<<16, []int{1, 2, 3}},
		{1<<48 | 2<<32 | 3<<16 | 4, []int{1, 2, 3, 4}},
	}

	for i, tc := range cases {
		ix0, ix1, ix2, ix3 := ϡinstr(tc.ins).decodeLs()
		for j, arg := range []int{ix0, ix1, ix2, ix3} {
			exp := 0
			if j < len(tc.args) {
				exp = tc.args[j]
			}
			if exp != arg {
				t.Errorf("%d: arg %d: want %d, got %d", i, j, exp, arg)
			}
		}
	}
}

func TestEncodeInstrValid(t *testing.T) {
	cases := []struct {
		op   ϡop
		args []int
		out  []uint64
	}{
		{ϡopExit, nil, []uint64{0}},
		{ϡopCall, nil, []uint64{1 << 58}},
		{ϡopCallA, []int{4}, []uint64{2<<58 | 1<<48 | 4<<32}},
		{ϡopPush, []int{1, 2}, []uint64{12<<58 | 2<<48 | 1<<32 | 2<<16}},
		{ϡopPush, []int{1, 2, 3}, []uint64{12<<58 | 3<<48 | 1<<32 | 2<<16 | 3}},
		{ϡopPush, []int{1, 2, 3, 4}, []uint64{
			12<<58 | 4<<48 | 1<<32 | 2<<16 | 3,
			4 << 48,
		}},
		{ϡopPush, []int{1, 2, 3, 4, 5}, []uint64{
			12<<58 | 5<<48 | 1<<32 | 2<<16 | 3,
			4<<48 | 5<<32,
		}},
		{ϡopPush, []int{1, 2, 3, 4, 5, 6}, []uint64{
			12<<58 | 6<<48 | 1<<32 | 2<<16 | 3,
			4<<48 | 5<<32 | 6<<16,
		}},
		{ϡopPush, []int{1, 2, 3, 4, 5, 6, 7}, []uint64{
			12<<58 | 7<<48 | 1<<32 | 2<<16 | 3,
			4<<48 | 5<<32 | 6<<16 | 7,
		}},
		{ϡopPush, []int{1, 2, 3, 4, 5, 6, 7, 8}, []uint64{
			12<<58 | 8<<48 | 1<<32 | 2<<16 | 3,
			4<<48 | 5<<32 | 6<<16 | 7,
			8 << 48,
		}},
		{ϡopPush, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, []uint64{
			12<<58 | 9<<48 | 1<<32 | 2<<16 | 3,
			4<<48 | 5<<32 | 6<<16 | 7,
			8<<48 | 9<<32,
		}},
		{ϡopPush, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, []uint64{
			12<<58 | 10<<48 | 1<<32 | 2<<16 | 3,
			4<<48 | 5<<32 | 6<<16 | 7,
			8<<48 | 9<<32 | 10<<16,
		}},
		{ϡopmax - 1, nil, []uint64{uint64(ϡopmax-1) << 58}},
	}
	for i, tc := range cases {
		got, err := ϡencodeInstr(tc.op, tc.args)

		if err != nil {
			t.Errorf("%d: got error %v", i, err)
			continue
		}
		if len(got) != len(tc.out) {
			t.Errorf("%d: want %d instructions, got %d", i, len(tc.out), len(got))
			continue
		}

		for j, want := range tc.out {
			if want != uint64(got[j]) {
				t.Errorf(`%d: instruction %d: want
1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef
%64b, got
%64b`, i, j, want, got[j])
			}
		}
	}
}

func TestEncodeInstrLimits(t *testing.T) {
	tooManyArgs := make([]int, ϡnMask+1)
	maxArgs := make([]int, ϡnMask)

	cases := []struct {
		op   ϡop
		args []int
		err  string
	}{
		{ϡopmax - 1, nil, ""},
		{ϡopmax, nil, "invalid op value"},
		{ϡopReturn, []int{1 << ϡlBits}, "argument value too big"},
		{ϡopReturn, []int{1<<ϡlBits - 1}, ""},
		{ϡopReturn, tooManyArgs, "too many arguments"},
		{ϡopReturn, maxArgs, ""},
	}

	for i, tc := range cases {
		_, err := ϡencodeInstr(tc.op, tc.args)

		if (err == nil) != (tc.err == "") {
			t.Errorf("%d: want error? %t, got %v", i, tc.err == "", err)
			continue
		}

		if err != nil && err.Error() != tc.err {
			t.Errorf("%d: want %q, got %q", i, tc.err, err)
		}
	}
}
