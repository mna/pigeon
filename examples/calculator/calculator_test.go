package main

import "testing"

var validCases = map[string]int{
	"0":   0,
	"1":   1,
	"-1":  -1,
	"10":  10,
	"-10": -10,

	"(0)":   0,
	"(1)":   1,
	"(-1)":  -1,
	"(10)":  10,
	"(-10)": -10,

	"1+1":   2,
	"1-1":   0,
	"1*1":   1,
	"1/1":   1,
	"1 + 1": 2,
	"1 - 1": 0,
	"1 * 1": 1,
	"1 / 1": 1,

	"1+0":   1,
	"1-0":   1,
	"1*0":   0,
	"1 + 0": 1,
	"1 - 0": 1,
	"1 * 0": 0,

	"1\n+\t2\r\n +\n3\n": 6,
	"(2) * 3":            6,

	" 1 + 2 - 3 * 4 / 5 ":       1,
	" 1 + (2 - 3) * 4 / 5 ":     1,
	" (1 + 2 - 3) * 4 / 5 ":     0,
	" 1 + 2 - (3 * 4) / 5 ":     1,
	" 18 + 3 - 27 * (-18 / -3)": -141,
}

func TestValidCases(t *testing.T) {
	for tc, exp := range validCases {
		got, err := Parse("", []byte(tc))
		if err != nil {
			t.Errorf("%q: want no error, got %v", tc, err)
			continue
		}
		goti, ok := got.(int)
		if !ok {
			t.Errorf("%q: want type %T, got %T", tc, exp, got)
			continue
		}
		if exp != goti {
			t.Errorf("%q: want %d, got %d", tc, exp, goti)
		}
	}
}

var invalidCases = map[string]string{
	"":        "1:1 (0): no match found",
	"(":       "1:1 (0): no match found",
	")":       "1:1 (0): no match found",
	"()":      "1:1 (0): no match found",
	"+":       "1:1 (0): no match found",
	"-":       "1:1 (0): no match found",
	"*":       "1:1 (0): no match found",
	"/":       "1:1 (0): no match found",
	"+1":      "1:1 (0): no match found",
	"*1":      "1:1 (0): no match found",
	"/1":      "1:1 (0): no match found",
	"1/0":     "1:4 (3): rule Term: runtime error: integer divide by zero",
	"1+":      "1:1 (0): no match found",
	"1-":      "1:1 (0): no match found",
	"1*":      "1:1 (0): no match found",
	"1/":      "1:1 (0): no match found",
	"1 (+ 2)": "1:1 (0): no match found",
	"1 (2)":   "1:1 (0): no match found",
	"\xfe":    "1:1 (0): invalid encoding",
}

func TestInvalidCases(t *testing.T) {
	for tc, exp := range invalidCases {
		got, err := Parse("", []byte(tc))
		if err == nil {
			t.Errorf("%q: want error, got none (%v)", tc, got)
			continue
		}
		el, ok := err.(*errList)
		if !ok {
			t.Errorf("%q: want error type %T, got %T", tc, &errList{}, err)
			continue
		}
		for _, e := range *el {
			if _, ok := e.(*ParserError); !ok {
				t.Errorf("%q: want all individual errors to be %T, got %T (%[3]v)", tc, &ParserError{}, e)
			}
		}
		if exp != err.Error() {
			t.Errorf("%q: want \n%s\n, got \n%s\n", tc, exp, err)
		}
	}
}
