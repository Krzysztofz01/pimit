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

			var c color.Color = nil

			for xIndex := 0; xIndex < srcWidth; xIndex += 1 {
				c = src.At(xIndex, yIndex)
				d(xIndex, yIndex, c)
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

			var (
				c   color.Color = nil
				err error       = nil
			)

			for xIndex := 0; xIndex < srcWidth; xIndex += 1 {
				select {
				case <-ctx.Done():
					return
				default:
				}

				c = src.At(xIndex, yIndex)
				if err = d(xIndex, yIndex, c); err != nil {
					errt.Set(fmt.Errorf("pimit: delegate function failed on x=%d y=%d with: %w", xIndex, yIndex, err))
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

			var c color.Color = nil

			for xIndex := 0; xIndex < srcWidth; xIndex += 1 {
				c = src.At(xIndex, yIndex)
				c = d(xIndex, yIndex, c)
				src.Set(xIndex, yIndex, c)
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

			var (
				c   color.Color = nil
				err error       = nil
			)

			for xIndex := 0; xIndex < srcWidth; xIndex += 1 {
				select {
				case <-ctx.Done():
					return
				default:
				}

				c = src.At(xIndex, yIndex)
				c, err = d(xIndex, yIndex, c)

				if err != nil {
					errt.Set(fmt.Errorf("pimit: delegate function failed on x=%d y=%d with: %w", xIndex, yIndex, err))
					cancel()
					return
				}

				src.Set(xIndex, yIndex, c)
			}
		}(y)
	}

	wg.Wait()
	return errt.Err()
}

// Perform a parallel iteration of the pixels of the provided image. For each pixel, execute the delegate function
// allowing you to read the color and coordinates, the delegate return color will be set at the given coordinates.
// This changes will be applied to the passed image instance. The integer parameter is the number of clustes into
// which the image will be devided. Each cluster is then iterated in a separate goroutine.
func ParallelDistributedReadWrite(src draw.Image, c int, d ReadWriteDelegate) {
	if src == nil {
		panic("pimit: the provided image reference is nil")
	}

	if c <= 0 {
		panic("pimit: the provided negative or zero distribution cluster size is invalid")
	}

	if d == nil {
		panic("pimit: the provided access delegate function is nil")
	}

	width := src.Bounds().Dx()
	height := src.Bounds().Dy()

	pCount := width * height
	cCount := pCount / c
	cLeft := pCount % c

	wg := &sync.WaitGroup{}

	for offsetFactor := 0; offsetFactor < c; offsetFactor += 1 {
		cOffset := cCount * offsetFactor
		cLength := cCount
		if offsetFactor+1 == c {
			cLength += cLeft
		}

		wg.Add(1)
		go func(offset, length int) {
			defer wg.Done()

			var (
				xIndex int         = 0
				yIndex int         = 0
				c      color.Color = nil
			)

			for innerOffset := 0; innerOffset < length; innerOffset += 1 {
				xIndex = (offset + innerOffset) % width
				yIndex = (offset + innerOffset - xIndex) / width

				c = src.At(xIndex, yIndex)
				c = d(xIndex, yIndex, c)

				src.Set(xIndex, yIndex, c)
			}
		}(cOffset, cLength)
	}

	wg.Wait()
}

// Perform a parallel iteration of the pixels of the provided image. For each pixel, execute the delegate function
// allowing you to read the color and coordinates, the delegate return color will be set at the given coordinates.
// This changes will be applied to the passed image instance. The integer parameter is the number of clustes into
// which the image will be devided. Each cluster is then iterated in a separate goroutine. The iteration will break
// after the first error occurs and the error will be returned.
func ParallelDistributedReadWriteE(src draw.Image, c int, d ReadWriteErrorableDelegate) error {
	if src == nil {
		panic("pimit: the provided image reference is nil")
	}

	if c <= 0 {
		panic("pimit: the provided negative or zero distribution cluster size is invalid")
	}

	if d == nil {
		panic("pimit: the provided access delegate function is nil")
	}

	width := src.Bounds().Dx()
	height := src.Bounds().Dy()

	pCount := width * height
	cCount := pCount / c
	cLeft := pCount % c

	wg := &sync.WaitGroup{}

	errt := NewErrorTrap()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for offsetFactor := 0; offsetFactor < c; offsetFactor += 1 {
		cOffset := cCount * offsetFactor
		cLength := cCount
		if offsetFactor+1 == c {
			cLength += cLeft
		}

		wg.Add(1)
		go func(offset, length int) {
			defer wg.Done()

			var (
				xIndex int         = 0
				yIndex int         = 0
				c      color.Color = nil
				err    error       = nil
			)

			for innerOffset := 0; innerOffset < length; innerOffset += 1 {
				select {
				case <-ctx.Done():
					return
				default:
				}

				xIndex = (offset + innerOffset) % width
				yIndex = (offset + innerOffset - xIndex) / width

				c = src.At(xIndex, yIndex)

				c, err = d(xIndex, yIndex, c)
				if err != nil {
					errt.Set(fmt.Errorf("pimit: delegate function failed on x=%d y=%d with: %w", xIndex, yIndex, err))
					cancel()
					return
				}

				src.Set(xIndex, yIndex, c)
			}
		}(cOffset, cLength)
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

			var c color.Color = nil

			for xIndex := 0; xIndex < srcWidth; xIndex += 1 {
				c = src.At(xIndex, yIndex)
				c = d(xIndex, yIndex, c)

				dst.Set(xIndex, yIndex, c)
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

			var (
				c   color.Color = nil
				err error       = nil
			)

			for xIndex := 0; xIndex < srcWidth; xIndex += 1 {
				select {
				case <-ctx.Done():
					return
				default:
				}

				c = src.At(xIndex, yIndex)
				c, err = d(xIndex, yIndex, c)

				if err != nil {
					errt.Set(fmt.Errorf("pimit: delegate function failed on x=%d y=%d with: %w", xIndex, yIndex, err))
					cancel()
					return
				}

				dst.Set(xIndex, yIndex, c)
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
