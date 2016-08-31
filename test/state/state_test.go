package main

import "testing"

var cases = map[string]int{
	"abce":         1,
    "abcd":         3,
}

func TestState(t *testing.T) {
	for tc, exp := range cases {
		got, err := Parse("", []byte(tc), Memoize(true))
		
		if err != nil {
			t.Errorf(err.Error())
		}
		if got != exp {
			t.Errorf("%q: want %v, got %v", tc, exp, got)
		}
	}
}
