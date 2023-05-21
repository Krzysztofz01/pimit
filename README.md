# (P)arallel (Im)age (It)eration

[![Go Reference](https://pkg.go.dev/badge/github.com/Krzysztofz01/pimit.svg)](https://pkg.go.dev/github.com/Krzysztofz01/pimit)
[![Go Report Card](https://goreportcard.com/badge/github.com/Krzysztofz01/pimit)](https://goreportcard.com/report/github.com/Krzysztofz01/pimit)
![GitHub](https://img.shields.io/github/license/Krzysztofz01/pimit)
![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/Krzysztofz01/pimit?include_prereleases)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/Krzysztofz01/pimit)

A minimalist library that adds functionality to wrap logic for concurrent iteration over images. The library contains various types of functions that allow parallel iteration over images for reading or writing. It is also possible to indicate whether the concurrent iteration is to be performed against columns, rows or clusters. Thanks to this library, it is possible to clean up the code fragments in which you perform image processing by separating the parts of the code related to iteration and concurrent operations. Thanks to the concurrent operation of these iterators, it is possible to increase the performance of the algorithm at the expense of a small increase in memory consumption.

## Installation
```
go get -u github.com/Krzysztofz01/pimit
```

## Documentation

[https://pkg.go.dev/github.com/Krzysztofz01/pimit](https://pkg.go.dev/github.com/Krzysztofz01/pimit)

## Example
A quick example of iteration over the image with current color printing. The "After" is additionaly handling all the concurrency work and is performing 2x faster on average.

### Before
```go
image := CreateExampleImage()
height := image.Bounds().Dy()
width := image.Bounds().Dx()

for y := 0; y < height; y += 1 {
    for x := 0; x < width; x += 1 {
        color := image.At(xIndex, yIndex)
        fmt.Print(color)
    }
} 
```

### After
```go
image := CreateExampleImage()

ParallelColumnColorRead(image, func(c color.Color) {
    fmt.Print(c)
})
```

A quick example of making the picture black and white using parallel row iteration.
```go
image := CreateExampleImage()

ParallelRowColorReadWrite(img, func(c color.Color) color.Color {
    rgb, _ := c.(color.RGBA)
    value := uint8(0.299*float32(rgb.R) + 0.587*float32(rgb.G) + 0.114*float32(rgb.B))

    return color.RGBA{value, value, value, 0xff}
})
```