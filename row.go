package pimit

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"sync"
)

// Perform a parallel iteration of the indexes according to the width and height indicated by the parameters. For
// each combination of indexes excute the passed access function. Each row is iterated in a separate goroutine.
//
// Deprecated: Use ParallelIndices instead.
func ParallelRowIndices(w, h int, a IndicesAccess) {
	if w <= 0 {
		panic("pimit: the provided nagative or zero width is invalid")
	}

	if h <= 0 {
		panic("pimit: the provided negative or zero height is invalid")
	}

	wg := sync.WaitGroup{}
	wg.Add(h)

	for y := 0; y < h; y += 1 {
		go func(yIndex int) {
			defer wg.Done()

			for xIndex := 0; xIndex < w; xIndex += 1 {
				a(xIndex, yIndex)
			}
		}(y)
	}

	wg.Wait()
}

// Perform a parallel reading of the pixels of the passed image. For each pixel, execute the passed access
// function allowing you to read the color and coordinates. Every row is iterated in a separate goroutine.
//
// Deprecated: Use ParallelRead function instead.
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
	wg.Add(height)

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
//
// Deprecated: Use ParallelRead instead.
func ParallelRowColorRead(i image.Image, a ReadColorAccess) {
	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	ParallelRowRead(i, func(_, _ int, color color.Color) { a(color) })
}

// Perform a parallel reading of the pixels of the passed image. For each pixel, execute the passed access
// function allowing you to read the color and coordinates. Every row is iterated in a separate goroutine.
// Errors that occur in the function will be caught and the first one will be returned by the function.
//
// Deprecated: Used ParallelReadE instead.
func ParallelRowReadE(i image.Image, a ReadAccessE) error {
	if i == nil {
		panic("pimit: the provided image reference is nil")
	}

	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	width := i.Bounds().Dx()
	height := i.Bounds().Dy()

	wg := sync.WaitGroup{}
	wg.Add(height)

	iterationErrors := make(chan error, height)

	for y := 0; y < height; y += 1 {
		go func(yIndex int, errCh chan error) {
			defer wg.Done()

			for xIndex := 0; xIndex < width; xIndex += 1 {
				color := i.At(xIndex, yIndex)

				if err := a(xIndex, yIndex, color); err != nil {
					errCh <- fmt.Errorf("pimit: access function failed on x:%d y:%d: %w", xIndex, yIndex, err)
					return
				}
			}
		}(y, iterationErrors)
	}

	wg.Wait()

	var err error = nil
	if len(iterationErrors) > 0 {
		err = <-iterationErrors
		close(iterationErrors)
	}

	return err
}

// Perform a parallel reading of the pixels of the passed image. For each pixel, execute the passed access
// function allowing you to read only the color. Every row is iterated in a separate goroutine. Errors
// that occur in the function will be caught and the first one will be returned by the function.
//
// Deprecated: Use ParallelReadE instead.
func ParallelRowColorReadE(i image.Image, a ReadColorAccessE) error {
	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	return ParallelRowReadE(i, func(_, _ int, color color.Color) error { return a(color) })
}

// Perform parallel reading and editing of pixels of the passed image. For each pixel, execute the passed access
// function allowing to read the color and coordinates, which will return the color that the pixel should take
// after this operation. The changes will be applied to the passed image instance. Every row is iterated in a
// separate goroutine.
//
// Deprecated: Use ParallelReadWrite instead.
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
	wg.Add(height)

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
//
// Deprecated: Use ParallelReadWrite instead.
func ParallelRowColorReadWrite(i draw.Image, a ReadWriteColorAccess) {
	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	ParallelRowReadWrite(i, func(_, _ int, color color.Color) color.Color { return a(color) })
}

// Perform parallel reading and editing of pixels of the passed image. For each pixel, execute the passed access
// function allowing to read the color and coordinates, which will return the color that the pixel should take
// after this operation. The changes will be applied to the passed image instance. Every row is iterated in a
// separate goroutine. Errors that occur in the function will be caught and the first one will be returned by the function.
//
// Deprecated: Use ParallelReadWriteE instead.
func ParallelRowReadWriteE(i draw.Image, a ReadWriteAccessE) error {
	if i == nil {
		panic("pimit: the provided image reference is nil")
	}

	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	width := i.Bounds().Dx()
	height := i.Bounds().Dy()

	wg := sync.WaitGroup{}
	wg.Add(height)

	iterationErrors := make(chan error, height)

	for y := 0; y < height; y += 1 {
		go func(yIndex int, errCh chan error) {
			defer wg.Done()

			for xIndex := 0; xIndex < width; xIndex += 1 {
				originalColor := i.At(xIndex, yIndex)
				modifiedColor, err := a(xIndex, yIndex, originalColor)
				if err != nil {
					errCh <- fmt.Errorf("pimit: access function failed on x:%d y:%d: %w", xIndex, yIndex, err)
					return
				}

				i.Set(xIndex, yIndex, modifiedColor)
			}
		}(y, iterationErrors)
	}

	wg.Wait()

	var err error = nil
	if len(iterationErrors) > 0 {
		err = <-iterationErrors
		close(iterationErrors)
	}

	return err
}

// Perform parallel reading and editing of pixels of the passed image. For each pixel, execute the passed access
// function allowing to read only the color, which will return the color that the pixel should take after this
// operation. The changes will be applied to the passed image instance.  Every row is iterated in a separate
// goroutine. Errors that occur in the function will be caught and the first one will be returned by the function.
//
// Deprecated: Use ParallelReadWriteE instead.
func ParallelRowColorReadWriteE(i draw.Image, a ReadWriteColorAccessE) error {
	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	return ParallelRowReadWriteE(i, func(_, _ int, color color.Color) (color.Color, error) { return a(color) })
}

// Perform parallel reading and editing of pixels of the passed image. For each pixel, execute the passed access
// function allowing to read the color and coordinates, which will return the color that the pixel should take
// after this operation. The changes will be applied to a new image instance returned from this function. Every
// row is iterated in a separate goroutine.
//
// Deprecated: Use ParallelReadWriteNew instead.
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
	wg.Add(height)

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
//
// Deprecated: Use ParallelReadWriteNew instead.
func ParallelRowColorReadWriteNew(i image.Image, a ReadWriteColorAccess) image.Image {
	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	return ParallelRowReadWriteNew(i, func(_, _ int, color color.Color) color.Color { return a(color) })
}

// Perform parallel reading and editing of pixels of the passed image. For each pixel, execute the passed access
// function allowing to read the color and coordinates, which will return the color that the pixel should take
// after this operation. The changes will be applied to a new image instance returned from this function. Every
// row is iterated in a separate goroutine. Errors that occur in the function will be caught and the first one
// will be returned by the function.
//
// Deprecated: Use ParallelReadWriteNewE instead.
func ParallelRowReadWriteNewE(i image.Image, a ReadWriteAccessE) (image.Image, error) {
	if i == nil {
		panic("pimit: the provided image reference is nil")
	}

	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	width := i.Bounds().Dx()
	height := i.Bounds().Dy()

	wg := sync.WaitGroup{}
	wg.Add(height)

	iterationErrors := make(chan error, height)

	outputImage := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y += 1 {
		go func(yIndex int, errCh chan error) {
			defer wg.Done()

			for xIndex := 0; xIndex < width; xIndex += 1 {
				color := i.At(xIndex, yIndex)

				modifiedColor, err := a(xIndex, yIndex, color)
				if err != nil {
					errCh <- fmt.Errorf("pimit: access function failed on x:%d y:%d: %w", xIndex, yIndex, err)
					return
				}

				outputImage.Set(xIndex, yIndex, modifiedColor)
			}
		}(y, iterationErrors)
	}

	wg.Wait()

	var err error = nil
	if len(iterationErrors) > 0 {
		err = <-iterationErrors
		close(iterationErrors)

		return nil, err
	}

	return outputImage, err
}

// Perform parallel reading and editing of pixels of the passed image. For each pixel, execute the passed access
// function allowing to read only the color, which will return the color that the pixel should take after this
// operation. The changes will be applied to the passed image instance. The changes will be applied to a new image
// instance returned from this function. Every row is iterated in a separate goroutine. Errors that occur in
// the function will be caught and the first one will be returned by the function.
//
// Deprecated: Use ParallelReadWriteNewE instead.
func ParallelRowColorReadWriteNewE(i image.Image, a ReadWriteColorAccessE) (image.Image, error) {
	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	return ParallelRowReadWriteNewE(i, func(_, _ int, color color.Color) (color.Color, error) { return a(color) })
}
