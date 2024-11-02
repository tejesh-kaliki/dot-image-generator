package dotimagegenerator

import (
	"image"
	"image/color"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
)

func GetImageBlockCount(size image.Rectangle, blockSize uint) BlockCount {
	if blockSize == 0 {
		blockSize = 1
	}
	width := size.Dx() / int(blockSize)
	height := size.Dy() / int(blockSize)
	return BlockCount{width, height}
}

func GetBlockAverageColor(pixels []color.Color, width, height uint) color.Color {
	totalPixels := width * height
	if totalPixels == 0 {
		return color.Transparent
	}

	totalRed, totalBlue, totalGreen, totalAlpha := uint(0), uint(0), uint(0), uint(0)
	for _, pixel := range pixels {
		r, g, b, a := pixel.RGBA()

		// Need to right shift by 8 and add because pixel.RGBA() does left shift
		totalRed += uint(r >> 8)
		totalGreen += uint(g >> 8)
		totalBlue += uint(b >> 8)
		totalAlpha += uint(a >> 8)
	}

	red := uint8(totalRed / totalPixels)
	green := uint8(totalGreen / totalPixels)
	blue := uint8(totalBlue / totalPixels)
	alpha := uint8(totalAlpha / totalPixels)

	return color.RGBA{red, green, blue, alpha}
}

func GetBlockPixels(imageData image.Image, offsetX, offsetY int, width, height uint) []color.Color {
	colors := make([]color.Color, width*height)
	for i := 0; i < int(width); i++ {
		for j := 0; j < int(height); j++ {
			colors[i*int(width)+j] = imageData.At(offsetX+i, offsetY+j)
		}
	}
	return colors
}

func ComputeDotImageColors(imageData image.Image, blockSize uint) image.Image {
	if blockSize == 0 {
		blockSize = 1
	}
	blockCount := GetImageBlockCount(imageData.Bounds(), blockSize)
	newImage := image.NewRGBA(blockCount.GetImageRect())

	for x := 0; x < newImage.Rect.Dx(); x++ {
		for y := 0; y < newImage.Bounds().Dy(); y++ {
			colors := GetBlockPixels(imageData, x*int(blockSize), y*int(blockSize), blockSize, blockSize)

			pixelColor := GetBlockAverageColor(colors, blockSize, blockSize)
			newImage.Set(x, y, pixelColor)
		}
	}

	return newImage
}

func drawCircleWithTopLeft(gc draw2d.GraphicContext, topLeft image.Point, radius float64, c color.Color) {
	gc.SetFillColor(c)
	arcCenterX := float64(topLeft.X) + radius
	arcCenterY := float64(topLeft.Y) + radius
	draw2dkit.Circle(gc, arcCenterX, arcCenterY, radius)
	gc.Close()
	gc.Fill()
}

func CreateDotStyleImage(dotImageData image.Image, blockSize uint) image.Image {
	imageSize := image.Rect(0, 0, dotImageData.Bounds().Dx()*int(blockSize), dotImageData.Bounds().Dy()*int(blockSize))
	dest := image.NewRGBA(imageSize)
	gc := draw2dimg.NewGraphicContext(dest)

	gc.SetFillColor(color.White)
	draw2dkit.Rectangle(gc, 0, 0, float64(dest.Bounds().Dx()), float64(dest.Bounds().Dy()))
	gc.Close()
	gc.Fill()

	gap := 2
	radius := float64(blockSize)/2 - float64(gap)
	for i := 0; i < dotImageData.Bounds().Dx(); i++ {
		for j := 0; j < dotImageData.Bounds().Dy(); j++ {
			topLeft := image.Pt(i*int(blockSize)+gap, j*int(blockSize)+gap)
			drawCircleWithTopLeft(gc, topLeft, radius, dotImageData.At(i, j))
		}
	}

	return dest
}
