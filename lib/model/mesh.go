package model

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/lib/modelspec"
)

type Mesh struct {
	vertexCount int
	vertices    []mgl64.Vec3
}

func NewMesh(spec *modelspec.ModelSpecification) *Mesh {
	vertexAttributes, totalAttributeSize, vertices := constructMeshVertexAttributes(
		spec.TriIndices,
		spec.TriIndicesStride,
		spec.PositionSourceData,
		spec.NormalSourceData,
		spec.ColorSourceData,
		spec.TextureSourceData,
	)

	configureMeshVertexAttributes(vertexAttributes, totalAttributeSize)

	return &Mesh{
		vertexCount: len(vertexAttributes) / totalAttributeSize,
		vertices:    vertices,
	}
}

func (m *Mesh) Vertices() []mgl64.Vec3 {
	return m.vertices
}

func constructMeshVertexAttributes(
	triIndices []int,
	triIndicesStride int,
	positionSourceData []mgl32.Vec3,
	normalSourceData []mgl32.Vec3,
	colorSourceData []mgl32.Vec3,
	textureSourceData []mgl32.Vec2,
) ([]float32, int, []mgl64.Vec3) {
	vertexAttributes := []float32{}

	totalAttributeSize := len(positionSourceData[0]) + len(normalSourceData[0]) + len(textureSourceData[0])

	colorPresent := colorSourceData != nil
	if colorPresent {
		totalAttributeSize += len(colorSourceData[0])
	}

	var vertices []mgl64.Vec3

	// TODO: i'm still ordering vertex attributes by the face order, rather than keeping the original exported source order
	// this current way will repeat data since i explicity store data for every vertex, rather than using indicies for lookup
	// in the future, i should refactor this to store the data in source data order then use an index buffer for VAO creation

	// triIndicies format: position, normal, texture, color
	for i := 0; i < len(triIndices); i += triIndicesStride {
		// TODO: we are assuming this ordering of position, normal, texture but this is not
		// necessarily the case. it depends on the <input> elements are ordered in the collada file
		position := positionSourceData[triIndices[i]]
		normal := normalSourceData[triIndices[i+1]]
		texture := textureSourceData[triIndices[i+2]]

		vertexAttributes = append(vertexAttributes, position.X(), position.Y(), position.Z())
		vertexAttributes = append(vertexAttributes, normal.X(), normal.Y(), normal.Z())
		vertexAttributes = append(vertexAttributes, texture.X(), texture.Y())

		if colorPresent {
			color := colorSourceData[triIndices[i+3]]
			vertexAttributes = append(vertexAttributes, color.X(), color.Y(), color.Z())
		}
		vertices = append(vertices, mgl64.Vec3{float64(position.X()), float64(position.Y()), float64(position.Z())})
	}

	return vertexAttributes, totalAttributeSize, vertices
}

// lays out the vertex atrributes for:
// 0 - position         vec3
// 1 - normal           vec3
// 2 - texture coord    vec2
// 3 - color            vec3
func configureMeshVertexAttributes(vertexAttributes []float32, totalAttributeSize int) {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexAttributes)*4, gl.Ptr(vertexAttributes), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(totalAttributeSize)*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(totalAttributeSize)*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, int32(totalAttributeSize)*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	gl.VertexAttribPointer(3, 3, gl.FLOAT, false, int32(totalAttributeSize)*4, gl.PtrOffset(8*4))
	gl.EnableVertexAttribArray(3)
}
