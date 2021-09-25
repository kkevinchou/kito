package model

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/kkevinchou/kito/lib/modelspec"
)

const (
	maxWeights int = 3
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
	var animation *Animation
	vao := configureVAO()
	mesh := NewMesh(spec)
	if spec.Animation != nil {
		animation = NewAnimation(spec)
	}
	configureIndexBuffer(mesh.vertexCount)

	return &Model{
		Mesh:      mesh,
		Animation: animation,

		vao: vao,
	}
}

func (m *Model) VAO() uint32 {
	return m.vao
}
func (m *Model) VertexCount() int {
	return m.Mesh.vertexCount
}

func NewMeshedModel(spec *modelspec.ModelSpecification) *Model {
	vao := configureVAO()
	mesh := NewMesh(spec)
	configureIndexBuffer(mesh.vertexCount)

	return &Model{
		Mesh: mesh,
		vao:  vao,
	}
}

// TODO: NewPlaceholderModel is meant for the server to create models without
// the related GL calls. For now this is just to allow simualtion of animations
// on the backend. this may not actually be needed
func NewPlaceholderModel(spec *modelspec.ModelSpecification) *Model {
	// if we're on the server and it's a static model, there will be no root joint
	// skip the animation stuff
	var animation *Animation
	if spec.Root != nil {
		animation = NewJointOnlyAnimation(spec)
	}
	return &Model{
		Animation: animation,
	}
}

func configureVAO() uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	return vao
}

func configureIndexBuffer(vertexCount int) {
	// super inefficient, the benefit of an EBO is that we don't have to store duplicate vertices.
	// since we haven't actually removed duplicate vertices yet and our EBO isn't intelligently
	// pointing to older vertices we do the dumb thing of keeping all the vertices and pointing
	// to each in order (we store the vertices in draw order)
	indices := []uint32{}
	for i := 0; i < vertexCount; i++ {
		indices = append(indices, uint32(i))
	}

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
}
