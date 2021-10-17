package collision

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
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
		cc := e.GetComponentContainer()

		if cc.ColliderComponent.CapsuleCollider != nil {
			capsule := cc.ColliderComponent.CapsuleCollider.Transform(cc.TransformComponent.Position)
			cc.ColliderComponent.TransformedCapsuleCollider = &capsule
		} else if cc.ColliderComponent.TriMeshCollider != nil {
			transformMatrix := mgl64.Translate3D(cc.TransformComponent.Position.X(), cc.TransformComponent.Position.Y(), cc.TransformComponent.Position.Z())
			triMesh := cc.ColliderComponent.TriMeshCollider.Transform(transformMatrix)
			cc.ColliderComponent.TransformedTriMeshCollider = &triMesh

		}

		e.GetComponentContainer().ColliderComponent.ContactManifolds = nil
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
					contactManifolds := collision.CheckCollisionCapsuleTriMesh(*e1cc.ColliderComponent.TransformedCapsuleCollider, *e2cc.ColliderComponent.TransformedTriMeshCollider)
					if contactManifolds != nil {
						e1cc.ColliderComponent.ContactManifolds = contactManifolds
						// fmt.Printf(
						// 	"[%d] collision ids[%d, %d]\n    capsule transform %v\n    capBottom %v\n    point %v\n    separating vector%v\n",
						// 	s.world.GetSingleton().CommandFrame,
						// 	e1.GetID(),
						// 	e2.GetID(),
						// 	transformComponent.Position,
						// 	capsule.Bottom,
						// 	contactManifold.Contacts[0].Point,
						// 	contactManifold.Contacts[0].SeparatingVector,
						// )
					}
				}
			}
		}
	}

}
