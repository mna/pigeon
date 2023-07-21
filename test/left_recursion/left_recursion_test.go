package leftrecursion_test

import (
	"testing"

	"github.com/mna/pigeon/test/left_recursion/standart/leftrecursion"
	"github.com/mna/pigeon/test/left_recursion/standart/withoutleftrecursion"

	optimizedleftrecursion "github.com/mna/pigeon/test/left_recursion/optimized/leftrecursion"
	optimizedwithoutleftrecursion "github.com/mna/pigeon/test/left_recursion/optimized/withoutleftrecursion"
)

func TestLeftRecursionParse(t *testing.T) {
	t.Parallel()

	type want struct {
		expr string
	}

	tests := []struct {
		name string
		expr string
		want want
	}{
		{
			name: "Complex",
			expr: "7+10/2*-4+5*3%6-8*6",
			want: want{expr: "(((7+((10/2)*(-4)))+((5*3)%6))-(8*6))"},
		},
		{
			name: "Simple",
			expr: "2*1+7",
			want: want{expr: "((2*1)+7)"},
		},
		{
			name: "Simple revers",
			expr: "2+1*7",
			want: want{expr: "(2+(1*7))"},
		},
		{
			name: "Same operations",
			expr: "2+1+7",
			want: want{expr: "((2+1)+7)"},
		},
		{
			name: "Start with unary minus",
			expr: "-2+1",
			want: want{expr: "((-2)+1)"},
		},
		{
			name: "unary minus between + and *",
			expr: "2+-7*-1",
			want: want{expr: "(2+((-7)*(-1)))"},
		},
	}

	for _, testCase := range tests {
		testCase := testCase

		setOptionsLR := map[string][]leftrecursion.Option{
			"memoize": {leftrecursion.Memoize(true)},
			"-":       {},
		}
		for nameOptionsLR, optionsLR := range setOptionsLR {
			optionsLR := optionsLR

			t.Run(
				testCase.name+" default(recursion). Options: "+nameOptionsLR,
				func(t *testing.T) {
					t.Parallel()

					resLR, err := leftrecursion.Parse(
						"", []byte(testCase.expr), optionsLR...)
					if err != nil {
						t.Fatalf(
							"for input %q got error: %s, but expect to parse without errors",
							testCase.expr, err)
					}
					exprLR, ok := resLR.(string)
					if !ok {
						t.FailNow()
					}
					if exprLR != testCase.want.expr {
						t.Fatalf(
							"for input %q\ngot result: %q,\nbut expect: %q",
							testCase.expr, exprLR, testCase.want.expr)
					}
				})
		}

		setOptions := map[string][]withoutleftrecursion.Option{
			"memoize": {withoutleftrecursion.Memoize(true)},
			"-":       {},
		}
		for nameOptions, options := range setOptions {
			options := options

			t.Run(testCase.name+" default(without recursion). Options: "+nameOptions, func(t *testing.T) {
				t.Parallel()

				res, err := withoutleftrecursion.Parse("", []byte(testCase.expr), options...)
				if err != nil {
					t.Fatalf(
						"for input %q got error: %s, but expect to parse without errors",
						testCase.expr, err)
				}
				expr, ok := res.(string)
				if !ok {
					t.FailNow()
				}
				if expr != testCase.want.expr {
					t.Fatalf(
						"for input %q\ngot result: %q,\nbut expect: %q",
						testCase.expr, expr, testCase.want.expr)
				}
			})
		}

		t.Run(testCase.name+" optimized(recursion)", func(t *testing.T) {
			t.Parallel()

			resLR, err := optimizedleftrecursion.Parse("", []byte(testCase.expr))
			if err != nil {
				t.Fatalf(
					"for input %q got error: %s, but expect to parse without errors",
					testCase.expr, err)
			}
			exprLR, ok := resLR.(string)
			if !ok {
				t.FailNow()
			}
			if exprLR != testCase.want.expr {
				t.Fatalf(
					"for input %q\ngot result: %q,\nbut expect: %q",
					testCase.expr, exprLR, testCase.want.expr)
			}
		})

		t.Run(testCase.name+" optimized(without recursion)", func(t *testing.T) {
			t.Parallel()

			res, err := optimizedwithoutleftrecursion.Parse("", []byte(testCase.expr))
			if err != nil {
				t.Fatalf(
					"for input %q got error: %s, but expect to parse without errors",
					testCase.expr, err)
			}
			expr, ok := res.(string)
			if !ok {
				t.FailNow()
			}
			if expr != testCase.want.expr {
				t.Fatalf(
					"for input %q\ngot result: %q,\nbut expect: %q",
					testCase.expr, expr, testCase.want.expr)
			}
		})
	}
}

func FuzzLeftRecursionParse(f *testing.F) {
	chars := []byte("0123456789+-/*%")

	f.Fuzz(func(t *testing.T, bytes []byte) {
		data := make([]byte, 0, len(bytes))
		for _, b := range bytes {
			data = append(data, chars[int(b)%len(chars)])
		}
		resLR, errLR := leftrecursion.Parse("", data)
		res, err := withoutleftrecursion.Parse("", data)
		if err != nil || errLR != nil {
			if err == nil || errLR == nil {
				t.Fatalf(
					"for input %q\ngot error: %q,\nbut expect: %q",
					data, errLR, err)
			}
			return
		}
		exprLR, okLR := resLR.(string)
		if !okLR {
			t.FailNow()
		}
		expr, ok := res.(string)
		if !ok {
			t.FailNow()
		}
		if expr != exprLR {
			t.Fatalf(
				"for input %q\ngot result: %q,\nbut expect: %q",
				data, exprLR, expr)
		}
	})
}

func FuzzLeftRecursionParseMemoize(f *testing.F) {
	chars := []byte("0123456789+-/*%")

	f.Fuzz(func(t *testing.T, bytes []byte) {
		data := make([]byte, 0, len(bytes))
		for _, b := range bytes {
			data = append(data, chars[int(b)%len(chars)])
		}

		resLR, errLR := leftrecursion.Parse(
			"", data, leftrecursion.Memoize(true))
		res, err := withoutleftrecursion.Parse(
			"", data, withoutleftrecursion.Memoize(true))
		if err != nil || errLR != nil {
			if err == nil || errLR == nil {
				t.Fatalf(
					"for input %q\ngot error: %q,\nbut expect: %q",
					data, errLR, err)
			}
			return
		}
		exprLR, okLR := resLR.(string)
		if !okLR {
			t.FailNow()
		}
		expr, ok := res.(string)
		if !ok {
			t.FailNow()
		}
		if expr != exprLR {
			t.Fatalf(
				"for input %q\ngot result: %q,\nbut expect: %q",
				data, exprLR, expr)
		}
	})
}

func FuzzLeftRecursionParseOptimized(f *testing.F) {
	chars := []byte("0123456789+-/*%")

	f.Fuzz(func(t *testing.T, bytes []byte) {
		data := make([]byte, 0, len(bytes))
		for _, b := range bytes {
			data = append(data, chars[int(b)%len(chars)])
		}
		resLR, errLR := optimizedleftrecursion.Parse("", data)
		res, err := optimizedwithoutleftrecursion.Parse("", data)
		if err != nil || errLR != nil {
			if err == nil || errLR == nil {
				t.Fatalf(
					"for input %q\ngot error: %q,\nbut expect: %q",
					data, errLR, err)
			}
			return
		}
		exprLR, okLR := resLR.(string)
		if !okLR {
			t.FailNow()
		}
		expr, ok := res.(string)
		if !ok {
			t.FailNow()
		}
		if expr != exprLR {
			t.Fatalf(
				"for input %q\ngot result: %q,\nbut expect: %q",
				data, exprLR, expr)
		}
	})
}
