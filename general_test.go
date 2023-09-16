package pimit

import (
	"errors"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestParallelIndicesShouldPanicOnInvalidWidth(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelIndices(0, 2, func(xIndex, yIndex int) {})
	})

	assert.Panics(t, func() {
		ParallelIndices(-2, 2, func(xIndex, yIndex int) {})
	})
}

func TestParallelIndicesShouldPanicOnInvalidHeight(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelIndices(2, 0, func(xIndex, yIndex int) {})
	})

	assert.Panics(t, func() {
		ParallelIndices(2, -2, func(xIndex, yIndex int) {})
	})
}

func TestParallelIndicesShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	width := 10
	height := 15

	ParallelIndices(width, height, func(xIndex, yIndex int) {
		assert.GreaterOrEqual(t, xIndex, 0)
		assert.Less(t, xIndex, width)

		assert.GreaterOrEqual(t, yIndex, 0)
		assert.Less(t, yIndex, height)
	})
}

func TestParallelReadShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelRead(nil, func(_, _ int, _ color.Color) {
		})
	})
}

func TestParallelReadShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImage()

	assert.Panics(t, func() {
		ParallelRead(img, nil)
	})
}

func TestParallelReadShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImage()

	exR, exG, exB, exA := color.White.RGBA()

	ParallelRead(img, func(xIndex, yIndex int, c color.Color) {
		assert.GreaterOrEqual(t, xIndex, 0)
		assert.Less(t, xIndex, img.Bounds().Dx())

		assert.GreaterOrEqual(t, yIndex, 0)
		assert.Less(t, yIndex, img.Bounds().Dy())

		acR, acG, acB, acA := c.RGBA()

		assert.Equal(t, exR, acR)
		assert.Equal(t, exG, acG)
		assert.Equal(t, exB, acB)
		assert.Equal(t, exA, acA)
	})
}

func TestParallelReadEShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelReadE(nil, func(_, _ int, _ color.Color) error {
			return nil
		})
	})
}

func TestParallelReadEShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImage()

	assert.Panics(t, func() {
		ParallelReadE(img, nil)
	})
}

func TestParallelReadEShouldReturnErrorOnAccessError(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImage()

	err := ParallelReadE(img, func(xIndex, yIndex int, c color.Color) error {
		return errors.New("pimit-test: test errror")
	})

	assert.NotNil(t, err)
}

func TestPrallelReadEShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImage()

	exR, exG, exB, exA := color.White.RGBA()

	err := ParallelReadE(img, func(xIndex, yIndex int, c color.Color) error {
		assert.GreaterOrEqual(t, xIndex, 0)
		assert.Less(t, xIndex, img.Bounds().Dx())

		assert.GreaterOrEqual(t, yIndex, 0)
		assert.Less(t, yIndex, img.Bounds().Dy())

		acR, acG, acB, acA := c.RGBA()

		assert.Equal(t, exR, acR)
		assert.Equal(t, exG, acG)
		assert.Equal(t, exB, acB)
		assert.Equal(t, exA, acA)
		return nil
	})

	assert.Nil(t, err)
}

func TestParallelReadWriteShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelReadWrite(nil, func(_, _ int, c color.Color) color.Color {
			return c
		})
	})
}

func TestParallelReadWriteShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteDrawableImage()

	assert.Panics(t, func() {
		ParallelReadWrite(img, nil)
	})
}

func TestPrallelReadWriteShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteDrawableImage()

	exR, exG, exB, exA := color.White.RGBA()

	ParallelReadWrite(img, func(xIndex, yIndex int, c color.Color) color.Color {
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

func TestParallelReadWriteShouldAccessPixelsOnce(t *testing.T) {
	defer goleak.VerifyNone(t)

	image := mockWhiteDrawableImage()

	rBlack, gBlack, bBlack, aBlack := color.Black.RGBA()
	rWhite, gWhite, bWhite, aWhite := color.White.RGBA()

	ParallelReadWrite(image, func(_, _ int, col color.Color) color.Color {
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

func TestParallelReadWriteEShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelReadWriteE(nil, func(_, _ int, c color.Color) (color.Color, error) {
			return c, nil
		})
	})
}

func TestParallelReadWriteEShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteDrawableImage()

	assert.Panics(t, func() {
		ParallelReadWriteE(img, nil)
	})
}

func TestParallelReadWriteEShouldReturnErrorOnAccessError(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteDrawableImage()

	err := ParallelReadWriteE(img, func(xIndex, yIndex int, c color.Color) (color.Color, error) {
		return c, errors.New("pimit-test: test errror")
	})

	assert.NotNil(t, err)
}

func TestPrallelReadWriteEShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteDrawableImage()

	exR, exG, exB, exA := color.White.RGBA()

	ParallelReadWriteE(img, func(xIndex, yIndex int, c color.Color) (color.Color, error) {
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

func TestParallelReadWriteEShouldAccessPixelsOnce(t *testing.T) {
	defer goleak.VerifyNone(t)

	image := mockWhiteDrawableImage()

	rBlack, gBlack, bBlack, aBlack := color.Black.RGBA()
	rWhite, gWhite, bWhite, aWhite := color.White.RGBA()

	err := ParallelReadWriteE(image, func(_, _ int, col color.Color) (color.Color, error) {
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

func TestParallelReadWriteNewShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelReadWriteNew(nil, func(_, _ int, c color.Color) color.Color {
			return c
		})
	})
}

func TestParallelReadWriteNewShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteDrawableImage()

	assert.Panics(t, func() {
		ParallelReadWriteNew(img, nil)
	})
}

func TestPrallelReadWriteNewShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteDrawableImage()

	exR, exG, exB, exA := color.White.RGBA()

	actualImage := ParallelReadWriteNew(img, func(xIndex, yIndex int, c color.Color) color.Color {
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

func TestParallelReadWriteNewShouldAccessPixelsOnce(t *testing.T) {
	defer goleak.VerifyNone(t)

	image := mockWhiteDrawableImage()

	rBlack, gBlack, bBlack, aBlack := color.Black.RGBA()
	rWhite, gWhite, bWhite, aWhite := color.White.RGBA()

	actualImage := ParallelReadWriteNew(image, func(_, _ int, col color.Color) color.Color {
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

func TestParallelReadWriteNewEShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelReadWriteNewE(nil, func(_, _ int, c color.Color) (color.Color, error) {
			return c, nil
		})
	})
}

func TestParallelReadWriteNewEShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteDrawableImage()

	assert.Panics(t, func() {
		ParallelReadWriteNewE(img, nil)
	})
}

func TestParallelReadWriteNewEShouldReturnErrorOnAccessError(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteDrawableImage()

	modifiedImg, err := ParallelReadWriteNewE(img, func(xIndex, yIndex int, c color.Color) (color.Color, error) {
		return c, errors.New("pimit-test: test errror")
	})

	assert.NotNil(t, err)
	assert.Nil(t, modifiedImg)
}

func TestPrallelReadWriteNewEShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteDrawableImage()

	exR, exG, exB, exA := color.White.RGBA()

	actualImage, err := ParallelReadWriteNewE(img, func(xIndex, yIndex int, c color.Color) (color.Color, error) {
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

	assert.Nil(t, err)

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

func TestParallelReadWriteENewShouldAccessPixelsOnce(t *testing.T) {
	defer goleak.VerifyNone(t)

	image := mockWhiteDrawableImage()

	rBlack, gBlack, bBlack, aBlack := color.Black.RGBA()
	rWhite, gWhite, bWhite, aWhite := color.White.RGBA()

	actualImage, err := ParallelReadWriteNewE(image, func(_, _ int, col color.Color) (color.Color, error) {
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
