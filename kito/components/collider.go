package components

import (
	"github.com/kkevinchou/kito/lib/collision"
	"github.com/kkevinchou/kito/lib/collision/collider"
)

// func (c *CollisionInstance) String() string {
// 	var result = fmt.Sprintf("{ EntityID: %d ", 69)
// 	for i, contact := range c.Contacts {
// 		result += fmt.Sprintf("[ %v ]", contact)
// 		if i < len(c.Contacts)-1 {
// 			result += ", "
// 		}
// 	}
// 	result += " }"
// 	return result
// }

type ColliderComponent struct {
	ContactCandidates []*collision.Contact

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
