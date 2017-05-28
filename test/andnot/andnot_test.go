package main

import "testing"

// ABs must end in Bs, CDs must end in Ds
var cases = map[string]string{
	"":             "1:1 (0): no match found",
	"a":            "1:1 (0): no match found",
	"b":            "",
	"ab":           "",
	"ba":           "1:1 (0): no match found",
	"aab":          "",
	"bba":          "1:1 (0): no match found",
	"aabbaba":      "1:1 (0): no match found",
	"bbaabaaabbbb": "",
	"abc":          "1:1 (0): no match found",
	"c":            "1:1 (0): no match found",
	"d":            "",
	"cd":           "",
	"dc":           "1:1 (0): no match found",
	"dcddcc":       "1:1 (0): no match found",
	"dcddccdd":     "",
}

var casesFailureTracking = map[string]string{
	"":        `1:1 (0): no match found, expected: [ \t\n\r], [ab], [cd]`,
	"a":       `1:2 (1): no match found, expected: [ab]`,
	"ba":      `1:3 (2): no match found, expected: [ab]`,
	"bba":     `1:4 (3): no match found, expected: [ab]`,
	"aabbaba": `1:8 (7): no match found, expected: [ab]`,
	"abc":     `1:3 (2): no match found, expected: [ \t\n\r], [ab] or EOF`,
	"c":       `1:2 (1): no match found, expected: [cd]`,
	"dc":      `1:3 (2): no match found, expected: [cd]`,
	"dcddcc":  `1:7 (6): no match found, expected: [cd]`,
}

func TestAndNot(t *testing.T) {
	config := []struct {
		failureTracking bool
	}{
		{
			failureTracking: false,
		},
		{
			failureTracking: true,
		},
	}
	for _, conf := range config {
		for tc, exp := range cases {
			_, err := Parse("", []byte(tc), FailureTracking(conf.failureTracking))
			var got string
			if err != nil {
				got = err.Error()
			}
			if conf.failureTracking {
				if failTrackExp, ok := casesFailureTracking[tc]; ok {
					exp = failTrackExp
				}
			}
			if got != exp {
				t.Errorf("%q: want %v, got %v", tc, exp, got)
			}
		}
	}
}
