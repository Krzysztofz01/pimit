package pimit

import "image/color"

type ReadAccess = func(xIndex, yIndex int, color color.Color)

type ReadColorAccess = func(color color.Color)

type ReadWriteAccess = func(xIndex, yIndex int, color color.Color) color.Color

type ReadWriteColorAccess = func(color color.Color) color.Color
