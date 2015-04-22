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

func TestOptionsReset(t *testing.T) {
	v := &ϡvm{}
	opts := []Option{Memoize(true), Debug(true), Recover(true)}
	flds := []*bool{&v.memoize, &v.debug, &v.recover}
	for i, opt := range opts {
		old := opt(v)
		if !(*flds[i]) {
			t.Errorf("%d: on set, want true, got false", i)
		}
		old(v)
		if *flds[i] {
			t.Errorf("%d: on reset, want false, got true", i)
		}
	}
}
