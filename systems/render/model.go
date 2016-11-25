package render

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kkevinchou/ant/lib/models"
)

func (r *RenderSystem) renderModel(model *models.Model) {
	color := []float32{0.5, 0.5, 0.5}
	for _, face := range model.Faces {
		gl.Begin(gl.POLYGON)
		for _, vertex := range face.Verticies {
			normal := vertex.Normal
			v := vertex.Vertex
			gl.Normal3f(float32(normal.X), float32(normal.Y), float32(normal.Z))
			gl.Color3f(color[0], color[1], color[2])
			gl.Vertex3f(float32(v.X), float32(v.Y), float32(v.Z))
		}
		gl.End()
	}
}
