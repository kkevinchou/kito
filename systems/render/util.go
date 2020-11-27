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

var skyboxVertices []float32 = []float32{
	// positions
	-1.0, 1.0, -1.0,
	-1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,

	-1.0, -1.0, 1.0,
	-1.0, -1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,
	-1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0,

	1.0, -1.0, -1.0,
	1.0, -1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, -1.0,
	1.0, -1.0, -1.0,

	-1.0, -1.0, 1.0,
	-1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	1.0, -1.0, 1.0,
	-1.0, -1.0, 1.0,

	-1.0, 1.0, -1.0,
	1.0, 1.0, -1.0,
	1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,
	-1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0,

	-1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	1.0, -1.0, 1.0,
}

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

func newCubeMap(file string) uint32 {
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
	gl.Enable(gl.TEXTURE_CUBE_MAP)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, texture)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	for i := 0; i < 6; i++ {
		gl.TexImage2D(
			gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i),
			0,
			gl.RGBA,
			int32(rgba.Rect.Size().X),
			int32(rgba.Rect.Size().Y),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(rgba.Pix),
		)
	}

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

func drawSkyBox(textureMap map[string]uint32, x, y, z float32, edgeLength float32, lightningEnabled bool) {
	if !lightningEnabled {
		gl.Disable(gl.LIGHTING)
	}
	gl.Enable(gl.TEXTURE_2D)
	gl.Color4f(1, 1, 1, 1)

	gl.BindTexture(gl.TEXTURE_2D, textureMap["front"])
	gl.Begin(gl.QUADS)

	// // FRONT
	gl.Normal3f(0, 0, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x-(edgeLength/2), y+edgeLength, z+(edgeLength/2))
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x-(edgeLength/2), y+0, z+(edgeLength/2))
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+(edgeLength/2), y+0, z+(edgeLength/2))
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+(edgeLength/2), y+edgeLength, z+(edgeLength/2))

	gl.End()

	gl.BindTexture(gl.TEXTURE_2D, textureMap["back"])
	gl.Begin(gl.QUADS)

	// BACK
	gl.Normal3f(0, 0, -1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+(edgeLength/2), y+edgeLength, z-(edgeLength/2))
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+(edgeLength/2), y+0, z-(edgeLength/2))
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x-(edgeLength/2), y+0, z-(edgeLength/2))
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x-(edgeLength/2), y+edgeLength, z-(edgeLength/2))

	gl.End()

	gl.BindTexture(gl.TEXTURE_2D, textureMap["top"])
	gl.Begin(gl.QUADS)

	// TOP
	gl.Normal3f(0, 1, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x-(edgeLength/2), y+edgeLength, z-(edgeLength/2))
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x-(edgeLength/2), y+edgeLength, z+(edgeLength/2))
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+(edgeLength/2), y+edgeLength, z+(edgeLength/2))
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+(edgeLength/2), y+edgeLength, z-(edgeLength/2))

	gl.End()

	gl.BindTexture(gl.TEXTURE_2D, textureMap["bottom"])
	gl.Begin(gl.QUADS)

	// BOTTOM
	gl.Normal3f(0, -1, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x-(edgeLength/2), y+0, z+(edgeLength/2))
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x-(edgeLength/2), y+0, z-(edgeLength/2))
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+(edgeLength/2), y+0, z-(edgeLength/2))
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+(edgeLength/2), y+0, z+(edgeLength/2))

	gl.End()

	gl.BindTexture(gl.TEXTURE_2D, textureMap["right"])
	gl.Begin(gl.QUADS)

	// RIGHT
	gl.Normal3f(1, 0, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+(edgeLength/2), y+edgeLength, z+(edgeLength/2))
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+(edgeLength/2), y+0, z+(edgeLength/2))
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+(edgeLength/2), y+0, z-(edgeLength/2))
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+(edgeLength/2), y+edgeLength, z-(edgeLength/2))

	gl.End()

	gl.BindTexture(gl.TEXTURE_2D, textureMap["left"])
	gl.Begin(gl.QUADS)

	// LEFT
	gl.Normal3f(-1, 0, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x-(edgeLength/2), y+edgeLength, z-(edgeLength/2))
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x-(edgeLength/2), y+0, z-(edgeLength/2))
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x-(edgeLength/2), y+0, z+(edgeLength/2))
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x-(edgeLength/2), y+edgeLength, z+(edgeLength/2))

	gl.End()

	gl.Disable(gl.TEXTURE_2D)
	if !lightningEnabled {
		gl.Enable(gl.LIGHTING)
	}

}

func drawCube(texture uint32, x, y, z, edgeLength float32, lightningEnabled bool) {
	if !lightningEnabled {
		gl.Disable(gl.LIGHTING)
	}
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.Color4f(1, 1, 1, 1)

	gl.Begin(gl.QUADS)

	// // FRONT
	gl.Normal3f(0, 0, 1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x-(edgeLength/2), y+edgeLength, z+(edgeLength/2))
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x-(edgeLength/2), y+0, z+(edgeLength/2))
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+(edgeLength/2), y+0, z+(edgeLength/2))
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+(edgeLength/2), y+edgeLength, z+(edgeLength/2))

	// BACK
	gl.Normal3f(0, 0, -1)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+(edgeLength/2), y+edgeLength, z-(edgeLength/2))
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+(edgeLength/2), y+0, z-(edgeLength/2))
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x-(edgeLength/2), y+0, z-(edgeLength/2))
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x-(edgeLength/2), y+edgeLength, z-(edgeLength/2))

	// TOP
	gl.Normal3f(0, 1, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x-(edgeLength/2), y+edgeLength, z-(edgeLength/2))
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x-(edgeLength/2), y+edgeLength, z+(edgeLength/2))
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+(edgeLength/2), y+edgeLength, z+(edgeLength/2))
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+(edgeLength/2), y+edgeLength, z-(edgeLength/2))

	// BOTTOM
	gl.Normal3f(0, -1, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x-(edgeLength/2), y+0, z+(edgeLength/2))
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x-(edgeLength/2), y+0, z-(edgeLength/2))
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+(edgeLength/2), y+0, z-(edgeLength/2))
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+(edgeLength/2), y+0, z+(edgeLength/2))

	// RIGHT
	gl.Normal3f(1, 0, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x+(edgeLength/2), y+edgeLength, z+(edgeLength/2))
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x+(edgeLength/2), y+0, z+(edgeLength/2))
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x+(edgeLength/2), y+0, z-(edgeLength/2))
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x+(edgeLength/2), y+edgeLength, z-(edgeLength/2))

	// LEFT
	gl.Normal3f(-1, 0, 0)
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(x-(edgeLength/2), y+edgeLength, z-(edgeLength/2))
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(x-(edgeLength/2), y+0, z-(edgeLength/2))
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(x-(edgeLength/2), y+0, z+(edgeLength/2))
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(x-(edgeLength/2), y+edgeLength, z+(edgeLength/2))

	gl.End()
	gl.Disable(gl.TEXTURE_2D)
	if !lightningEnabled {
		gl.Enable(gl.LIGHTING)
	}
}

func drawCubeMap(texture uint32, x, y, z, edgeLength float32, lightningEnabled bool) {
	if !lightningEnabled {
		gl.Disable(gl.LIGHTING)
	}
	gl.Enable(gl.TEXTURE_CUBE_MAP)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, texture)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(skyboxVertices)*4, gl.Ptr(&skyboxVertices), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
}
