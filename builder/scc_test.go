package builder_test

import (
	"testing"

	"github.com/mna/pigeon/builder"
	"github.com/mna/pigeon/testutils"
)

func TestStronglyConnectedComponents(t *testing.T) { //nolint:funlen
	t.Parallel()

	type want struct {
		sccs []map[string]struct{}
	}

	tests := []struct {
		name  string
		graph map[string]map[string]struct{}
		want  want
	}{
		{
			name: "Simple",
			graph: map[string]map[string]struct{}{
				"1": {"2": {}},
				"2": {"1": {}},
			},
			want: want{sccs: []map[string]struct{}{
				{"2": {}, "1": {}},
			}},
		},
		{
			name: "Without scc",
			graph: map[string]map[string]struct{}{
				"1": {"2": {}},
			},
			want: want{sccs: []map[string]struct{}{
				{"2": {}},
				{"1": {}},
			}},
		},
		{
			name: "One element",
			graph: map[string]map[string]struct{}{
				"1": {},
			},
			want: want{sccs: []map[string]struct{}{
				{"1": {}},
			}},
		},
		{
			name: "One element with loop",
			graph: map[string]map[string]struct{}{
				"1": {"1": {}},
			},
			want: want{sccs: []map[string]struct{}{
				{"1": {}},
			}},
		},
		{
			name: "Wiki 1",
			graph: map[string]map[string]struct{}{
				"1": {"2": {}},
				"2": {"3": {}},
				"3": {"1": {}},
				"4": {"2": {}, "3": {}, "6": {}},
				"5": {"3": {}, "7": {}},
				"6": {"4": {}, "5": {}},
				"7": {"5": {}},
				"8": {"6": {}, "7": {}, "8": {}},
			},
			want: want{sccs: []map[string]struct{}{
				{"2": {}, "3": {}, "1": {}},
				{"5": {}, "7": {}},
				{"4": {}, "6": {}},
				{"8": {}},
			}},
		},
		{
			name: "Wiki 2",
			graph: map[string]map[string]struct{}{
				"1": {"2": {}, "6": {}},
				"2": {"6": {}, "4": {}},
				"3": {"9": {}, "4": {}, "8": {}},
				"4": {"1": {}, "7": {}},
				"5": {"9": {}, "8": {}},
				"6": {"1": {}, "4": {}, "7": {}},
				"7": {"1": {}},
				"8": {"5": {}, "3": {}},
				"9": {"8": {}},
			},
			want: want{sccs: []map[string]struct{}{
				{"1": {}, "2": {}, "4": {}, "6": {}, "7": {}},
				{"3": {}, "5": {}, "9": {}, "8": {}},
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
			sccs := builder.StronglyConnectedComponents(vertices, testCase.graph)
			if !testutils.ElementsMatch(sccs, testCase.want.sccs) {
				t.Fatalf("Result %v, expected %v", sccs, testCase.want.sccs)
			}
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
		graph map[string]map[string]struct{}
		scc   map[string]struct{}
		start string
		want  want
	}{
		{
			name: "Wiki 1 1",
			graph: map[string]map[string]struct{}{
				"1": {"2": {}},
				"2": {"3": {}},
				"3": {"1": {}},
				"4": {"2": {}, "3": {}, "6": {}},
				"5": {"3": {}, "7": {}},
				"6": {"4": {}, "5": {}},
				"7": {"5": {}},
				"8": {"6": {}, "7": {}, "8": {}},
			},
			scc:   map[string]struct{}{"2": {}, "3": {}, "1": {}},
			start: "3",
			want:  want{paths: [][]string{{"3", "1", "2", "3"}}},
		},
		{
			name: "Wiki 1 2",
			graph: map[string]map[string]struct{}{
				"1": {"2": {}},
				"2": {"3": {}},
				"3": {"1": {}},
				"4": {"2": {}, "3": {}, "6": {}},
				"5": {"3": {}, "7": {}},
				"6": {"4": {}, "5": {}},
				"7": {"5": {}},
				"8": {"6": {}, "7": {}, "8": {}},
			},
			scc:   map[string]struct{}{"5": {}, "7": {}},
			start: "5",
			want:  want{paths: [][]string{{"5", "7", "5"}}},
		},
		{
			name: "Wiki 2",
			graph: map[string]map[string]struct{}{
				"1": {"2": {}, "6": {}},
				"2": {"6": {}, "4": {}},
				"3": {"9": {}, "4": {}, "8": {}},
				"4": {"1": {}, "7": {}},
				"5": {"9": {}, "8": {}},
				"6": {"1": {}, "4": {}, "7": {}},
				"7": {"1": {}},
				"8": {"5": {}, "3": {}},
				"9": {"8": {}},
			},
			scc: map[string]struct{}{
				"1": {}, "2": {}, "4": {}, "6": {}, "7": {},
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
			graph: map[string]map[string]struct{}{
				"1": {"2": {}},
				"2": {"3": {}},
				"3": {"1": {}, "2": {}},
			},
			scc: map[string]struct{}{
				"1": {}, "2": {}, "3": {},
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
			paths, err := builder.FindCyclesInSCC(
				testCase.graph, testCase.scc, testCase.start)
			if err != nil {
				t.FailNow()
			}
			if !testutils.ElementsMatch(paths, testCase.want.paths) {
				t.Fatalf("Result %v, expected %v", paths, testCase.want.paths)
			}
		})
	}
}
