[![Join the chat at https://gitter.im/manyminds/tmx](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/manyminds/tmx)

[![GoDoc](https://godoc.org/github.com/manyminds/tmx?status.svg)](https://godoc.org/github.com/manyminds/tmx)
[![Build Status](https://travis-ci.org/manyminds/tmx.svg?branch=master)](https://travis-ci.org/manyminds/tmx) 
[![Coverage Status](https://coveralls.io/repos/manyminds/tmx/badge.svg?branch=master&service=github)](https://coveralls.io/github/manyminds/tmx?branch=master) 
[![Go Report Card](http://goreportcard.com/badge/manyminds/tmx)](http://goreportcard.com/report/manyminds/tmx)
# TMX Map File Loader

This repository aims to provide go support for maps that are saved according the the [TMX Map Format](http://doc.mapeditor.org/reference/tmx-map-format)

## Support

This library currently supports loading of base64 encoded tile maps with either gzip, zlip or no compression.

## Usage

Currently the library only provides functionality to load maps, in the future it should provide utility functions
to make using tmx files even more convenient. 

## Renderer

To generate a preview image of your tilemap you can use the Render function: 

```go
  testfile := "example.tmx"

  reader, err := os.Open(testfile)

  m, err := tmx.NewMap(reader)
  if err != nil {
    log.Fatal(err)
  }

  canvas := tmx.NewImageCanvasFromMap(*m)

  renderer := tmx.NewRenderer(*m, canvas)
  err = renderer.Render()
  if err != nil {
    log.Fatal(err)
  }

  target, err := os.Create("result.png")
  if err != nil {
    target, err = os.Open("result.png")
    if err != nil {
      log.Fatal(err)
    }
  }
  defer target.Close()

  err = png.Encode(target, canvas.Image())
  if err != nil {
    log.Fatal(err)
  }
```

The renderer is still a work in progress and currently only renders tiles and layers. 