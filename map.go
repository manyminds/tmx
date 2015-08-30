package tmx

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

//Map contains map information
type Map struct {
	Version         string        `xml:"title,attr"`
	Orientation     string        `xml:"orientation,attr"`
	Width           int           `xml:"width,attr"`
	Height          int           `xml:"height,attr"`
	TileWidth       int           `xml:"tilewidth,attr"`
	TileHeight      int           `xml:"tileheight,attr"`
	BackgroundColor hexcolor      `xml:"backgroundcolor,attr"`
	RenderOrder     string        `xml:"renderorder,attr"`
	Properties      []Property    `xml:"properties>property"`
	Tilesets        []Tileset     `xml:"tileset"`
	Layers          []Layer       `xml:"layer"`
	ObjectGroups    []ObjectGroup `xml:"objectgroup"`
	//since tileset loading sucks so much and uses relative paths
	//we store the original filename for this map if possible
	filename string
}

//GetTilesetForGID returns the correct tileset for a given gid
func (m Map) GetTilesetForGID(gid GID) (*Tileset, error) {
	if gid == 0 {
		return nil, nil
	}

	for i, tileset := range m.Tilesets {
		if gid >= tileset.FirstGID && gid < tileset.FirstGID+GID(tileset.GetNumTiles()) {
			return &m.Tilesets[i], nil
		}
	}

	return nil, fmt.Errorf("Invalid GID %d given.", gid)
}

// NewMap creates a new map from a given io.Reader
func NewMap(f io.Reader) (*Map, error) {
	var target Map
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	filename := ""
	if f, ok := f.(*os.File); ok {
		filename = filepath.Dir(f.Name()) + string(os.PathSeparator)
	}

	err = xml.Unmarshal(data, &target)
	if err != nil {
		return nil, err
	}

	target.filename = filename

	for key, layer := range target.Layers {
		if e := layer.Data.loadEncodedTiles(); e != nil {
			return nil, e
		}

		target.Layers[key] = layer
	}

	return &target, nil
}
