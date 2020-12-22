package render

import (
	"image"
	"image/draw"
	"log"
	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/lib/animation"
	"github.com/kkevinchou/kito/lib/math/vector"
	"github.com/kkevinchou/kito/lib/shaders"
)

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

func drawSkyBox(sb *SkyBox, shader *shaders.Shader, textureMap map[string]uint32, modelMatrix, viewMatrix, projectionMatrix mgl32.Mat4) {
	textures := []uint32{textureMap["front"], textureMap["top"], textureMap["left"], textureMap["right"], textureMap["bottom"], textureMap["back"]}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindVertexArray(sb.VAO())
	shader.Use()
	shader.SetUniformInt("skyboxTexture", 0)
	shader.SetUniformMat4("model", modelMatrix)
	shader.SetUniformMat4("view", viewMatrix)
	shader.SetUniformMat4("projection", projectionMatrix)
	for i := 0; i < 6; i++ {
		gl.BindTexture(gl.TEXTURE_2D, textures[i])
		gl.DrawArrays(gl.TRIANGLES, int32(i*6), 6)
	}
}

func drawQuad(q *Quad, shader *shaders.Shader, modelMatrix, viewMatrix, projectionMatrix mgl32.Mat4, viewerPosition vector.Vector3) {
	shader.Use()
	shader.SetUniformMat4("model", modelMatrix)
	shader.SetUniformMat4("view", viewMatrix)
	shader.SetUniformMat4("projection", projectionMatrix)
	shader.SetUniformVec3("viewPos", mgl32.Vec3{float32(viewerPosition.X), float32(viewerPosition.Y), float32(viewerPosition.Z)})
	gl.BindVertexArray(q.VAO())
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func drawMesh(mesh *animation.Mesh, shader *shaders.Shader, modelMatrix, viewMatrix, projectionMatrix mgl32.Mat4, viewerPosition vector.Vector3) {
	shader.Use()
	shader.SetUniformMat4("model", modelMatrix)
	shader.SetUniformMat4("view", viewMatrix)
	shader.SetUniformMat4("projection", projectionMatrix)
	shader.SetUniformVec3("viewPos", mgl32.Vec3{float32(viewerPosition.X), float32(viewerPosition.Y), float32(viewerPosition.Z)})
	gl.BindVertexArray(mesh.VAO())

	gl.DrawElements(gl.TRIANGLES, int32(mesh.VertexCount()), gl.UNSIGNED_INT, nil)
}

func createModelMatrix(scaleMatrix, rotationMatrix, translationMatrix mgl32.Mat4) mgl32.Mat4 {
	return translationMatrix.Mul4(rotationMatrix).Mul4(scaleMatrix)
}

// func RenderNoiseMap(noiseMap [][]float64, xOffset, zOffset, edgeLength float32) {
// 	var amplitude float32 = 5
// 	var lineOffset float32 = 0.01
// 	// gl.Disable(gl.LIGHTING)

// 	gl.Normal3f(0, 1, 0)
// 	for y := 0; y < len(noiseMap[0])-1; y++ {
// 		gl.Begin(gl.TRIANGLE_STRIP)
// 		for x := 0; x < len(noiseMap); x++ {
// 			// val := noiseMap[x][y]
// 			// var r, g, b float32 = float32(val), float32(val), float32(val)
// 			// if y%2 != x%2 {
// 			// 	r, g, b = 1, 1, 1
// 			// }
// 			// drawQuad(float32(x), float32(val)*10, float32(y), 5, r, g, b, false)
// 			gl.Color3f(0.11, 0.31, 0.29)
// 			gl.Vertex3f(float32(x)*edgeLength+xOffset, float32(noiseMap[x][y])*amplitude, float32(y)*edgeLength+zOffset)
// 			gl.Vertex3f(float32(x)*edgeLength+xOffset, float32(noiseMap[x][y+1])*amplitude, float32(y+1)*edgeLength+zOffset)
// 		}
// 		gl.End()
// 	}

// 	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
// 	gl.Normal3f(0, 1, 0)
// 	for y := 0; y < len(noiseMap[0])-1; y++ {
// 		gl.Begin(gl.TRIANGLE_STRIP)
// 		for x := 0; x < len(noiseMap); x++ {
// 			// val := noiseMap[x][y]
// 			// var r, g, b float32 = float32(val), float32(val), float32(val)
// 			// if y%2 != x%2 {
// 			// 	r, g, b = 1, 1, 1
// 			// }
// 			// drawQuad(float32(x), float32(val)*10, float32(y), 5, r, g, b, false)
// 			gl.Color3f(0, 0, 0)
// 			gl.Vertex3f(float32(x)*edgeLength+xOffset, float32(noiseMap[x][y])*amplitude+lineOffset, float32(y)*edgeLength+zOffset)
// 			gl.Vertex3f(float32(x)*edgeLength+xOffset, float32(noiseMap[x][y+1])*amplitude+lineOffset, float32(y+1)*edgeLength+zOffset)
// 		}
// 		gl.End()
// 	}
// 	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
// 	// gl.Enable(gl.LIGHTING)
// }

// func RenderNavMesh(navMesh *pathing.NavMesh) {
// 	polygons := navMesh.Polygons()
// 	for i, polygon := range polygons {
// 		p0 := polygon.Points()[0]
// 		p1 := polygon.Points()[1]
// 		p2 := polygon.Points()[2]

// 		a := vector.Vector3{X: p1.X - p0.X, Y: p1.Y - p0.Y, Z: p1.Z - p0.Z}
// 		b := vector.Vector3{X: p2.X - p0.X, Y: p2.Y - p0.Y, Z: p2.Z - p0.Z}

// 		// normal always points "up"
// 		normal := a.Cross(b).Normalize()
// 		if normal.Y < 0 {
// 			normal.Y = -normal.Y
// 		}

// 		color := make([]float32, 3)

// 		gl.Begin(gl.POLYGON)
// 		gl.Normal3f(float32(normal.X), float32(normal.Y), float32(normal.Z))

// 		for _, point := range polygon.Points() {
// 			if i%2 == 0 {
// 				color[0], color[1], color[2] = 0, 0, 0
// 			} else {
// 				color[0], color[1], color[2] = 1, 1, 1
// 			}
// 			gl.Color3f(color[0], color[1], color[2])
// 			gl.Vertex3f(float32(point.X), float32(point.Y), float32(point.Z))
// 		}
// 		gl.End()
// 	}
// }

// func drawLine(start, end vector.Vector3) {
// 	gl.LineWidth(2.5)
// 	gl.Color3f(1.0, 0.0, 0.0)
// 	gl.Begin(gl.LINES)
// 	gl.Vertex3f(float32(start.X), float32(start.Y), float32(start.Z))
// 	gl.Vertex3f(float32(end.X), float32(end.Y), float32(end.Z))
// 	gl.End()
// }
