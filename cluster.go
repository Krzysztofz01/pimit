package pimit

import "image"

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

func ParallelClusterDistributedReadWrite(i image.Image, c int, a ReadWriteAccess) {
	// TODO: Implementation
}

func ParallelClusterLimitedReadWrite(i image.Image, c, l int, a ReadWriteAccess) {
	// TODO: Implementation
}
