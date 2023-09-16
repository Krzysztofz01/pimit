package pimit

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"sync"
)

type (
	IndicesDelegate            = func(x, y int)
	ReadDelegate               = func(x, y int, c color.Color)
	ReadErrorableDelegate      = func(x, y int, c color.Color) error
	ReadWriteDelegate          = func(x, y int, c color.Color) color.Color
	ReadWriteErrorableDelegate = func(x, y int, c color.Color) (color.Color, error)
)

// Perform a parallel iteration of the indexes according to the width and height provided via the parameters.
// Execute the delegate for each indexes combination. Each row is iterated in a separate goroutine.
func ParallelIndices(w, h int, d IndicesDelegate) {
	if w <= 0 {
		panic("pimit: the provided nagative or zero width is invalid")
	}

	if h <= 0 {
		panic("pimit: the provided negative or zero height is invalid")
	}

	wg := &sync.WaitGroup{}

	for y := 0; y < h; y += 1 {
		wg.Add(1)
		go func(yIndex int) {
			defer wg.Done()

			for xIndex := 0; xIndex < w; xIndex += 1 {
				d(xIndex, yIndex)
			}
		}(y)
	}

	wg.Wait()
}

// Perform a parallel iteration of the pixels of the provided image. For each pixel, execute the delegate function
// allowing you to read the color and coordinates. Each row is iterated in a separate goroutine.
func ParallelRead(src image.Image, d ReadDelegate) {
	if src == nil {
		panic("pimit: the provided image reference is nil")
	}

	if d == nil {
		panic("pimit: the provided access delegate function is nil")
	}

	srcWidth := src.Bounds().Dx()
	srcHeight := src.Bounds().Dy()
	wg := &sync.WaitGroup{}

	for y := 0; y < srcHeight; y += 1 {
		wg.Add(1)
		go func(yIndex int) {
			defer wg.Done()

			for xIndex := 0; xIndex < srcWidth; xIndex += 1 {
				color := src.At(xIndex, yIndex)
				d(xIndex, yIndex, color)
			}
		}(y)
	}

	wg.Wait()
}

// Perform a parallel iteration of the pixels of the provided image. For each pixel, execute the delegate function
// allowing you to read the color and coordinates. Each row is iterated in a separate goroutine. The iteration will
// break after the first error occurs and the error will be returned.
func ParallelReadE(src image.Image, d ReadErrorableDelegate) error {
	if src == nil {
		panic("pimit: the provided image reference is nil")
	}

	if d == nil {
		panic("pimit: the provided access delegate function is nil")
	}

	srcWidth := src.Bounds().Dx()
	srcHeight := src.Bounds().Dy()
	wg := &sync.WaitGroup{}

	errt := NewErrorTrap()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for y := 0; y < srcHeight; y += 1 {
		wg.Add(1)
		go func(yIndex int) {
			defer wg.Done()

			for xIndex := 0; xIndex < srcWidth; xIndex += 1 {
				select {
				case <-ctx.Done():
					return
				default:
				}

				color := src.At(xIndex, yIndex)
				if err := d(xIndex, yIndex, color); err != nil {
					errt.Set(fmt.Errorf("pimit: delegate function failed om x:%d y:%d with: %w", xIndex, yIndex, err))
					cancel()
					return
				}
			}
		}(y)
	}

	wg.Wait()
	return errt.Err()
}

// Perform a parallel iteration of the pixels of the provided image. For each pixel, execute the delegate function
// allowing you to read the color and coordinates, the delegate return color will be set at the given coordinates.
// This changes will be applied to the passed image instance. Consider using ParallelReadWriteNew if you want to
// avoid changes to the original image at the expense of additional allocations. Each row is iterated in a separate
// goroutine.
func ParallelReadWrite(src draw.Image, d ReadWriteDelegate) {
	if src == nil {
		panic("pimit: the provided image reference is nil")
	}

	if d == nil {
		panic("pimit: the provided access delegate function is nil")
	}

	srcWidth := src.Bounds().Dx()
	srcHeight := src.Bounds().Dy()
	wg := &sync.WaitGroup{}

	for y := 0; y < srcHeight; y += 1 {
		wg.Add(1)
		go func(yIndex int) {
			defer wg.Done()

			for xIndex := 0; xIndex < srcWidth; xIndex += 1 {
				color := src.At(xIndex, yIndex)
				color = d(xIndex, yIndex, color)
				src.Set(xIndex, yIndex, color)
			}
		}(y)
	}

	wg.Wait()
}

// Perform a parallel iteration of the pixels of the provided image. For each pixel, execute the delegate function
// allowing you to read the color and coordinates, the delegate return color will be set at the given coordinates.
// This changes will be applied to the passed image instance. Consider using ParallelReadWriteNewE if you want to
// avoid changes to the original image at the expense of additional allocations. Each row is iterated in a separate
// goroutine. The iteration will break after the first error occurs and the error will be returned.
func ParallelReadWriteE(src draw.Image, d ReadWriteErrorableDelegate) error {
	if src == nil {
		panic("pimit: the provided image reference is nil")
	}

	if d == nil {
		panic("pimit: the provided access delegate function is nil")
	}

	srcWidth := src.Bounds().Dx()
	srcHeight := src.Bounds().Dy()
	wg := &sync.WaitGroup{}

	errt := NewErrorTrap()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for y := 0; y < srcHeight; y += 1 {
		wg.Add(1)
		go func(yIndex int) {
			defer wg.Done()

			for xIndex := 0; xIndex < srcWidth; xIndex += 1 {
				select {
				case <-ctx.Done():
					return
				default:
				}

				color := src.At(xIndex, yIndex)
				color, err := d(xIndex, yIndex, color)

				if err != nil {
					errt.Set(fmt.Errorf("pimit: delegate function failed om x:%d y:%d with: %w", xIndex, yIndex, err))
					cancel()
					return
				}

				src.Set(xIndex, yIndex, color)
			}
		}(y)
	}

	wg.Wait()
	return errt.Err()
}

// Perform a parallel iteration of the pixels of the provided image. For each pixel, execute the delegate function
// allowing you to read the color and coordinates, the delegate return color will be set at the given coordinates.
// This changes will be applied to a new image instance which internaly uses the NRGBA color space and is returned
// by the function. Each row is iterated in a separate goroutine.
func ParallelReadWriteNew(src image.Image, d ReadWriteDelegate) draw.Image {
	if src == nil {
		panic("pimit: the provided image reference is nil")
	}

	if d == nil {
		panic("pimit: the provided access delegate function is nil")
	}

	srcWidth := src.Bounds().Dx()
	srcHeight := src.Bounds().Dy()
	dst := image.NewNRGBA(image.Rect(0, 0, srcWidth, srcHeight))
	wg := &sync.WaitGroup{}

	for y := 0; y < srcHeight; y += 1 {
		wg.Add(1)
		go func(yIndex int) {
			defer wg.Done()

			for xIndex := 0; xIndex < srcWidth; xIndex += 1 {
				color := src.At(xIndex, yIndex)
				color = d(xIndex, yIndex, color)

				dst.Set(xIndex, yIndex, color)
			}
		}(y)
	}

	wg.Wait()
	return dst
}

// Perform a parallel iteration of the pixels of the provided image. For each pixel, execute the delegate function
// allowing you to read the color and coordinates, the delegate return color will be set at the given coordinates.
// This changes will be applied to a new image instance which internaly uses the NRGBA color space and is returned
// by the function. Each row is iterated in a separate goroutine. The iteration will break after the first error
// occurs and the error will be returned.
func ParallelReadWriteNewE(src image.Image, d ReadWriteErrorableDelegate) (draw.Image, error) {
	if src == nil {
		panic("pimit: the provided image reference is nil")
	}

	if d == nil {
		panic("pimit: the provided access delegate function is nil")
	}

	srcWidth := src.Bounds().Dx()
	srcHeight := src.Bounds().Dy()
	dst := image.NewNRGBA(image.Rect(0, 0, srcWidth, srcHeight))
	wg := &sync.WaitGroup{}

	errt := NewErrorTrap()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for y := 0; y < srcHeight; y += 1 {
		wg.Add(1)
		go func(yIndex int) {
			defer wg.Done()

			for xIndex := 0; xIndex < srcWidth; xIndex += 1 {
				select {
				case <-ctx.Done():
					return
				default:
				}

				color := src.At(xIndex, yIndex)
				color, err := d(xIndex, yIndex, color)

				if err != nil {
					errt.Set(fmt.Errorf("pimit: delegate function failed om x:%d y:%d with: %w", xIndex, yIndex, err))
					cancel()
					return
				}

				dst.Set(xIndex, yIndex, color)
			}
		}(y)
	}

	wg.Wait()

	if err := errt.Err(); err != nil {
		return nil, err
	} else {
		return dst, nil
	}
}
