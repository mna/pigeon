package leftrecursionthrownrecover_test

import (
	"errors"
	"reflect"
	"testing"

	leftrecursionthrownrecover "github.com/mna/pigeon/test/left_recursion_thrownrecover"
)

func TestLeftRecursionWithThrowAndRecover(t *testing.T) {
	t.Parallel()

	type want struct {
		captures any
		errors   []string
	}

	cases := []struct {
		name       string
		entrypoint string
		input      string
		want       want
	}{
		// Case 01: Recover multiple labels
		{
			name:       "Case 01: Recover multiple labels[correct]",
			entrypoint: "case01",
			input:      "123",
			want:       want{captures: "123"},
		},
		{
			name:       "Case 01: Recover multiple labels[second character is not a number]",
			entrypoint: "case01",
			input:      "1a3",
			want: want{
				captures: "1?3",
				errors: []string{
					"1:2 (1): rule ErrNonNumber: expecting a number",
				},
			},
		},
		{
			name:       "Case 01: Recover multiple labels[third character is not a number]",
			entrypoint: "case01",
			input:      "11+3",
			want: want{
				captures: "11?3",
				errors: []string{
					"1:3 (2): rule ErrNonNumber: expecting a number",
				},
			},
		},

		// Case 02: Throw a undefined label
		{
			name:       "Case 02: Throw a undefined label",
			entrypoint: "case02",
			input:      "",
			want: want{
				captures: nil,
				errors: []string{
					"1:1 (0): rule case02: Throwed undefined label",
				},
			},
		},

		// Case 03: Nested Recover
		{
			name:       "Case 03: Nested Recover[correct]",
			entrypoint: "case03",
			input:      "123",
			want:       want{captures: "123"},
		},
		{
			name:       "Case 03: Nested Recover[second character is lower case char]",
			entrypoint: "case03",
			input:      "1a3",
			want: want{
				captures: "1<3",
				errors: []string{
					"1:2 (1): rule ErrAlphaInner03: expecting a number, got lower case char",
				},
			},
		},
		{
			name:       "Case 03: Nested Recover[third character is upper case char]",
			entrypoint: "case03",
			input:      "11A3",
			want: want{
				captures: "11>3",
				errors: []string{
					"1:3 (2): rule ErrAlphaOuter03: expecting a number, got upper case char",
				},
			},
		},
		{
			name:       "Case 03: Nested Recover[fourth character is non-char]",
			entrypoint: "case03",
			input:      "111+3",
			want: want{
				captures: "111?3",
				errors: []string{
					"1:4 (3): rule ErrOtherOuter03: expecting a number, got a non-char",
				},
			},
		},

		// Case 04: Nested Recover, which fails in inner recover
		{
			name:       "Case 04: Nested Recover, which fails in inner recover[correct]",
			entrypoint: "case04",
			input:      "123",
			want:       want{captures: "123"},
		},
		{
			name:       "Case 04: Nested Recover, which fails in inner recover[second character is lower case char]",
			entrypoint: "case04",
			input:      "1a3",
			want: want{
				captures: "1x3",
				errors: []string{
					"1:2 (1): rule ErrAlphaOuter04: expecting a number, got a char",
				},
			},
		},
		{
			name:       "Case 04: Nested Recover, which fails in inner recover[third character is upper case char]",
			entrypoint: "case04",
			input:      "11A3",
			want: want{
				captures: "11x3",
				errors: []string{
					"1:3 (2): rule ErrAlphaOuter04: expecting a number, got a char",
				},
			},
		},
		{
			name:       "Case 04: Nested Recover, which fails in inner recover[fourth character is non-char]",
			entrypoint: "case04",
			input:      "111+3",
			want: want{
				captures: "111?3",
				errors: []string{
					"1:4 (3): rule ErrOtherOuter04: expecting a number, got a non-char",
				},
			},
		},
	}
	for _, testCase := range cases {
		testCase := testCase

		setOptions := map[string][]leftrecursionthrownrecover.Option{
			"memoize": {
				leftrecursionthrownrecover.Memoize(true),
				leftrecursionthrownrecover.Entrypoint(testCase.entrypoint),
			},
			"-": {
				leftrecursionthrownrecover.Entrypoint(testCase.entrypoint),
			},
		}
		for nameOptions, options := range setOptions {
			options := options

			t.Run(testCase.name+". Options: "+nameOptions, func(t *testing.T) {
				t.Parallel()

				got, err := leftrecursionthrownrecover.Parse(
					"", []byte(testCase.input), options...)
				if testCase.want.errors == nil && err != nil {
					t.Fatalf(
						"for input %q got error: %s, but expect to parse without errors",
						testCase.input, err)
				}
				if testCase.want.errors != nil && err == nil {
					t.Fatalf(
						"for input %q got no error, but expect to parse with errors: %s",
						testCase.input, testCase.want.errors)
				}
				if !reflect.DeepEqual(got, testCase.want.captures) {
					t.Errorf(
						"for input %q want %s, got %s",
						testCase.input, testCase.want.captures, got)
				}
				if err != nil {
					var errorLister leftrecursionthrownrecover.ErrorLister
					if !errors.As(err, &errorLister) {
						t.FailNow()
					}
					list := errorLister.Errors()
					if len(list) != len(testCase.want.errors) {
						t.Errorf(
							"for input %q want %d error(s), got %d",
							testCase.input, len(testCase.want.errors), len(list))
						t.Logf("expected errors:\n")
						for _, ee := range testCase.want.errors {
							t.Logf("- %s\n", ee)
						}
						t.Logf("got errors:\n")
						for _, ee := range list {
							t.Logf("- %s\n", ee)
						}
						t.FailNow()
					}
					for index, err := range list {
						var parserError leftrecursionthrownrecover.ParserError
						if !errors.As(err, &parserError) {
							t.FailNow()
						}
						if parserError.Error() != testCase.want.errors[index] {
							t.Errorf(
								"for input %q want %dth error to be %s, got %s",
								testCase.input, index+1,
								testCase.want.errors[index], parserError)
						}
					}
				}
			})
		}
	}
}
