package render

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kkevinchou/ant/lib/models"
)

func (r *RenderSystem) renderModel(model *models.Model) {
	for _, face := range model.Faces {
		gl.Color3f(float32(face.Material.Diffuse.X), float32(face.Material.Diffuse.Y), float32(face.Material.Diffuse.Z))
		gl.Begin(gl.POLYGON)
		for _, vertex := range face.Verticies {
			normal := vertex.Normal
			v := vertex.Vertex
			gl.Normal3f(float32(normal.X), float32(normal.Y), float32(normal.Z))
			gl.Vertex3f(float32(v.X), float32(v.Y)+1, float32(v.Z))
		}
		gl.End()
	}
}
