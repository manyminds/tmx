package tmx

// Tileset entry describes a complete tileset
type Tileset struct {
	FirstGID   GID        `xml:"firstgid,attr"`
	Source     string     `xml:"source,attr"`
	Name       string     `xml:"name,attr"`
	TileWidth  int        `xml:"tilewidth,attr"`
	TileHeight int        `xml:"tileheight,attr"`
	Spacing    int        `xml:"spacing,attr"`
	Margin     int        `xml:"margin,attr"`
	Properties []Property `xml:"properties>property"`
	Image      Image      `xml:"image"`
	Tiles      []Tile     `xml:"tile"`
}

//GetTileByID returns special tile information
func (t Tileset) GetTileByID(tileID uint32) *Tile {
	for _, t := range t.Tiles {
		if t.ID == tileID {
			return &t
		}
	}

	return nil
}

//GetFilename returns the filename for this tileset
func (t Tileset) GetFilename() string {
	if t.Source != "" {
		return t.Source
	}

	return t.Image.Source
}

//GetNumTiles returns the number of tiles of this tileset
func (t Tileset) GetNumTiles() int {
	return t.GetNumTilesX() * t.GetNumTilesY()
}

//GetNumTilesX returns the number of tiles in x direction
func (t Tileset) GetNumTilesX() int {
	return t.Image.Width / t.TileWidth
}

//GetNumTilesY returns the number of tiles in y direction
func (t Tileset) GetNumTilesY() int {
	return t.Image.Height / t.TileHeight
}

// Image refers to the image of one tile or the tileset
type Image struct {
	Source string `xml:"source,attr"`
	Trans  string `xml:"trans,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

// Tile refers to one tile in the tileset
type Tile struct {
	ID        uint32     `xml:"id,attr"`
	Image     Image      `xml:"image"`
	Animation *Animation `xml:"animation"`
}

//Animation references an animated tile
type Animation struct {
	Frames        []*Frame `xml:"frame"`
	animationTime int64
	totalDuration int64
	currentFrame  int
}

//GetFrame returns the frame that will be drawn
func (a Animation) GetFrame() *Frame {
	if len(a.Frames) == 0 {
		return nil
	}

	return a.Frames[a.currentFrame]
}

//Update the animation with elapsedTime since last update
func (a *Animation) Update(elapsedTime int64) {
	if a.totalDuration == 0 {
		for i, f := range a.Frames {
			a.totalDuration += f.Duration
			a.Frames[i].endTime = a.totalDuration
		}
	}

	if len(a.Frames) > 1 {
		a.animationTime += elapsedTime
		if a.animationTime >= a.totalDuration {
			a.animationTime = a.animationTime % a.totalDuration
			a.currentFrame = 0
		}

		for a.animationTime >= a.GetFrame().endTime {
			a.currentFrame++
		}
	}
}

//Frame is one frame of an animation
type Frame struct {
	//Duration is given in milliseconds
	Duration int64 `xml:"duration,attr"`
	TileID   int   `xml:"tileid,attr"`
	endTime  int64
}
