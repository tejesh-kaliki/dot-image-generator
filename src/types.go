package dotimagegenerator

import "image"

type BlockCount struct {
	X int
	Y int
}

func (bc BlockCount) GetImageRect() image.Rectangle {
	return image.Rect(0, 0, bc.X, bc.Y)
}
