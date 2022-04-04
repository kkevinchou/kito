package model

import (
	"github.com/kkevinchou/kito/lib/modelspec"
)

type Model struct {
	Animations map[string]*Animation
	meshes     []*Mesh
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
