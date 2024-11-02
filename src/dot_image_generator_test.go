package dotimagegenerator_test

import (
	"image"
	"image/color"
	"testing"

	dotimage "github.com/tejesh-kaliki/dot-image-generator/src"
)

var RedColor = color.RGBA{255, 0, 0, 255}

func TestGetImageBounds(t *testing.T) {
	testCases := []struct {
		Name      string
		Size      image.Rectangle
		BlockSize uint
		Want      dotimage.BlockCount
	}{
		{"empty rectangle returns empty", image.Rect(0, 0, 0, 0), 1, dotimage.BlockCount{0, 0}},
		{"zero block size defaults to 1", image.Rect(0, 0, 5, 5), 0, dotimage.BlockCount{5, 5}},
		{"square with non-divisible size", image.Rect(0, 0, 10, 10), 4, dotimage.BlockCount{2, 2}},
		{"square with divisible size", image.Rect(0, 0, 15, 15), 3, dotimage.BlockCount{5, 5}},
		{"rectangle with non-divisible size", image.Rect(0, 0, 10, 13), 4, dotimage.BlockCount{2, 3}},
		{"rectangle with non-divisible height", image.Rect(0, 0, 15, 13), 3, dotimage.BlockCount{5, 4}},
		{"rectangle with non-divisible width", image.Rect(0, 0, 13, 15), 3, dotimage.BlockCount{4, 5}},
		{"rectangle with divisible bounds", image.Rect(0, 0, 15, 18), 3, dotimage.BlockCount{5, 6}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			got := dotimage.GetImageBlockCount(testCase.Size, testCase.BlockSize)
			assertEqual(t, got, testCase.Want)
		})
	}
}

func TestDotImageGenerator(t *testing.T) {
	testCases := []struct {
		Name      string
		InputSize image.Rectangle
		BlockSize uint
		Color     color.Color
	}{
		{"plain transparent color", image.Rect(0, 0, 12, 12), 4, color.Transparent},
		{"plain white color", image.Rect(0, 0, 12, 12), 4, color.White},
		{"plain red color", image.Rect(0, 0, 12, 12), 4, RedColor},
		{"empty image does not fail", image.Rect(0, 0, 0, 0), 4, RedColor},
		{"size less than block size", image.Rect(0, 0, 3, 3), 4, RedColor},
		{"empty block size does not fail", image.Rect(0, 0, 4, 4), 0, RedColor},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			imageData := image.NewRGBA(testCase.InputSize)
			for x := 0; x < testCase.InputSize.Dx(); x++ {
				for y := 0; y < testCase.InputSize.Dy(); y++ {
					imageData.Set(x, y, testCase.Color)
				}
			}

			dotImageData := dotimage.ComputeDotImageColors(imageData, testCase.BlockSize)

			bc := dotimage.GetImageBlockCount(testCase.InputSize, testCase.BlockSize)
			for i := 0; i < bc.X; i++ {
				for j := 0; j < bc.Y; j++ {
					assertColorEqual(t, dotImageData.At(i, j), testCase.Color)
				}
			}
		})
	}
}

func TestGroupAverageColor(t *testing.T) {
	testCases := []struct {
		Name   string
		Width  uint
		Height uint
		Colors []color.Color
		Want   color.Color
	}{
		{"0 width gives transparent", 0, 2, []color.Color{}, color.Transparent},
		{"0 height gives transparent", 2, 0, []color.Color{}, color.Transparent},
		{"single block returns same color: Black", 1, 1, []color.Color{color.Black}, color.Black},
		{"single block returns same color: Transparent", 1, 1, []color.Color{color.Transparent}, color.Transparent},
		{"single block returns same color: Red", 1, 1, []color.Color{RedColor}, RedColor},
		{"white and black gives grey", 2, 1, []color.Color{color.White, color.Black}, color.Gray{127}},
		{"red and white gives light red", 2, 1, []color.Color{RedColor, color.White}, color.RGBA{255, 127, 127, 255}},
		{"4 same colors gives same color", 2, 2, []color.Color{RedColor, RedColor, RedColor, RedColor}, RedColor},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			got := dotimage.GetBlockAverageColor(testCase.Colors, testCase.Width, testCase.Height)
			assertColorEqual(t, got, testCase.Want)
		})
	}
}

func assertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()

	if got != want {
		t.Errorf("Values does not match: got %+v want %+v", got, want)
	}
}

func assertColorEqual(t *testing.T, got, want color.Color) {
	t.Helper()

	r1, g1, b1, a1 := got.RGBA()
	r2, g2, b2, a2 := want.RGBA()

	if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
		t.Errorf("Values does not match: got %+v want %+v", got, want)
	}
}
