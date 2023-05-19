package pimit

import (
	"image"
	"image/color"
	"image/draw"
	"sync"
)

// Perform a parallel reading of the pixels of the passed image. For each pixel, execute the passed access
// function allowing you to read the color and coordinates.
func ParallelColumnRead(i image.Image, a ReadAccess) {
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

	for x := 0; x < width; x += 1 {
		go func(xIndex int) {
			defer wg.Done()

			for yIndex := 0; yIndex < height; yIndex += 1 {
				color := i.At(xIndex, yIndex)
				a(xIndex, yIndex, color)
			}
		}(x)
	}

	wg.Wait()
}

// Perform a parallel reading of the pixels of the passed image. For each pixel, execute the passed access
// function allowing you to read only the color.
func ParallelColumnColorRead(i image.Image, a ReadColorAccess) {
	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	ParallelColumnRead(i, func(_, _ int, color color.Color) { a(color) })
}

// Perform parallel reading and editing of pixels of the passed image. For each pixel, execute the passed access
// function allowing to read the color and coordinates, which will return the color that the pixel should take
// after this operation. The changes will be applied to the passed image instance.
func ParallelColumnReadWrite(i draw.Image, a ReadWriteAccess) {
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

	for x := 0; x < width; x += 1 {
		go func(xIndex int) {
			defer wg.Done()

			for yIndex := 0; yIndex < height; yIndex += 1 {
				originalColor := i.At(xIndex, yIndex)
				modifiedColor := a(xIndex, yIndex, originalColor)

				i.Set(xIndex, yIndex, modifiedColor)
			}
		}(x)
	}

	wg.Wait()
}

// Perform parallel reading and editing of pixels of the passed image. For each pixel, execute the passed access
// function allowing to read only the color, which will return the color that the pixel should take after this
// operation. The changes will be applied to the passed image instance.
func ParallelColumnColorReadWrite(i draw.Image, a ReadWriteColorAccess) {
	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	ParallelColumnReadWrite(i, func(_, _ int, color color.Color) color.Color { return a(color) })
}

// Perform parallel reading and editing of pixels of the passed image. For each pixel, execute the passed access
// function allowing to read the color and coordinates, which will return the color that the pixel should take
// after this operation. The changes will be applied to a new image instance returned from this function.
func ParallelColumnReadWriteNew(i image.Image, a ReadWriteAccess) image.Image {
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

	for x := 0; x < width; x += 1 {
		go func(xIndex int) {
			defer wg.Done()

			for yIndex := 0; yIndex < height; yIndex += 1 {
				originalColor := i.At(xIndex, yIndex)
				modifiedColor := a(xIndex, yIndex, originalColor)

				outputImage.Set(xIndex, yIndex, modifiedColor)
			}
		}(x)
	}

	wg.Wait()
	return outputImage
}

// Perform parallel reading and editing of pixels of the passed image. For each pixel, execute the passed access
// function allowing to read only the color, which will return the color that the pixel should take after this
// operation. The changes will be applied to the passed image instance. The changes will be applied to a new image
// instance returned from this function.
func ParallelColumnColorReadWriteNew(i image.Image, a ReadWriteColorAccess) image.Image {
	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	return ParallelColumnReadWriteNew(i, func(_, _ int, color color.Color) color.Color { return a(color) })
}
