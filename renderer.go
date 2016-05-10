package tmx

import (
	"errors"
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
	timer  *timer
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
	t := createTimer()
	t.Start()
	return &fullRenderer{m: m, canvas: c, loader: locator, tf: tf, timer: t}
}

//Render will generate a preview image of the tmx map provided
func (r *fullRenderer) Render() error {
	elapsed := r.timer.GetElapsedTime() / (1000 * 1000)
	canvas := tilemap{subject: r.m}
	canvas.renderBackground(r)
	canvas.updateIdentities(elapsed)
	err := canvas.renderLayer(r)
	if err != nil {
		return err
	}
	r.timer.UpdateTime()

	return nil
}

type tilemap struct {
	subject Map
}

type subImager interface {
	SubImage(r image.Rectangle) image.Image
}

func (t tilemap) renderBackground(r *fullRenderer) {
	color := t.subject.BackgroundColor
	r.canvas.FillRect(color, r.canvas.Bounds())
}

func (t *tilemap) updateIdentities(elapsedTime int64) {
	for i, ts := range t.subject.Tilesets {
		for j, tile := range ts.Tiles {
			if tile.Animation != nil {
				t.subject.Tilesets[i].Tiles[j].Animation.Update(elapsedTime)
			}
		}
	}
}

func (t *tilemap) renderLayer(r *fullRenderer) error {
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

			tileID := int(dt.GID - tileset.FirstGID)
			tile := tileset.GetTileByID(uint32(tileID))
			if tile != nil {
				if tile.Animation != nil {
					tileID = tile.Animation.GetFrame().TileID
				}
			}

			tx := tileID % tileset.GetNumTilesX()
			ty := tileID / tileset.GetNumTilesX()
			tx *= t.subject.TileWidth
			ty *= t.subject.TileHeight

			tileBounds := image.Rect(tx, ty, tx+t.subject.TileWidth, ty+t.subject.TileHeight)
			x := (i % l.Width) * t.subject.TileWidth
			y := (i / l.Width) * t.subject.TileHeight

			bounds := image.Rect(x, y, x+t.subject.TileWidth, y+t.subject.TileWidth)

			if relativeCanvas, ok := r.canvas.(RelativeCanvas); ok {
				relativeCanvas.Draw(tileBounds, bounds, tileset.GetFilename())
			} else if imgCanvas, ok := r.canvas.(ImageCanvas); ok {
				tilesetgfx, err := r.loader.LocateResource(filepath.Clean(t.subject.filename + tileset.Image.Source))
				if err != nil {
					return errors.New("invalid tileset path")
				}
				ptileset, ok := tilesetgfx.(subImager)
				if !ok {
					return errors.New("invalid image type given")
				}

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

				imgCanvas.Draw(tile, bounds)
			}
		}
	}

	return nil
}
