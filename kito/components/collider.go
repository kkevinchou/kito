package components

import (
	"github.com/kkevinchou/kito/lib/collision"
	"github.com/kkevinchou/kito/lib/collision/collider"
)

type ColliderComponent struct {
	CapsuleCollider  *collider.Capsule
	TriMeshCollider  *collider.TriMesh
	ContactManifolds []*collision.ContactManifold

	// stores the transformed collider (e.g. if the entity moves)
	TransformedCapsuleCollider *collider.Capsule
	TransformedTriMeshCollider *collider.TriMesh
}

func (c *ColliderComponent) AddToComponentContainer(container *ComponentContainer) {
	container.ColliderComponent = c
}
