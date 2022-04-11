package model

import (
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/lib/modelspec"
)

type Model struct {
	meshes    []*Mesh
	modelSpec *modelspec.ModelSpecification
}

func NewModel(spec *modelspec.ModelSpecification) *Model {
	var meshes []*Mesh
	for _, ms := range spec.Meshes {
		meshes = append(meshes, NewMesh(ms))
	}

	m := &Model{
		modelSpec: spec,
		meshes:    meshes,
	}

	if utils.IsClient() {
		m.prepare()
	}

	return m
}

func (m *Model) RootJoint() *modelspec.JointSpec {
	return m.modelSpec.RootJoint
}

func (m *Model) Animations() map[string]*modelspec.AnimationSpec {
	return m.modelSpec.Animations
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

func (m *Model) prepare() {
	for _, mesh := range m.meshes {
		mesh.prepare(m.modelSpec.Textures)
	}
}

func (m *Model) MeshChunks() []*MeshChunk {
	var meshChunks []*MeshChunk
	for _, mesh := range m.Meshes() {
		meshChunks = append(meshChunks, mesh.MeshChunks()...)
	}
	return meshChunks
}
