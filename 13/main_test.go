package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSolve76(t *testing.T) {
	h, c := solve(76)
	assert.Equal(t, 3, h)
	assert.Equal(t, 76, c)
}

func TestSolve156(t *testing.T) {
	h, c := solve(156)
	assert.Equal(t, 4, h)
	assert.Equal(t, 156, c)
}

func TestF(t *testing.T) {
	assert.Equal(t, 76, f(3, 2, 3))
	assert.Equal(t, 156, f(1, 1, 4))
}
