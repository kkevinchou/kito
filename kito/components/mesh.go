package components

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/lib/modelspec"
	"github.com/kkevinchou/kito/lib/textures"
)

type MeshComponent struct {
	// should probably store this data in a separate component
	ModelVAO         uint32
	ModelVertexCount int
	Texture          *textures.Texture
	ShaderProgram    string
	Scale            mgl64.Mat4
	Orientation      mgl64.Mat4
	Material         *modelspec.EffectSpec
}

func (c *MeshComponent) GetMeshComponent() *MeshComponent {
	return c
}

func (c *MeshComponent) AddToComponentContainer(container *ComponentContainer) {
	container.MeshComponent = c
}
