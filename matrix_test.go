package pimit

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestGetMatrixSizeShouldTellIfValidAndReturnSize(t *testing.T) {
	defer goleak.VerifyNone(t)

	cases := []struct {
		matrixWidth    int
		matrixHeight   int
		expectedOk     bool
		expectedWidth  int
		expectedHeight int
	}{
		{0, 0, false, 0, 0},
		{0, 1, false, 0, 0},
		{1, 0, false, 0, 0},
		{1, 1, false, 0, 0},
		{1, 2, false, 0, 0},
		{1, 1, true, 1, 1},
		{2, 1, true, 2, 1},
		{1, 2, true, 1, 2},
		{3, 4, true, 3, 4},
		{4, 4, true, 4, 4},
	}

	for _, c := range cases {
		matrix := make([][]int, c.matrixWidth)
		for mx := 0; mx < c.expectedWidth; mx += 1 {
			matrix[mx] = make([]int, c.matrixHeight)
		}

		actualWidth, actualHeight, actualOk := getMatrixSize(matrix)

		assert.Equal(t, c.expectedOk, actualOk)
		assert.Equal(t, c.expectedWidth, actualWidth)
		assert.Equal(t, c.expectedHeight, actualHeight)
	}
}

func TestGetMatrixSizeShouldReturnFalseOnInconsistentSizes(t *testing.T) {
	defer goleak.VerifyNone(t)

	matrix := make([][]bool, 2)
	matrix[0] = make([]bool, 2)
	matrix[1] = make([]bool, 1)

	_, _, ok := getMatrixSize(matrix)

	assert.False(t, ok)
}

func TestParallelMatrixReadWriteShouldPanicOnNilMatrix(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelMatrixReadWrite(nil, func(_, _ int, value bool) bool {
			return value
		})
	})
}

func TestParallelMatrixReadWriteShouldPanicOnInvalidMatrix(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		matrix := mockSpecificMatrix(1, 0, true)

		ParallelMatrixReadWrite(matrix, func(_, _ int, value bool) bool {
			return value
		})
	})
}

func TestParallelMatrixReadWriteShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	matrix := mockSpecificMatrix(2, 3, true)

	assert.Panics(t, func() {
		ParallelMatrixReadWrite(matrix, nil)
	})
}

func TestParallelMatrixReadWriteShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	matWidth := 2
	matHeight := 3

	matrix := mockSpecificMatrix(matWidth, matHeight, true)

	ParallelMatrixReadWrite(matrix, func(_, _ int, value bool) bool {
		assert.True(t, value)
		return !value
	})

	expectedMatrix := mockSpecificMatrix(matWidth, matHeight, false)

	for x := 0; x < matWidth; x += 1 {
		for y := 0; y < matHeight; y += 1 {
			assert.Equal(t, expectedMatrix[x][y], matrix[x][y])
		}
	}
}

func TestParallelMatrixReadWriteShouldAccessPixelsOnce(t *testing.T) {
	defer goleak.VerifyNone(t)

	matWidth := 2
	matHeight := 3

	matrix := mockSpecificMatrix(matWidth, matHeight, true)

	ParallelMatrixReadWrite(matrix, func(_, _ int, value bool) bool {
		return !value
	})

	expectedMatrix := mockSpecificMatrix(matWidth, matHeight, false)

	for x := 0; x < matWidth; x += 1 {
		for y := 0; y < matHeight; y += 1 {
			assert.Equal(t, expectedMatrix[x][y], matrix[x][y])
		}
	}
}

func TestParallelMatrixReadWriteEShouldPanicOnNilMatrix(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ParallelMatrixReadWriteE(nil, func(_, _ int, value bool) (bool, error) {
			return value, nil
		})
	})
}

func TestParallelMatrixReadWriteEShouldPanicOnInvalidMatrix(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		matrix := mockSpecificMatrix(1, 0, true)

		ParallelMatrixReadWriteE(matrix, func(_, _ int, value bool) (bool, error) {
			return value, nil
		})
	})
}

func TestParallelMatrixReadWriteEShouldPanicOnNilAccessFunc(t *testing.T) {
	defer goleak.VerifyNone(t)

	matrix := mockSpecificMatrix(2, 3, true)

	assert.Panics(t, func() {
		ParallelMatrixReadWriteE(matrix, nil)
	})
}

func TestParallelMatrixReadWriteEShouldCorrectlyIterate(t *testing.T) {
	defer goleak.VerifyNone(t)

	matWidth := 2
	matHeight := 3

	matrix := mockSpecificMatrix(matWidth, matHeight, true)

	err := ParallelMatrixReadWriteE(matrix, func(_, _ int, value bool) (bool, error) {
		assert.True(t, value)
		return !value, nil
	})

	assert.Nil(t, err)

	expectedMatrix := mockSpecificMatrix(matWidth, matHeight, false)

	for x := 0; x < matWidth; x += 1 {
		for y := 0; y < matHeight; y += 1 {
			assert.Equal(t, expectedMatrix[x][y], matrix[x][y])
		}
	}
}

func TestParallelMatrixReadWriteEShouldAccessPixelsOnce(t *testing.T) {
	defer goleak.VerifyNone(t)

	matWidth := 2
	matHeight := 3

	matrix := mockSpecificMatrix(matWidth, matHeight, true)

	err := ParallelMatrixReadWriteE(matrix, func(_, _ int, value bool) (bool, error) {
		return !value, nil
	})

	assert.Nil(t, err)

	expectedMatrix := mockSpecificMatrix(matWidth, matHeight, false)

	for x := 0; x < matWidth; x += 1 {
		for y := 0; y < matHeight; y += 1 {
			assert.Equal(t, expectedMatrix[x][y], matrix[x][y])
		}
	}
}

func TestParallelMatrixReadWriteEShouldReturnErrorOnAccessError(t *testing.T) {
	defer goleak.VerifyNone(t)

	matrix := mockSpecificMatrix(2, 3, true)

	err := ParallelMatrixReadWriteE(matrix, func(_, _ int, value bool) (bool, error) {
		return value, errors.New("pimit-test: test errror")
	})

	assert.NotNil(t, err)
}
