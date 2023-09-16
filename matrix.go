package pimit

import (
	"context"
	"fmt"
	"sync"
)

// Perform a parallel iteration of the values of the provided matrix represented as a two-dimentional generic slice.
// For each entry, execute the delegate function allowing you to read the values and coordinates, the delegate return
// value will be set at the given coordinates. This changes will be applied to the passed two-dimentional slice instance.
// Each column is iterated in a separate goroutine.
func ParallelMatrixReadWrite[T any](m [][]T, d func(x, y int, value T) T) {
	if m == nil {
		panic("pimit: the provided matrix slice reference is nil")
	}

	width, height, ok := getMatrixSize(m)
	if !ok {
		panic("pimit: the provided matrix slice has inconsistent lengths")
	}

	if d == nil {
		panic("pimit: the provided access delegate function is nil")
	}

	wg := &sync.WaitGroup{}

	for x := 0; x < width; x += 1 {
		wg.Add(1)
		go func(xIndex int) {
			defer wg.Done()

			var value T = *new(T)

			for yIndex := 0; yIndex < height; yIndex += 1 {
				value = m[xIndex][yIndex]
				value = d(xIndex, yIndex, value)

				m[xIndex][yIndex] = value
			}
		}(x)
	}

	wg.Wait()
}

// Perform a parallel iteration of the values of the provided matrix represented as a two-dimentional generic slice.
// For each entry, execute the delegate function allowing you to read the values and coordinates, the delegate return
// value will be set at the given coordinates. This changes will be applied to the passed two-dimentional slice instance.
// Each column is iterated in a separate goroutine. The iteration will break after the first error occurs and the error
// will be returned.
func ParallelMatrixReadWriteE[T any](m [][]T, d func(x, y int, value T) (T, error)) error {
	if m == nil {
		panic("pimit: the provided matrix slice reference is nil")
	}

	width, height, ok := getMatrixSize(m)
	if !ok {
		panic("pimit: the provided matrix slice has inconsistent lengths")
	}

	if d == nil {
		panic("pimit: the provided access delegate function is nil")
	}

	wg := &sync.WaitGroup{}

	errt := NewErrorTrap()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for x := 0; x < width; x += 1 {
		wg.Add(1)
		go func(xIndex int) {
			defer wg.Done()

			var (
				value T     = *new(T)
				err   error = nil
			)

			for yIndex := 0; yIndex < height; yIndex += 1 {
				select {
				case <-ctx.Done():
					return
				default:
				}

				value = m[xIndex][yIndex]

				value, err = d(xIndex, yIndex, value)
				if err != nil {
					errt.Set(fmt.Errorf("pimit: delegate function failed on x=%d y=%d with: %w", xIndex, yIndex, err))
					cancel()
					return
				}

				m[xIndex][yIndex] = value
			}
		}(x)
	}

	wg.Wait()
	return errt.Err()
}

func getMatrixSize[T any](m [][]T) (int, int, bool) {
	width := len(m)
	if width == 0 {
		return 0, 0, false
	}

	height := len(m[0])
	if height == 0 {
		return 0, 0, false
	}

	for x := 1; x < width; x += 1 {
		if len(m[x]) != height {
			return 0, 0, false
		}
	}

	return width, height, true
}
