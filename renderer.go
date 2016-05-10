package tmx

import (
	"image"
	"path/filepath"
)

//Renderer renders
//the given Map on a provided Canvas
type Renderer interface {
	Render() error
}

type fullRenderer struct {
	canvas Canvas
	m      Map
	loader ResourceLocator
	tf     TileFlipper
}

//NewRenderer lets you draw the map on a custom canvas
//with a default FilesystemLocator
func NewRenderer(m Map, c Canvas) Renderer {
	return NewRendererWithResourceLocator(m, c, NewLazyResourceLocator(FilesystemLocator{}))
}

//NewRendererWithResourceLocator return a new renderer
func NewRendererWithResourceLocator(m Map, c Canvas, locator ResourceLocator) Renderer {
	return NewRendererWithResourceLocatorAndTileFlipper(m, c, locator, &imagingFlipper{})
}

//NewRendererWithResourceLocatorAndTileFlipper allows you to specify
//a custom canvas, locator and TileFlipper
func NewRendererWithResourceLocatorAndTileFlipper(
	m Map,
	c Canvas,
	locator ResourceLocator,
	tf TileFlipper,
) Renderer {
	return &fullRenderer{m: m, canvas: c, loader: locator, tf: tf}
}

//Render will generate a preview image of the tmx map provided
func (r fullRenderer) Render() error {
	canvas := tilemap{subject: r.m}
	canvas.renderBackground(r)
	err := canvas.renderLayer(r)
	if err != nil {
		return err
	}

	return nil
}

type tilemap struct {
	subject Map
}

type subImager interface {
	SubImage(r image.Rectangle) image.Image
}

func (t tilemap) renderBackground(r fullRenderer) {
	color := t.subject.BackgroundColor
	r.canvas.FillRect(color, r.canvas.Bounds())
}

func (t *tilemap) renderLayer(r fullRenderer) error {
	for _, l := range t.subject.Layers {
		if !l.IsVisible() {
			continue
		}

		for i, dt := range l.Data.DataTiles {
			tileset, err := t.subject.GetTilesetForGID(dt.GID)
			if err != nil {
				continue
			}

			if tileset == nil {
				continue
			}

			tx := int(dt.GID-tileset.FirstGID) % tileset.GetNumTilesX()
			ty := int(dt.GID-tileset.FirstGID) / tileset.GetNumTilesX()
			tx *= t.subject.TileWidth
			ty *= t.subject.TileHeight

			tilesetgfx, err := r.loader.LocateResource(filepath.Clean(t.subject.filename + tileset.Image.Source))
			if err != nil {
				panic("invalid tileset path")
			}

			ptileset, ok := tilesetgfx.(subImager)
			if !ok {
				panic("invalid image type given")
			}

			tileBounds := image.Rect(tx, ty, tx+t.subject.TileWidth, ty+t.subject.TileHeight)
			tile := ptileset.SubImage(tileBounds)

			if dt.DiagonalFlip {
				tile = r.tf.FlipDiagonal(tile)
			}

			if dt.HorizontalFlip {
				tile = r.tf.FlipHorizontal(tile)
			}

			if dt.VerticalFlip {
				tile = r.tf.FlipVertical(tile)
			}

			x := (i % l.Width) * t.subject.TileWidth
			y := (i / l.Height) * t.subject.TileHeight

			bounds := image.Rect(x, y, x+t.subject.TileWidth, y+t.subject.TileWidth)

			r.canvas.Draw(tile, bounds)
		}
	}

	return nil
}
