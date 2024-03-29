package pimit

import (
	"errors"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestOffsetToIndexShouldPanicOnNegativeOffset(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		offsetToIndex(-2, 4)
	})
}

func TestOffsetToIndexShouldPanicOnInvalidWidth(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		offsetToIndex(0, 0)
	})

	assert.Panics(t, func() {
		offsetToIndex(0, -2)
	})
}

func TestOffsetToIndexShouldCalculate(t *testing.T) {
	defer goleak.VerifyNone(t)

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

func TestParallelClusterDistributedReadWriteShouldAccessPixelsOnce(t *testing.T) {
	defer goleak.VerifyNone(t)

	cases := []struct {
		width    int
		height   int
		clusters int
	}{
		{2, 2, 2},
		{2, 2, 2},
		{2, 3, 2},
		{3, 2, 2},
		{3, 3, 2},
		{2, 2, 3},
		{2, 2, 3},
		{2, 3, 3},
		{3, 2, 3},
		{3, 3, 3},
	}

	rBlack, gBlack, bBlack, aBlack := color.Black.RGBA()
	rWhite, gWhite, bWhite, aWhite := color.White.RGBA()

	for _, c := range cases {
		image := mockSpecificDrawableImage(c.width, c.height, color.White)

		ParallelClusterDistributedReadWrite(image, c.clusters, func(xIndex, yIndex int, col color.Color) color.Color {
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

		for x := 0; x < c.width; x += 1 {
			for y := 0; y < c.height; y += 1 {
				acR, acG, acB, acA := image.At(x, y).RGBA()

				assert.Equal(t, rBlack, acR)
				assert.Equal(t, gBlack, acG)
				assert.Equal(t, bBlack, acB)
				assert.Equal(t, aBlack, acA)
			}
		}
	}
}

func TestParallelClusterDistributedReadWriteShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelClusterDistributedReadWrite(nil, 2, func(_, _ int, c color.Color) color.Color {
			return c
		})
	})
}

func TestParallelClusterDistributedReadWriteShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteDrawableImage()

	assert.Panics(t, func() {
		ParallelClusterDistributedReadWrite(img, 2, nil)
	})
}

func TestParallelClusterDistributedReadWriteShouldPanicOnInvalidClusterCount(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteDrawableImage()

	assert.Panics(t, func() {
		ParallelClusterDistributedReadWrite(img, 0, func(_, _ int, c color.Color) color.Color {
			return c
		})
	})

	assert.Panics(t, func() {
		ParallelClusterDistributedReadWrite(img, -2, func(_, _ int, c color.Color) color.Color {
			return c
		})
	})
}

func TestParallelClusterDistributedReadWriteShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteDrawableImage()

	exR, exG, exB, exA := color.White.RGBA()

	ParallelClusterDistributedReadWrite(img, 2, func(xIndex, yIndex int, c color.Color) color.Color {
		assert.GreaterOrEqual(t, xIndex, 0)
		assert.Less(t, xIndex, img.Bounds().Dx())

		assert.GreaterOrEqual(t, yIndex, 0)
		assert.Less(t, yIndex, img.Bounds().Dy())

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

func TestParallelClusterDistributedReadWriteEShouldAccessPixelsOnce(t *testing.T) {
	defer goleak.VerifyNone(t)

	cases := []struct {
		width    int
		height   int
		clusters int
	}{
		{2, 2, 2},
		{2, 2, 2},
		{2, 3, 2},
		{3, 2, 2},
		{3, 3, 2},
		{2, 2, 3},
		{2, 2, 3},
		{2, 3, 3},
		{3, 2, 3},
		{3, 3, 3},
	}

	rBlack, gBlack, bBlack, aBlack := color.Black.RGBA()
	rWhite, gWhite, bWhite, aWhite := color.White.RGBA()

	for _, c := range cases {
		image := mockSpecificDrawableImage(c.width, c.height, color.White)

		err := ParallelClusterDistributedReadWriteE(image, c.clusters, func(xIndex, yIndex int, col color.Color) (color.Color, error) {
			rCurrent, gCurrent, bCurrent, aCurrent := col.RGBA()

			if rCurrent == rBlack && gCurrent == gBlack && bCurrent == bBlack && aCurrent == aBlack {
				return color.White, nil
			}

			if rCurrent == rWhite && gCurrent == gWhite && bCurrent == bWhite && aCurrent == aWhite {
				return color.Black, nil
			}

			assert.FailNow(t, "This should never happen")
			return col, nil
		})

		assert.Nil(t, err)

		for x := 0; x < c.width; x += 1 {
			for y := 0; y < c.height; y += 1 {
				acR, acG, acB, acA := image.At(x, y).RGBA()

				assert.Equal(t, rBlack, acR)
				assert.Equal(t, gBlack, acG)
				assert.Equal(t, bBlack, acB)
				assert.Equal(t, aBlack, acA)
			}
		}
	}
}

func TestParallelClusterDistributedReadWriteEShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelClusterDistributedReadWriteE(nil, 2, func(_, _ int, c color.Color) (color.Color, error) {
			return c, nil
		})
	})
}

func TestParallelClusterDistributedReadWriteEShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteDrawableImage()

	assert.Panics(t, func() {
		ParallelClusterDistributedReadWriteE(img, 2, nil)
	})
}

func TestParallelClusterDistributedReadWriteEShouldPanicOnInvalidClusterCount(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteDrawableImage()

	assert.Panics(t, func() {
		ParallelClusterDistributedReadWriteE(img, 0, func(_, _ int, c color.Color) (color.Color, error) {
			return c, nil
		})
	})

	assert.Panics(t, func() {
		ParallelClusterDistributedReadWriteE(img, -2, func(_, _ int, c color.Color) (color.Color, error) {
			return c, nil
		})
	})
}

func TestParallelClusterDistributedReadWriteEShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteDrawableImage()

	exR, exG, exB, exA := color.White.RGBA()

	ParallelClusterDistributedReadWriteE(img, 2, func(xIndex, yIndex int, c color.Color) (color.Color, error) {
		assert.GreaterOrEqual(t, xIndex, 0)
		assert.Less(t, xIndex, img.Bounds().Dx())

		assert.GreaterOrEqual(t, yIndex, 0)
		assert.Less(t, yIndex, img.Bounds().Dy())

		acR, acG, acB, acA := c.RGBA()

		assert.Equal(t, exR, acR)
		assert.Equal(t, exG, acG)
		assert.Equal(t, exB, acB)
		assert.Equal(t, exA, acA)

		return color.Black, nil
	})

	expectedImage := mockBlackDrawableImage()

	assert.Equal(t, expectedImage.Bounds(), img.Bounds())

	for x := 0; x < expectedImage.Bounds().Dx(); x += 1 {
		for y := 0; y < expectedImage.Bounds().Dy(); y += 1 {
			assert.Equal(t, expectedImage.At(x, y), img.At(x, y))
		}
	}
}

func TestParallelClusterDistributedReadWriteEShouldReturnErrorOnAccessError(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteDrawableImage()

	err := ParallelClusterDistributedReadWriteE(img, 2, func(xIndex, yIndex int, c color.Color) (color.Color, error) {
		return c, errors.New("pimit-test: test errror")
	})

	assert.NotNil(t, err)
}
