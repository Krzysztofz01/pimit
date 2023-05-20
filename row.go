package pimit

import (
	"image"
	"image/color"
	"image/draw"
	"sync"
)

// Perform a parallel reading of the pixels of the passed image. For each pixel, execute the passed access
// function allowing you to read the color and coordinates. Every row is iterated in a separate goroutine.
func ParallelRowRead(i image.Image, a ReadAccess) {
	if i == nil {
		panic("pimit: the provided image reference is nil")
	}

	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	width := i.Bounds().Dx()
	height := i.Bounds().Dy()

	wg := sync.WaitGroup{}
	wg.Add(width)

	for y := 0; y < height; y += 1 {
		go func(yIndex int) {
			defer wg.Done()

			for xIndex := 0; xIndex < width; xIndex += 1 {
				color := i.At(xIndex, yIndex)
				a(xIndex, yIndex, color)
			}
		}(y)
	}

	wg.Wait()
}

// Perform a parallel reading of the pixels of the passed image. For each pixel, execute the passed access
// function allowing you to read only the color. Every row is iterated in a separate goroutine.
func ParallelRowColorRead(i image.Image, a ReadColorAccess) {
	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	ParallelRowRead(i, func(_, _ int, color color.Color) { a(color) })
}

// Perform parallel reading and editing of pixels of the passed image. For each pixel, execute the passed access
// function allowing to read the color and coordinates, which will return the color that the pixel should take
// after this operation. The changes will be applied to the passed image instance. Every row is iterated in a
// separate goroutine.
func ParallelRowReadWrite(i draw.Image, a ReadWriteAccess) {
	if i == nil {
		panic("pimit: the provided image reference is nil")
	}

	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	width := i.Bounds().Dx()
	height := i.Bounds().Dy()

	wg := sync.WaitGroup{}
	wg.Add(width)

	for y := 0; y < height; y += 1 {
		go func(yIndex int) {
			defer wg.Done()

			for xIndex := 0; xIndex < width; xIndex += 1 {
				originalColor := i.At(xIndex, yIndex)
				modifiedColor := a(xIndex, yIndex, originalColor)

				i.Set(xIndex, yIndex, modifiedColor)
			}
		}(y)
	}

	wg.Wait()
}

// Perform parallel reading and editing of pixels of the passed image. For each pixel, execute the passed access
// function allowing to read only the color, which will return the color that the pixel should take after this
// operation. The changes will be applied to the passed image instance.  Every row is iterated in a separate
// goroutine.
func ParallelRowColorReadWrite(i draw.Image, a ReadWriteColorAccess) {
	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	ParallelRowReadWrite(i, func(_, _ int, color color.Color) color.Color { return a(color) })
}

// Perform parallel reading and editing of pixels of the passed image. For each pixel, execute the passed access
// function allowing to read the color and coordinates, which will return the color that the pixel should take
// after this operation. The changes will be applied to a new image instance returned from this function. Every
// row is iterated in a separate goroutine.
func ParallelRowReadWriteNew(i image.Image, a ReadWriteAccess) image.Image {
	if i == nil {
		panic("pimit: the provided image reference is nil")
	}

	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	width := i.Bounds().Dx()
	height := i.Bounds().Dy()

	wg := sync.WaitGroup{}
	wg.Add(width)

	outputImage := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y += 1 {
		go func(yIndex int) {
			defer wg.Done()

			for xIndex := 0; xIndex < width; xIndex += 1 {
				originalColor := i.At(xIndex, yIndex)
				modifiedColor := a(xIndex, yIndex, originalColor)

				outputImage.Set(xIndex, yIndex, modifiedColor)
			}
		}(y)
	}

	wg.Wait()
	return outputImage
}

// Perform parallel reading and editing of pixels of the passed image. For each pixel, execute the passed access
// function allowing to read only the color, which will return the color that the pixel should take after this
// operation. The changes will be applied to the passed image instance. The changes will be applied to a new image
// instance returned from this function. Every row is iterated in a separate goroutine.
func ParallelRowColorReadWriteNew(i image.Image, a ReadWriteColorAccess) image.Image {
	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	return ParallelRowReadWriteNew(i, func(_, _ int, color color.Color) color.Color { return a(color) })
}
