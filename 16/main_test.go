package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSolveSample1(t *testing.T) {
	sol := solve(2, []*group{
		&group{employees: 1, floors: []int{0, 1}},
		&group{employees: 1, floors: []int{1}},
	})

	assert.Equal(t, 1, sol)
}

func TestSolveSample2(t *testing.T) {
	sol := solve(1, []*group{
		&group{employees: 3, floors: []int{0}},
		&group{employees: 2, floors: []int{0}},
	})

	assert.Equal(t, 5, sol)
}

func TestSolveSample3(t *testing.T) {
	sol := solve(10, []*group{
		&group{employees: 2, floors: []int{4}},
		&group{employees: 3, floors: []int{8, 4, 5, 2, 0, 7, 6}},
	})

	assert.Equal(t, 2, sol)
}

func makeInts(n int) []int {
	var ints []int
	for i := 0; i < n; i++ {
		ints = append(ints, i)
	}
	return ints
}
