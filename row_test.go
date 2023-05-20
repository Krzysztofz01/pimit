package pimit

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParallelRowReadShouldPanicOnNilImage(t *testing.T) {
	assert.Panics(t, func() {
		ParallelRowRead(nil, func(_, _ int, _ color.Color) {
		})
	})
}

func TestParallelRowReadShouldPanicOnNilAccessFunc(t *testing.T) {
	img := mockWhiteImage()

	assert.Panics(t, func() {
		ParallelRowRead(img, nil)
	})
}

func TestPrallelRowReadShouldCorrectlyIterate(t *testing.T) {
	img := mockWhiteImage()

	exR, exG, exB, exA := color.White.RGBA()

	ParallelRowRead(img, func(xIndex, yIndex int, c color.Color) {
		assert.GreaterOrEqual(t, xIndex, 0)
		assert.GreaterOrEqual(t, yIndex, 0)

		acR, acG, acB, acA := c.RGBA()

		assert.Equal(t, exR, acR)
		assert.Equal(t, exG, acG)
		assert.Equal(t, exB, acB)
		assert.Equal(t, exA, acA)
	})
}

func TestParallelRowColorReadShouldPanicOnNilImage(t *testing.T) {
	assert.Panics(t, func() {
		ParallelRowColorRead(nil, func(_ color.Color) {
		})
	})
}

func TestParallelRowColorReadShouldPanicOnNilAccessFunc(t *testing.T) {
	img := mockWhiteImage()

	assert.Panics(t, func() {
		ParallelRowColorRead(img, nil)
	})
}

func TestPrallelRowColorReadShouldCorrectlyIterate(t *testing.T) {
	img := mockWhiteImage()

	exR, exG, exB, exA := color.White.RGBA()

	ParallelRowColorRead(img, func(c color.Color) {
		acR, acG, acB, acA := c.RGBA()

		assert.Equal(t, exR, acR)
		assert.Equal(t, exG, acG)
		assert.Equal(t, exB, acB)
		assert.Equal(t, exA, acA)
	})
}

func TestParallelRowReadWriteShouldPanicOnNilImage(t *testing.T) {
	assert.Panics(t, func() {
		ParallelRowReadWrite(nil, func(_, _ int, c color.Color) color.Color {
			return c
		})
	})
}

func TestParallelRowReadWriteShouldPanicOnNilAccessFunc(t *testing.T) {
	img := mockWhiteDrawableImage()

	assert.Panics(t, func() {
		ParallelRowReadWrite(img, nil)
	})
}

func TestPrallelRowReadWriteShouldCorrectlyIterate(t *testing.T) {
	img := mockWhiteDrawableImage()

	exR, exG, exB, exA := color.White.RGBA()

	ParallelRowReadWrite(img, func(xIndex, yIndex int, c color.Color) color.Color {
		assert.GreaterOrEqual(t, xIndex, 0)
		assert.GreaterOrEqual(t, yIndex, 0)
		acR, acG, acB, acA := c.RGBA()

		assert.Equal(t, exR, acR)
		assert.Equal(t, exG, acG)
		assert.Equal(t, exB, acB)
		assert.Equal(t, exA, acA)

		return color.Black
	})

	expectedImage := mockBlackDrawableImage()

	assert.Equal(t, expectedImage.Bounds(), img.Bounds())

	for x := 0; x < expectedImage.Bounds().Dx(); x += 1 {
		for y := 0; y < expectedImage.Bounds().Dy(); y += 1 {
			assert.Equal(t, expectedImage.At(x, y), img.At(x, y))
		}
	}
}

func TestParallelRowColorReadWriteShouldPanicOnNilImage(t *testing.T) {
	assert.Panics(t, func() {
		ParallelRowColorReadWrite(nil, func(c color.Color) color.Color {
			return c
		})
	})
}

func TestParallelRowColorReadWriteShouldPanicOnNilAccessFunc(t *testing.T) {
	img := mockWhiteDrawableImage()

	assert.Panics(t, func() {
		ParallelRowColorReadWrite(img, nil)
	})
}

func TestPrallelRowColorReadWriteShouldCorrectlyIterate(t *testing.T) {
	img := mockWhiteDrawableImage()

	exR, exG, exB, exA := color.White.RGBA()

	ParallelRowColorReadWrite(img, func(c color.Color) color.Color {
		acR, acG, acB, acA := c.RGBA()

		assert.Equal(t, exR, acR)
		assert.Equal(t, exG, acG)
		assert.Equal(t, exB, acB)
		assert.Equal(t, exA, acA)

		return color.Black
	})

	expectedImage := mockBlackDrawableImage()

	assert.Equal(t, expectedImage.Bounds(), img.Bounds())

	for x := 0; x < expectedImage.Bounds().Dx(); x += 1 {
		for y := 0; y < expectedImage.Bounds().Dy(); y += 1 {
			assert.Equal(t, expectedImage.At(x, y), img.At(x, y))
		}
	}
}

func TestParallelRowReadWriteNewShouldPanicOnNilImage(t *testing.T) {
	assert.Panics(t, func() {
		ParallelRowReadWriteNew(nil, func(_, _ int, c color.Color) color.Color {
			return c
		})
	})
}

func TestParallelRowReadWriteNewShouldPanicOnNilAccessFunc(t *testing.T) {
	img := mockWhiteDrawableImage()

	assert.Panics(t, func() {
		ParallelRowReadWriteNew(img, nil)
	})
}

func TestPrallelRowReadWriteNewShouldCorrectlyIterate(t *testing.T) {
	img := mockWhiteDrawableImage()

	exR, exG, exB, exA := color.White.RGBA()

	actualImage := ParallelRowReadWriteNew(img, func(xIndex, yIndex int, c color.Color) color.Color {
		assert.GreaterOrEqual(t, xIndex, 0)
		assert.GreaterOrEqual(t, yIndex, 0)
		acR, acG, acB, acA := c.RGBA()

		assert.Equal(t, exR, acR)
		assert.Equal(t, exG, acG)
		assert.Equal(t, exB, acB)
		assert.Equal(t, exA, acA)

		return color.Black
	})

	expectedImage := mockBlackDrawableImage()

	assert.Equal(t, expectedImage.Bounds(), actualImage.Bounds())

	for x := 0; x < expectedImage.Bounds().Dx(); x += 1 {
		for y := 0; y < expectedImage.Bounds().Dy(); y += 1 {
			exR, exG, exB, exA := expectedImage.At(x, y).RGBA()
			acR, acG, acB, acA := actualImage.At(x, y).RGBA()

			assert.Equal(t, exR, acR)
			assert.Equal(t, exG, acG)
			assert.Equal(t, exB, acB)
			assert.Equal(t, exA, acA)
		}
	}
}

func TestParallelRowColorReadWriteNewShouldPanicOnNilImage(t *testing.T) {
	assert.Panics(t, func() {
		ParallelRowColorReadWriteNew(nil, func(c color.Color) color.Color {
			return c
		})
	})
}

func TestParallelRowColorReadWriteNewShouldPanicOnNilAccessFunc(t *testing.T) {
	img := mockWhiteDrawableImage()

	assert.Panics(t, func() {
		ParallelRowColorReadWriteNew(img, nil)
	})
}

func TestPrallelRowColorReadWriteNewShouldCorrectlyIterate(t *testing.T) {
	img := mockWhiteDrawableImage()

	exR, exG, exB, exA := color.White.RGBA()

	actualImage := ParallelRowColorReadWriteNew(img, func(c color.Color) color.Color {
		acR, acG, acB, acA := c.RGBA()

		assert.Equal(t, exR, acR)
		assert.Equal(t, exG, acG)
		assert.Equal(t, exB, acB)
		assert.Equal(t, exA, acA)

		return color.Black
	})

	expectedImage := mockBlackDrawableImage()

	assert.Equal(t, expectedImage.Bounds(), actualImage.Bounds())

	for x := 0; x < expectedImage.Bounds().Dx(); x += 1 {
		for y := 0; y < expectedImage.Bounds().Dy(); y += 1 {
			exR, exG, exB, exA := expectedImage.At(x, y).RGBA()
			acR, acG, acB, acA := actualImage.At(x, y).RGBA()

			assert.Equal(t, exR, acR)
			assert.Equal(t, exG, acG)
			assert.Equal(t, exB, acB)
			assert.Equal(t, exA, acA)
		}
	}
}
