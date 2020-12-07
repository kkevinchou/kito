package render

import (
	"image"
	"image/draw"
	"log"
	"os"

	"github.com/go-gl/gl/v4.6-compatibility/gl"
	"github.com/kkevinchou/kito/lib/math/vector"
	"github.com/kkevinchou/kito/lib/pathing"
)

func RenderNoiseMap() {
	var amplitude float32 = 5
	var edgeLength float32 = 1
	var lineOffset float32 = 0.01
	// gl.Disable(gl.LIGHTING)

	gl.Normal3f(0, 1, 0)
	for y := 0; y < len(noiseMap[0])-1; y++ {
		gl.Begin(gl.TRIANGLE_STRIP)
		for x := 0; x < len(noiseMap); x++ {
			// val := noiseMap[x][y]
			// var r, g, b float32 = float32(val), float32(val), float32(val)
			// if y%2 != x%2 {
			// 	r, g, b = 1, 1, 1
			// }
			// drawQuad(float32(x), float32(val)*10, float32(y), 5, r, g, b, false)
			gl.Color3f(0.75, 0.75, 0.75)
			gl.Vertex3f(float32(x)*edgeLength, float32(noiseMap[x][y])*amplitude, float32(y)*edgeLength)
			gl.Vertex3f(float32(x)*edgeLength, float32(noiseMap[x][y+1])*amplitude, float32(y+1)*edgeLength)
		}
		gl.End()
	}

	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	gl.Normal3f(0, 1, 0)
	for y := 0; y < len(noiseMap[0])-1; y++ {
		gl.Begin(gl.TRIANGLE_STRIP)
		for x := 0; x < len(noiseMap); x++ {
			// val := noiseMap[x][y]
			// var r, g, b float32 = float32(val), float32(val), float32(val)
			// if y%2 != x%2 {
			// 	r, g, b = 1, 1, 1
			// }
			// drawQuad(float32(x), float32(val)*10, float32(y), 5, r, g, b, false)
			gl.Color3f(0, 0, 0)
			gl.Vertex3f(float32(x)*edgeLength, float32(noiseMap[x][y])*amplitude+lineOffset, float32(y)*edgeLength)
			gl.Vertex3f(float32(x)*edgeLength, float32(noiseMap[x][y+1])*amplitude+lineOffset, float32(y+1)*edgeLength)
		}
		gl.End()
	}
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	// gl.Enable(gl.LIGHTING)
}

func RenderNavMesh(navMesh *pathing.NavMesh) {
	polygons := navMesh.Polygons()
	for i, polygon := range polygons {
		p0 := polygon.Points()[0]
		p1 := polygon.Points()[1]
		p2 := polygon.Points()[2]

		a := vector.Vector3{X: p1.X - p0.X, Y: p1.Y - p0.Y, Z: p1.Z - p0.Z}
		b := vector.Vector3{X: p2.X - p0.X, Y: p2.Y - p0.Y, Z: p2.Z - p0.Z}

		// normal always points "up"
		normal := a.Cross(b).Normalize()
		if normal.Y < 0 {
			normal.Y = -normal.Y
		}

		color := make([]float32, 3)

		gl.Begin(gl.POLYGON)
		gl.Normal3f(float32(normal.X), float32(normal.Y), float32(normal.Z))

		for _, point := range polygon.Points() {
			if i%2 == 0 {
				color[0], color[1], color[2] = 0, 0, 0
			} else {
				color[0], color[1], color[2] = 1, 1, 1
			}
			gl.Color3f(color[0], color[1], color[2])
			gl.Vertex3f(float32(point.X), float32(point.Y), float32(point.Z))
		}
		gl.End()
	}
}

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
