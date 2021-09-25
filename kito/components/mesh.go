package components

import "github.com/kkevinchou/kito/lib/textures"

type MeshComponent struct {
	// should probably store this data in a separate component
	ModelVAO         uint32
	ModelVertexCount int
	Texture          *textures.Texture
}

func (c *MeshComponent) GetMeshComponent() *MeshComponent {
	return c
}

func (c *MeshComponent) AddToComponentContainer(container *ComponentContainer) {
	container.MeshComponent = c
}
