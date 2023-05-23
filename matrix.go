package pimit

import (
	"fmt"
	"sync"
)

// Perform parallel reading and editing of the value of the passed generic matrix, representing the image or
// its properties. For each pixel, execute the passed access function allowing to read the value and coordinates,
// which will return the value that the given matrix position should take after this operation. The changes
// will be applied to the passed matrix slice instance. Every column is iterated in a separate goroutine.
func ParallelMatrixReadWrite[T any](m [][]T, a func(xIndex, yIndex int, value T) T) {
	if m == nil {
		panic("pimit: the provided matrix slice reference is nil")
	}

	if !isMatrixSizeValid(m) {
		panic("pimit: the provided matrix slice reference has a invalid size")
	}

	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	width, height := getMatrixDimensions(m)

	wg := sync.WaitGroup{}
	wg.Add(width)

	for x := 0; x < width; x += 1 {
		go func(xIndex int) {
			defer wg.Done()

			for yIndex := 0; yIndex < height; yIndex += 1 {
				value := m[xIndex][yIndex]
				modifiedValue := a(xIndex, yIndex, value)

				m[xIndex][yIndex] = modifiedValue
			}
		}(x)
	}

	wg.Wait()
}

// Perform parallel reading and editing of the value of the passed generic matrix, representing the image or
// its properties. For each pixel, execute the passed access function allowing to read the value and coordinates,
// which will return the value that the given matrix position should take after this operation. The changes
// will be applied to the passed matrix slice instance. Every column is iterated in a separate goroutine. Errors
// that occur in the function will be caught and the first one will be returned by the function.
func ParallelMatrixReadWriteE[T any](m [][]T, a func(xIndex, yIndex int, value T) (T, error)) error {
	if m == nil {
		panic("pimit: the provided matrix slice reference is nil")
	}

	if !isMatrixSizeValid(m) {
		panic("pimit: the provided matrix slice reference has a invalid size")
	}

	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	width, height := getMatrixDimensions(m)

	wg := sync.WaitGroup{}
	wg.Add(width)

	iterationErrors := make(chan error, width)

	for x := 0; x < width; x += 1 {
		go func(xIndex int, errCh chan error) {
			defer wg.Done()

			for yIndex := 0; yIndex < height; yIndex += 1 {
				value := m[xIndex][yIndex]
				modifiedValue, err := a(xIndex, yIndex, value)
				if err != nil {
					errCh <- fmt.Errorf("pimit: access function failed on x:%d y:%d: %w", xIndex, yIndex, err)
					return
				}

				m[xIndex][yIndex] = modifiedValue
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

func isMatrixSizeValid[T any](m [][]T) bool {
	width := len(m)
	if width == 0 {
		return false
	}

	height := len(m[0])
	if height == 0 {
		return false
	}

	for x := 1; x < width; x += 1 {
		if len(m[x]) != height {
			return false
		}
	}

	return true
}

func getMatrixDimensions[T any](m [][]T) (int, int) {
	return len(m), len(m[0])
}
