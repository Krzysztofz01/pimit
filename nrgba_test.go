package pimit

import (
	"errors"
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestParallelNrgbaReadShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelNrgbaRead(nil, func(x, y int, r, g, b, a uint8) {})
	})
}

func TestParallelNrgbaReadShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageNrgba()

	assert.Panics(t, func() {
		ParallelNrgbaRead(img, nil)
	})
}

func TestParallelNrgbaReadShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageNrgba()

	exR, exG, exB, exA := uint8(255), uint8(255), uint8(255), uint8(255)

	ParallelNrgbaRead(img, func(xIndex int, yIndex int, acR, acG, acB, acA uint8) {
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

func TestParallelNrgbaReadEShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelNrgbaReadE(nil, func(x, y int, r, g, b, a uint8) error {
			return nil
		})
	})
}

func TestParallelNrgbaReadEShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageNrgba()

	assert.Panics(t, func() {
		ParallelNrgbaReadE(img, nil)
	})
}

func TestParallelNrgbaReadEShouldReturnErrorOnAccessError(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageNrgba()

	err := ParallelNrgbaReadE(img, func(x, y int, r, g, b, a uint8) error {
		return errors.New("pimit-test: test errror")
	})

	assert.NotNil(t, err)
}

func TestParallelNrgbaReadEShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageNrgba()

	exR, exG, exB, exA := uint8(255), uint8(255), uint8(255), uint8(255)

	err := ParallelNrgbaReadE(img, func(xIndex, yIndex int, acR, acG, acB, acA uint8) error {
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

func TestParallelNrgbaReadWriteShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelNrgbaReadWrite(nil, func(x int, y int, r uint8, g uint8, b uint8, a uint8) (uint8, uint8, uint8, uint8) {
			return r, g, b, a
		})
	})
}

func TestParallelNrgbaReadWriteShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageNrgba()

	assert.Panics(t, func() {
		ParallelNrgbaReadWrite(img, nil)
	})
}

func TestParallelNrgbaReadWriteShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageNrgba()

	exR, exG, exB, exA := uint8(255), uint8(255), uint8(255), uint8(255)

	ParallelNrgbaReadWrite(img, func(xIndex, yIndex int, acR, acG, acB, acA uint8) (uint8, uint8, uint8, uint8) {
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

	expectedImage := mockBlackImageNrgba()

	assert.Equal(t, expectedImage.Bounds(), img.Bounds())

	for x := 0; x < expectedImage.Bounds().Dx(); x += 1 {
		for y := 0; y < expectedImage.Bounds().Dy(); y += 1 {
			assert.Equal(t, expectedImage.NRGBAAt(x, y), img.NRGBAAt(x, y))
		}
	}
}

func TestParallelNrgbaReadWriteShouldAccessPixelsOnce(t *testing.T) {
	defer goleak.VerifyNone(t)

	image := mockWhiteImageNrgba()

	rBlack, gBlack, bBlack, aBlack := uint8(0), uint8(0), uint8(0), uint8(0)
	rWhite, gWhite, bWhite, aWhite := uint8(255), uint8(255), uint8(255), uint8(255)

	ParallelNrgbaReadWrite(image, func(x, y int, rCurrent, gCurrent, bCurrent, aCurrent uint8) (uint8, uint8, uint8, uint8) {
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
			c := image.NRGBAAt(x, y)

			assert.Equal(t, rBlack, c.R)
			assert.Equal(t, gBlack, c.G)
			assert.Equal(t, bBlack, c.B)
			assert.Equal(t, aBlack, c.A)
		}
	}
}

func TestParallelNrgbaReadWriteEShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelNrgbaReadWriteE(nil, func(x, y int, r, g, b, a uint8) (uint8, uint8, uint8, uint8, error) {
			return r, g, b, a, nil
		})
	})
}

func TestParallelNrgbaReadWriteEShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageNrgba()

	assert.Panics(t, func() {
		ParallelNrgbaReadWriteE(img, nil)
	})
}

func TestParallelNrgbaReadWriteEShouldReturnErrorOnAccessError(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageNrgba()

	err := ParallelNrgbaReadWriteE(img, func(x, y int, r, g, b, a uint8) (uint8, uint8, uint8, uint8, error) {
		return r, g, b, a, errors.New("pimit-test: test errror")
	})

	assert.NotNil(t, err)
}

func TestParallelNrgbaReadWriteEShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageNrgba()

	exR, exG, exB, exA := uint8(255), uint8(255), uint8(255), uint8(255)

	ParallelNrgbaReadWriteE(img, func(xIndex, yIndex int, acR, acG, acB, acA uint8) (uint8, uint8, uint8, uint8, error) {
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

	expectedImage := mockBlackImageNrgba()

	assert.Equal(t, expectedImage.Bounds(), img.Bounds())

	for x := 0; x < expectedImage.Bounds().Dx(); x += 1 {
		for y := 0; y < expectedImage.Bounds().Dy(); y += 1 {
			assert.Equal(t, expectedImage.NRGBAAt(x, y), img.NRGBAAt(x, y))
		}
	}
}

func TestParallelNrgbaReadWriteEShouldAccessPixelsOnce(t *testing.T) {
	defer goleak.VerifyNone(t)

	image := mockWhiteImageNrgba()

	rBlack, gBlack, bBlack, aBlack := uint8(0), uint8(0), uint8(0), uint8(0)
	rWhite, gWhite, bWhite, aWhite := uint8(255), uint8(255), uint8(255), uint8(255)

	err := ParallelNrgbaReadWriteE(image, func(xIndex, yIndex int, rCurrent, gCurrent, bCurrent, aCurrent uint8) (uint8, uint8, uint8, uint8, error) {
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
			c := image.NRGBAAt(x, y)

			assert.Equal(t, rBlack, c.R)
			assert.Equal(t, gBlack, c.G)
			assert.Equal(t, bBlack, c.B)
			assert.Equal(t, aBlack, c.A)
		}
	}
}

func TestParallelNrgbaReadWriteNewShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelNrgbaReadWriteNew(nil, func(x int, y int, r uint8, g uint8, b uint8, a uint8) (uint8, uint8, uint8, uint8) {
			return r, g, b, a
		})
	})
}

func TestParallelNrgbaReadWriteNewShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageNrgba()

	assert.Panics(t, func() {
		ParallelNrgbaReadWriteNew(img, nil)
	})
}

func TestParallelNrgbaReadWriteNewShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageNrgba()

	exR, exG, exB, exA := uint8(255), uint8(255), uint8(255), uint8(255)

	actualImage := ParallelNrgbaReadWriteNew(img, func(xIndex, yIndex int, acR, acG, acB, acA uint8) (uint8, uint8, uint8, uint8) {
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

	expectedImage := mockBlackImageNrgba()

	assert.Equal(t, expectedImage.Bounds(), actualImage.Bounds())

	for x := 0; x < expectedImage.Bounds().Dx(); x += 1 {
		for y := 0; y < expectedImage.Bounds().Dy(); y += 1 {
			assert.Equal(t, expectedImage.NRGBAAt(x, y), actualImage.NRGBAAt(x, y))
		}
	}
}

func TestParallelNrgbaReadWriteNewShouldAccessPixelsOnce(t *testing.T) {
	defer goleak.VerifyNone(t)

	image := mockWhiteImageNrgba()

	rBlack, gBlack, bBlack, aBlack := uint8(0), uint8(0), uint8(0), uint8(0)
	rWhite, gWhite, bWhite, aWhite := uint8(255), uint8(255), uint8(255), uint8(255)

	actualImage := ParallelNrgbaReadWriteNew(image, func(xIndex, yIndex int, rCurrent, gCurrent, bCurrent, aCurrent uint8) (uint8, uint8, uint8, uint8) {
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
			c := actualImage.NRGBAAt(x, y)

			assert.Equal(t, rBlack, c.R)
			assert.Equal(t, gBlack, c.G)
			assert.Equal(t, bBlack, c.B)
			assert.Equal(t, aBlack, c.A)
		}
	}
}

func TestParallelNrgbaReadWriteNewEShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelNrgbaReadWriteNewE(nil, func(x, y int, r, g, b, a uint8) (uint8, uint8, uint8, uint8, error) {
			return r, g, b, a, nil
		})
	})
}

func TestParallelNrgbaReadWriteNewEShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageNrgba()

	assert.Panics(t, func() {
		ParallelNrgbaReadWriteNewE(img, nil)
	})
}

func TestParallelNrgbaReadWriteNewEShouldReturnErrorOnAccessError(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageNrgba()

	modifiedImg, err := ParallelNrgbaReadWriteNewE(img, func(x, y int, r, g, b, a uint8) (uint8, uint8, uint8, uint8, error) {
		return r, g, b, a, errors.New("pimit-test: test errror")
	})

	assert.NotNil(t, err)
	assert.Nil(t, modifiedImg)
}

func TestParallelNrgbaReadWriteNewEShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockWhiteImageNrgba()

	exR, exG, exB, exA := uint8(255), uint8(255), uint8(255), uint8(255)

	actualImage, err := ParallelNrgbaReadWriteNewE(img, func(xIndex, yIndex int, acR, acG, acB, acA uint8) (uint8, uint8, uint8, uint8, error) {
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

	expectedImage := mockBlackImageNrgba()

	assert.Equal(t, expectedImage.Bounds(), actualImage.Bounds())

	for x := 0; x < expectedImage.Bounds().Dx(); x += 1 {
		for y := 0; y < expectedImage.Bounds().Dy(); y += 1 {
			assert.Equal(t, expectedImage.NRGBAAt(x, y), actualImage.NRGBAAt(x, y))
		}
	}
}

func TestParallelNrgbaReadWriteENewShouldAccessPixelsOnce(t *testing.T) {
	defer goleak.VerifyNone(t)

	image := mockWhiteImageNrgba()

	rBlack, gBlack, bBlack, aBlack := uint8(0), uint8(0), uint8(0), uint8(0)
	rWhite, gWhite, bWhite, aWhite := uint8(255), uint8(255), uint8(255), uint8(255)

	actualImage, err := ParallelNrgbaReadWriteNewE(image, func(xIndex, yIndex int, rCurrent, gCurrent, bCurrent, aCurrent uint8) (uint8, uint8, uint8, uint8, error) {
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
			c := actualImage.NRGBAAt(x, y)

			assert.Equal(t, rBlack, c.R)
			assert.Equal(t, gBlack, c.G)
			assert.Equal(t, bBlack, c.B)
			assert.Equal(t, aBlack, c.A)
		}
	}
}

func mockWhiteImageNrgba() *image.NRGBA {
	width, height := 5, 6

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y += 1 {
		for x := 0; x < width; x += 1 {
			img.Set(x, y, color.White)
		}
	}

	return img
}

func mockBlackImageNrgba() *image.NRGBA {
	width, height := 5, 6

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y += 1 {
		for x := 0; x < width; x += 1 {
			img.Set(x, y, color.Black)
		}
	}

	return img
}
