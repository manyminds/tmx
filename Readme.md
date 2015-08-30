[![Build Status](https://travis-ci.org/manyminds/tmx.svg?branch=master)](https://travis-ci.org/manyminds/tmx)

# TMX Map File Loader

This repository aims to provide go support for maps that are saved according the the [http://doc.mapeditor.org/reference/tmx-map-format](TMX Map Format).

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

  renderer := tmx.NewRendererWithCanvas(*m, canvas)
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

  err = png.Encode(target, canvas.Image())
  if err != nil {
    log.Fatal(err)
  }
```

The renderer is still a work in progress and currently only renders tiles and layers. 