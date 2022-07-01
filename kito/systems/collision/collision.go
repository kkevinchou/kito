package collision

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/managers/player"
	"github.com/kkevinchou/kito/kito/netsync"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/lib/collision"
)

type World interface {
	GetSingleton() *singleton.Singleton
	GetPlayerEntity() entities.Entity
	QueryEntity(componentFlags int) []entities.Entity
	GetPlayer() *player.Player
	GetEntityByID(id int) entities.Entity
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
	collisionDetected := true
	iterationMax := 10
	count := 0

	for collisionDetected && count < iterationMax {
		collisionDetected = s.findCollisionCandidates()
		if collisionDetected {
			s.resolveCollisions()
		}
		count += 1
	}
}

func (s *CollisionSystem) resolveCollisions() {
	d := directory.GetDirectory()
	playerManager := d.PlayerManager()

	var players []*player.Player
	if utils.IsClient() {
		players = []*player.Player{s.world.GetPlayer()}
	} else {
		players = playerManager.GetPlayers()
	}

	for _, player := range players {
		entity := s.world.GetEntityByID(player.EntityID)
		if entity == nil {
			continue
		}
		netsync.ResolveControllerCollision(entity)
	}
}

func (s *CollisionSystem) findCollisionCandidates() bool {
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

		e.GetComponentContainer().ColliderComponent.ContactCandidates = nil
		handledCollisions[e.GetID()] = map[int]bool{}
	}

	checkPairs := [][]entities.Entity{}
	if utils.IsClient() {
		player := s.world.GetPlayerEntity()
		for _, e2 := range s.world.QueryEntity(components.ComponentFlagCollider | components.ComponentFlagTransform) {
			checkPairs = append(checkPairs, []entities.Entity{player, e2})
		}
	} else {
		for _, e1 := range s.world.QueryEntity(components.ComponentFlagCollider | components.ComponentFlagTransform) {
			for _, e2 := range s.world.QueryEntity(components.ComponentFlagCollider | components.ComponentFlagTransform) {
				checkPairs = append(checkPairs, []entities.Entity{e1, e2})
			}
		}
	}

	collisionFound := false
	for _, pair := range checkPairs {
		e1 := pair[0]
		e2 := pair[1]
		if handledCollisions[e1.GetID()][e2.GetID()] || handledCollisions[e2.GetID()][e1.GetID()] {
			continue
		}
		if s.collide(e1, e2) {
			collisionFound = true
		}
		handledCollisions[e1.GetID()][e2.GetID()] = true
		handledCollisions[e2.GetID()][e1.GetID()] = true
	}

	return collisionFound
}

func (s *CollisionSystem) collide(e1 entities.Entity, e2 entities.Entity) bool {
	// don't check an entity against itself or if we've already computed collisions
	if e1.GetID() == e2.GetID() {
		return false
	}

	isSupportedCollision, capsuleEntity, triMeshEntity := isCapsuleTriMeshCollision(e1, e2)
	if !isSupportedCollision {
		return false
	}

	contacts := collision.CheckCollisionCapsuleTriMesh(
		*capsuleEntity.GetComponentContainer().ColliderComponent.TransformedCapsuleCollider,
		*triMeshEntity.GetComponentContainer().ColliderComponent.TransformedTriMeshCollider,
	)

	if contacts != nil {
		for _, contact := range contacts {
			triEntityID := triMeshEntity.GetID()
			contact.EntityID = &triEntityID
		}

		capsuleEntity.GetComponentContainer().ColliderComponent.ContactCandidates = append(capsuleEntity.GetComponentContainer().ColliderComponent.ContactCandidates, contacts...)
		return true
	}
	return false
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
