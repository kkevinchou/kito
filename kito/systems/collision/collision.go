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
		for _, e2 := range s.world.QueryEntity(components.ComponentFlagCollider | components.ComponentFlagTransform) {
			if handledCollisions[player.GetID()][e2.GetID()] || handledCollisions[e2.GetID()][player.GetID()] {
				continue
			}
			s.collide(player, e2, handledCollisions)
			handledCollisions[player.GetID()][e2.GetID()] = true
			handledCollisions[e2.GetID()][player.GetID()] = true
		}
	} else {
		for _, e1 := range s.world.QueryEntity(components.ComponentFlagCollider | components.ComponentFlagTransform) {
			for _, e2 := range s.world.QueryEntity(components.ComponentFlagCollider | components.ComponentFlagTransform) {
				if handledCollisions[e1.GetID()][e2.GetID()] || handledCollisions[e2.GetID()][e1.GetID()] {
					continue
				}
				s.collide(e1, e2, handledCollisions)
				handledCollisions[e1.GetID()][e2.GetID()] = true
				handledCollisions[e2.GetID()][e1.GetID()] = true
			}
		}
	}
}

func (s *CollisionSystem) collide(e1 entities.Entity, e2 entities.Entity, handledCollisions map[int]map[int]bool) {
	// don't check an entity against itself or if we've already computed collisions
	if e1.GetID() == e2.GetID() {
		return
	}

	isSupportedCollision, capsuleEntity, triMeshEntity := isCapsuleTriMeshCollision(e1, e2)
	if !isSupportedCollision {
		return
	}

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
}

func isCapsuleTriMeshCollision(e1, e2 entities.Entity) (bool, entities.Entity, entities.Entity) {
	e1cc := e1.GetComponentContainer()
	e2cc := e2.GetComponentContainer()

	if e1cc.ColliderComponent.CapsuleCollider != nil {
		if e2cc.ColliderComponent.TriMeshCollider != nil {
			return true, e1, e2
		}
	}

	if e2cc.ColliderComponent.CapsuleCollider != nil {
		if e1cc.ColliderComponent.TriMeshCollider != nil {
			return true, e2, e1
		}
	}

	return false, nil, nil
}
