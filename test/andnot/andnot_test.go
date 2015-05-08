package main

import "testing"

// ABs must end in Bs, CDs must end in Ds
var cases = map[string]string{
	"":             `1:0 (0): rule _: expected [ \t\n\r], got ""`,
	"a":            `1:1 (1): rule AB: expected [ab], got ""`,
	"b":            "",
	"ab":           "",
	"ba":           `1:2 (2): rule AB: expected [ab], got ""`,
	"aab":          "",
	"bba":          `1:3 (3): rule AB: expected [ab], got ""`,
	"aabbaba":      `1:7 (7): rule AB: expected [ab], got ""`,
	"bbaabaaabbbb": "",
	"abc":          `1:3 (2): rule AB: expected [ab], got "c"`,
	"c":            `1:1 (1): rule CD: expected [cd], got ""`,
	"d":            "",
	"cd":           "",
	"dc":           `1:2 (2): rule CD: expected [cd], got ""`,
	"dcddcc":       `1:6 (6): rule CD: expected [cd], got ""`,
	"dcddccdd":     "",
}

func TestAndNot(t *testing.T) {
	for tc, exp := range cases {
		_, err := Parse("", []byte(tc))
		var got string
		if err != nil {
			got = err.Error()
		}
		if got != exp {
			t.Errorf("%q: want %v, got %v", tc, exp, got)
		}
	}
}
