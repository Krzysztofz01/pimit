# (P)arallel (Im)age (It)eration

[![Go Reference](https://pkg.go.dev/badge/github.com/Krzysztofz01/pimit.svg)](https://pkg.go.dev/github.com/Krzysztofz01/pimit)
[![Go Report Card](https://goreportcard.com/badge/github.com/Krzysztofz01/pimit)](https://goreportcard.com/report/github.com/Krzysztofz01/pimit)
![GitHub](https://img.shields.io/github/license/Krzysztofz01/pimit)
![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/Krzysztofz01/pimit?include_prereleases)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/Krzysztofz01/pimit)

A minimalist library that adds concurrent image pixel iteration functionality wrapped in a convenient and intuitive API. The main idea is that the functions take as a parameter the image whose pixels are to be iterated over, and a function, which is a delegate, that will be executed on each pixel. The library contains a number of functions, some allow only reading, some allow reading as well as editing. Some of the functions make changes to the original image and some create a new instance. Error propagation and iteration interrupts are also possible. In general, for each row of pixels, iteration takes place in a separate goroutine, but it is also possible to choose a function that allows you to select the appropriate number of goroutines. Pimit also allows you to perform iterations on matrices, which are represented as two-dimensional generic slices `[][]T`.

The library includes a general API that works on universal types like `image.Image` and `color.Color`, it is more convenient but less efficient and performs more memory allocations, despite this, the operation is expected to perform **2x faster** on average.

The library also includes an API for specific color spaces, it works on image pointers and primitive types e.g.: `*image.RGBA` and `uint8`, these APIs are more verbose but also much more efficient and avoid additional memory allocations. The operation is expected to perform **10x faster** on average.

## Installation
```
go get -u github.com/Krzysztofz01/pimit
```

## Documentation

[https://pkg.go.dev/github.com/Krzysztofz01/pimit](https://pkg.go.dev/github.com/Krzysztofz01/pimit)

## Examples

### Read example: Count all the black pixels in the image.

#### Without **pimit** (no concurrecy and more code)
```golang
func CountBlackPixel(i image.Image) int {
    height := image.Bounds().Dy()
    width := image.Bounds().Dx()
    count := 0
    
    for y := 0; y < height; y += 1 {
        for x := 0; x < width; x += 1 {
            color := image.At(xIndex, yIndex)
            if color == color.Black {
                count += 1
            }
        }
    }
    
    return count
}
```

#### With **pimit**, using the general API
```golang
func CountBlackPixel(i image.Image) int {
    var count int32 = 0
    
    pimit.ParallelRead(i, func(x, y int, c color.Color) {
        if c == color.Black {
            atomic.AddInt32(&count, 1)
        }
    })
    
    return int(atomic.LoadInt32(count))
}
```

#### With **pimit**, using the specific API
```golang
func CountBlackPixel(i *image.RGBA) int {
    var count int32 = 0
    
    pimit.ParallelRgbaRead(i, func(x, y int, r, g, b, a uint8) {
        if 0 == a && a == b && b == c {
            atomic.AddInt32(&count, 1)
        }
    })
    
    return int(count)
}
```

### Read/write example: Image to grayscale converting.

#### Without **pimit** (no concurrecy and more code)
```golang
func ToGrayscale(i draw.Image) (image.Image) {
    height := image.Bounds().Dy()
    width := image.Bounds().Dx()
 
    for y := 0; y < height; y += 1 {
        for x := 0; x < width; x += 1 {
            c := image.At(xIndex, yIndex)
            rgba, _ := c.(color.RGBA)

            y := uint8(0.299*float32(rgb.R) + 0.587*float32(rgb.G) + 0.114*float32(rgb.B))
                        
            i.Set(xIndex, yIndex, color.RGBA{
                R: y,
                G: y,
                B: y,
                A: rgba.A,
            })
        }
    }
    
    return i
}
```

#### With **pimit**, using the general API
```golang
func ToGrayscale(i draw.Image) image.Image  {
    pimit.ParallelReadWrite(i, func(x, y int, c color.Color) color.Color {
        rgba, _ := c.(color.RGBA)
        y := uint8(0.299*float32(rgb.R) + 0.587*float32(rgb.G) + 0.114*float32(rgb.B))
                        
        return color.RGBA{
            R: y,
            G: y,
            B: y,
            A: rgba.A,
        })
    })
    
    return i
}
```

#### With **pimit**, using the specific API
```golang
func ToGrayscale(i *image.RGBA) image.Image {
    pimit.ParallelRgbaReadWrite(i, func(x, y int, r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
        y := uint8(0.299*float32(r) + 0.587*float32(g) + 0.114*float32(b))
        return y, y, y, a
    })
    
    return i
}
```