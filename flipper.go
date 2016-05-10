package tmx

import (
	"image"

	"github.com/disintegration/imaging"
)

//TileFlipper allows a custom implementation of
//flipping techniques
type TileFlipper interface {
	FlipHorizontal(image.Image) image.Image
	FlipVertical(image.Image) image.Image
	FlipDiagonal(image.Image) image.Image
}

type imagingFlipper struct {
}

func (i imagingFlipper) FlipHorizontal(tile image.Image) image.Image {
	return imaging.FlipH(tile)
}
func (i imagingFlipper) FlipVertical(tile image.Image) image.Image {
	return imaging.FlipV(tile)
}

func (i imagingFlipper) FlipDiagonal(tile image.Image) image.Image {
	return imaging.FlipH(imaging.Rotate270(tile))
}
