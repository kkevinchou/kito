package render

import (
	"image"
	"image/draw"
	"log"
	"os"

	"github.com/go-gl/gl/v4.6-compatibility/gl"
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

func drawQuad(x, y, z, edgeLength, r, g, b float32, lightningEnabled bool) {
	if !lightningEnabled {
		gl.Disable(gl.LIGHTING)
	}
	gl.Begin(gl.POLYGON)
	gl.Normal3f(0, 1, 0)

	gl.Color3f(r, g, b)
	gl.Vertex3f((x*edgeLength)-(edgeLength/2), y, (z*edgeLength)+(edgeLength/2))
	gl.Vertex3f((x*edgeLength)-(edgeLength/2), y, (z*edgeLength)-(edgeLength/2))
	gl.Vertex3f((x*edgeLength)+(edgeLength/2), y, (z*edgeLength)-(edgeLength/2))
	gl.Vertex3f((x*edgeLength)+(edgeLength/2), y, (z*edgeLength)+(edgeLength/2))
	gl.End()
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
