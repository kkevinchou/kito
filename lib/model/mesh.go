package model

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/lib/modelspec"
)

type Mesh struct {
	vertexCount        int
	vertices           []mgl64.Vec3
	vertexAttributes   []float32
	totalAttributeSize int
	material           *modelspec.EffectSpec
	pbrMaterial        *modelspec.PBRMaterial
}

func NewMesh(spec *modelspec.ModelSpecification) *Mesh {
	vertexAttributes, totalAttributeSize, vertices := constructMeshVertexAttributes(spec)

	return &Mesh{
		vertexCount:        len(vertexAttributes) / totalAttributeSize,
		vertices:           vertices,
		vertexAttributes:   vertexAttributes,
		totalAttributeSize: totalAttributeSize,
		material:           spec.EffectSpecData,
		pbrMaterial:        spec.Meshes[0].PBRMaterial,
	}
}

func (m *Mesh) Vertices() []mgl64.Vec3 {
	return m.vertices
}

func (m *Mesh) Material() *modelspec.EffectSpec {
	return m.material
}

func (m *Mesh) PBRMaterial() *modelspec.PBRMaterial {
	return m.pbrMaterial
}

func constructMeshVertexAttributes(
	spec *modelspec.ModelSpecification,
) ([]float32, int, []mgl64.Vec3) {
	var vertices []mgl64.Vec3
	vertexAttributes := []float32{}

	for _, mesh := range spec.Meshes {
		positionSourceData := mesh.PositionSourceData
		normalSourceData := mesh.NormalSourceData
		textureSourceData := mesh.TextureSourceData
		vertexAttributeIndices := mesh.VertexAttributeIndices
		vertexAttributesStride := mesh.VertexAttributesStride

		if mesh.VertexAttributesStride <= 0 {
			panic(fmt.Sprintf("unexpected stride value %d", mesh.VertexAttributesStride))
		}

		// triIndicies format: position, normal, texture, color
		for i := 0; i < len(vertexAttributeIndices); i += vertexAttributesStride {
			// TODO: we are assuming this ordering of position, normal, texture but this is not
			// necessarily the case. it depends on the <input> elements are ordered in the collada file
			position := positionSourceData[vertexAttributeIndices[i]]
			normal := normalSourceData[vertexAttributeIndices[i+1]]
			texture := textureSourceData[vertexAttributeIndices[i+2]]

			vertexAttributes = append(vertexAttributes, position.X(), position.Y(), position.Z())
			vertexAttributes = append(vertexAttributes, normal.X(), normal.Y(), normal.Z())
			vertexAttributes = append(vertexAttributes, texture.X(), texture.Y())

			vertices = append(vertices, mgl64.Vec3{float64(position.X()), float64(position.Y()), float64(position.Z())})
		}
	}

	totalAttributeSize := len(spec.Meshes[0].PositionSourceData[0]) + len(spec.Meshes[0].NormalSourceData[0]) + len(spec.Meshes[0].TextureSourceData[0])
	return vertexAttributes, totalAttributeSize, vertices
}

// lays out the vertex atrributes for:
// 0 - position         vec3
// 1 - normal           vec3
// 2 - texture coord    vec2
// 3 - color            vec3
func (m *Mesh) BindVertexAttributes() {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(m.vertexAttributes)*4, gl.Ptr(m.vertexAttributes), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(m.totalAttributeSize)*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(m.totalAttributeSize)*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, int32(m.totalAttributeSize)*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	// TODO: deprecate this, color not really useful
	// gl.VertexAttribPointer(3, 3, gl.FLOAT, false, int32(m.totalAttributeSize)*4, gl.PtrOffset(8*4))
	// gl.EnableVertexAttribArray(3)
}
