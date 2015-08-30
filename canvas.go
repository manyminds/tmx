package tmx

import (
	"image"
	"image/color"
	"image/draw"
)

//Canvas to draw on
type Canvas interface {
	Draw(what image.Image, where image.Rectangle)
	FillRect(what color.Color, where image.Rectangle)
	Bounds() image.Rectangle
}

//ImageCanvas is a sample renderer that renders
//on a image.RGBA to generate snapshots
type ImageCanvas struct {
	target *image.RGBA
}

//Bounds returns the canvas' bounds
func (i ImageCanvas) Bounds() image.Rectangle {
	return i.target.Bounds()
}

//Draw renders on the image.RGBA surface
func (i ImageCanvas) Draw(what image.Image, where image.Rectangle) {
	draw.Draw(
		i.target,
		where,
		what,
		what.Bounds().Min,
		draw.Over,
	)
}

//FillRect draws a rectangle on the canvas
func (i ImageCanvas) FillRect(what color.Color, where image.Rectangle) {
	draw.Draw(
		i.target,
		where,
		&image.Uniform{what},
		image.ZP,
		draw.Src,
	)
}

//Image returns the image that has been drawn
func (i ImageCanvas) Image() *image.RGBA {
	return i.target
}

//NewImageCanvasFromMap returns an image canvas with correct bounds
func NewImageCanvasFromMap(m Map) *ImageCanvas {
	bounds := image.Rect(0, 0, m.Width*m.TileWidth, m.Height*m.TileHeight)
	target := image.NewRGBA(bounds)
	ic := ImageCanvas{target: target}

	return &ic
}
