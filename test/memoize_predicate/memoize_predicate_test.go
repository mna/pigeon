package main

import (
    "testing"
    "fmt"
)
var cases = []string{
	"aa",
}

func TestMemoizePredicate(t *testing.T) {
	for _, tc:= range cases {
		exp, err := Parse("", []byte(tc), Memoize(false))
		
		if err != nil {
			t.Errorf(err.Error())
		}
        exp = fmt.Sprintf("%v", exp)

		got, err := Parse("", []byte(tc), Memoize(true))
		if err != nil {
			t.Errorf(err.Error())
		}
        got = fmt.Sprintf("%v", got)
        
		if got != exp {
			t.Errorf("%q: want %v, got %v", tc, exp, got)
		}
	}
}
