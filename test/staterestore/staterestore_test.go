package main

import "testing"

var cases = map[string]int{
	"f#\n": 1,
	"f\n":  1,
}

func TestStateRestore(t *testing.T) {
	for tc, exp := range cases {
		got, err := Parse("", []byte(tc))
		if err != nil {
			t.Errorf("%q: %v", tc, err)
			continue
		}
		if got != exp {
			t.Errorf("%q: want %v, got %v", tc, exp, got)
		}
	}
}
