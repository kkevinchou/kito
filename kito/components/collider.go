package components

import (
	"github.com/kkevinchou/kito/lib/collision"
	"github.com/kkevinchou/kito/lib/collision/collider"
)

type ColliderComponent struct {
	CapsuleCollider *collider.Capsule
	TriMeshCollider *collider.TriMesh
	ContactManifold *collision.ContactManifold
}

func (c *ColliderComponent) AddToComponentContainer(container *ComponentContainer) {
	container.ColliderComponent = c
}
