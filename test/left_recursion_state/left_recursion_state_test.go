package leftrecursionstate_test

import (
	"testing"

	optimizedleftrecursionstate "github.com/mna/pigeon/test/left_recursion_state/optimized"
	leftrecursionstate "github.com/mna/pigeon/test/left_recursion_state/standart"
)

func TestLeftRecursionWithState(t *testing.T) {
	t.Parallel()

	initCount := 100000

	type want struct {
		count int
	}

	tests := []struct {
		name string
		expr string
		want want
	}{
		{
			name: "atom",
			expr: "1",
			want: want{count: 3 + 15 + 63 + 127 + initCount},
		},
		{
			name: "factor",
			expr: "-1",
			want: want{count: 3 + 15 + 31 + 63 + 127 + initCount},
		},
		{
			name: "expr",
			expr: "1+1",
			want: want{count: 1 +
				(3 + 15 + 63 + 127) +
				(15 + 63 + 127) + initCount},
		},
		{
			name: "expr",
			expr: "1*1*1",
			want: want{count: 3 +
				7 +
				7 +
				(15 + 63 + 127) +
				(63 + 127) +
				(63 + 127) +
				+initCount},
		},
		{
			name: "invalid",
			expr: "**",
			want: want{count: initCount},
		},
	}

	for _, testCase := range tests {
		testCase := testCase

		setOptions := map[string][]leftrecursionstate.Option{
			"memoize": {
				leftrecursionstate.Memoize(true),
				leftrecursionstate.InitState("count", initCount),
			},
			"-": {
				leftrecursionstate.InitState("count", initCount),
			},
		}
		for nameOptions, options := range setOptions {
			options := options

			t.Run(testCase.name+" default. Options: "+nameOptions, func(t *testing.T) {
				t.Parallel()

				count, err := leftrecursionstate.Parse(
					"", []byte(testCase.expr), options...)
				if err != nil {
					t.Fatalf(
						"for input %q got error: %s, but expect to parse without errors",
						testCase.expr, err)
				}
				if count != testCase.want.count {
					t.Fatalf(
						"for input %q\ngot result: %d,\nbut expect: %d",
						testCase.expr, count, testCase.want.count)
				}
			})
		}

		t.Run(testCase.name+" optimized", func(t *testing.T) {
			t.Parallel()

			count, err := optimizedleftrecursionstate.Parse(
				"", []byte(testCase.expr),
				optimizedleftrecursionstate.InitState("count", initCount))
			if err != nil {
				t.Fatalf(
					"for input %q got error: %s, but expect to parse without errors",
					testCase.expr, err)
			}
			if count != testCase.want.count {
				t.Fatalf(
					"for input %q\ngot result: %q,\nbut expect: %q",
					testCase.expr, count, testCase.want.count)
			}
			if count != testCase.want.count {
				t.Fatalf(
					"for input %q\ngot result: %d,\nbut expect: %d",
					testCase.expr, count, testCase.want.count)
			}
		})
	}
}
