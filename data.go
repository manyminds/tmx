package tmx

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
)

//Data contains raw loaded tmx data
type Data struct {
	Encoding    string     `xml:"encoding,attr"`
	Compression string     `xml:"compression,attr"`
	RawData     []byte     `xml:",innerxml"`
	DataTiles   []DataTile `xml:"tile"`
}

//utility function to check flipping
func isHorizontallyFlipped(gid GID) bool {
	return gid&GIDHorizontalFlip != 0
}

//utility function to check flipping
func isVerticallyFlipped(gid GID) bool {
	return gid&GIDVerticalFlip != 0
}

//utility function to check flipping
func isDiagonallyFlipped(gid GID) bool {
	return gid&GIDDiagonalFlip != 0
}

// loadEncodedTiles loads all GID informations
// from RawData to `DataTiles`
func (d *Data) loadEncodedTiles() error {
	if len(d.RawData) == 0 {
		return nil
	}

	rawData := bytes.TrimSpace(d.RawData)
	bytesReader := bytes.NewReader(rawData)

	var reader io.Reader
	reader = bytesReader

	if d.Encoding == "base64" {
		reader = base64.NewDecoder(base64.StdEncoding, reader)
	}

	decodedData, err := decompress(reader, d.Compression)
	if err != nil {
		return err
	}

	// every 4 bytes is one tile
	if len(decodedData)%4 != 0 {
		return errors.New("Tile information []byte must consist solely of 32bit integers.")
	}

	d.DataTiles = make([]DataTile, len(decodedData)/4)

	for j := 0; j < len(decodedData); {
		gid := GID(decodedData[j]) +
			GID(decodedData[j+1])<<8 +
			GID(decodedData[j+2])<<16 +
			GID(decodedData[j+3])<<24

		tile := DataTile{}
		tile.HorizontalFlip = isHorizontallyFlipped(gid)
		tile.VerticalFlip = isVerticallyFlipped(gid)
		tile.DiagonalFlip = isDiagonallyFlipped(gid)

		//flip information must be cleared
		tile.GID = gid &^ GIDFlips

		d.DataTiles[j/4] = tile

		j += 4
	}

	return nil
}

// decompress input from `r` with the given `compression` standard
func decompress(r io.Reader, compression string) (data []byte, err error) {
	var compressionReader io.Reader
	switch compression {
	case "gzip":
		compressionReader, err = gzip.NewReader(r)
		if err != nil {
			return
		}
	case "zlib":
		compressionReader, err = zlib.NewReader(r)
		if err != nil {
			return
		}
	case "":
		compressionReader = r
	default:
		err = errors.New("Only zlib and gzip compressions are supported")
		return
	}

	return ioutil.ReadAll(compressionReader)
}
