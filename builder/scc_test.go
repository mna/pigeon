package builder_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/mna/pigeon/builder"
)

// isEmpty gets whether the specified object is considered empty or not.
func isEmpty(object interface{}) bool {
	// get nil case out of the way
	if object == nil {
		return true
	}

	objValue := reflect.ValueOf(object)

	switch objValue.Kind() {
	// collection types are empty when they have no element
	case reflect.Chan, reflect.Map, reflect.Slice:
		return objValue.Len() == 0
	// pointers are empty if nil or if the value they point to is empty
	case reflect.Ptr:
		if objValue.IsNil() {
			return true
		}
		deref := objValue.Elem().Interface()
		return isEmpty(deref)
	// for all other types, compare against the zero value
	// array types are empty when they match their zero-initialized state
	default:
		zero := reflect.Zero(objValue.Type())
		return reflect.DeepEqual(object, zero.Interface())
	}
}

// isList checks that the provided value is array or slice.
func isList(list interface{}) (ok bool) {
	kind := reflect.TypeOf(list).Kind()
	return kind == reflect.Array || kind == reflect.Slice
}

// diffLists diffs two arrays/slices and returns slices of elements that are only in A and only in B.
// If some element is present multiple times, each instance is counted separately (e.g. if something is 2x in A and
// 5x in B, it will be 0x in extraA and 3x in extraB). The order of items in both lists is ignored.
func diffLists(listA, listB interface{}) (extraA, extraB []interface{}) {
	aValue := reflect.ValueOf(listA)
	bValue := reflect.ValueOf(listB)

	aLen := aValue.Len()
	bLen := bValue.Len()

	// Mark indexes in bValue that we already used
	visited := make([]bool, bLen)
	for i := 0; i < aLen; i++ {
		element := aValue.Index(i).Interface()
		found := false
		for j := 0; j < bLen; j++ {
			if visited[j] {
				continue
			}
			if ObjectsAreEqual(bValue.Index(j).Interface(), element) {
				visited[j] = true
				found = true
				break
			}
		}
		if !found {
			extraA = append(extraA, element)
		}
	}

	for j := 0; j < bLen; j++ {
		if visited[j] {
			continue
		}
		extraB = append(extraB, bValue.Index(j).Interface())
	}

	return
}

// ObjectsAreEqual determines if two objects are considered equal.
//
// This function does no assertion of any kind.
func ObjectsAreEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}

// ElementsMatch asserts that the specified listA(array, slice...) is equal to specified
// listB(array, slice...) ignoring the order of the elements. If there are duplicate elements,
// the number of appearances of each of them in both lists should match.
//
// ElementsMatch([1, 3, 2, 3], [1, 3, 3, 2]).
func ElementsMatch(listA interface{}, listB interface{}) bool {
	if isEmpty(listA) && isEmpty(listB) {
		return true
	}

	if !isList(listA) || !isList(listB) {
		return false
	}

	extraA, extraB := diffLists(listA, listB)

	return len(extraA) == 0 && len(extraB) == 0
}

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
			sccs := builder.StronglyConnectedComponents(vertices, testCase.graph)
			if !ElementsMatch(sccs, testCase.want.sccs) {
				t.FailNow()
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
			paths := builder.FindCyclesInSCC(
				testCase.graph, testCase.scc, testCase.start)
			if !ElementsMatch(paths, testCase.want.paths) {
				t.FailNow()
			}
		})
	}
}
