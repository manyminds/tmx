package tmx

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	validator "gopkg.in/bluesuncorp/validator.v8"

	"github.com/manyminds/tmx/spec"
)

//Map contains map inspecion
type Map struct {
	Version         string             `xml:"title,attr"`
	Orientation     string             `xml:"orientation,attr"`
	Width           int                `validate:"gt=0" xml:"width,attr"`
	Height          int                `validate:"gt=0" xml:"height,attr"`
	TileWidth       int                `validate:"gt=0" xml:"tilewidth,attr"`
	TileHeight      int                `validate:"gt=0" xml:"tileheight,attr"`
	BackgroundColor hexcolor           `validate:"omitempty,rgb|hexcolor" xml:"backgroundcolor,attr"`
	RenderOrder     string             `xml:"renderorder,attr"`
	Properties      []spec.Property    `xml:"properties>property"`
	Tilesets        []spec.Tileset     `xml:"tileset"`
	Layers          []spec.Layer       `xml:"layer"`
	ObjectGroups    []spec.ObjectGroup `xml:"objectgroup"`
	//since tileset loading sucks so much and uses relative paths
	//we store the original filename for this map if possible
	filename string
}

//GetTilesetForGID returns the correct tileset for a given gid
func (m Map) GetTilesetForGID(gid spec.GID) (*spec.Tileset, error) {
	if gid == 0 {
		return nil, nil
	}

	for i, tileset := range m.Tilesets {
		if gid >= tileset.FirstGID && gid < tileset.FirstGID+spec.GID(tileset.GetNumTiles()) {
			return &m.Tilesets[i], nil
		}
	}

	return nil, fmt.Errorf("Invalid GID %d given.", gid)
}

// NewMap creates a new map from a given io.Reader
func NewMap(f io.Reader) (*Map, error) {
	return NewMapWithValidation(f, true)
}

//NewMapWithValidation lets you create a map and decide
//wether you want to validate it, or skip validation
//and risk potential errors
//this is particulary useful if you want to skip validation
//during runtime for performance benefits
func NewMapWithValidation(f io.Reader, validation bool) (*Map, error) {
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

	if validation {
		err = target.Validate()
		if err != nil {
			return nil, err
		}
	}

	target.filename = filename

	for key, layer := range target.Layers {
		if e := layer.Data.LoadEncodedTiles(); e != nil {
			return nil, e
		}

		target.Layers[key] = layer
	}

	return &target, nil
}

//Validate will semantically validate the given map
func (m Map) Validate() error {
	v := validator.New(&validator.Config{TagName: "validate"})

	return v.Struct(m)
}
