package collision

import (
	"fmt"
	"time"

	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/lib/collision"
)

type World interface {
	GetSingleton() *singleton.Singleton
	GetEntityByID(id int) (entities.Entity, error)
	GetPlayer() entities.Entity
}

type CollisionSystem struct {
	*base.BaseSystem
	world    World
	entities []entities.Entity
}

func NewCollisionSystem(world World) *CollisionSystem {
	return &CollisionSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}
}

func (s *CollisionSystem) RegisterEntity(entity entities.Entity) {
	componentContainer := entity.GetComponentContainer()

	if componentContainer.ColliderComponent != nil && componentContainer.TransformComponent != nil {
		s.entities = append(s.entities, entity)
	}
}

func (s *CollisionSystem) Update(delta time.Duration) {
	handledCollisions := map[int]map[int]bool{}
	for _, e := range s.entities {
		e.GetComponentContainer().ColliderComponent.ContactManifold = nil
		handledCollisions[e.GetID()] = map[int]bool{}
	}

	for _, e1 := range s.entities {
		e1cc := e1.GetComponentContainer()
		for _, e2 := range s.entities {
			e2cc := e2.GetComponentContainer()

			// don't check an entity against itself or if we've already computed collisions
			if e1.GetID() == e2.GetID() {
				continue
			}
			if handledCollisions[e1.GetID()][e2.GetID()] {
				continue
			}

			if e1cc.ColliderComponent.CapsuleCollider != nil {
				if e2cc.ColliderComponent.TriMeshCollider != nil {
					transformComponent := e1cc.TransformComponent
					capsule := e1cc.ColliderComponent.CapsuleCollider.Transform(transformComponent.Position)
					contactManifold := collision.CheckCollisionCapsuleTriMesh(capsule, *e2cc.ColliderComponent.TriMeshCollider)
					if contactManifold != nil {
						e1cc.ColliderComponent.ContactManifold = contactManifold
						fmt.Printf("collision detected %v, %f\n", contactManifold.Contacts[0].Point, contactManifold.Contacts[0].SeparatingVector)
					}
				}
			}
		}
	}

}
