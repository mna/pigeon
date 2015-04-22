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

func TestEncodeInstrValid(t *testing.T) {
	cases := []struct {
		op   ϡop
		args []int
		out  []uint64
	}{
		{ϡopExit, nil, []uint64{0}},
		{ϡopCall, nil, []uint64{1 << 58}},
		{ϡopCallA, []int{4}, []uint64{2<<58 | 1<<48 | 4<<32}},
		{ϡopPushL, []int{1, 2}, []uint64{13<<58 | 2<<48 | 1<<32 | 2<<16}},
		{ϡopPushL, []int{1, 2, 3}, []uint64{13<<58 | 3<<48 | 1<<32 | 2<<16 | 3}},
		{ϡopPushL, []int{1, 2, 3, 4}, []uint64{
			13<<58 | 4<<48 | 1<<32 | 2<<16 | 3,
			4 << 48,
		}},
		{ϡopPushL, []int{1, 2, 3, 4, 5}, []uint64{
			13<<58 | 5<<48 | 1<<32 | 2<<16 | 3,
			4<<48 | 5<<32,
		}},
		{ϡopPushL, []int{1, 2, 3, 4, 5, 6}, []uint64{
			13<<58 | 6<<48 | 1<<32 | 2<<16 | 3,
			4<<48 | 5<<32 | 6<<16,
		}},
		{ϡopPushL, []int{1, 2, 3, 4, 5, 6, 7}, []uint64{
			13<<58 | 7<<48 | 1<<32 | 2<<16 | 3,
			4<<48 | 5<<32 | 6<<16 | 7,
		}},
		{ϡopPushL, []int{1, 2, 3, 4, 5, 6, 7, 8}, []uint64{
			13<<58 | 8<<48 | 1<<32 | 2<<16 | 3,
			4<<48 | 5<<32 | 6<<16 | 7,
			8 << 48,
		}},
		{ϡopPushL, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, []uint64{
			13<<58 | 9<<48 | 1<<32 | 2<<16 | 3,
			4<<48 | 5<<32 | 6<<16 | 7,
			8<<48 | 9<<32,
		}},
		{ϡopPushL, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, []uint64{
			13<<58 | 10<<48 | 1<<32 | 2<<16 | 3,
			4<<48 | 5<<32 | 6<<16 | 7,
			8<<48 | 9<<32 | 10<<16,
		}},
		{ϡopMax - 1, nil, []uint64{uint64(ϡopMax-1) << 58}},
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

func TestEncodeInstrInvalid(t *testing.T) {
	// TODO..
}
