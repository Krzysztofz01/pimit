package pimit

import (
	"image"
	"image/draw"
	"sync"
)

func offsetToIndex(offset, width int) (int, int) {
	if offset < 0 {
		panic("pimit: invalid negative offset provided")
	}

	if width <= 0 {
		panic("pimit: invalid non positive width provided")
	}

	xIndex := offset % width
	yIndex := (offset - xIndex) / width

	return xIndex, yIndex
}

// Perform parallel reading and editing of pixels of the passed image. For each pixel, execute the passed access
// function allowing to read the color and coordinates, which will return the color that the pixel should take
// after this operation. The changes will be applied to the passed image instance. The passed integer is the number
// of clusters into which the image will be divided. Each cluster is then iterated through a separate goroutine.
func ParallelClusterDistributedReadWrite(i draw.Image, c int, a ReadWriteAccess) {
	if i == nil {
		panic("pimit: the provided image reference is nil")
	}

	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	if c <= 0 {
		panic("pimit: the provided negative or zero cluster size is invalid")
	}

	width := i.Bounds().Dx()
	height := i.Bounds().Dy()

	pixelCount := width * height

	clusterBaseLength := pixelCount / c
	clusterRemnantLength := pixelCount % c

	wg := sync.WaitGroup{}
	wg.Add(c)

	for offsetFactor := 0; offsetFactor < c; offsetFactor += 1 {
		targetClusterOffset := offsetFactor * clusterBaseLength
		targetClusterLength := clusterBaseLength
		if offsetFactor+1 == c {
			targetClusterLength += clusterRemnantLength
		}

		func(offset, length int) {
			defer wg.Done()

			for offsetIteration := 0; offsetIteration < length; offsetIteration += 1 {
				xIndex, yIndex := offsetToIndex(offset+offsetIteration, width)

				currentColor := i.At(xIndex, yIndex)
				modifiedColor := a(xIndex, yIndex, currentColor)

				i.Set(xIndex, yIndex, modifiedColor)

			}
		}(targetClusterOffset, targetClusterLength)
	}

	wg.Wait()
}

func ParallelClusterLimitedReadWrite(i image.Image, c, l int, a ReadWriteAccess) {
	// TODO: Implementation
}
