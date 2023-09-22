// Copied from https://github.com/stretchr/testify

// Copyright (c) 2012-2020 Mat Ryer, Tyler Bunnell and contributors. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in
// the THIRD-PARTY-NOTICES file.

package testutils_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/mna/pigeon/testutils"
)

func TestIsEmpty(t *testing.T) {
	t.Parallel()

	chWithValue := make(chan struct{}, 1)
	chWithValue <- struct{}{}

	tests := []struct {
		obj  interface{}
		want bool
	}{
		{obj: "", want: true},
		{obj: nil, want: true},
		{obj: []string{}, want: true},
		{obj: 0, want: true},
		{obj: int32(0), want: true},
		{obj: int64(0), want: true},
		{obj: false, want: true},
		{obj: map[string]string{}, want: true},
		{obj: new(time.Time), want: true},
		{obj: time.Time{}, want: true},
		{obj: make(chan struct{}), want: true},
		{obj: [1]int{}, want: true},
		{obj: "something", want: false},
		{obj: errors.New("something"), want: false},
		{obj: []string{"something"}, want: false},
		{obj: 1, want: false},
		{obj: true, want: false},
		{obj: map[string]string{"Hello": "World"}, want: false},
		{obj: chWithValue, want: false},
		{obj: [1]int{42}, want: false},
	}
	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("IsEmpty(%#v)", test.obj), func(t *testing.T) {
			t.Parallel()

			isEmpty := testutils.IsEmpty(test.obj)
			if isEmpty != test.want {
				t.Fatalf("IsEmpty(%#v) should return %v", test.obj, test.want)
			}
		})
	}
}

func TestVlidateEqualArgs(t *testing.T) {
	t.Parallel()

	if testutils.ValidateEqualArgs(func() {}, func() {}) == nil {
		t.Error("non-nil functions should error")
	}

	if testutils.ValidateEqualArgs(func() {}, func() {}) == nil {
		t.Error("non-nil functions should error")
	}

	if testutils.ValidateEqualArgs(nil, nil) != nil {
		t.Error("nil functions are equal")
	}
}

func TestEqual(t *testing.T) {
	t.Parallel()

	type myType string

	var m map[string]interface{}

	tests := []struct {
		expected interface{}
		actual   interface{}
		result   bool
		remark   string
	}{
		{"Hello World", "Hello World", true, ""},
		{123, 123, true, ""},
		{123.5, 123.5, true, ""},
		{[]byte("Hello World"), []byte("Hello World"), true, ""},
		{nil, nil, true, ""},
		{int32(123), int32(123), true, ""},
		{uint64(123), uint64(123), true, ""},
		{myType("1"), myType("1"), true, ""},
		{&struct{}{}, &struct{}{}, true, "pointer equality is based on equality of underlying value"},

		// Not expected to be equal
		{m["bar"], "something", false, ""},
		{myType("1"), myType("2"), false, ""},

		// A case that might be confusing, especially with numeric literals
		{10, uint(10), false, ""},
	}

	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("Equal(%#v, %#v)", test.expected, test.actual), func(t *testing.T) {
			t.Parallel()

			res := testutils.Equal(test.expected, test.actual)
			if res != test.result {
				t.Errorf(
					"Equal(%#v, %#v) should return %#v: %s",
					test.expected, test.actual, test.result, test.remark)
			}
		})
	}
}

func TestDiffLists(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		listA  interface{}
		listB  interface{}
		extraA []interface{}
		extraB []interface{}
	}{
		{
			name:   "equal empty",
			listA:  []string{},
			listB:  []string{},
			extraA: nil,
			extraB: nil,
		},
		{
			name:   "equal same order",
			listA:  []string{"hello", "world"},
			listB:  []string{"hello", "world"},
			extraA: nil,
			extraB: nil,
		},
		{
			name:   "equal different order",
			listA:  []string{"hello", "world"},
			listB:  []string{"world", "hello"},
			extraA: nil,
			extraB: nil,
		},
		{
			name:   "extra A",
			listA:  []string{"hello", "hello", "world"},
			listB:  []string{"hello", "world"},
			extraA: []interface{}{"hello"},
			extraB: nil,
		},
		{
			name:   "extra A twice",
			listA:  []string{"hello", "hello", "hello", "world"},
			listB:  []string{"hello", "world"},
			extraA: []interface{}{"hello", "hello"},
			extraB: nil,
		},
		{
			name:   "extra B",
			listA:  []string{"hello", "world"},
			listB:  []string{"hello", "hello", "world"},
			extraA: nil,
			extraB: []interface{}{"hello"},
		},
		{
			name:   "extra B twice",
			listA:  []string{"hello", "world"},
			listB:  []string{"hello", "hello", "world", "hello"},
			extraA: nil,
			extraB: []interface{}{"hello", "hello"},
		},
		{
			name:   "integers 1",
			listA:  []int{1, 2, 3, 4, 5},
			listB:  []int{5, 4, 3, 2, 1},
			extraA: nil,
			extraB: nil,
		},
		{
			name:   "integers 2",
			listA:  []int{1, 2, 1, 2, 1},
			listB:  []int{2, 1, 2, 1, 2},
			extraA: []interface{}{1},
			extraB: []interface{}{2},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actualExtraA, actualExtraB := testutils.DiffLists(
				test.listA, test.listB)
			if !testutils.Equal(test.extraA, actualExtraA) {
				t.Errorf(
					"extra A does not match for listA=%v listB=%v",
					test.listA, test.listB)
			}
			if !testutils.Equal(test.extraB, actualExtraB) {
				t.Errorf(
					"extra B does not match for listA=%v listB=%v",
					test.listA, test.listB)
			}
		})
	}
}

func TestObjectsAreEqual(t *testing.T) {
	t.Parallel()

	cases := []struct {
		expected interface{}
		actual   interface{}
		result   bool
	}{
		// cases that are expected to be equal
		{"Hello World", "Hello World", true},
		{123, 123, true},
		{123.5, 123.5, true},
		{[]byte("Hello World"), []byte("Hello World"), true},
		{nil, nil, true},

		// cases that are expected not to be equal
		{map[int]int{5: 10}, map[int]int{10: 20}, false},
		{'x', "x", false},
		{"x", 'x', false},
		{0, 0.1, false},
		{0.1, 0, false},
		{time.Now, time.Now, false},
		{func() {}, func() {}, false},
		{uint32(10), int32(10), false},
	}

	for _, test := range cases {
		test := test
		t.Run(fmt.Sprintf("ObjectsAreEqual(%#v, %#v)", test.expected, test.actual), func(t *testing.T) {
			t.Parallel()

			res := testutils.ObjectsAreEqual(test.expected, test.actual)
			if res != test.result {
				t.Errorf(
					"ObjectsAreEqual(%#v, %#v) should return %#v",
					test.expected, test.actual, test.result)
			}
		})
	}
}

func TestElementsMatch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expected interface{}
		actual   interface{}
		result   bool
	}{
		// matching
		{nil, nil, true},

		{nil, nil, true},
		{[]int{}, []int{}, true},
		{[]int{1}, []int{1}, true},
		{[]int{1, 1}, []int{1, 1}, true},
		{[]int{1, 2}, []int{1, 2}, true},
		{[]int{1, 2}, []int{2, 1}, true},
		{[2]int{1, 2}, [2]int{2, 1}, true},
		{[]string{"hello", "world"}, []string{"world", "hello"}, true},
		{[]string{"hello", "hello"}, []string{"hello", "hello"}, true},
		{[]string{"hello", "hello", "world"}, []string{"hello", "world", "hello"}, true},
		{[3]string{"hello", "hello", "world"}, [3]string{"hello", "world", "hello"}, true},
		{[]int{}, nil, true},

		// not matching
		{[]int{1}, []int{1, 1}, false},
		{[]int{1, 2}, []int{2, 2}, false},
		{[]string{"hello", "hello"}, []string{"hello"}, false},
	}

	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("ElementsMatch(%#v, %#v)", test.expected, test.actual), func(t *testing.T) {
			t.Parallel()

			res := testutils.ElementsMatch(test.actual, test.expected)
			if res != test.result {
				t.Errorf(
					"ElementsMatch(%#v, %#v) should return %v",
					test.actual, test.expected, test.result)
			}
		})
	}
}
