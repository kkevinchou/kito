package render

import (
	"github.com/go-gl/gl/v4.6-compatibility/gl"
	"github.com/kkevinchou/kito/lib/math/vector"
	"github.com/kkevinchou/kito/lib/models"
)

func (r *RenderSystem) renderModel(model *models.Model, position vector.Vector3) {
	for _, face := range model.Faces {
		gl.Begin(gl.POLYGON)
		gl.Color3f(float32(face.Material.Diffuse.X), float32(face.Material.Diffuse.Y), float32(face.Material.Diffuse.Z))
		for _, vertex := range face.Verticies {
			normal := vertex.Normal
			v := vertex.Vertex
			specReflection := []float32{
				float32(face.Material.Specular.X),
				float32(face.Material.Specular.Y),
				float32(face.Material.Specular.X),
				1.0,
			}
			gl.Materialfv(gl.FRONT, gl.SPECULAR, &specReflection[0])
			gl.Materiali(gl.FRONT, gl.SHININESS, 56)
			gl.Normal3f(float32(normal.X), float32(normal.Y), float32(normal.Z))
			gl.Vertex3f(float32(v.X+position.X), float32(v.Y+position.Y), float32(v.Z+position.Z))
		}
		gl.End()
	}
}
