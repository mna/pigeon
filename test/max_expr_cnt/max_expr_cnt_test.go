package main

import "testing"

func TestMaxExprCnt(t *testing.T) {
	_, err := Parse("", []byte("infinite parse"), MaxExpressions(5))
	if err == nil {
		t.Errorf("expected non nil error message for testing max expr cnt option.")
		t.Fail()
	}

	errs, ok := err.(errList)
	if !ok {
		t.Errorf("expected err %v to be of type errList but got type %T", err, err)
		t.Fail()
	}

	var found bool
	for _, err := range errs {
		pe, ok := err.(*parserError)
		if !ok {
			t.Errorf("expected err %v to be of type parserError but got type %T", err, err)
			t.Fail()
		}

		if pe.Inner == errMaxExprCnt {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected to find errMaxExprCnt %v in error list %v", errMaxExprCnt, errs)
		t.Fail()
	}
}
