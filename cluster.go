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

func ParallelClusterDistributedReadWrite(i draw.Image, c int, a ReadWriteAccess) {
	if i == nil {
		panic("pimit: the provided image reference is nil")
	}

	if a == nil {
		panic("pimit: the provided access function is nil")
	}

	width := i.Bounds().Dx()
	height := i.Bounds().Dy()

	pixelCount := width * height

	clusterLengthBase := pixelCount / c
	clusterLengthRemaining := pixelCount % c

	wg := sync.WaitGroup{}
	wg.Add(c)

	for offsetMultiplier := 0; offsetMultiplier < c; offsetMultiplier += 1 {
		targetClusterOffset := offsetMultiplier * clusterLengthBase
		targetClusterLength := clusterLengthBase
		if offsetMultiplier+1 == c {
			targetClusterLength = clusterLengthRemaining
		}

		go func(startingOffset, clusterLength int) {
			defer wg.Done()

			for innerOffset := 0; innerOffset < clusterLength; innerOffset += 1 {
				currentOffset := innerOffset + startingOffset

				xIndex, yIndex := offsetToIndex(currentOffset+startingOffset, width)

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
