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

func mockSpecificDrawableImage(width, height int, color color.Color) draw.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			img.Set(x, y, color)
		}
	}

	return img
}

func mockWhiteDrawableImage() draw.Image {
	return mockSpecificDrawableImage(mockImageWidth, mockImageHeight, color.White)
}

func mockWhiteImage() image.Image {
	return mockSpecificDrawableImage(mockImageWidth, mockImageHeight, color.White)
}

func mockBlackDrawableImage() draw.Image {
	return mockSpecificDrawableImage(mockImageWidth, mockImageHeight, color.Black)
}

func mockBlackImage() image.Image {
	return mockSpecificDrawableImage(mockImageWidth, mockImageHeight, color.Black)
}
