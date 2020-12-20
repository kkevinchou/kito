package animation

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/kkevinchou/kito/lib/loaders/collada"
)

type Mesh struct {
	vao         uint32
	vertexCount int
}

func NewMesh() *Mesh {
	return &Mesh{}
}

func NewMeshFromCollada(c *collada.Collada) *Mesh {
	vertexCount := len(c.VertexFloatsSourceData) / 3
	dataSize := 8

	vertexAttributes := constructVertexAttributes(
		c.VertexFloatsSourceData,
		c.NormalFloatsSourceData,
		c.TextureFloatsSourceData,
		vertexCount,
		dataSize,
	)

	return &Mesh{
		vao:         constructVAO(vertexAttributes, dataSize),
		vertexCount: vertexCount,
	}
}

func (m *Mesh) VAO() uint32 {
	return m.vao
}

func (m *Mesh) VertexCount() int {
	return m.vertexCount
}

func constructVertexAttributes(vertices []float32, normals []float32, texCoords []float32, vertexCount, dataSize int) []float32 {
	vertexAttributes := make([]float32, vertexCount*dataSize)

	for i := 0; i < vertexCount; i++ {
		threeStepIndex := i * 3
		// twoStepIndex := i * 2
		vertexAttributeStepIndex := i * 8
		vertexAttributes[vertexAttributeStepIndex] = vertices[threeStepIndex]
		vertexAttributes[vertexAttributeStepIndex+1] = vertices[threeStepIndex+1]
		vertexAttributes[vertexAttributeStepIndex+2] = vertices[threeStepIndex+2]

		vertexAttributes[vertexAttributeStepIndex+3] = normals[threeStepIndex]
		vertexAttributes[vertexAttributeStepIndex+4] = normals[threeStepIndex+1]
		vertexAttributes[vertexAttributeStepIndex+5] = normals[threeStepIndex+2]

		// vertexAttributes[strideStepIndex+6] = texCoords[twoStepIndex]
		// vertexAttributes[strideStepIndex+7] = texCoords[twoStepIndex+1]
		vertexAttributes[vertexAttributeStepIndex+6] = 0
		vertexAttributes[vertexAttributeStepIndex+7] = 0
	}

	return vertexAttributes
}

func constructVAO(vertexAttributes []float32, dataSize int) uint32 {
	var vbo, vao uint32
	gl.GenBuffers(1, &vbo)
	gl.GenVertexArrays(1, &vao)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexAttributes)*4, gl.Ptr(vertexAttributes), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(dataSize)*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(dataSize)*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, int32(dataSize)*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	return vao
}
