package model

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/kkevinchou/kito/lib/modelspec"
)

type Model struct {
	Mesh      *Mesh
	Animation *Animation

	vao uint32
}

// NewModel takes a ModelSpecification and performs the necessary OpenGL operations
// to pack the vertex and joint data into vertex buffers. It also holds animation
// key frame data for the animation system
func NewModel(spec *modelspec.ModelSpecification) *Model {
	mesh := NewMesh(spec)

	var animation *Animation
	if spec.Animation != nil {
		animation = NewAnimation(spec)
	}

	return &Model{
		Mesh:      mesh,
		Animation: animation,
	}
}

func (m *Model) VAO() uint32 {
	return m.vao
}

func (m *Model) VertexCount() int {
	return m.Mesh.vertexCount
}

func (m *Model) Bind() uint32 {
	vao := configureVAO()
	m.vao = vao

	m.Mesh.BindVertexAttributes()
	if m.Animation != nil {
		m.Animation.BindVertexAttributes()
	}

	configureIndexBuffer(m.Mesh.vertexCount)
	return vao
}

func configureVAO() uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	return vao
}

func configureIndexBuffer(vertexCount int) {
	indices := []uint32{}
	for i := 0; i < vertexCount; i++ {
		indices = append(indices, uint32(i))
	}

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
}
