package pimit

import (
	"context"
	"fmt"
	"image"
	"sync"
)

type (
	NrgbaReadDelegate               = func(x, y int, r, g, b, a uint8)
	NrgbaReadErrorableDelegate      = func(x, y int, r, g, b, a uint8) error
	NrgbaReadWriteDelegate          = func(x, y int, r, g, b, a uint8) (uint8, uint8, uint8, uint8)
	NrgbaReadWriteErrorableDelegate = func(x, y int, r, g, b, a uint8) (uint8, uint8, uint8, uint8, error)
)

// Perform a parallel iteration of the pixels of the provided NRGBA image. For each pixel, execute the delegate function
// allowing you to read the color (R, G, B and A as uint8) and coordinates. Each row is iterated in a separate goroutine.
func ParallelNrgbaRead(src *image.NRGBA, d NrgbaReadDelegate) {
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

			var (
				baseIndex  int   = 4 * yIndex * srcWidth
				r, g, b, a uint8 = 0, 0, 0, 0
			)

			for xIndex := 0; xIndex < srcWidth; xIndex += 1 {
				r = src.Pix[baseIndex+0]
				g = src.Pix[baseIndex+1]
				b = src.Pix[baseIndex+2]
				a = src.Pix[baseIndex+3]

				d(xIndex, yIndex, r, g, b, a)
				baseIndex += 4
			}
		}(y)
	}

	wg.Wait()
}

// Perform a parallel iteration of the pixels of the provided NRGBA image. For each pixel, execute the delegate function
// allowing you to read the color (R, G, B and A as uint8) and coordinates. Each row is iterated in a separate goroutine.
// The iteration will break after the first error occurs and the error will be returned.
func ParallelNrgbaReadE(src *image.NRGBA, d NrgbaReadErrorableDelegate) error {
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
				baseIndex  int   = 4 * yIndex * srcWidth
				r, g, b, a uint8 = 0, 0, 0, 0
				err        error = nil
			)

			for xIndex := 0; xIndex < srcWidth; xIndex += 1 {
				select {
				case <-ctx.Done():
					return
				default:
				}

				r = src.Pix[baseIndex+0]
				g = src.Pix[baseIndex+1]
				b = src.Pix[baseIndex+2]
				a = src.Pix[baseIndex+3]

				if err = d(xIndex, yIndex, r, g, b, a); err != nil {
					errt.Set(fmt.Errorf("pimit: delegate function failed on x=%d y=%d with: %w", xIndex, yIndex, err))
					cancel()
					return
				}

				baseIndex += 4
			}
		}(y)
	}

	wg.Wait()
	return errt.Err()
}

// Perform a parallel iteration of the pixels of the provided NRGBA image. For each pixel, execute the delegate function
// allowing you to read the color (R, G, B and A as uint8) and coordinates, the delegate return color will be set at the
// given coordinates. This changes will be applied to the passed image instance. Consider using ParallelReadWriteNew if
// you want to avoid changes to the original image at the expense of additional allocations. Each row is iterated in a
// separate goroutine.
func ParallelNrgbaReadWrite(src *image.NRGBA, d NrgbaReadWriteDelegate) {
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

			var (
				baseIndex  int   = 4 * yIndex * srcWidth
				r, g, b, a uint8 = 0, 0, 0, 0
			)

			for xIndex := 0; xIndex < srcWidth; xIndex += 1 {
				r = src.Pix[baseIndex+0]
				g = src.Pix[baseIndex+1]
				b = src.Pix[baseIndex+2]
				a = src.Pix[baseIndex+3]

				r, g, b, a = d(xIndex, yIndex, r, g, b, a)

				src.Pix[baseIndex+0] = r
				src.Pix[baseIndex+1] = g
				src.Pix[baseIndex+2] = b
				src.Pix[baseIndex+3] = a

				baseIndex += 4
			}
		}(y)
	}

	wg.Wait()
}

// Perform a parallel iteration of the pixels of the provided NRGBA image. For each pixel, execute the delegate function
// allowing you to read the color (R, G, B and A as uint8) and coordinates, the delegate return color will be set at the
// given coordinates. This changes will be applied to the passed image instance. Consider using ParallelReadWriteNewE if
// you want to avoid changes to the original image at the expense of additional allocations. Each row is iterated in a
// separate goroutine. The iteration will break after the first error occurs and the error will be returned.
func ParallelNrgbaReadWriteE(src *image.NRGBA, d NrgbaReadWriteErrorableDelegate) error {
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
				baseIndex  int   = 4 * yIndex * srcWidth
				r, g, b, a uint8 = 0, 0, 0, 0
				err        error = nil
			)

			for xIndex := 0; xIndex < srcWidth; xIndex += 1 {
				select {
				case <-ctx.Done():
					return
				default:
				}

				r = src.Pix[baseIndex+0]
				g = src.Pix[baseIndex+1]
				b = src.Pix[baseIndex+2]
				a = src.Pix[baseIndex+3]

				r, g, b, a, err = d(xIndex, yIndex, r, g, b, a)

				if err != nil {
					errt.Set(fmt.Errorf("pimit: delegate function failed on x=%d y=%d with: %w", xIndex, yIndex, err))
					cancel()
					return
				}

				src.Pix[baseIndex+0] = r
				src.Pix[baseIndex+1] = g
				src.Pix[baseIndex+2] = b
				src.Pix[baseIndex+3] = a

				baseIndex += 4
			}
		}(y)
	}

	wg.Wait()
	return errt.Err()
}

// Perform a parallel iteration of the pixels of the provided NRGBA image. For each pixel, execute the delegate function
// allowing you to read the color (R, G, B and A as uint8) and coordinates, the delegate return color will be set at the
// given coordinates. This changes will be applied to a new image instance which internaly uses the NRGBA color space
// and is returned by the function. Each row is iterated in a separate goroutine.
func ParallelNrgbaReadWriteNew(src *image.NRGBA, d NrgbaReadWriteDelegate) *image.NRGBA {
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

			var (
				baseIndex  int   = 4 * yIndex * srcWidth
				r, g, b, a uint8 = 0, 0, 0, 0
			)

			for xIndex := 0; xIndex < srcWidth; xIndex += 1 {
				r = src.Pix[baseIndex+0]
				g = src.Pix[baseIndex+1]
				b = src.Pix[baseIndex+2]
				a = src.Pix[baseIndex+3]

				r, g, b, a = d(xIndex, yIndex, r, g, b, a)

				dst.Pix[baseIndex+0] = r
				dst.Pix[baseIndex+1] = g
				dst.Pix[baseIndex+2] = b
				dst.Pix[baseIndex+3] = a

				baseIndex += 4
			}
		}(y)
	}

	wg.Wait()
	return dst
}

// Perform a parallel iteration of the pixels of the provided NRGBA image. For each pixel, execute the delegate function
// allowing you to read the color (R, G, B and A as uint8) and coordinates, the delegate return color will be set at the
// given coordinates. This changes will be applied to a new image instance which internaly uses the NRGBA color space
// and is returned by the function. Each row is iterated in a separate goroutine. The iteration will break after the first
// error occurs and the error will be returned.
func ParallelNrgbaReadWriteNewE(src *image.NRGBA, d NrgbaReadWriteErrorableDelegate) (*image.NRGBA, error) {
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
				baseIndex  int   = 4 * yIndex * srcWidth
				r, g, b, a uint8 = 0, 0, 0, 0
				err        error = nil
			)

			for xIndex := 0; xIndex < srcWidth; xIndex += 1 {
				select {
				case <-ctx.Done():
					return
				default:
				}

				r = src.Pix[baseIndex+0]
				g = src.Pix[baseIndex+1]
				b = src.Pix[baseIndex+2]
				a = src.Pix[baseIndex+3]

				r, g, b, a, err = d(xIndex, yIndex, r, g, b, a)

				if err != nil {
					errt.Set(fmt.Errorf("pimit: delegate function failed on x=%d y=%d with: %w", xIndex, yIndex, err))
					cancel()
					return
				}

				dst.Pix[baseIndex+0] = r
				dst.Pix[baseIndex+1] = g
				dst.Pix[baseIndex+2] = b
				dst.Pix[baseIndex+3] = a

				baseIndex += 4
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
