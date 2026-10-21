package stateclone

import (
	"fmt"
	"slices"
	"testing"
)

var cases = map[string]int{
	"abce":           11,
	"abcd":           13,
	"abcd abce abcd": 17,
	"a":              10,
	"abcf":           15,
}

func TestState(t *testing.T) {
	for tc, exp := range cases {
		vals := make(values, 1)
		vals[0] = 10
		got, err := Parse("", []byte(tc), Memoize(false), InitState("vals", vals))
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		vals = got.(values)
		res := 0
		for _, v := range vals {
			res += v
		}
		if res != exp {
			t.Errorf("%q: want %v, got %v", tc, exp, got)
		}
	}
}

// Reusing parser state must not mix values between parses or change an
// earlier result, including after backtracking.
func TestStateIsolation(t *testing.T) {
	for seed := range 16 {
		t.Run(fmt.Sprint(seed), func(t *testing.T) {
			t.Parallel()
			input := make(values, 1, 8)
			input[0] = seed
			parse := func(initial values) values {
				t.Helper()
				got, err := Parse("", []byte("abcd abce abcd"), Memoize(false), InitState("vals", initial))
				if err != nil {
					t.Fatal(err)
				}
				result := got.(values)
				want := values{initial[0], 3, 1, 3}
				if !slices.Equal(result, want) {
					t.Fatalf("got %v, want %v", result, want)
				}
				return result
			}
			first := parse(input)
			for range 16 {
				parse(values{seed + 100})
			}
			if !slices.Equal(first, values{seed, 3, 1, 3}) {
				t.Fatalf("earlier result changed: %v", first)
			}
		})
	}
}
