package model

import (
	"github.com/kkevinchou/kito/lib/modelspec"
)

type Model struct {
	// Mesh *Mesh
	// Animation  *Animation
	Animations map[string]*Animation

	vao    uint32
	meshes []*Mesh
}

// NewModel takes a ModelSpecification and performs the necessary OpenGL operations
// to pack the vertex and joint data into vertex buffers. It also holds animation
// key frame data for the animation system
func NewModel(spec *modelspec.ModelSpecification) *Model {
	// mesh := NewMesh(spec.Meshes[0])
	var animations map[string]*Animation
	if spec.Animations != nil {
		animations = NewAnimations(spec)
	}

	var meshes []*Mesh
	for _, ms := range spec.Meshes {
		meshes = append(meshes, NewMesh(ms))
	}

	return &Model{
		// Mesh:       mesh,
		meshes:     meshes,
		Animations: animations,
	}
}

func (m *Model) Meshes() []*Mesh {
	return m.meshes
}

func (m *Model) Prepare() {
	for _, mesh := range m.meshes {
		mesh.Prepare()
	}
}

// func (m *Model) Bind() uint32 {
// 	vao := configureVAO()
// 	m.vao = vao

// 	m.Mesh.BindVertexAttributes()
// 	if len(m.Animations) > 0 {
// 		// TODO: is this safe?
// 		// assume vertex attributes of first animation is the same for all animations
// 		for _, animation := range m.Animations {
// 			animation.BindVertexAttributes()
// 			break
// 		}
// 	}

// 	configureIndexBuffer(m.Mesh.vertexCount)
// 	return vao
// }

// func configureVAO() uint32 {
// 	var vao uint32
// 	gl.GenVertexArrays(1, &vao)
// 	gl.BindVertexArray(vao)
// 	return vao
// }

// func configureIndexBuffer(vertexCount int) {
// 	indices := []uint32{}
// 	for i := 0; i < vertexCount; i++ {
// 		indices = append(indices, uint32(i))
// 	}

// 	var ebo uint32
// 	gl.GenBuffers(1, &ebo)
// 	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
// 	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
// }
