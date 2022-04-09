package components

import (
	"github.com/kkevinchou/kito/lib/collision"
	"github.com/kkevinchou/kito/lib/collision/collider"
)

type CollisionInstance struct {
	OtherEntityID    int
	ContactManifolds []*collision.ContactManifold
}

type ColliderComponent struct {
	CapsuleCollider    *collider.Capsule
	TriMeshCollider    *collider.TriMesh
	CollisionInstances []*CollisionInstance

	// stores the transformed collider (e.g. if the entity moves)
	TransformedCapsuleCollider *collider.Capsule
	TransformedTriMeshCollider *collider.TriMesh
}

func (c *ColliderComponent) AddToComponentContainer(container *ComponentContainer) {
	container.ColliderComponent = c
}

func (c *ColliderComponent) ComponentFlag() int {
	return ComponentFlagCollider
}
