package vm

import "testing"

func TestOptions(t *testing.T) {
	v := &ϡvm{}

	wantm, wantd, wantr := true, true, false
	v.setOptions([]Option{
		Memoize(true),
		Debug(true),
		Debug(false),
		Recover(false),
		Debug(true),
	})

	if v.memoize != wantm {
		t.Errorf("Memoize: want %t, got %t", wantm, v.memoize)
	}
	if v.debug != wantd {
		t.Errorf("Debug: want %t, got %t", wantd, v.debug)
	}
	if v.recover != wantr {
		t.Errorf("Recover: want %t, got %t", wantr, v.recover)
	}
}
