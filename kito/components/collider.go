package components

import "github.com/kkevinchou/kito/lib/collision/primitives"

type ColliderComponent struct {
	CapsuleCollider *primitives.Capsule
	TriMeshCollider *primitives.TriMesh
}

func (c *ColliderComponent) AddToComponentContainer(container *ComponentContainer) {
	container.ColliderComponent = c
}
