package main

import "testing"

var cases = map[string]string{
	`123455`: ``,
	`

    1

    2

    3

    `: ``,
	`

    1

    2

    x
    `: `2:0 (0): no match found`,
}

var casesFailureTracking = map[string]string{
	`

    1

    2

    x
    `: `7:5 (20): no match found, expected: [ \t\r\n], [0-9] or EOF`,
}

func TestErrorReporting(t *testing.T) {
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
