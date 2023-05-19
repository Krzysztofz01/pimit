package pimit

import (
	"image"
	"image/color"
	"image/draw"
)

const (
	mockImageWidth  = 5
	mockImageHeight = 5
)

func mockWhiteDrawableImage() draw.Image {
	img := image.NewRGBA(image.Rect(0, 0, mockImageWidth, mockImageHeight))

	for x := 0; x < mockImageWidth; x += 1 {
		for y := 0; y < mockImageHeight; y += 1 {
			img.Set(x, y, color.White)
		}
	}

	return img
}

func mockWhiteImage() image.Image {
	return mockWhiteDrawableImage()
}

func mockBlackDrawableImage() draw.Image {
	img := image.NewRGBA(image.Rect(0, 0, mockImageWidth, mockImageHeight))

	for x := 0; x < mockImageWidth; x += 1 {
		for y := 0; y < mockImageHeight; y += 1 {
			img.Set(x, y, color.Black)
		}
	}

	return img
}

func mockBlackImage() image.Image {
	return mockBlackDrawableImage()
}
