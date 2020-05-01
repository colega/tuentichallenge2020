package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	lines := []*inputLine{
		&inputLine{s: "a"},
		&inputLine{s: "b"},
		&inputLine{s: "c"},
		&inputLine{s: "d"},
	}
	lines = join(lines, 2)
	lines = join(lines, 0)
	lines = join(lines, 0)

	assert.Equal(t, "abcd", lines[0].s)
}

func TestInput_lolmaoLiteralToLolmaoList(t *testing.T) {
	t.Run("easy", func(t *testing.T) {
		inp := input{&inputLine{s: "easy"}}
		expected := "[as]"

		got, ok := inp.lolmaoLiteralToLolmaoList()
		require.True(t, ok)
		assert.Equal(t, expected, got.String())
	})
	t.Run("two chars", func(t *testing.T) {
		inp := input{&inputLine{s: "xx"}}
		expected := "[]"

		got, ok := inp.lolmaoLiteralToLolmaoList()
		require.True(t, ok)
		assert.Equal(t, expected, got.String())
	})

	t.Run("easy two lines", func(t *testing.T) {
		inp := input{
			&inputLine{s: "easy"},
			&inputLine{s: "easy"},
		}
		expected := "[asy\neas]"

		got, ok := inp.lolmaoLiteralToLolmaoList()
		require.True(t, ok)
		assert.Equal(t, expected, got.String())
	})

	t.Run("first line empty", func(t *testing.T) {
		inp := input{
			&inputLine{s: ""},
			&inputLine{s: "two"},
		}
		expected := "[tw]"

		got, ok := inp.lolmaoLiteralToLolmaoList()
		require.True(t, ok)
		assert.Equal(t, expected, got.String())
	})

	t.Run("last line empty", func(t *testing.T) {
		inp := input{
			&inputLine{s: "one"},
			&inputLine{s: "two"},
			&inputLine{s: ""},
		}
		expected := "[ne\ntwo]"

		got, ok := inp.lolmaoLiteralToLolmaoList()
		require.True(t, ok)
		assert.Equal(t, expected, got.String())
	})

	t.Run("one char on first line", func(t *testing.T) {
		inp := input{
			&inputLine{s: "x"},
			&inputLine{s: ""},
		}
		expected := "[]"

		got, ok := inp.lolmaoLiteralToLolmaoList()
		require.True(t, ok)
		assert.Equal(t, expected, got.String())
	})

	t.Run("one char on last line", func(t *testing.T) {
		inp := input{
			&inputLine{s: ""},
			&inputLine{s: "x"},
		}
		expected := "[]"

		got, ok := inp.lolmaoLiteralToLolmaoList()
		require.True(t, ok)
		assert.Equal(t, expected, got.String())
	})

	t.Run("three empty lines", func(t *testing.T) {
		inp := input{
			&inputLine{s: ""},
			&inputLine{s: ""},
			&inputLine{s: ""},
		}
		expected := "[]"

		got, ok := inp.lolmaoLiteralToLolmaoList()
		require.True(t, ok)
		assert.Equal(t, expected, got.String())
	})

	t.Run("two empty lines", func(t *testing.T) {
		inp := input{
			&inputLine{s: ""},
			&inputLine{s: ""},
		}

		_, ok := inp.lolmaoLiteralToLolmaoList()
		require.False(t, ok)
	})
	t.Run("single char", func(t *testing.T) {
		inp := input{
			&inputLine{s: "x"},
		}

		_, ok := inp.lolmaoLiteralToLolmaoList()
		require.False(t, ok)
	})

}

func TestBracketHereWouldBeUnclosed1(t *testing.T) {
	inp := []byte("[5[b3,]][]")
	assert.False(t, bracketHereWouldBeUnclosed(inp, 1, 1))
}

func TestBracketHereWouldBeUnclosed2(t *testing.T) {
	inp := []byte("[,,,[N,][]")
	assert.False(t, bracketHereWouldBeUnclosed(inp, 1, 4))
}
