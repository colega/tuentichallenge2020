package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSolveSample1(t *testing.T) {
	got := solve("aardwolf0000", []add{
		{0, 224, 1},
		{0, 192, 2},
	})

	expected := []string{
		"aardwolf0000 0: 00000000",
		"aardwolf0000 1: 72080df5",
		"aardwolf0000 2: 2a2927c9",
	}

	assert.Equal(t, expected, got)
}
func TestSolveSample2(t *testing.T) {
	got := solve("admiral0000", []add{
		{1, 227, 1},
		{2, 232, 2},
		{2, 46, 3},
		{0, 169, 4},
	})

	expected := []string{
		"admiral0000 0: d202ef8d",
		"admiral0000 1: 78daa13d",
		"admiral0000 2: 24c31377",
		"admiral0000 3: 2f36d283",
		"admiral0000 4: b5670765",
	}

	assert.Equal(t, expected, got)
}

func TestAntelope(t *testing.T) {
	got := solve("antelope0000", []add{
		{336237519, 75, 1},
		{3136973822, 193, 2},
	})

	expected := []string{
		"antelope0000 0: 3c9e51d9",
		"antelope0000 1: 14dfc526",
		"antelope0000 2: 23fff51c",
	}

	assert.Equal(t, expected, got)
}
