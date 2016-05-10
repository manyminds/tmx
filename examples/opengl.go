package main

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/manyminds/tmx"
)

type chipset struct {
	handle   uint32
	width    int
	height   int
	tilesize int
}

const screenWidth = 640
const screenHeight = 480

func init() {
	runtime.LockOSThread()
}

type openGLCanvas struct {
	width, height int
	sets          map[string]chipset
	scaleX        float32
	scaleY        float32
}

func (o *openGLCanvas) Draw(tile image.Rectangle, where image.Rectangle, tileset string) {
	if _, ok := o.sets[tileset]; !ok {
		o.sets[tileset] = newChipset(tileset, tile.Max.X-tile.Min.X)
	}

	c := o.sets[tileset]
	gl.BindTexture(gl.TEXTURE_2D, c.handle)

	// Texture coords
	fts := float32(c.tilesize)
	tileWidthPixels := fts / float32(c.width)
	tileHeightPixels := fts / float32(c.height)
	startX := (float32(tile.Min.X) / fts) * tileWidthPixels
	startY := (float32(tile.Min.Y) / fts) * tileHeightPixels
	endX := startX + tileWidthPixels
	endY := startY + tileHeightPixels

	// Draw coords
	drawX := float32(where.Min.X)
	drawY := float32(where.Min.Y)

	gl.Begin(gl.QUADS)
	{
		gl.TexCoord2f(startX, startY)
		gl.Vertex3f(drawX*o.scaleX, drawY*o.scaleY, 0)

		gl.TexCoord2f(startX, endY)
		gl.Vertex3f(drawX*o.scaleX, (drawY+fts)*o.scaleY, 0)

		gl.TexCoord2f(endX, endY)
		gl.Vertex3f((drawX+fts)*o.scaleX, (drawY+fts)*o.scaleY, 0)

		gl.TexCoord2f(endX, startY)
		gl.Vertex3f((drawX+fts)*o.scaleX, (drawY)*o.scaleY, 0)
	}
	gl.End()
}

func (o openGLCanvas) FillRect(what color.Color, where image.Rectangle) {
	return
	drawX := float32(where.Min.X) * o.scaleX
	drawY := float32(where.Min.Y) * o.scaleY
	endX := float32(where.Max.X) * o.scaleX
	endY := float32(where.Max.Y) * o.scaleY
	r, g, b, a := what.RGBA()
	gl.Color4f(float32(r)/0xFF, float32(g)/0xFF, float32(b)/0xFF, float32(a)/0xFF)
	gl.Begin(gl.QUADS)
	{
		gl.Vertex3f(drawX, drawY, 0)
		gl.Vertex3f(drawX, endX, 0)
		gl.Vertex3f(endX, endY, 0)
		gl.Vertex3f(endX, drawY, 0)
	}
	gl.End()
}

func (o openGLCanvas) Bounds() image.Rectangle {
	return image.Rectangle{Min: image.ZP, Max: image.Pt(o.width, o.height)}
}

func newOpenGLCanvas(width, height int, scaleX, scaleY float32) tmx.Canvas {
	return &openGLCanvas{width: width, height: height, sets: map[string]chipset{}, scaleX: scaleX, scaleY: scaleY}
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	fp, err := os.Open("example.tmx")
	if err != nil {
		panic(err)
	}

	m, err := tmx.NewMap(fp)
	if err != nil {
		panic(err)
	}

	var monitor *glfw.Monitor
	window, err := glfw.CreateWindow(screenWidth, screenHeight, "Map Renderer", monitor, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	width, height := window.GetFramebufferSize()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(width), float64(height), 0, -1, 1)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	canvas := newOpenGLCanvas(width, height, float32(width)/float32(screenWidth), float32(height)/float32(screenHeight))
	renderer := tmx.NewRenderer(*m, canvas)
	fps := 0
	startTime := time.Now().UnixNano()
	for !window.ShouldClose() {
		renderer.Render()
		fps++
		if time.Now().UnixNano()-startTime > 1000*1000*1000 {
			log.Println(fps)
			startTime = time.Now().UnixNano()
			fps = 0
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func newChipset(file string, tilesize int) chipset {
	imgFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("texture %q not found on disk: %v\n", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.ZP, draw.Src)

	var texture uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0, gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return chipset{width: rgba.Bounds().Dx(), height: rgba.Bounds().Dy(), handle: texture, tilesize: tilesize}
}
