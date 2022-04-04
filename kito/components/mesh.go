package components

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/lib/model"
)

type MeshComponent struct {
	// should probably store this data in a separate component
	// ModelVAO         uint32
	// ModelVertexCount int
	Scale       mgl64.Mat4
	Orientation mgl64.Mat4
	// Material    *modelspec.EffectSpec
	// PBRMaterial *modelspec.PBRMaterial

	// Texture *textures.Texture
	Model *model.Model
}

func (c *MeshComponent) GetMeshComponent() *MeshComponent {
	return c
}

func (c *MeshComponent) AddToComponentContainer(container *ComponentContainer) {
	container.MeshComponent = c
}
