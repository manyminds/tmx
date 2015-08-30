// Package tmx provides an TMX file loader
// specification: http://doc.mapeditor.org/reference/tmx-map-format/
package tmx

import "database/sql"

const (
	//GIDHorizontalFlip for horizontal flipped tiles
	GIDHorizontalFlip = 0x80000000
	//GIDVerticalFlip for vertical flipped tiles
	GIDVerticalFlip = 0x40000000
	//GIDDiagonalFlip for diagonally flipped tile
	GIDDiagonalFlip = 0x20000000
	//GIDFlips removes all flipping informations
	GIDFlips = GIDHorizontalFlip | GIDVerticalFlip | GIDDiagonalFlip
)

//GID is a global tile id
type GID uint32

//DataTile datatile
type DataTile struct {
	GID GID `xml:"gid,attr"`
	//HorizontalFlip true if the tile should be rendered with flip in x dir
	HorizontalFlip bool
	//VerticalFlip true if the tile should be rendered with flip in y dir
	VerticalFlip bool
	//DiagonalFlip true if it should be rendered flipped diagonally, can be combined
	//with `HorizontalFlip` or `VerticalFlip`
	DiagonalFlip bool
}

//Layer represents one layer of the map.
type Layer struct {
	Name       string       `xml:"name,attr"`
	Opacity    float32      `xml:"opacity,attr"`
	Visible    sql.NullBool `xml:"visible,attr"`
	Properties []Property   `xml:"properties>property"`
	Data       Data         `xml:"data"`
	Width      int          `xml:"width,attr"`
	Height     int          `xml:"height,attr"`
}

//IsVisible returns true if the layer is visible, false otherwise
func (l Layer) IsVisible() bool {
	if !l.Visible.Valid {
		return true
	}

	return l.Visible.Bool
}

//ObjectGroup is a group of objects
type ObjectGroup struct {
	Name       string       `xml:"name,attr"`
	Color      string       `xml:"color,attr"`
	Opacity    float32      `xml:"opacity,attr"`
	Visible    sql.NullBool `xml:"visible,attr"`
	Properties []Property   `xml:"properties>property"`
	Objects    []Object     `xml:"object"`
}

//IsVisible returns true if the object group is visible, false otherwise
func (o ObjectGroup) IsVisible() bool {
	if !o.Visible.Valid {
		return true
	}

	return o.Visible.Bool
}

// Object is an object
type Object struct {
	Name       string       `xml:"name,attr"`
	Type       string       `xml:"type,attr"`
	X          int          `xml:"x,attr"`
	Y          int          `xml:"y,attr"`
	Width      int          `xml:"width,attr"`
	Height     int          `xml:"height,attr"`
	GID        int          `xml:"gid,attr"`
	Visible    sql.NullBool `xml:"visible,attr"`
	Polygons   []Polygon    `xml:"polygon"`
	PolyLines  []PolyLine   `xml:"polyline"`
	Properties []Property   `xml:"properties>property"`
}

//IsVisible returns true if the object is visible, false otherwise
func (o Object) IsVisible() bool {
	if !o.Visible.Valid {
		return true
	}

	return o.Visible.Bool
}

//Polygon loads a polygon from tmx
type Polygon struct {
	Points string `xml:"points,attr"`
}

//PolyLine loads a polyline from tmx
type PolyLine struct {
	Points string `xml:"points,attr"`
}

//Property can be set on tiles
type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}
