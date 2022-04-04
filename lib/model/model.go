package model

import (
	"github.com/kkevinchou/kito/lib/modelspec"
)

type Model struct {
	animations map[string]*modelspec.AnimationSpec
	meshes     []*Mesh
	rootJoint  *modelspec.JointSpec
}

// NewModel takes a ModelSpecification and performs the necessary OpenGL operations
// to pack the vertex and joint data into vertex buffers. It also holds animation
// key frame data for the animation system
func NewModel(spec *modelspec.ModelSpecification) *Model {
	var meshes []*Mesh
	for _, ms := range spec.Meshes {
		meshes = append(meshes, NewMesh(ms))
	}

	return &Model{
		meshes:     meshes,
		animations: spec.Animations,
		rootJoint:  spec.RootJoint,
	}
}

func (m *Model) RootJoint() *modelspec.JointSpec {
	return m.rootJoint
}

func (m *Model) Animations() map[string]*modelspec.AnimationSpec {
	return m.animations
}

func (m *Model) Meshes() []*Mesh {
	return m.meshes
}

func (m *Model) Vertices() []modelspec.Vertex {
	var vertices []modelspec.Vertex
	for _, mesh := range m.meshes {
		meshVerts := mesh.Vertices()
		vertices = append(vertices, meshVerts...)
	}
	return vertices
}

func (m *Model) Prepare() {
	for _, mesh := range m.meshes {
		mesh.Prepare()
	}
}
