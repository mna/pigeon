package issue18

import "testing"

var cases = []struct {
	name string
	in   string
	exp  string
}{
	{
		name: "numbers",
		in:   `123455`,
		exp:  ``,
	}, {
		name: "numbers with whitespace",
		in: `

    1

    2

    3

	`,
		exp: ``,
	}, {
		name: "unexpected character",
		in: `

    1

    2

    x
	`,
		exp: `7:5 (20): no match found, expected: [ \t\r\n], [0-9] or EOF`,
	},
}

func TestErrorReporting(t *testing.T) {
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse("", []byte(tc.in))
			var got string
			if err != nil {
				got = err.Error()
			}

			if got != tc.exp {
				t.Errorf("%q: want %v, got %v", tc.name, tc.exp, got)
			}
		})
	}
}
