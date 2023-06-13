package builder_test

import (
	"testing"

	"github.com/mna/pigeon/builder"
	"github.com/stretchr/testify/require"
)

func TestStronglyConnectedComponents(t *testing.T) { //nolint:funlen
	t.Parallel()

	type want struct {
		sccs []map[string]bool
	}

	tests := []struct {
		name  string
		graph map[string]map[string]bool
		want  want
	}{
		{
			name: "Simple",
			graph: map[string]map[string]bool{
				"1": {"2": true},
				"2": {"1": true},
			},
			want: want{sccs: []map[string]bool{
				{"2": true, "1": true},
			}},
		},
		{
			name: "Without scc",
			graph: map[string]map[string]bool{
				"1": {"2": true},
			},
			want: want{sccs: []map[string]bool{
				{"2": true},
				{"1": true},
			}},
		},
		{
			name: "One element",
			graph: map[string]map[string]bool{
				"1": {},
			},
			want: want{sccs: []map[string]bool{
				{"1": true},
			}},
		},
		{
			name: "One element with loop",
			graph: map[string]map[string]bool{
				"1": {"1": true},
			},
			want: want{sccs: []map[string]bool{
				{"1": true},
			}},
		},
		{
			name: "Wiki 1",
			graph: map[string]map[string]bool{
				"1": {"2": true},
				"2": {"3": true},
				"3": {"1": true},
				"4": {"2": true, "3": true, "6": true},
				"5": {"3": true, "7": true},
				"6": {"4": true, "5": true},
				"7": {"5": true},
				"8": {"6": true, "7": true, "8": true},
			},
			want: want{sccs: []map[string]bool{
				{"2": true, "3": true, "1": true},
				{"5": true, "7": true},
				{"4": true, "6": true},
				{"8": true},
			}},
		},
		{
			name: "Wiki 2",
			graph: map[string]map[string]bool{
				"1": {"2": true, "6": true},
				"2": {"6": true, "4": true},
				"3": {"9": true, "4": true, "8": true},
				"4": {"1": true, "7": true},
				"5": {"9": true, "8": true},
				"6": {"1": true, "4": true, "7": true},
				"7": {"1": true},
				"8": {"5": true, "3": true},
				"9": {"8": true},
			},
			want: want{sccs: []map[string]bool{
				{"1": true, "2": true, "4": true, "6": true, "7": true},
				{"3": true, "5": true, "9": true, "8": true},
			}},
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			vertices := make([]string, 0, len(testCase.graph))
			for k := range testCase.graph {
				vertices = append(vertices, k)
			}
			require.ElementsMatch(t, builder.StronglyConnectedComponents(
				vertices, testCase.graph), testCase.want.sccs)
		})
	}
}

func TestFindCyclesInSCC(t *testing.T) { //nolint:funlen
	t.Parallel()

	type want struct {
		paths [][]string
	}

	tests := []struct {
		name  string
		graph map[string]map[string]bool
		scc   map[string]bool
		start string
		want  want
	}{
		{
			name: "Wiki 1 1",
			graph: map[string]map[string]bool{
				"1": {"2": true},
				"2": {"3": true},
				"3": {"1": true},
				"4": {"2": true, "3": true, "6": true},
				"5": {"3": true, "7": true},
				"6": {"4": true, "5": true},
				"7": {"5": true},
				"8": {"6": true, "7": true, "8": true},
			},
			scc:   map[string]bool{"2": true, "3": true, "1": true},
			start: "3",
			want:  want{paths: [][]string{{"3", "1", "2", "3"}}},
		},
		{
			name: "Wiki 1 2",
			graph: map[string]map[string]bool{
				"1": {"2": true},
				"2": {"3": true},
				"3": {"1": true},
				"4": {"2": true, "3": true, "6": true},
				"5": {"3": true, "7": true},
				"6": {"4": true, "5": true},
				"7": {"5": true},
				"8": {"6": true, "7": true, "8": true},
			},
			scc:   map[string]bool{"5": true, "7": true},
			start: "5",
			want:  want{paths: [][]string{{"5", "7", "5"}}},
		},
		{
			name: "Wiki 2",
			graph: map[string]map[string]bool{
				"1": {"2": true, "6": true},
				"2": {"6": true, "4": true},
				"3": {"9": true, "4": true, "8": true},
				"4": {"1": true, "7": true},
				"5": {"9": true, "8": true},
				"6": {"1": true, "4": true, "7": true},
				"7": {"1": true},
				"8": {"5": true, "3": true},
				"9": {"8": true},
			},
			scc: map[string]bool{
				"1": true, "2": true, "4": true, "6": true, "7": true,
			},
			start: "1",
			want: want{paths: [][]string{
				{"1", "2", "6", "1"},
				{"1", "2", "6", "4", "1"},
				{"1", "2", "6", "4", "7", "1"},
				{"1", "2", "6", "7", "1"},
				{"1", "2", "4", "1"},
				{"1", "2", "4", "7", "1"},
				{"1", "6", "1"},
				{"1", "6", "7", "1"},
				{"1", "6", "4", "7", "1"},
				{"1", "6", "4", "1"},
			}},
		},
		{
			name: "loop in loop",
			graph: map[string]map[string]bool{
				"1": {"2": true},
				"2": {"3": true},
				"3": {"1": true, "2": true},
			},
			scc: map[string]bool{
				"1": true, "2": true, "3": true,
			},
			start: "1",
			want: want{paths: [][]string{
				{"1", "2", "3", "1"},
				{"1", "2", "3", "2"},
			}},
		},
	}
	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			require.ElementsMatch(t, builder.FindCyclesInSCC(
				testCase.graph, testCase.scc, testCase.start),
				testCase.want.paths)
		})
	}
}
