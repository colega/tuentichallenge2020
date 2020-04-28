package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSolve(t *testing.T) {
	for _, tc := range []struct {
		numbers  []int
		expected int
	}{
		{[]int{2, 1}, 0},
		{[]int{5, 4, 2}, 2},

		{[]int{2}, 1},
		{[]int{7, 1, 3, 6}, 1},
	} {
		t.Run(fmt.Sprintf("%v", tc.numbers), func(t *testing.T) {
			f := solve(tc.numbers)
			assert.Equal(t, tc.expected, f)
		})
	}
}
