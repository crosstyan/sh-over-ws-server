package utils

import (
	"testing"

	"golang.org/x/exp/slices"
)

type testCase[T any] struct {
	lhs []T
	rhs []T
	el  T
}

func TestDeleteIfAll1(t *testing.T) {
	cases := []testCase[int]{
		testCase[int]{lhs: []int{1, 2, 1, 3, 4, 5, 1}, rhs: []int{2, 3, 4, 5}, el: 1},
		testCase[int]{lhs: []int{5, 5, 2, 1, 5, 5, 5}, rhs: []int{2, 1}, el: 5},
		testCase[int]{lhs: []int{1, 2, 3, 4, 5}, rhs: []int{1, 2, 3, 4, 5}, el: 6},
		testCase[int]{lhs: []int{1, 2, 3, 3, 4, 5, 5}, rhs: []int{1, 2, 4, 5, 5}, el: 3},
		testCase[int]{lhs: []int{}, rhs: []int{}, el: 1},
		testCase[int]{lhs: []int{1, 2, 3, 4, 5}, rhs: []int{1, 3, 4, 5}, el: 2},
		testCase[int]{lhs: []int{1, 1, 1, 1, 1}, rhs: []int{}, el: 1},
	}

	string_cases := []testCase[string]{
		testCase[string]{lhs: []string{"a", "b", "a", "c", "d", "e", "a"}, rhs: []string{"b", "c", "d", "e"}, el: "a"},
		testCase[string]{lhs: []string{"e", "e", "b", "a", "e", "e", "e"}, rhs: []string{"b", "a"}, el: "e"},
		testCase[string]{lhs: []string{"a", "b", "c", "d", "e"}, rhs: []string{"a", "b", "c", "d", "e"}, el: "f"},
		testCase[string]{lhs: []string{"a", "b", "c", "c", "d", "e", "e"}, rhs: []string{"a", "b", "d", "e", "e"}, el: "c"},
	}

	for _, c := range cases {
		res := DeleteIfAll(
			c.lhs,
			c.el,
			func(a, b int) bool {
				return a == b
			})
		if !slices.Equal(res, c.rhs) {
			t.Errorf("DeleteIfAll failed, got: %v, want: %v.", res, c.rhs)
		}
	}

	for _, c := range string_cases {
		res := DeleteIfAll(
			c.lhs,
			c.el,
			func(a, b string) bool {
				return a == b
			})
		if !slices.Equal(res, c.rhs) {
			t.Errorf("DeleteIfAll failed, got: %v, want: %v.", res, c.rhs)
		}
	}
}

func TestDeleteIfOnce(t *testing.T) {
	cases := []testCase[int]{
		testCase[int]{lhs: []int{1, 2, 1, 3, 4, 5, 1}, rhs: []int{2, 1, 3, 4, 5, 1}, el: 1},
		testCase[int]{lhs: []int{5, 5, 2, 1, 5, 5, 5}, rhs: []int{5, 2, 1, 5, 5, 5}, el: 5},
		testCase[int]{lhs: []int{1, 2, 3, 4, 5}, rhs: []int{1, 2, 3, 4, 5}, el: 6},
		testCase[int]{lhs: []int{1, 2, 3, 3, 4, 5, 5}, rhs: []int{1, 2, 3, 4, 5, 5}, el: 3},
		testCase[int]{lhs: []int{}, rhs: []int{}, el: 1},
		testCase[int]{lhs: []int{1, 2, 3, 4, 5}, rhs: []int{1, 3, 4, 5}, el: 2},
		testCase[int]{lhs: []int{1, 1, 1, 1, 1}, rhs: []int{1, 1, 1, 1}, el: 1},
	}
	for _, c := range cases {
		res := DeleteIfOnce(
			c.lhs,
			c.el,
			func(a, b int) bool {
				return a == b
			})
		if !slices.Equal(res, c.rhs) {
			t.Errorf("DeleteIfOnce failed, got: %v, want: %v.", res, c.rhs)
		}
	}
}
