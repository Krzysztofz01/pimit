package pimit

import (
	"errors"
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestParallelRgbaReadShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelRgbaRead(nil, func(x, y int, r, g, b, a uint8) {})
	})
}

func TestParallelRgbaReadShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageRgba()

	assert.Panics(t, func() {
		ParallelRgbaRead(img, nil)
	})
}

func TestParallelRgbaReadShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageRgba()

	exR, exG, exB, exA := uint8(255), uint8(255), uint8(255), uint8(255)

	ParallelRgbaRead(img, func(xIndex int, yIndex int, acR, acG, acB, acA uint8) {
		assert.GreaterOrEqual(t, xIndex, 0)
		assert.Less(t, xIndex, img.Bounds().Dx())

		assert.GreaterOrEqual(t, yIndex, 0)
		assert.Less(t, yIndex, img.Bounds().Dy())

		assert.Equal(t, exR, acR)
		assert.Equal(t, exG, acG)
		assert.Equal(t, exB, acB)
		assert.Equal(t, exA, acA)
	})
}

func TestParallelRgbaReadEShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelRgbaReadE(nil, func(x, y int, r, g, b, a uint8) error {
			return nil
		})
	})
}

func TestParallelRgbaReadEShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageRgba()

	assert.Panics(t, func() {
		ParallelRgbaReadE(img, nil)
	})
}

func TestParallelRgbaReadEShouldReturnErrorOnAccessError(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageRgba()

	err := ParallelRgbaReadE(img, func(x, y int, r, g, b, a uint8) error {
		return errors.New("pimit-test: test errror")
	})

	assert.NotNil(t, err)
}

func TestParallelRgbaReadEShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageRgba()

	exR, exG, exB, exA := uint8(255), uint8(255), uint8(255), uint8(255)

	err := ParallelRgbaReadE(img, func(xIndex, yIndex int, acR, acG, acB, acA uint8) error {
		assert.GreaterOrEqual(t, xIndex, 0)
		assert.Less(t, xIndex, img.Bounds().Dx())

		assert.GreaterOrEqual(t, yIndex, 0)
		assert.Less(t, yIndex, img.Bounds().Dy())

		assert.Equal(t, exR, acR)
		assert.Equal(t, exG, acG)
		assert.Equal(t, exB, acB)
		assert.Equal(t, exA, acA)
		return nil
	})

	assert.Nil(t, err)
}

func TestParallelRgbaReadWriteShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelRgbaReadWrite(nil, func(x int, y int, r uint8, g uint8, b uint8, a uint8) (uint8, uint8, uint8, uint8) {
			return r, g, b, a
		})
	})
}

func TestParallelRgbaReadWriteShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageRgba()

	assert.Panics(t, func() {
		ParallelRgbaReadWrite(img, nil)
	})
}

func TestParallelRgbaReadWriteShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageRgba()

	exR, exG, exB, exA := uint8(255), uint8(255), uint8(255), uint8(255)

	ParallelRgbaReadWrite(img, func(xIndex, yIndex int, acR, acG, acB, acA uint8) (uint8, uint8, uint8, uint8) {
		assert.GreaterOrEqual(t, xIndex, 0)
		assert.Less(t, xIndex, img.Bounds().Dx())

		assert.GreaterOrEqual(t, yIndex, 0)
		assert.Less(t, yIndex, img.Bounds().Dy())

		assert.Equal(t, exR, acR)
		assert.Equal(t, exG, acG)
		assert.Equal(t, exB, acB)
		assert.Equal(t, exA, acA)

		return 0, 0, 0, 255
	})

	expectedImage := mockBlackImageRgba()

	assert.Equal(t, expectedImage.Bounds(), img.Bounds())

	for x := 0; x < expectedImage.Bounds().Dx(); x += 1 {
		for y := 0; y < expectedImage.Bounds().Dy(); y += 1 {
			assert.Equal(t, expectedImage.RGBAAt(x, y), img.RGBAAt(x, y))
		}
	}
}

func TestParallelRgbaReadWriteShouldAccessPixelsOnce(t *testing.T) {
	defer goleak.VerifyNone(t)

	image := mockWhiteImageRgba()

	rBlack, gBlack, bBlack, aBlack := uint8(0), uint8(0), uint8(0), uint8(0)
	rWhite, gWhite, bWhite, aWhite := uint8(255), uint8(255), uint8(255), uint8(255)

	ParallelRgbaReadWrite(image, func(x, y int, rCurrent, gCurrent, bCurrent, aCurrent uint8) (uint8, uint8, uint8, uint8) {
		if rCurrent == rBlack && gCurrent == gBlack && bCurrent == bBlack && aCurrent == aBlack {
			return 255, 255, 255, 255
		}

		if rCurrent == rWhite && gCurrent == gWhite && bCurrent == bWhite && aCurrent == aWhite {
			return 0, 0, 0, 0
		}

		assert.FailNow(t, "This should never happen")
		return rCurrent, gCurrent, bCurrent, aCurrent
	})

	for x := 0; x < image.Bounds().Dx(); x += 1 {
		for y := 0; y < image.Bounds().Dy(); y += 1 {
			c := image.RGBAAt(x, y)

			assert.Equal(t, rBlack, c.R)
			assert.Equal(t, gBlack, c.G)
			assert.Equal(t, bBlack, c.B)
			assert.Equal(t, aBlack, c.A)
		}
	}
}

func TestParallelRgbaReadWriteEShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelRgbaReadWriteE(nil, func(x, y int, r, g, b, a uint8) (uint8, uint8, uint8, uint8, error) {
			return r, g, b, a, nil
		})
	})
}

func TestParallelRgbaReadWriteEShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageRgba()

	assert.Panics(t, func() {
		ParallelRgbaReadWriteE(img, nil)
	})
}

func TestParallelRgbaReadWriteEShouldReturnErrorOnAccessError(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageRgba()

	err := ParallelRgbaReadWriteE(img, func(x, y int, r, g, b, a uint8) (uint8, uint8, uint8, uint8, error) {
		return r, g, b, a, errors.New("pimit-test: test errror")
	})

	assert.NotNil(t, err)
}

func TestParallelRgbaReadWriteEShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageRgba()

	exR, exG, exB, exA := uint8(255), uint8(255), uint8(255), uint8(255)

	ParallelRgbaReadWriteE(img, func(xIndex, yIndex int, acR, acG, acB, acA uint8) (uint8, uint8, uint8, uint8, error) {
		assert.GreaterOrEqual(t, xIndex, 0)
		assert.Less(t, xIndex, img.Bounds().Dx())

		assert.GreaterOrEqual(t, yIndex, 0)
		assert.Less(t, yIndex, img.Bounds().Dy())

		assert.Equal(t, exR, acR)
		assert.Equal(t, exG, acG)
		assert.Equal(t, exB, acB)
		assert.Equal(t, exA, acA)

		return 0, 0, 0, 255, nil
	})

	expectedImage := mockBlackImageRgba()

	assert.Equal(t, expectedImage.Bounds(), img.Bounds())

	for x := 0; x < expectedImage.Bounds().Dx(); x += 1 {
		for y := 0; y < expectedImage.Bounds().Dy(); y += 1 {
			assert.Equal(t, expectedImage.RGBAAt(x, y), img.RGBAAt(x, y))
		}
	}
}

func TestParallelRgbaReadWriteEShouldAccessPixelsOnce(t *testing.T) {
	defer goleak.VerifyNone(t)

	image := mockWhiteImageRgba()

	rBlack, gBlack, bBlack, aBlack := uint8(0), uint8(0), uint8(0), uint8(0)
	rWhite, gWhite, bWhite, aWhite := uint8(255), uint8(255), uint8(255), uint8(255)

	err := ParallelRgbaReadWriteE(image, func(xIndex, yIndex int, rCurrent, gCurrent, bCurrent, aCurrent uint8) (uint8, uint8, uint8, uint8, error) {
		if rCurrent == rBlack && gCurrent == gBlack && bCurrent == bBlack && aCurrent == aBlack {
			return 255, 255, 255, 255, nil
		}

		if rCurrent == rWhite && gCurrent == gWhite && bCurrent == bWhite && aCurrent == aWhite {
			return 0, 0, 0, 0, nil
		}

		assert.FailNow(t, "This should never happen")
		return rCurrent, gCurrent, bCurrent, aCurrent, nil

	})

	assert.Nil(t, err)

	for x := 0; x < image.Bounds().Dx(); x += 1 {
		for y := 0; y < image.Bounds().Dy(); y += 1 {
			c := image.RGBAAt(x, y)

			assert.Equal(t, rBlack, c.R)
			assert.Equal(t, gBlack, c.G)
			assert.Equal(t, bBlack, c.B)
			assert.Equal(t, aBlack, c.A)
		}
	}
}

func TestParallelRgbaReadWriteNewShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelRgbaReadWriteNew(nil, func(x int, y int, r uint8, g uint8, b uint8, a uint8) (uint8, uint8, uint8, uint8) {
			return r, g, b, a
		})
	})
}

func TestParallelRgbaReadWriteNewShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageRgba()

	assert.Panics(t, func() {
		ParallelRgbaReadWriteNew(img, nil)
	})
}

func TestParallelRgbaReadWriteNewShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageRgba()

	exR, exG, exB, exA := uint8(255), uint8(255), uint8(255), uint8(255)

	actualImage := ParallelRgbaReadWriteNew(img, func(xIndex, yIndex int, acR, acG, acB, acA uint8) (uint8, uint8, uint8, uint8) {
		assert.GreaterOrEqual(t, xIndex, 0)
		assert.Less(t, xIndex, img.Bounds().Dx())

		assert.GreaterOrEqual(t, yIndex, 0)
		assert.Less(t, yIndex, img.Bounds().Dy())

		assert.Equal(t, exR, acR)
		assert.Equal(t, exG, acG)
		assert.Equal(t, exB, acB)
		assert.Equal(t, exA, acA)

		return 0, 0, 0, 255
	})

	expectedImage := mockBlackImageRgba()

	assert.Equal(t, expectedImage.Bounds(), actualImage.Bounds())

	for x := 0; x < expectedImage.Bounds().Dx(); x += 1 {
		for y := 0; y < expectedImage.Bounds().Dy(); y += 1 {
			assert.Equal(t, expectedImage.RGBAAt(x, y), actualImage.RGBAAt(x, y))
		}
	}
}

func TestParallelRgbaReadWriteNewShouldAccessPixelsOnce(t *testing.T) {
	defer goleak.VerifyNone(t)

	image := mockWhiteImageRgba()

	rBlack, gBlack, bBlack, aBlack := uint8(0), uint8(0), uint8(0), uint8(0)
	rWhite, gWhite, bWhite, aWhite := uint8(255), uint8(255), uint8(255), uint8(255)

	actualImage := ParallelRgbaReadWriteNew(image, func(xIndex, yIndex int, rCurrent, gCurrent, bCurrent, aCurrent uint8) (uint8, uint8, uint8, uint8) {
		if rCurrent == rBlack && gCurrent == gBlack && bCurrent == bBlack && aCurrent == aBlack {
			return 255, 255, 255, 255
		}

		if rCurrent == rWhite && gCurrent == gWhite && bCurrent == bWhite && aCurrent == aWhite {
			return 0, 0, 0, 0
		}

		assert.FailNow(t, "This should never happen")
		return rCurrent, gCurrent, bCurrent, aCurrent
	})

	for x := 0; x < image.Bounds().Dx(); x += 1 {
		for y := 0; y < image.Bounds().Dy(); y += 1 {
			c := actualImage.RGBAAt(x, y)

			assert.Equal(t, rBlack, c.R)
			assert.Equal(t, gBlack, c.G)
			assert.Equal(t, bBlack, c.B)
			assert.Equal(t, aBlack, c.A)
		}
	}
}

func TestParallelRgbaReadWriteNewEShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelRgbaReadWriteNewE(nil, func(x, y int, r, g, b, a uint8) (uint8, uint8, uint8, uint8, error) {
			return r, g, b, a, nil
		})
	})
}

func TestParallelRgbaReadWriteNewEShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageRgba()

	assert.Panics(t, func() {
		ParallelRgbaReadWriteNewE(img, nil)
	})
}

func TestParallelRgbaReadWriteNewEShouldReturnErrorOnAccessError(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageRgba()

	modifiedImg, err := ParallelRgbaReadWriteNewE(img, func(x, y int, r, g, b, a uint8) (uint8, uint8, uint8, uint8, error) {
		return r, g, b, a, errors.New("pimit-test: test errror")
	})

	assert.NotNil(t, err)
	assert.Nil(t, modifiedImg)
}

func TestParallelRgbaReadWriteNewEShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageRgba()

	exR, exG, exB, exA := uint8(255), uint8(255), uint8(255), uint8(255)

	actualImage, err := ParallelRgbaReadWriteNewE(img, func(xIndex, yIndex int, acR, acG, acB, acA uint8) (uint8, uint8, uint8, uint8, error) {
		assert.GreaterOrEqual(t, xIndex, 0)
		assert.Less(t, xIndex, img.Bounds().Dx())

		assert.GreaterOrEqual(t, yIndex, 0)
		assert.Less(t, yIndex, img.Bounds().Dy())

		assert.Equal(t, exR, acR)
		assert.Equal(t, exG, acG)
		assert.Equal(t, exB, acB)
		assert.Equal(t, exA, acA)

		return 0, 0, 0, 255, nil
	})

	assert.Nil(t, err)

	expectedImage := mockBlackImageRgba()

	assert.Equal(t, expectedImage.Bounds(), actualImage.Bounds())

	for x := 0; x < expectedImage.Bounds().Dx(); x += 1 {
		for y := 0; y < expectedImage.Bounds().Dy(); y += 1 {
			assert.Equal(t, expectedImage.RGBAAt(x, y), actualImage.RGBAAt(x, y))
		}
	}
}

func TestParallelRgbaReadWriteENewShouldAccessPixelsOnce(t *testing.T) {
	defer goleak.VerifyNone(t)

	image := mockWhiteImageRgba()

	rBlack, gBlack, bBlack, aBlack := uint8(0), uint8(0), uint8(0), uint8(0)
	rWhite, gWhite, bWhite, aWhite := uint8(255), uint8(255), uint8(255), uint8(255)

	actualImage, err := ParallelRgbaReadWriteNewE(image, func(xIndex, yIndex int, rCurrent, gCurrent, bCurrent, aCurrent uint8) (uint8, uint8, uint8, uint8, error) {
		if rCurrent == rBlack && gCurrent == gBlack && bCurrent == bBlack && aCurrent == aBlack {
			return 255, 255, 255, 255, nil
		}

		if rCurrent == rWhite && gCurrent == gWhite && bCurrent == bWhite && aCurrent == aWhite {
			return 0, 0, 0, 0, nil
		}

		assert.FailNow(t, "This should never happen")
		return rCurrent, gCurrent, bCurrent, aCurrent, nil
	})

	assert.Nil(t, err)

	for x := 0; x < image.Bounds().Dx(); x += 1 {
		for y := 0; y < image.Bounds().Dy(); y += 1 {
			c := actualImage.RGBAAt(x, y)

			assert.Equal(t, rBlack, c.R)
			assert.Equal(t, gBlack, c.G)
			assert.Equal(t, bBlack, c.B)
			assert.Equal(t, aBlack, c.A)
		}
	}
}

func mockWhiteImageRgba() *image.RGBA {
	width, height := 5, 6

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y += 1 {
		for x := 0; x < width; x += 1 {
			img.Set(x, y, color.White)
		}
	}

	return img
}

func mockBlackImageRgba() *image.RGBA {
	width, height := 5, 6

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y += 1 {
		for x := 0; x < width; x += 1 {
			img.Set(x, y, color.Black)
		}
	}

	return img
}
