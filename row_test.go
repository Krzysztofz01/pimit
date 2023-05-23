package pimit

import (
	"errors"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParallelRowIndicesShouldPanicOnInvalidWidth(t *testing.T) {
	assert.Panics(t, func() {
		ParallelRowIndices(0, 2, func(xIndex, yIndex int) {})
	})

	assert.Panics(t, func() {
		ParallelRowIndices(-2, 2, func(xIndex, yIndex int) {})
	})
}

func TestParallelRowIndicesShouldPanicOnInvalidHeight(t *testing.T) {
	assert.Panics(t, func() {
		ParallelRowIndices(2, 0, func(xIndex, yIndex int) {})
	})

	assert.Panics(t, func() {
		ParallelRowIndices(2, -2, func(xIndex, yIndex int) {})
	})
}

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

func TestParallelRowReadEShouldPanicOnNilImage(t *testing.T) {
	assert.Panics(t, func() {
		ParallelRowReadE(nil, func(_, _ int, _ color.Color) error {
			return nil
		})
	})
}

func TestParallelRowReadEShouldPanicOnNilAccessFunc(t *testing.T) {
	img := mockWhiteImage()

	assert.Panics(t, func() {
		ParallelRowReadE(img, nil)
	})
}

func TestParallelRowReadEShouldReturnErrorOnAccessError(t *testing.T) {
	img := mockWhiteImage()

	err := ParallelRowReadE(img, func(xIndex, yIndex int, c color.Color) error {
		return errors.New("pimit-test: test errror")
	})

	assert.NotNil(t, err)
}

func TestPrallelRowReadEShouldCorrectlyIterate(t *testing.T) {
	img := mockWhiteImage()

	exR, exG, exB, exA := color.White.RGBA()

	err := ParallelRowReadE(img, func(xIndex, yIndex int, c color.Color) error {
		assert.GreaterOrEqual(t, xIndex, 0)
		assert.GreaterOrEqual(t, yIndex, 0)

		acR, acG, acB, acA := c.RGBA()

		assert.Equal(t, exR, acR)
		assert.Equal(t, exG, acG)
		assert.Equal(t, exB, acB)
		assert.Equal(t, exA, acA)
		return nil
	})

	assert.Nil(t, err)
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

func TestParallelRowColorReadEShouldPanicOnNilImage(t *testing.T) {
	assert.Panics(t, func() {
		ParallelRowColorReadE(nil, func(_ color.Color) error {
			return nil
		})
	})
}

func TestParallelRowColorReadEShouldPanicOnNilAccessFunc(t *testing.T) {
	img := mockWhiteImage()

	assert.Panics(t, func() {
		ParallelRowColorReadE(img, nil)
	})
}

func TestParallelRowColorReadEShouldReturnErrorOnAccessError(t *testing.T) {
	img := mockWhiteImage()

	err := ParallelRowColorReadE(img, func(c color.Color) error {
		return errors.New("pimit-test: test errror")
	})

	assert.NotNil(t, err)
}

func TestPrallelRowColorReadEShouldCorrectlyIterate(t *testing.T) {
	img := mockWhiteImage()

	exR, exG, exB, exA := color.White.RGBA()

	err := ParallelRowColorReadE(img, func(c color.Color) error {
		acR, acG, acB, acA := c.RGBA()

		assert.Equal(t, exR, acR)
		assert.Equal(t, exG, acG)
		assert.Equal(t, exB, acB)
		assert.Equal(t, exA, acA)
		return nil
	})

	assert.Nil(t, err)
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

func TestParallelRowReadWriteShouldAccessPixelsOnce(t *testing.T) {
	image := mockWhiteDrawableImage()

	rBlack, gBlack, bBlack, aBlack := color.Black.RGBA()
	rWhite, gWhite, bWhite, aWhite := color.White.RGBA()

	ParallelRowReadWrite(image, func(_, _ int, col color.Color) color.Color {
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

func TestParallelRowReadWriteEShouldPanicOnNilImage(t *testing.T) {
	assert.Panics(t, func() {
		ParallelRowReadWriteE(nil, func(_, _ int, c color.Color) (color.Color, error) {
			return c, nil
		})
	})
}

func TestParallelRowReadWriteEShouldPanicOnNilAccessFunc(t *testing.T) {
	img := mockWhiteDrawableImage()

	assert.Panics(t, func() {
		ParallelRowReadWriteE(img, nil)
	})
}

func TestParallelRowReadWriteEShouldReturnErrorOnAccessError(t *testing.T) {
	img := mockWhiteDrawableImage()

	err := ParallelRowReadWriteE(img, func(xIndex, yIndex int, c color.Color) (color.Color, error) {
		return c, errors.New("pimit-test: test errror")
	})

	assert.NotNil(t, err)
}

func TestPrallelRowReadWriteEShouldCorrectlyIterate(t *testing.T) {
	img := mockWhiteDrawableImage()

	exR, exG, exB, exA := color.White.RGBA()

	ParallelRowReadWriteE(img, func(xIndex, yIndex int, c color.Color) (color.Color, error) {
		assert.GreaterOrEqual(t, xIndex, 0)
		assert.GreaterOrEqual(t, yIndex, 0)
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

func TestParallelRowReadWriteEShouldAccessPixelsOnce(t *testing.T) {
	image := mockWhiteDrawableImage()

	rBlack, gBlack, bBlack, aBlack := color.Black.RGBA()
	rWhite, gWhite, bWhite, aWhite := color.White.RGBA()

	err := ParallelRowReadWriteE(image, func(_, _ int, col color.Color) (color.Color, error) {
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

func TestParallelRowColorReadWriteShouldAccessPixelsOnce(t *testing.T) {
	image := mockWhiteDrawableImage()

	rBlack, gBlack, bBlack, aBlack := color.Black.RGBA()
	rWhite, gWhite, bWhite, aWhite := color.White.RGBA()

	ParallelRowColorReadWrite(image, func(col color.Color) color.Color {
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

func TestParallelRowColorReadWriteEShouldPanicOnNilImage(t *testing.T) {
	assert.Panics(t, func() {
		ParallelRowColorReadWriteE(nil, func(c color.Color) (color.Color, error) {
			return c, nil
		})
	})
}

func TestParallelRowColorReadWriteEShouldPanicOnNilAccessFunc(t *testing.T) {
	img := mockWhiteDrawableImage()

	assert.Panics(t, func() {
		ParallelRowColorReadWriteE(img, nil)
	})
}

func TestParallelRowColorReadWriteEShouldReturnErrorOnAccessError(t *testing.T) {
	img := mockWhiteDrawableImage()

	err := ParallelRowColorReadWriteE(img, func(c color.Color) (color.Color, error) {
		return c, errors.New("pimit-test: test errror")
	})

	assert.NotNil(t, err)
}

func TestPrallelRowColorReadWriteEShouldCorrectlyIterate(t *testing.T) {
	img := mockWhiteDrawableImage()

	exR, exG, exB, exA := color.White.RGBA()

	err := ParallelRowColorReadWriteE(img, func(c color.Color) (color.Color, error) {
		acR, acG, acB, acA := c.RGBA()

		assert.Equal(t, exR, acR)
		assert.Equal(t, exG, acG)
		assert.Equal(t, exB, acB)
		assert.Equal(t, exA, acA)

		return color.Black, nil
	})

	assert.Nil(t, err)

	expectedImage := mockBlackDrawableImage()

	assert.Equal(t, expectedImage.Bounds(), img.Bounds())

	for x := 0; x < expectedImage.Bounds().Dx(); x += 1 {
		for y := 0; y < expectedImage.Bounds().Dy(); y += 1 {
			assert.Equal(t, expectedImage.At(x, y), img.At(x, y))
		}
	}
}

func TestParallelRowColorReadWriteEShouldAccessPixelsOnce(t *testing.T) {
	image := mockWhiteDrawableImage()

	rBlack, gBlack, bBlack, aBlack := color.Black.RGBA()
	rWhite, gWhite, bWhite, aWhite := color.White.RGBA()

	err := ParallelRowColorReadWriteE(image, func(col color.Color) (color.Color, error) {
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

func TestParallelRowReadWriteNewShouldAccessPixelsOnce(t *testing.T) {
	image := mockWhiteDrawableImage()

	rBlack, gBlack, bBlack, aBlack := color.Black.RGBA()
	rWhite, gWhite, bWhite, aWhite := color.White.RGBA()

	actualImage := ParallelRowReadWriteNew(image, func(_, _ int, col color.Color) color.Color {
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

func TestParallelRowReadWriteNewEShouldPanicOnNilImage(t *testing.T) {
	assert.Panics(t, func() {
		ParallelRowReadWriteNewE(nil, func(_, _ int, c color.Color) (color.Color, error) {
			return c, nil
		})
	})
}

func TestParallelRowReadWriteNewEShouldPanicOnNilAccessFunc(t *testing.T) {
	img := mockWhiteDrawableImage()

	assert.Panics(t, func() {
		ParallelRowReadWriteNewE(img, nil)
	})
}

func TestParallelRowReadWriteNewEShouldReturnErrorOnAccessError(t *testing.T) {
	img := mockWhiteDrawableImage()

	modifiedImg, err := ParallelRowReadWriteNewE(img, func(xIndex, yIndex int, c color.Color) (color.Color, error) {
		return c, errors.New("pimit-test: test errror")
	})

	assert.NotNil(t, err)
	assert.Nil(t, modifiedImg)
}

func TestPrallelRowReadWriteNewEShouldCorrectlyIterate(t *testing.T) {
	img := mockWhiteDrawableImage()

	exR, exG, exB, exA := color.White.RGBA()

	actualImage, err := ParallelRowReadWriteNewE(img, func(xIndex, yIndex int, c color.Color) (color.Color, error) {
		assert.GreaterOrEqual(t, xIndex, 0)
		assert.GreaterOrEqual(t, yIndex, 0)
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

func TestParallelRowReadWriteENewShouldAccessPixelsOnce(t *testing.T) {
	image := mockWhiteDrawableImage()

	rBlack, gBlack, bBlack, aBlack := color.Black.RGBA()
	rWhite, gWhite, bWhite, aWhite := color.White.RGBA()

	actualImage, err := ParallelRowReadWriteNewE(image, func(_, _ int, col color.Color) (color.Color, error) {
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

func TestParallelRowColorReadWriteNewShouldAccessPixelsOnce(t *testing.T) {
	image := mockWhiteDrawableImage()

	rBlack, gBlack, bBlack, aBlack := color.Black.RGBA()
	rWhite, gWhite, bWhite, aWhite := color.White.RGBA()

	actualImage := ParallelRowColorReadWriteNew(image, func(col color.Color) color.Color {
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

func TestParallelRowColorReadWriteNewEShouldPanicOnNilImage(t *testing.T) {
	assert.Panics(t, func() {
		ParallelRowColorReadWriteNewE(nil, func(c color.Color) (color.Color, error) {
			return c, nil
		})
	})
}

func TestParallelRowColorReadWriteNewEShouldPanicOnNilAccessFunc(t *testing.T) {
	img := mockWhiteDrawableImage()

	assert.Panics(t, func() {
		ParallelRowColorReadWriteNewE(img, nil)
	})
}

func TestParallelRowColorReadWriteNewEShouldReturnErrorOnAccessError(t *testing.T) {
	img := mockWhiteDrawableImage()

	modifiedImg, err := ParallelRowColorReadWriteNewE(img, func(c color.Color) (color.Color, error) {
		return c, errors.New("pimit-test: test errror")
	})

	assert.NotNil(t, err)
	assert.Nil(t, modifiedImg)
}

func TestPrallelRowColorReadWriteNewEShouldCorrectlyIterate(t *testing.T) {
	img := mockWhiteDrawableImage()

	exR, exG, exB, exA := color.White.RGBA()

	actualImage, err := ParallelRowColorReadWriteNewE(img, func(c color.Color) (color.Color, error) {
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

func TestParallelRowColorReadWriteNewEShouldAccessPixelsOnce(t *testing.T) {
	image := mockWhiteDrawableImage()

	rBlack, gBlack, bBlack, aBlack := color.Black.RGBA()
	rWhite, gWhite, bWhite, aWhite := color.White.RGBA()

	actualImage, err := ParallelRowColorReadWriteNewE(image, func(col color.Color) (color.Color, error) {
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
