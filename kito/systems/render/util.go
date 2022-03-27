package render

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/lib/collision/collider"
	"github.com/kkevinchou/kito/lib/font"
	utils "github.com/kkevinchou/kito/lib/libutils"
	"github.com/kkevinchou/kito/lib/shaders"
	"github.com/kkevinchou/kito/lib/textures"
)

func drawModel(viewerContext ViewerContext, lightContext LightContext, shadowMap *ShadowMap, shader *shaders.ShaderProgram, meshComponent *components.MeshComponent, animationComponent *components.AnimationComponent, modelMatrix mgl64.Mat4, modelRotationMatrix mgl64.Mat4) {
	shader.Use()
	shader.SetUniformMat4("model", utils.Mat4F64ToF32(modelMatrix))
	shader.SetUniformMat4("modelRotationMatrix", utils.Mat4F64ToF32(modelRotationMatrix))
	shader.SetUniformMat4("view", utils.Mat4F64ToF32(viewerContext.InverseViewMatrix))
	shader.SetUniformMat4("projection", utils.Mat4F64ToF32(viewerContext.ProjectionMatrix))
	shader.SetUniformVec3("viewPos", utils.Vec3F64ToF32(viewerContext.Position))
	shader.SetUniformFloat("shadowDistance", float32(shadowMap.ShadowDistance()))
	shader.SetUniformVec3("directionalLightDir", utils.Vec3F64ToF32(lightContext.DirectionalLightDir))
	shader.SetUniformMat4("lightSpaceMatrix", utils.Mat4F64ToF32(lightContext.LightSpaceMatrix))
	shader.SetUniformInt("shadowMap", 31)

	if meshComponent.Material != nil && meshComponent.Material.DiffuseColor != nil {
		shader.SetUniformInt("materialHasDiffuseColor", 1)
		shader.SetUniformVec3("materialDiffuseColor", *meshComponent.Material.DiffuseColor)
	} else {
		shader.SetUniformInt("materialHasDiffuseColor", 0)
	}

	if meshComponent.PBRMaterial != nil {
		shader.SetUniformInt("hasPBRMaterial", 1)
		shader.SetUniformVec4("pbrBaseColorFactor", meshComponent.PBRMaterial.PBRMetallicRoughness.BaseColorFactor)
	} else {
		shader.SetUniformInt("hasPBRMaterial", 0)
	}

	if animationComponent != nil {
		animationTransforms := animationComponent.AnimationTransforms
		for i := 0; i < len(animationTransforms); i++ {
			shader.SetUniformMat4(fmt.Sprintf("jointTransforms[%d]", i), animationTransforms[i])
		}
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, meshComponent.Texture.ID)
	gl.ActiveTexture(gl.TEXTURE31)
	gl.BindTexture(gl.TEXTURE_2D, shadowMap.DepthTexture())
	gl.BindVertexArray(meshComponent.ModelVAO)

	gl.DrawElements(gl.TRIANGLES, int32(meshComponent.ModelVertexCount), gl.UNSIGNED_INT, nil)
}

func drawTriMeshCollider(viewerContext ViewerContext, lightContext LightContext, shader *shaders.ShaderProgram, triMeshCollider *collider.TriMesh) {
	var vertices []float32

	for _, triangle := range triMeshCollider.Triangles {
		for _, point := range triangle.Points {
			vertices = append(vertices, float32(point.X()), float32(point.Y()), float32(point.Z()))
		}
	}

	var vbo, vao uint32
	gl.GenBuffers(1, &vbo)
	gl.GenVertexArrays(1, &vao)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(vao)
	shader.Use()
	shader.SetUniformMat4("model", mgl32.Ident4())
	shader.SetUniformMat4("view", utils.Mat4F64ToF32(viewerContext.InverseViewMatrix))
	shader.SetUniformMat4("projection", utils.Mat4F64ToF32(viewerContext.ProjectionMatrix))
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)))
}

func drawCapsuleCollider(viewerContext ViewerContext, lightContext LightContext, shader *shaders.ShaderProgram, capsuleCollider *collider.Capsule, billboardModelMatrix mgl64.Mat4) {
	radius := float32(capsuleCollider.Radius)
	top := float32(capsuleCollider.Top.Y()) + radius
	bottom := float32(capsuleCollider.Bottom.Y()) - radius

	vertices := []float32{
		-radius, bottom, 0,
		radius, bottom, 0,
		radius, top, 0,
		radius, top, 0,
		-radius, top, 0,
		-radius, bottom, 0,
	}

	// for i := 0; i < len(vertices); i += 3 {
	// 	vertices[i] = vertices[i] + float32(capsuleCollider.Bottom.X())
	// 	vertices[i+2] = vertices[i+2] + float32(capsuleCollider.Bottom.Z())
	// }

	var vbo, vao uint32
	gl.GenBuffers(1, &vbo)
	gl.GenVertexArrays(1, &vao)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(vao)
	shader.Use()
	shader.SetUniformMat4("model", utils.Mat4F64ToF32(billboardModelMatrix))
	shader.SetUniformMat4("view", utils.Mat4F64ToF32(viewerContext.InverseViewMatrix))
	shader.SetUniformMat4("projection", utils.Mat4F64ToF32(viewerContext.ProjectionMatrix))
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)))
}

func drawSkyBox(viewerContext ViewerContext, sb *SkyBox, shader *shaders.ShaderProgram, frontTexture, topTexture, leftTexture, rightTexture, bottomTexture, backTexture *textures.Texture) {
	textures := []*textures.Texture{frontTexture, topTexture, leftTexture, rightTexture, bottomTexture, backTexture}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindVertexArray(sb.VAO())
	shader.Use()
	shader.SetUniformInt("skyboxTexture", 0)
	shader.SetUniformMat4("model", mgl32.Ident4())
	shader.SetUniformMat4("view", utils.Mat4F64ToF32(viewerContext.Orientation.Mat4().Inv()))
	shader.SetUniformMat4("projection", utils.Mat4F64ToF32(viewerContext.ProjectionMatrix))
	for i := 0; i < 6; i++ {
		gl.BindTexture(gl.TEXTURE_2D, textures[i].ID)
		gl.DrawArrays(gl.TRIANGLES, int32(i*6), 6)
	}
}

// drawText draws text at an x,y position that represents a fractional placement (0 -> 1)
// drawText expects the glyphs within `font` to be of equal width and height
func drawText(shader *shaders.ShaderProgram, font font.Font, text string, x, y float32) {
	var vertices []float32

	// assuming the height of all glyphs are equal - may not be the case in the future
	var glyphHeight float32
	for _, glyph := range font.Glyphs {
		glyphHeight = float32(glyph.Height)
		break
	}

	// convert porportion to pixel value
	x = x * float32(settings.Width)
	y = float32(settings.Height)*(1-y) - float32(glyphHeight)

	var xOffset float32
	var yOffset float32

	textureID := font.TextureID
	for _, c := range text {
		stringChar := string(c)

		if stringChar == "\n" {
			xOffset = 0
			yOffset++
			continue
		}

		glyph := font.Glyphs[stringChar]
		if _, ok := font.Glyphs[stringChar]; !ok {
			panic(fmt.Sprintf("glyph %s not found in font", stringChar))
		}

		width := float32(glyph.Width)
		height := float32(glyph.Height)

		textureX := float32(glyph.TextureCoords.X())
		textureY := float32(glyph.TextureCoords.Y())
		widthTextureCoord := (float32(glyph.Width) / float32(font.TotalWidth))
		heightTextureCoord := (float32(glyph.Height) / float32(font.TotalHeight))

		var characterVertices []float32 = []float32{
			xOffset * width, -(yOffset * glyphHeight), -5, textureX, textureY,
			(xOffset + 1) * width, -(yOffset * glyphHeight), -5, textureX + widthTextureCoord, textureY,
			(xOffset + 1) * width, height - (yOffset * glyphHeight), -5, textureX + widthTextureCoord, heightTextureCoord,

			(xOffset + 1) * width, height - (yOffset * glyphHeight), -5, textureX + widthTextureCoord, heightTextureCoord,
			xOffset * width, height - (yOffset * glyphHeight), -5, textureX, heightTextureCoord,
			xOffset * width, -(yOffset * glyphHeight), -5, textureX, textureY,
		}

		xOffset += 1
		vertices = append(vertices, characterVertices...)
	}

	// offset based on passed in x, y position which is constant across all characters
	for i := 0; i < len(vertices); i += 5 {
		vertices[i] = vertices[i] + x
		vertices[i+1] = vertices[i+1] + y
	}

	var vbo, vao uint32
	gl.GenBuffers(1, &vbo)
	gl.GenVertexArrays(1, &vao)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(vao)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, textureID)

	shader.Use()
	shader.SetUniformMat4("model", mgl32.Ident4())
	shader.SetUniformMat4("view", mgl32.Ident4())
	shader.SetUniformMat4("projection", mgl32.Ortho(0, float32(settings.Width), 0, float32(settings.Height), 1, 100))

	numCharacters := len(text)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(numCharacters*6))
}

// // drawHUDTextureToQuad does a shitty perspective based rendering of a flat texture
// func drawTexture(shader *shaders.ShaderProgram, texture uint32, x, y, width, height float32) {
// 	// convert porportion to pixel value
// 	x = x * float32(settings.Width)
// 	y = float32(settings.Height)*(1-y) - height

// 	// texture coords top left = 0,0 | bottom right = 1,1
// 	var vertices []float32 = []float32{
// 		0, 0, -5, 0.0, 0.0,
// 		width, 0, -5, 1.0, 0.0,
// 		width, height, -5, 1.0, 1.0,
// 		width, height, -5, 1.0, 1.0,
// 		0, height, -5, 0.0, 1.0,
// 		0, 0, -5, 0.0, 0.0,
// 	}

// 	for i := 0; i < len(vertices); i += 5 {
// 		vertices[i] = vertices[i] + x
// 		vertices[i+1] = vertices[i+1] + y
// 	}

// 	var vbo, vao uint32
// 	gl.GenBuffers(1, &vbo)
// 	gl.GenVertexArrays(1, &vao)

// 	gl.BindVertexArray(vao)
// 	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
// 	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

// 	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
// 	gl.EnableVertexAttribArray(0)

// 	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
// 	gl.EnableVertexAttribArray(1)

// 	gl.BindVertexArray(vao)
// 	gl.ActiveTexture(gl.TEXTURE0)
// 	gl.BindTexture(gl.TEXTURE_2D, texture)

// 	shader.Use()
// 	// shader.SetUniformMat4("model", mgl32.Translate3D(1.2, 0.8, -2))
// 	shader.SetUniformMat4("model", mgl32.Ident4())
// 	shader.SetUniformMat4("view", mgl32.Ident4())
// 	shader.SetUniformMat4("projection", mgl32.Ortho(0, float32(settings.Width), 0, float32(settings.Height), 1, 100))

// 	gl.DrawArrays(gl.TRIANGLES, 0, 6)
// }

// drawHUDTextureToQuad does a shitty perspective based rendering of a flat texture
func drawHUDTextureToQuad(viewerContext ViewerContext, shader *shaders.ShaderProgram, texture uint32, hudScale float32) {
	// texture coords top left = 0,0 | bottom right = 1,1
	var vertices []float32 = []float32{
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
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(vao)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	shader.Use()
	shader.SetUniformMat4("model", mgl32.Translate3D(1.2, 0.8, -2))
	shader.SetUniformMat4("view", mgl32.Ident4())
	shader.SetUniformMat4("projection", utils.Mat4F64ToF32(viewerContext.ProjectionMatrix))

	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func createModelMatrix(scaleMatrix, rotationMatrix, translationMatrix mgl64.Mat4) mgl64.Mat4 {
	return translationMatrix.Mul4(rotationMatrix).Mul4(scaleMatrix)
}

func resetGLRenderSettings() {
	gl.BindVertexArray(0)
	gl.UseProgram(0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.CullFace(gl.BACK)
}
