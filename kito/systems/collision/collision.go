package collision

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/lib/collision"
)

type World interface {
	GetSingleton() *singleton.Singleton
	GetPlayerEntity() entities.Entity
	QueryEntity(componentFlags int) []entities.Entity
}

type CollisionSystem struct {
	*base.BaseSystem
	world World
}

func NewCollisionSystem(world World) *CollisionSystem {
	return &CollisionSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}
}

func (s *CollisionSystem) Update(delta time.Duration) {
	handledCollisions := map[int]map[int]bool{}

	// initialize collision state
	for _, e := range s.world.QueryEntity(components.ComponentFlagCollider | components.ComponentFlagTransform) {
		cc := e.GetComponentContainer()

		if cc.ColliderComponent.CapsuleCollider != nil {
			capsule := cc.ColliderComponent.CapsuleCollider.Transform(cc.TransformComponent.Position)
			cc.ColliderComponent.TransformedCapsuleCollider = &capsule
		} else if cc.ColliderComponent.TriMeshCollider != nil {
			transformMatrix := mgl64.Translate3D(cc.TransformComponent.Position.X(), cc.TransformComponent.Position.Y(), cc.TransformComponent.Position.Z())
			triMesh := cc.ColliderComponent.TriMeshCollider.Transform(transformMatrix)
			cc.ColliderComponent.TransformedTriMeshCollider = &triMesh
		}

		e.GetComponentContainer().ColliderComponent.CollisionInstances = nil
		handledCollisions[e.GetID()] = map[int]bool{}
	}

	if utils.IsClient() {
		player := s.world.GetPlayerEntity()
		s.collide(player, handledCollisions)
	} else {
		for _, e1 := range s.world.QueryEntity(components.ComponentFlagCollider | components.ComponentFlagTransform) {
			s.collide(e1, handledCollisions)
		}
	}
}

func (s *CollisionSystem) collide(e1 entities.Entity, handledCollisions map[int]map[int]bool) {
	for _, e2 := range s.world.QueryEntity(components.ComponentFlagCollider | components.ComponentFlagTransform) {
		// don't check an entity against itself or if we've already computed collisions
		if e1.GetID() == e2.GetID() {
			continue
		}
		if handledCollisions[e1.GetID()][e2.GetID()] || handledCollisions[e2.GetID()][e1.GetID()] {
			continue
		}

		isSupportedCollision, capsuleEntity, triMeshEntity := isCapsulteTriMeshCollision(e1, e2)
		if isSupportedCollision {
			contactManifolds := collision.CheckCollisionCapsuleTriMesh(*capsuleEntity.GetComponentContainer().ColliderComponent.TransformedCapsuleCollider, *triMeshEntity.GetComponentContainer().ColliderComponent.TransformedTriMeshCollider)
			if contactManifolds != nil {
				capsuleEntity.GetComponentContainer().ColliderComponent.CollisionInstances = append(
					capsuleEntity.GetComponentContainer().ColliderComponent.CollisionInstances,
					&components.CollisionInstance{
						OtherEntityID:    capsuleEntity.GetID(),
						ContactManifolds: contactManifolds,
					},
				)
			}
			// TODO: decide how we want to handle avoiding double calculating collisions
			handledCollisions[capsuleEntity.GetID()][triMeshEntity.GetID()] = true
			handledCollisions[triMeshEntity.GetID()][capsuleEntity.GetID()] = true
		}
	}
}

func isCapsulteTriMeshCollision(e1, e2 entities.Entity) (bool, entities.Entity, entities.Entity) {
	e1cc := e1.GetComponentContainer()
	e2cc := e2.GetComponentContainer()

	if e1cc.ColliderComponent.CapsuleCollider != nil {
		if e2cc.ColliderComponent.TriMeshCollider != nil {
			return true, e1, e2
		}
	}

	if e1cc.ColliderComponent.CapsuleCollider != nil {
		if e2cc.ColliderComponent.TriMeshCollider != nil {
			return true, e1, e2
		}
	}

	return false, nil, nil
}
