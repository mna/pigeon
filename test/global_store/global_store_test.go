package globalstore

import "testing"

func TestGlobalStore(t *testing.T) {
	cases := []struct {
		initial  int
		input    string
		expected int
	}{
		{
			input:    "i",
			expected: 1,
		},
		{
			initial:  -1,
			input:    "ii",
			expected: 1,
		},
		{
			input:    "d",
			expected: -1,
		},
		{
			initial:  10,
			input:    "iizddiii",
			expected: 1,
		},
	}
	for _, test := range cases {
		got, err := Parse("", []byte(test.input), GlobalStore("initial", test.initial))
		if err != nil {
			t.Fatalf("for input %q got error: %s", test.input, err)
		}
		goti := got.(int)
		if goti != test.expected {
			t.Errorf("for input %q want %d, got %d", test.input, test.expected, goti)
		}
	}
}

func TestChoiceStatistics(t *testing.T) {
	for _, tracked := range []bool{false, true} {
		name := "default"
		if tracked {
			name = "requested"
		}
		t.Run(name, func(t *testing.T) {
			opts := []Option{GlobalStore("initial", 0)}
			var stats Stats
			if tracked {
				opts = append(opts, Statistics(&stats, "no match"))
			}
			p := newParser("", []byte("iid"), opts...)
			got, err := p.parse(g)
			if err != nil || got != 1 {
				t.Fatalf("parse returned %v, %v; want 1, nil", got, err)
			}
			if p.ExprCnt == 0 {
				t.Fatal("expression counting must remain enabled")
			}
			if !tracked {
				if len(p.ChoiceAltCnt) != 0 {
					t.Fatal("unrequested choice statistics were collected")
				}
				return
			}
			if len(stats.ChoiceAltCnt) != 1 {
				t.Fatalf("got statistics for %d choices, want 1", len(stats.ChoiceAltCnt))
			}
			for _, counts := range stats.ChoiceAltCnt {
				if counts["1"] != 2 || counts["2"] != 1 || counts["3"] != 0 || counts["no match"] != 1 {
					t.Fatalf("unexpected choice counts: %v", counts)
				}
			}
		})
	}
}
