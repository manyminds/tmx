package tmx

import (
	"image"
	"image/color"
	"image/draw"
)

//Canvas to draw on
type Canvas interface {
	FillRect(what color.Color, where image.Rectangle)
	Bounds() image.Rectangle
}

//ImageCanvas will always draw the exact image to where
type ImageCanvas interface {
	Canvas
	Draw(what image.Image, where image.Rectangle)
}

//RelativeCanvas allows custom cropping of tiles
type RelativeCanvas interface {
	Draw(tile image.Rectangle, where image.Rectangle, f FlipMode, tileset string)
}

//ImgCanvas is a sample renderer that renders
//on a image.RGBA to generate snapshots
type ImgCanvas struct {
	target *image.RGBA
}

//Bounds returns the canvas' bounds
func (i ImgCanvas) Bounds() image.Rectangle {
	return i.target.Bounds()
}

//Draw renders on the image.RGBA surface
func (i ImgCanvas) Draw(what image.Image, where image.Rectangle) {
	draw.Draw(
		i.target,
		where,
		what,
		what.Bounds().Min,
		draw.Over,
	)
}

//FillRect draws a rectangle on the canvas
func (i ImgCanvas) FillRect(what color.Color, where image.Rectangle) {
	draw.Draw(
		i.target,
		where,
		&image.Uniform{what},
		image.ZP,
		draw.Src,
	)
}

//Image returns the image that has been drawn
func (i ImgCanvas) Image() *image.RGBA {
	return i.target
}

//NewImageCanvasFromMap returns an image canvas with correct bounds
func NewImageCanvasFromMap(m Map) *ImgCanvas {
	bounds := image.Rect(0, 0, m.Width*m.TileWidth, m.Height*m.TileHeight)
	target := image.NewRGBA(bounds)
	ic := ImgCanvas{target: target}

	return &ic
}
