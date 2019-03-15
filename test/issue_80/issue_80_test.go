package issue80

import "testing"

var cases = map[string]interface{}{
	"12345": 12345,
	"asdf": "asdf",
}

func TestIgnoreComments(t *testing.T) {

	for tc, exp := range cases {
		ret, err := Parse("", []byte(tc))
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}

		switch expt := exp.(type) {
		case int:
			got, ok := ret.(int)
			if !ok {
				t.Fatalf("incorrect output type %T for case %v, expected %T",
				ret, tc, expt)
			}

			if (expt != got) {
				t.Fatalf("incorrect output got %v, expected %v", got, exp)
			}
		case string:
			got, ok := ret.(string)
			if !ok {
				t.Fatalf("incorrect output type %T for case %v, expected %T",
				ret, tc, expt)
			}

			if (expt != got) {
				t.Fatalf("incorrect output got %v, expected %v", got, exp)
			}
		}
	}
}


