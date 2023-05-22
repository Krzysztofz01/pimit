package pimit

import "image/color"

type accessFuncConstraint interface {
	ReadAccess | ReadColorAccess | ReadWriteAccess | ReadWriteColorAccess | ReadAccessE | ReadColorAccessE | ReadWriteAccessE | ReadWriteColorAccessE
}

// Iteration access function. Allows you to read the coordinates of a pixel and its color.
type ReadAccess = func(xIndex, yIndex int, color color.Color)

// Iteration access function. Allows you to read the pixel color.
type ReadColorAccess = func(color color.Color)

// Iteration access function. Allows you to read the coordinates of a pixel and its color. Returns the color to be set.
type ReadWriteAccess = func(xIndex, yIndex int, color color.Color) color.Color

// Iteration access function. Allows you to read the pixel color. Returns the color to be set.
type ReadWriteColorAccess = func(color color.Color) color.Color

// Iteration access function. Allows you to read the coordinates of a pixel and its color. Errors can be passed down to the iterator.
type ReadAccessE = func(xIndex, yIndex int, color color.Color) error

// Iteration access function. Allows you to read the pixel color. Errors can be passed down to the iterator.
type ReadColorAccessE = func(color color.Color) error

// Iteration access function. Allows you to read the coordinates of a pixel and its color. Returns the color to be set. Errors can be passed down to the iterator.
type ReadWriteAccessE = func(xIndex, yIndex int, color color.Color) (color.Color, error)

// Iteration access function. Allows you to read the pixel color. Returns the color to be set. Errors can be passed down to the iterator.
type ReadWriteColorAccessE = func(color color.Color) (color.Color, error)
