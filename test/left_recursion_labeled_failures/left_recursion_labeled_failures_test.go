package leftrecursionlabeledfailures_test

import (
	"reflect"
	"testing"

	leftrecursionlabeledfailures "github.com/mna/pigeon/test/left_recursion_labeled_failures"
)

func TestLeftRecursionWithLabeledFailures(t *testing.T) {
	t.Parallel()

	type want struct {
		captures []string
		errors   []string
	}

	cases := []struct {
		name  string
		input string
		want  want
	}{
		// Test cases from reference implementation peglabel:
		// https://github.com/sqmedeiros/lpeglabel/blob/976b38458e0bba58ca748e96b53afd9ee74a1d1d/README.md#relabel-syntax
		// https://github.com/sqmedeiros/lpeglabel/blame/976b38458e0bba58ca748e96b53afd9ee74a1d1d/README.md#L418-L440
		{
			name:  "correct",
			input: "one,two",
			want:  want{captures: []string{"one", "two"}},
		},
		{
			name:  "missing commas",
			input: "one two three",
			want: want{
				captures: []string{"one", "two", "three"},
				errors: []string{
					"1:4 (3): rule ErrComma: expecting ','",
					"1:8 (7): rule ErrComma: expecting ','",
				},
			},
		},
		{
			name:  "missing id and incorrect ids",
			input: "1,\n two, \n3,",
			want: want{
				captures: []string{"NONE", "two", "NONE", "NONE"},
				errors: []string{
					"1:1 (0): rule ErrID: expecting an identifier",
					"2:6 (8): rule ErrID: expecting an identifier",
					// is line 3, col 2 in peglabel, pigeon increments the position
					// behind the last character of the input if !. is matched
					"3:3 (12): rule ErrID: expecting an identifier",
				},
			},
		},
		{
			name:  "missing comma, id and incorrect id",
			input: "one\n two123, \nthree,",
			want: want{
				captures: []string{"one", "two", "three", "NONE"},
				errors: []string{
					// is line 2, col 1 in peglabel, in pigeon, if a \n causes
					// an error, this is at col 0
					"2:0 (3): rule ErrComma: expecting ','",
					"2:5 (8): rule ErrComma: expecting ','",
					// is line 3, col 6 in peglabel, pigeon increments the position
					// behind the last character of the input if !. is matched
					"3:7 (20): rule ErrID: expecting an identifier",
				},
			},
		},
		// Additional test cases
		{
			name:  "empty",
			input: "",
			want: want{
				captures: []string{"NONE"},
				errors: []string{
					"1:1 (0): rule ErrID: expecting an identifier",
				},
			},
		},
		{
			name:  "incorrect id",
			input: "1",
			want: want{
				captures: []string{"NONE"},
				errors:   []string{"1:1 (0): rule ErrID: expecting an identifier"},
			},
		},
		{
			name:  "incorrect ids",
			input: "1,2",
			want: want{
				captures: []string{"NONE", "NONE"},
				errors: []string{
					"1:1 (0): rule ErrID: expecting an identifier",
					"1:3 (2): rule ErrID: expecting an identifier",
				},
			},
		},
	}
	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got, err := leftrecursionlabeledfailures.Parse(
				"", []byte(testCase.input))
			if testCase.want.errors == nil && err != nil {
				t.Fatalf(
					"for input %q got error: %s, but expect to parse without errors",
					testCase.input, err)
			}
			if !reflect.DeepEqual(got, testCase.want.captures) {
				t.Errorf(
					"for input %q want %s, got %s",
					testCase.input, testCase.want.captures, got)
			}
			if err != nil {
				errorLister, ok := err.(leftrecursionlabeledfailures.ErrorLister)
				if !ok {
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
					pe, ok := err.(leftrecursionlabeledfailures.ParserError)
					if !ok {
						t.FailNow()
					}
					if pe.Error() != testCase.want.errors[index] {
						t.Errorf(
							"for input %q want %dth error to be %s, got %s",
							testCase.input, index+1, testCase.want.errors[index], pe)
					}
				}
			}
		})
	}
}
