package pimit

import "image/color"

// Iteration access function. Allows you to read the coordinates of a pixel and its color.
type ReadAccess = func(xIndex, yIndex int, color color.Color)

// Iteration access function. Allows you to read the pixel color.
type ReadColorAccess = func(color color.Color)

// Iteration access function. Allows you to read the coordinates of a pixel and its color. Returns the color to be set.
type ReadWriteAccess = func(xIndex, yIndex int, color color.Color) color.Color

// Iteration access function. Allows you to read the pixel color. Returns the color to be set.
type ReadWriteColorAccess = func(color color.Color) color.Color
