package render

import (
	"errors"
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl64"
)

func initializeShadowMap(width, height int32) (uint32, uint32, error) {
	var depthMapFBO uint32
	gl.GenFramebuffers(1, &depthMapFBO)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT,
		width, height, 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	gl.BindFramebuffer(gl.FRAMEBUFFER, depthMapFBO)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, texture, 0)
	gl.DrawBuffer(gl.NONE)
	gl.ReadBuffer(gl.NONE)

	defer gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		return 999, 999, errors.New("failed to initialize shadow map frame buffer - who knows why?")
	}

	return depthMapFBO, texture, nil
}

func calculateFrustumVertices(near, far, fovy, aspectRatio float64) []mgl64.Vec3 {
	halfY := 2 * math.Tan(mgl64.DegToRad(fovy/2))
	halfX := aspectRatio * halfY
	_ = halfX
	return nil

}
