package pimit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOffsetToIndexShouldPanicOnNegativeOffset(t *testing.T) {
	assert.Panics(t, func() {
		offsetToIndex(-2, 4)
	})
}

func TestOffsetToIndexShouldPanicOnInvalidWidth(t *testing.T) {
	assert.Panics(t, func() {
		offsetToIndex(0, 0)
	})

	assert.Panics(t, func() {
		offsetToIndex(0, -2)
	})
}

func TestOffsetToIndexShouldCalculate(t *testing.T) {
	cases := map[struct {
		offset int
		width  int
	}]struct {
		xIndex int
		yIndex int
	}{
		{0, 3}: {0, 0},
		{1, 3}: {1, 0},
		{2, 3}: {2, 0},
		{3, 3}: {0, 1},
		{4, 3}: {1, 1},
		{0, 4}: {0, 0},
		{1, 4}: {1, 0},
		{2, 4}: {2, 0},
		{3, 4}: {3, 0},
		{4, 4}: {0, 1},
	}

	for input, expected := range cases {
		xIndexActual, yIndexActual := offsetToIndex(input.offset, input.width)

		assert.Equal(t, expected.xIndex, xIndexActual)
		assert.Equal(t, expected.yIndex, yIndexActual)
	}
}
