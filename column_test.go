package pimit

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParallelColumnReadShouldPanicOnNilImage(t *testing.T) {
	assert.Panics(t, func() {
		ParallelColumnRead(nil, func(_, _ int, _ color.Color) {
		})
	})
}

func TestParallelColumnReadShouldPanicOnNilAccessFunc(t *testing.T) {
	img := mockWhiteImage()

	assert.Panics(t, func() {
		ParallelColumnRead(img, nil)
	})
}

func TestPrallelColumnReadShouldCorrectlyIterate(t *testing.T) {
	img := mockWhiteImage()

	exR, exG, exB, exA := color.White.RGBA()

	ParallelColumnRead(img, func(xIndex, yIndex int, c color.Color) {
		assert.GreaterOrEqual(t, xIndex, 0)
		assert.GreaterOrEqual(t, yIndex, 0)

		acR, acG, acB, acA := c.RGBA()

		assert.Equal(t, exR, acR)
		assert.Equal(t, exG, acG)
		assert.Equal(t, exB, acB)
		assert.Equal(t, exA, acA)
	})
}

func TestParallelColumnColorReadShouldPanicOnNilImage(t *testing.T) {
	assert.Panics(t, func() {
		ParallelColumnColorRead(nil, func(_ color.Color) {
		})
	})
}

func TestParallelColumnColorReadShouldPanicOnNilAccessFunc(t *testing.T) {
	img := mockWhiteImage()

	assert.Panics(t, func() {
		ParallelColumnColorRead(img, nil)
	})
}

func TestPrallelColumnColorReadShouldCorrectlyIterate(t *testing.T) {
	img := mockWhiteImage()

	exR, exG, exB, exA := color.White.RGBA()

	ParallelColumnColorRead(img, func(c color.Color) {
		acR, acG, acB, acA := c.RGBA()

		assert.Equal(t, exR, acR)
		assert.Equal(t, exG, acG)
		assert.Equal(t, exB, acB)
		assert.Equal(t, exA, acA)
	})
}

func TestParallelColumnReadWriteShouldPanicOnNilImage(t *testing.T) {
	assert.Panics(t, func() {
		ParallelColumnReadWrite(nil, func(_, _ int, c color.Color) color.Color {
			return c
		})
	})
}

func TestParallelColumnReadWriteShouldPanicOnNilAccessFunc(t *testing.T) {
	img := mockWhiteDrawableImage()

	assert.Panics(t, func() {
		ParallelColumnReadWrite(img, nil)
	})
}

func TestPrallelColumnReadWriteShouldCorrectlyIterate(t *testing.T) {
	img := mockWhiteDrawableImage()

	exR, exG, exB, exA := color.White.RGBA()

	ParallelColumnReadWrite(img, func(xIndex, yIndex int, c color.Color) color.Color {
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

func TestParallelColumnReadWriteShouldAccessPixelsOnce(t *testing.T) {
	image := mockWhiteDrawableImage()

	rBlack, gBlack, bBlack, aBlack := color.Black.RGBA()
	rWhite, gWhite, bWhite, aWhite := color.White.RGBA()

	ParallelColumnReadWrite(image, func(_, _ int, col color.Color) color.Color {
		rCurrent, gCurrent, bCurrent, aCurrent := col.RGBA()

		if rCurrent == rBlack && gCurrent == gBlack && bCurrent == bBlack && aCurrent == aBlack {
			return color.White
		}

		if rCurrent == rWhite && gCurrent == gWhite && bCurrent == bWhite && aCurrent == aWhite {
			return color.Black
		}

		assert.FailNow(t, "This should never happen")
		return col
	})

	for x := 0; x < image.Bounds().Dx(); x += 1 {
		for y := 0; y < image.Bounds().Dy(); y += 1 {
			acR, acG, acB, acA := image.At(x, y).RGBA()

			assert.Equal(t, rBlack, acR)
			assert.Equal(t, gBlack, acG)
			assert.Equal(t, bBlack, acB)
			assert.Equal(t, aBlack, acA)
		}
	}
}

func TestParallelColumnColorReadWriteShouldPanicOnNilImage(t *testing.T) {
	assert.Panics(t, func() {
		ParallelColumnColorReadWrite(nil, func(c color.Color) color.Color {
			return c
		})
	})
}

func TestParallelColumnColorReadWriteShouldPanicOnNilAccessFunc(t *testing.T) {
	img := mockWhiteDrawableImage()

	assert.Panics(t, func() {
		ParallelColumnColorReadWrite(img, nil)
	})
}

func TestPrallelColumnColorReadWriteShouldCorrectlyIterate(t *testing.T) {
	img := mockWhiteDrawableImage()

	exR, exG, exB, exA := color.White.RGBA()

	ParallelColumnColorReadWrite(img, func(c color.Color) color.Color {
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

func TestParallelColumnColorReadWriteShouldAccessPixelsOnce(t *testing.T) {
	image := mockWhiteDrawableImage()

	rBlack, gBlack, bBlack, aBlack := color.Black.RGBA()
	rWhite, gWhite, bWhite, aWhite := color.White.RGBA()

	ParallelColumnColorReadWrite(image, func(col color.Color) color.Color {
		rCurrent, gCurrent, bCurrent, aCurrent := col.RGBA()

		if rCurrent == rBlack && gCurrent == gBlack && bCurrent == bBlack && aCurrent == aBlack {
			return color.White
		}

		if rCurrent == rWhite && gCurrent == gWhite && bCurrent == bWhite && aCurrent == aWhite {
			return color.Black
		}

		assert.FailNow(t, "This should never happen")
		return col
	})

	for x := 0; x < image.Bounds().Dx(); x += 1 {
		for y := 0; y < image.Bounds().Dy(); y += 1 {
			acR, acG, acB, acA := image.At(x, y).RGBA()

			assert.Equal(t, rBlack, acR)
			assert.Equal(t, gBlack, acG)
			assert.Equal(t, bBlack, acB)
			assert.Equal(t, aBlack, acA)
		}
	}
}

func TestParallelColumnReadWriteNewShouldPanicOnNilImage(t *testing.T) {
	assert.Panics(t, func() {
		ParallelColumnReadWriteNew(nil, func(_, _ int, c color.Color) color.Color {
			return c
		})
	})
}

func TestParallelColumnReadWriteNewShouldPanicOnNilAccessFunc(t *testing.T) {
	img := mockWhiteDrawableImage()

	assert.Panics(t, func() {
		ParallelColumnReadWriteNew(img, nil)
	})
}

func TestPrallelColumnReadWriteNewShouldCorrectlyIterate(t *testing.T) {
	img := mockWhiteDrawableImage()

	exR, exG, exB, exA := color.White.RGBA()

	actualImage := ParallelColumnReadWriteNew(img, func(xIndex, yIndex int, c color.Color) color.Color {
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

func TestParallelColumnReadWriteNewShouldAccessPixelsOnce(t *testing.T) {
	image := mockWhiteDrawableImage()

	rBlack, gBlack, bBlack, aBlack := color.Black.RGBA()
	rWhite, gWhite, bWhite, aWhite := color.White.RGBA()

	actualImage := ParallelColumnReadWriteNew(image, func(_, _ int, col color.Color) color.Color {
		rCurrent, gCurrent, bCurrent, aCurrent := col.RGBA()

		if rCurrent == rBlack && gCurrent == gBlack && bCurrent == bBlack && aCurrent == aBlack {
			return color.White
		}

		if rCurrent == rWhite && gCurrent == gWhite && bCurrent == bWhite && aCurrent == aWhite {
			return color.Black
		}

		assert.FailNow(t, "This should never happen")
		return col
	})

	for x := 0; x < image.Bounds().Dx(); x += 1 {
		for y := 0; y < image.Bounds().Dy(); y += 1 {
			acR, acG, acB, acA := actualImage.At(x, y).RGBA()

			assert.Equal(t, rBlack, acR)
			assert.Equal(t, gBlack, acG)
			assert.Equal(t, bBlack, acB)
			assert.Equal(t, aBlack, acA)
		}
	}
}

func TestParallelColumnColorReadWriteNewShouldPanicOnNilImage(t *testing.T) {
	assert.Panics(t, func() {
		ParallelColumnColorReadWriteNew(nil, func(c color.Color) color.Color {
			return c
		})
	})
}

func TestParallelColumnColorReadWriteNewShouldPanicOnNilAccessFunc(t *testing.T) {
	img := mockWhiteDrawableImage()

	assert.Panics(t, func() {
		ParallelColumnColorReadWriteNew(img, nil)
	})
}

func TestPrallelColumnColorReadWriteNewShouldCorrectlyIterate(t *testing.T) {
	img := mockWhiteDrawableImage()

	exR, exG, exB, exA := color.White.RGBA()

	actualImage := ParallelColumnColorReadWriteNew(img, func(c color.Color) color.Color {
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

func TestParallelColumnColorReadWriteNewShouldAccessPixelsOnce(t *testing.T) {
	image := mockWhiteDrawableImage()

	rBlack, gBlack, bBlack, aBlack := color.Black.RGBA()
	rWhite, gWhite, bWhite, aWhite := color.White.RGBA()

	actualImage := ParallelColumnColorReadWriteNew(image, func(col color.Color) color.Color {
		rCurrent, gCurrent, bCurrent, aCurrent := col.RGBA()

		if rCurrent == rBlack && gCurrent == gBlack && bCurrent == bBlack && aCurrent == aBlack {
			return color.White
		}

		if rCurrent == rWhite && gCurrent == gWhite && bCurrent == bWhite && aCurrent == aWhite {
			return color.Black
		}

		assert.FailNow(t, "This should never happen")
		return col
	})

	for x := 0; x < image.Bounds().Dx(); x += 1 {
		for y := 0; y < image.Bounds().Dy(); y += 1 {
			acR, acG, acB, acA := actualImage.At(x, y).RGBA()

			assert.Equal(t, rBlack, acR)
			assert.Equal(t, gBlack, acG)
			assert.Equal(t, bBlack, acB)
			assert.Equal(t, aBlack, acA)
		}
	}
}
