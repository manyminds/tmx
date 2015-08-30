package tmx

import (
	"image"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"

	//import for gif support
	_ "image/gif"
	//import for jpeg support
	_ "image/jpeg"
	//import for png support
	_ "image/png"
)

//Renderer renders a tmx to the given canvas
type Renderer struct {
	canvas Canvas
	m      Map
}

//NewRendererWithCanvas lets you draw the map on a custom canvas
func NewRendererWithCanvas(m Map, c Canvas) *Renderer {
	return &Renderer{m: m, canvas: c}
}

//Render will generate a preview image of the tmx map provided
func (r Renderer) Render() error {
	canvas := tilemap{subject: r.m, tilesets: map[string]image.Image{}}
	canvas.renderBackground(r)
	err := canvas.renderLayer(r)
	if err != nil {
		return err
	}

	return nil
}

type tilemap struct {
	subject  Map
	tilesets map[string]image.Image
}

func (t tilemap) renderBackground(r Renderer) {
	color := t.subject.BackgroundColor
	r.canvas.FillRect(color, r.canvas.Bounds())
}

func (t *tilemap) renderLayer(r Renderer) error {
	for _, tileset := range t.subject.Tilesets {
		path := tileset.Image.Source

		tileset, err := loadImage(filepath.Clean(t.subject.filename + path))

		if err != nil {
			return err
		}

		t.tilesets[path] = tileset
	}

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

			tilesetgfx, found := t.tilesets[tileset.Image.Source]
			if !found {
				panic("invalid tileset path")
			}

			ptileset, ok := tilesetgfx.(*image.Paletted)
			if !ok {
				panic("invalid image type given")
			}

			tileBounds := image.Rect(tx, ty, tx+t.subject.TileWidth, ty+t.subject.TileHeight)
			tile := ptileset.SubImage(tileBounds)

			if dt.DiagonalFlip {
				tile = imaging.Rotate270(tile)
				tile = imaging.FlipH(tile)
			}

			if dt.HorizontalFlip {
				tile = imaging.FlipH(tile)
			}

			if dt.VerticalFlip {
				tile = imaging.FlipV(tile)
			}

			x := (i % l.Width) * t.subject.TileWidth
			y := (i / l.Height) * t.subject.TileHeight

			bounds := image.Rect(x, y, x+t.subject.TileWidth, y+t.subject.TileWidth)

			r.canvas.Draw(tile, bounds)
		}
	}

	return nil
}

func loadImage(src string) (image.Image, error) {
	file, err := os.Open(src)
	if err != nil {
		return nil, err
	}

	data, _, err := image.Decode(file)
	return data, err
}
