package render

import (
	"image"
	"image/draw"
	"log"
	"math"
	"os"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kkevinchou/kito/lib/math/vector"
)

func drawLine(start, end vector.Vector3) {
	gl.LineWidth(2.5)
	gl.Color3f(1.0, 0.0, 0.0)
	gl.Begin(gl.LINES)
	gl.Vertex3f(float32(start.X), float32(start.Y), float32(start.Z))
	gl.Vertex3f(float32(end.X), float32(end.Y), float32(end.Z))
	gl.End()
}

func drawFloor() {
	width := 21
	height := 21
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			x := (i - int(math.Floor(float64(width)/2))) * floorPanelDimension
			y := (j - int(math.Floor(float64(height)/2))) * floorPanelDimension
			drawFloorPanel(float32(x), float32(y), (i+j)%2 == 0)
		}
	}
}

func newTexture(file string) uint32 {
	imgFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("texture %q not found on disk: %v\n", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

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
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture
}

func drawFloorPanel(x, z float32, black bool) {
	color := make([]float32, 3)
	if black {
		color[0] = 0
		color[1] = 0
		color[2] = 0
	} else {
		color[0] = 1
		color[1] = 1
		color[2] = 1
	}

	gl.Begin(gl.QUADS)

	halfDimension := float32(floorPanelDimension) / 2
	gl.Normal3f(0, 1, 0)
	gl.Color3f(color[0], color[1], color[2])
	gl.Vertex3f(x-halfDimension, 0, z-halfDimension)
	gl.Color3f(color[0], color[1], color[2])
	gl.Vertex3f(x-halfDimension, 0, z+halfDimension)
	gl.Color3f(color[0], color[1], color[2])
	gl.Vertex3f(x+halfDimension, 0, z+halfDimension)
	gl.Color3f(color[0], color[1], color[2])
	gl.Vertex3f(x+halfDimension, 0, z-halfDimension)

	gl.End()
}

func drawQuad(texture uint32, x, y, z float32) {
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.Color4f(1, 1, 1, 1)

	gl.Begin(gl.QUADS)

	// // FRONT
	gl.Normal3f(0, 0, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+-0.5, y+1, z+0.5)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+-0.5, y+0, z+0.5)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+0.5, y+0, z+0.5)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+0.5, y+1, z+0.5)

	// BACK
	gl.Normal3f(0, 0, -1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+0.5, y+1, z+-0.5)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+0.5, y+0, z+-0.5)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+-0.5, y+0, z+-0.5)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+-0.5, y+1, z+-0.5)

	// TOP
	gl.Normal3f(0, 1, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+-0.5, y+1, z+-0.5)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+-0.5, y+1, z+0.5)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+0.5, y+1, z+0.5)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+0.5, y+1, z+-0.5)

	// BOTTOM
	gl.Normal3f(0, -1, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+-0.5, y+0, z+0.5)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+-0.5, y+0, z+-0.5)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+0.5, y+0, z+-0.5)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+0.5, y+0, z+0.5)

	// RIGHT
	gl.Normal3f(1, 0, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+0.5, y+1, z+0.5)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+0.5, y+0, z+0.5)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+0.5, y+0, z-0.5)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+0.5, y+1, z-0.5)

	// LEFT
	gl.Normal3f(-1, 0, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+-0.5, y+1, z+-0.5)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+-0.5, y+0, z+-0.5)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+-0.5, y+0, z+0.5)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+-0.5, y+1, z+0.5)

	gl.End()
	gl.Disable(gl.TEXTURE_2D)
}

func drawQuad2(texture uint32, x, y, z, size float32) {
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.Color4f(1, 1, 1, 1)

	gl.Begin(gl.QUADS)

	// // FRONT
	gl.Normal3f(0, 0, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+-size, y+(size*2), z+size)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+-size, y+0, z+size)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+size, y+0, z+size)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+size, y+(size*2), z+size)

	// BACK
	gl.Normal3f(0, 0, -1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+size, y+(size*2), z+-size)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+size, y+0, z+-size)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+-size, y+0, z+-size)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+-size, y+(size*2), z+-size)

	// TOP
	gl.Normal3f(0, 1, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+-size, y+(size*2), z+-size)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+-size, y+(size*2), z+size)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+size, y+(size*2), z+size)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+size, y+(size*2), z+-size)

	// BOTTOM
	gl.Normal3f(0, -1, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+-size, y+0, z+size)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+-size, y+0, z+-size)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+size, y+0, z+-size)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+size, y+0, z+size)

	// RIGHT
	gl.Normal3f(1, 0, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+size, y+(size*2), z+size)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+size, y+0, z+size)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+size, y+0, z-size)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+size, y+(size*2), z-size)

	// LEFT
	gl.Normal3f(-1, 0, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+-size, y+(size*2), z+-size)
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+-size, y+0, z+-size)
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+-size, y+0, z+size)
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+-size, y+(size*2), z+size)

	gl.End()
	gl.Disable(gl.TEXTURE_2D)
}
