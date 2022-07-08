package components

import (
	"github.com/kkevinchou/kito/lib/collision"
	"github.com/kkevinchou/kito/lib/collision/collider"
)

type ColliderComponent struct {
	SkipMovementResolution bool

	// some field that marks which entities it collided with in the current frame
	Contacts map[int]*collision.Contact

	CapsuleCollider *collider.Capsule
	TriMeshCollider *collider.TriMesh

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
