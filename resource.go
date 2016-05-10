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
	LocateResource(filepath string) (image.Image, error)
}

//ResourceManager allows better memory handling
//and cleanup of resources
type ResourceManager interface {
	UnsetResource(filepath string)
}

//FilesystemLocator loads files simply from the filesystem
//it supports png, jpeg and non animated gifs
type FilesystemLocator struct {
}

//LocateResource to implement ResourceLocator interface
func (f FilesystemLocator) LocateResource(filepath string) (image.Image, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	data, _, err := image.Decode(file)
	return data, err
}

type lazyLocator struct {
	parent   ResourceLocator
	tilesets map[string]image.Image
}

func (l *lazyLocator) LocateResource(filepath string) (image.Image, error) {
	cached, ok := l.tilesets[filepath]
	if ok {
		return cached, nil
	}

	data, err := l.parent.LocateResource(filepath)
	if err != nil {
		return nil, err
	}

	l.tilesets[filepath] = data

	return data, nil
}

func (l *lazyLocator) UnsetResource(filepath string) {
	delete(l.tilesets, filepath)
}

//NewLazyResourceLocator wraps a ResourceLocator and caches results
func NewLazyResourceLocator(l ResourceLocator) ResourceLocator {
	return &lazyLocator{parent: l, tilesets: map[string]image.Image{}}
}
