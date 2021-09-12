package render

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"os"

	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/lib/animation"
	"github.com/kkevinchou/kito/lib/shaders"
	"github.com/kkevinchou/kito/lib/textures"
	"github.com/kkevinchou/kito/lib/utils"
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

	// is vertically flipped if directly read into opengl texture
	nrgba := imaging.FlipV(img)

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}

	draw.Draw(rgba, rgba.Bounds(), nrgba, image.Point{0, 0}, draw.Src)

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

func drawSkyBox(sb *SkyBox, shader *shaders.ShaderProgram, frontTexture, topTexture, leftTexture, rightTexture, bottomTexture, backTexture *textures.Texture, viewMatrix, projectionMatrix mgl32.Mat4) {
	textures := []*textures.Texture{frontTexture, topTexture, leftTexture, rightTexture, bottomTexture, backTexture}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindVertexArray(sb.VAO())
	shader.Use()
	shader.SetUniformInt("skyboxTexture", 0)
	shader.SetUniformMat4("model", mgl32.Ident4())
	shader.SetUniformMat4("view", viewMatrix)
	shader.SetUniformMat4("projection", projectionMatrix)
	for i := 0; i < 6; i++ {
		gl.BindTexture(gl.TEXTURE_2D, textures[i].ID)
		gl.DrawArrays(gl.TRIANGLES, int32(i*6), 6)
	}
}

// drawHUDTextureToQuad does a shitty perspective based rendering of a flat texture
func drawHUDTextureToQuad(shader *shaders.ShaderProgram, texture uint32, projectionMatrix mgl32.Mat4, hudScale float32) {
	// texture coords top left = 0,0 | bottom right = 1,1
	var skyboxVertices []float32 = []float32{
		// front
		-1 * hudScale, -1 * hudScale, 0, 0.0, 0.0,
		1 * hudScale, -1 * hudScale, 0, 1.0, 0.0,
		1 * hudScale, 1 * hudScale, 0, 1.0, 1.0,
		1 * hudScale, 1 * hudScale, 0, 1.0, 1.0,
		-1 * hudScale, 1 * hudScale, 0, 0.0, 1.0,
		-1 * hudScale, -1 * hudScale, 0, 0.0, 0.0,
	}

	var vbo, vao uint32
	gl.GenBuffers(1, &vbo)
	gl.GenVertexArrays(1, &vao)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(skyboxVertices)*4, gl.Ptr(skyboxVertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(vao)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	shader.Use()
	shader.SetUniformMat4("model", mgl32.Translate3D(1.2, 0.8, -2))
	shader.SetUniformMat4("view", mgl32.Ident4())
	shader.SetUniformMat4("projection", projectionMatrix)
	shader.SetUniformInt("depthMap", 0)

	gl.DrawArrays(gl.TRIANGLES, 0, 6)

}

func drawTextureToQuad(shader *shaders.ShaderProgram, texture uint32, modelMatrix, viewMatrix, projectionMatrix mgl32.Mat4) {
	// texture coords top left = 0,0 | bottom right = 1,1
	var skyboxVertices []float32 = []float32{
		// front
		-5, -5, -5, 0.0, 0.0,
		5, -5, -5, 1.0, 0.0,
		5, 5, -5, 1.0, 1.0,
		5, 5, -5, 1.0, 1.0,
		-5, 5, -5, 0.0, 1.0,
		-5, -5, -5, 0.0, 0.0,
	}

	var vbo, vao uint32
	gl.GenBuffers(1, &vbo)
	gl.GenVertexArrays(1, &vao)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(skyboxVertices)*4, gl.Ptr(skyboxVertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(vao)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	shader.Use()
	shader.SetUniformMat4("model", modelMatrix)
	shader.SetUniformMat4("view", viewMatrix)
	shader.SetUniformMat4("projection", projectionMatrix)
	shader.SetUniformInt("depthMap", 0)

	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func drawMesh(mesh Mesh, shader *shaders.ShaderProgram, modelMatrix, viewMatrix, projectionMatrix mgl32.Mat4, cameraPosition mgl32.Vec3, lightMVPMatrix mgl64.Mat4, texture uint32, directionalLightDir mgl64.Vec3, shadowDistance float64) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	shader.Use()
	shader.SetUniformMat4("model", modelMatrix)
	shader.SetUniformMat4("view", viewMatrix)
	shader.SetUniformMat4("projection", projectionMatrix)
	shader.SetUniformVec3("viewPos", cameraPosition)
	shader.SetUniformFloat("shadowDistance", float32(shadowDistance))
	shader.SetUniformVec3("directionalLightDir", utils.Vec3F64ToF32(directionalLightDir))
	shader.SetUniformMat4("lightSpaceMatrix", utils.Mat4F64ToF32(lightMVPMatrix))
	gl.BindVertexArray(mesh.GetVAO())
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func drawAnimatedMesh(mesh *animation.Mesh, animationTransforms map[int]mgl32.Mat4, texture *textures.Texture, shader *shaders.ShaderProgram, modelMatrix, viewMatrix, projectionMatrix mgl32.Mat4, cameraPosition mgl32.Vec3) {
	shader.Use()
	shader.SetUniformMat4("model", modelMatrix)
	shader.SetUniformMat4("view", viewMatrix)
	shader.SetUniformMat4("projection", projectionMatrix)
	shader.SetUniformVec3("viewPos", mgl32.Vec3{float32(cameraPosition.X()), float32(cameraPosition.Y()), float32(cameraPosition.Z())})

	for i := 0; i < len(animationTransforms); i++ {
		shader.SetUniformMat4(fmt.Sprintf("jointTransforms[%d]", i), animationTransforms[i])
	}
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture.ID)
	gl.BindVertexArray(mesh.VAO())

	gl.DrawElements(gl.TRIANGLES, int32(mesh.VertexCount()), gl.UNSIGNED_INT, nil)
}

func createModelMatrix(scaleMatrix, rotationMatrix, translationMatrix mgl32.Mat4) mgl32.Mat4 {
	return translationMatrix.Mul4(rotationMatrix).Mul4(scaleMatrix)
}

func resetGLRenderSettings() {
	gl.BindVertexArray(0)
	gl.UseProgram(0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.CullFace(gl.BACK)
}
