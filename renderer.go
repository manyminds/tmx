package tmx

import (
	"image"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"

	"image/color"
	"image/draw"
	//import for gif support
	_ "image/gif"
	//import for jpeg support
	_ "image/jpeg"
	//import for png support
	_ "image/png"
)

//Render will generate a preview image of the tmx map provided
func Render(m Map) (*image.RGBA, error) {
	bounds := image.Rect(0, 0, m.Width*m.TileWidth, m.Height*m.TileHeight)
	target := image.NewRGBA(bounds)

	canvas := tilemap{subject: m, canvas: target, tilesets: map[string]image.Image{}}
	canvas.renderBackground()
	err := canvas.renderLayer()
	if err != nil {
		return nil, err
	}

	return canvas.canvas, nil
}

type tilemap struct {
	canvas   *image.RGBA
	subject  Map
	tilesets map[string]image.Image
}

func (t tilemap) renderBackground() {
	color := getBackgroundColor(t.subject.BackgroundColor)

	draw.Draw(
		t.canvas,
		t.canvas.Rect,
		&image.Uniform{color},
		image.ZP,
		draw.Src,
	)
}

func (t *tilemap) renderLayer() error {
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

			draw.Draw(
				t.canvas,
				bounds,
				tile,
				tile.Bounds().Min,
				draw.Over,
			)
		}
	}

	return nil
}

func getBackgroundColor(background string) color.RGBA {
	d := color.RGBA{128, 128, 128, 255}
	if background == "" {
		return d
	}

	data := []byte(strings.ToLower(background))[1:]

	if len(data) != 6 {
		return d
	}

	r, err := strconv.ParseInt(string(data[0:2]), 16, 0)
	if err != nil {
		return d
	}

	g, err := strconv.ParseInt(string(data[2:4]), 16, 0)
	if err != nil {
		return d
	}

	b, err := strconv.ParseInt(string(data[4:6]), 16, 0)
	if err != nil {
		return d
	}

	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
}

func loadImage(src string) (image.Image, error) {
	file, err := os.Open(src)
	if err != nil {
		return nil, err
	}

	data, _, err := image.Decode(file)
	return data, err
}
