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
	ID    uint32 `xml:"id,attr"`
	Image Image  `xml:"image"`
}
