package tmx

import (
	"image"
	"os"

	//import for gif support
	_ "image/gif"
	//import for jpeg support
	_ "image/jpeg"
	//import for png support
	_ "image/png"
)

//ResourceLocator can be implemented to
//load resources differently than from filesystem
type ResourceLocator interface {
	LoadResource(filepath string) (image.Image, error)
}

//FilesystemLocator loads files simply from the filesystem
//it supports png, jpeg and non animated gifs
type FilesystemLocator struct {
}

//LoadResource to implement ResourceLocator interface
func (f FilesystemLocator) LoadResource(filepath string) (image.Image, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	data, _, err := image.Decode(file)
	return data, err
}
