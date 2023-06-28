package leftrecursion

import (
	"testing"
)

func TestLeftRecursion(t *testing.T) {
	t.Parallel()

	data := "7+10/2*-4+5*3%6-8*6"
	res, err := Parse("", []byte(data))
	if err != nil {
		t.Fatalf(
			"for input %q got error: %s, but expect to parse without errors",
			data, err)
	}
	str, ok := res.(string)
	if !ok {
		t.FailNow()
	}
	want := "(((7+((10/2)*(-4)))+((5*3)%6))-(8*6))"
	if str != want {
		t.Fatalf(
			"for input %q\ngot result: %q,\nbut expect: %q", data, str, want)
	}
}
