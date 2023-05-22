package pimit

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"sync"
)

// Perform a parallel reading of the pixels of the passed image. For each pixel, execute the passed access
// function allowing you to read the color and coordinates. Every column is iterated in a separate goroutine.
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
// function allowing you to read only the color. Every column is iterated in a separate goroutine.
func ParallelColumnColorRead(i image.Image, a ReadColorAccess) {
	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	ParallelColumnRead(i, func(_, _ int, color color.Color) { a(color) })
}

// Perform a parallel reading of the pixels of the passed image. For each pixel, execute the passed access
// function allowing you to read the color and coordinates. Every column is iterated in a separate goroutine.
// Errors that occur in the function will be caught and the first one will be returned by the function.
func ParallelColumnReadE(i image.Image, a ReadAccessE) error {
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

	iterationErrors := make(chan error, width)

	for x := 0; x < width; x += 1 {
		go func(xIndex int, errCh chan error) {
			defer wg.Done()

			for yIndex := 0; yIndex < height; yIndex += 1 {
				color := i.At(xIndex, yIndex)

				if err := a(xIndex, yIndex, color); err != nil {
					errCh <- fmt.Errorf("pimit: access function failed on x:%d y:%d: %w", xIndex, yIndex, err)
					return
				}
			}
		}(x, iterationErrors)
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
// function allowing you to read only the color. Every column is iterated in a separate goroutine. Errors
// that occur in the function will be caught and the first one will be returned by the function.
func ParallelColumnColorReadE(i image.Image, a ReadColorAccessE) error {
	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	return ParallelColumnReadE(i, func(_, _ int, color color.Color) error { return a(color) })
}

// Perform parallel reading and editing of pixels of the passed image. For each pixel, execute the passed access
// function allowing to read the color and coordinates, which will return the color that the pixel should take
// after this operation. The changes will be applied to the passed image instance. Every column is iterated in a
// separate goroutine.
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
// operation. The changes will be applied to the passed image instance.  Every column is iterated in a separate
// goroutine.
func ParallelColumnColorReadWrite(i draw.Image, a ReadWriteColorAccess) {
	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	ParallelColumnReadWrite(i, func(_, _ int, color color.Color) color.Color { return a(color) })
}

// Perform parallel reading and editing of pixels of the passed image. For each pixel, execute the passed access
// function allowing to read the color and coordinates, which will return the color that the pixel should take
// after this operation. The changes will be applied to the passed image instance. Every column is iterated in a
// separate goroutine. Errors that occur in the function will be caught and the first one will be returned by the function.
func ParallelColumnReadWriteE(i draw.Image, a ReadWriteAccessE) error {
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

	iterationErrors := make(chan error, width)

	for x := 0; x < width; x += 1 {
		go func(xIndex int, errCh chan error) {
			defer wg.Done()

			for yIndex := 0; yIndex < height; yIndex += 1 {
				color := i.At(xIndex, yIndex)

				modifiedColor, err := a(xIndex, yIndex, color)
				if err != nil {
					errCh <- fmt.Errorf("pimit: access function failed on x:%d y:%d: %w", xIndex, yIndex, err)
					return
				}

				i.Set(xIndex, yIndex, modifiedColor)
			}
		}(x, iterationErrors)
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
// operation. The changes will be applied to the passed image instance.  Every column is iterated in a separate
// goroutine. Errors that occur in the function will be caught and the first one will be returned by the function.
func ParallelColumnColorReadWriteE(i draw.Image, a ReadWriteColorAccessE) error {
	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	return ParallelColumnReadWriteE(i, func(_, _ int, color color.Color) (color.Color, error) { return a(color) })
}

// Perform parallel reading and editing of pixels of the passed image. For each pixel, execute the passed access
// function allowing to read the color and coordinates, which will return the color that the pixel should take
// after this operation. The changes will be applied to a new image instance returned from this function. Every
// column is iterated in a separate goroutine.
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
// instance returned from this function. Every column is iterated in a separate goroutine.
func ParallelColumnColorReadWriteNew(i image.Image, a ReadWriteColorAccess) image.Image {
	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	return ParallelColumnReadWriteNew(i, func(_, _ int, color color.Color) color.Color { return a(color) })
}

// Perform parallel reading and editing of pixels of the passed image. For each pixel, execute the passed access
// function allowing to read the color and coordinates, which will return the color that the pixel should take
// after this operation. The changes will be applied to a new image instance returned from this function. Every
// column is iterated in a separate goroutine. Errors that occur in the function will be caught and the first one
// will be returned by the function.
func ParallelColumnReadWriteNewE(i image.Image, a ReadWriteAccessE) (image.Image, error) {
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

	iterationErrors := make(chan error, width)

	outputImage := image.NewNRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x += 1 {
		go func(xIndex int, errCh chan error) {
			defer wg.Done()

			for yIndex := 0; yIndex < height; yIndex += 1 {
				color := i.At(xIndex, yIndex)

				modifiedColor, err := a(xIndex, yIndex, color)
				if err != nil {
					errCh <- fmt.Errorf("pimit: access function failed on x:%d y:%d: %w", xIndex, yIndex, err)
					return
				}

				outputImage.Set(xIndex, yIndex, modifiedColor)
			}
		}(x, iterationErrors)
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
// instance returned from this function. Every column is iterated in a separate goroutine. Errors that occur in
// the function will be caught and the first one will be returned by the function.
func ParallelColumnColorReadWriteNewE(i image.Image, a ReadWriteColorAccessE) (image.Image, error) {
	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	return ParallelColumnReadWriteNewE(i, func(_, _ int, color color.Color) (color.Color, error) { return a(color) })
}
