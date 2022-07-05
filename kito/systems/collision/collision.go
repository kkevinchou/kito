package collision

import (
	"sort"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/managers/player"
	"github.com/kkevinchou/kito/kito/netsync"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/lib/collision"
)

const (
	// the maximum number of times a distinct entity can have their collision resolved
	// this presents the collision resolution phase to go on forever
	resolveCountMax = 10
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
	entityPairs := [][]entities.Entity{}
	if utils.IsClient() {
		player := s.world.GetPlayerEntity()
		for _, e2 := range s.world.QueryEntity(components.ComponentFlagCollider | components.ComponentFlagTransform) {
			entityPairs = append(entityPairs, []entities.Entity{player, e2})
		}
	} else {
		for _, e1 := range s.world.QueryEntity(components.ComponentFlagCollider | components.ComponentFlagTransform) {
			for _, e2 := range s.world.QueryEntity(components.ComponentFlagCollider | components.ComponentFlagTransform) {
				entityPairs = append(entityPairs, []entities.Entity{e1, e2})
			}
		}
	}

	// 1. collect pairs of entities that are colliding, sorted by separating vector
	// 2. perform collision resolution for any colliding entities
	// 3. this can cause more collisions, repeat until no more further detected collisions, or we hit the max

	resolveCount := map[int]int{}
	reachedMaxCount := false
	for !reachedMaxCount {
		collisionCandidates := s.collectSortedCollisionCandidates(entityPairs)
		if len(collisionCandidates) == 0 {
			break
		}

		resolvedEntities := s.resolveCollisions(collisionCandidates)
		for _, id := range resolvedEntities {
			resolveCount[id] += 1
			if resolveCount[id] > resolveCountMax {
				reachedMaxCount = true
			}
		}
	}
}

func (s *CollisionSystem) resolveCollisions(contacts []*collision.Contact) []int {
	seen := map[int]any{}
	for _, contact := range contacts {
		if _, ok := seen[*contact.EntityID]; ok {
			continue
		}
		entity := s.world.GetEntityByID(*contact.EntityID)
		sourceEntity := s.world.GetEntityByID(*contact.SourceEntityID)
		netsync.ResolveControllerCollision(entity, sourceEntity, contact)

		seen[*contact.EntityID] = true
		seen[*contact.SourceEntityID] = true
	}

	var resolvedEntities []int
	for id, _ := range seen {
		resolvedEntities = append(resolvedEntities, id)
	}
	return resolvedEntities
}

// collectSortedCollisionCandidates collects all potential collisions that can occur in the frame.
// these are "candidates" in that if we resolve some of the collisions in the list, some will be
// invalidated
func (s *CollisionSystem) collectSortedCollisionCandidates(entityPairs [][]entities.Entity) []*collision.Contact {
	seen := map[int]map[int]bool{}

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

		seen[e.GetID()] = map[int]bool{}
	}

	var allContacts []*collision.Contact
	for _, pair := range entityPairs {
		e1 := pair[0]
		e2 := pair[1]
		if seen[e1.GetID()][e2.GetID()] || seen[e2.GetID()][e1.GetID()] {
			continue
		}
		contacts := s.collide(e1, e2)
		// if len(contacts) > 0 {
		// 	if e2.GetID() == 80002 {
		// 		fmt.Println(1)
		// 	}
		// 	if e1.GetID() == 80002 {
		// 		fmt.Println(2)
		// 	}
		// }
		allContacts = append(allContacts, contacts...)

		seen[e1.GetID()][e2.GetID()] = true
		seen[e2.GetID()][e1.GetID()] = true
	}
	sort.Sort(contactsBySeparatingDistance(allContacts))

	return allContacts
}

func (s *CollisionSystem) collide(e1 entities.Entity, e2 entities.Entity) []*collision.Contact {
	if e1.GetID() == e2.GetID() {
		return nil
	}

	if ok, capsuleEntity, triMeshEntity := isCapsuleTriMeshCollision(e1, e2); ok {
		contacts := collision.CheckCollisionCapsuleTriMesh(
			*capsuleEntity.GetComponentContainer().ColliderComponent.TransformedCapsuleCollider,
			*triMeshEntity.GetComponentContainer().ColliderComponent.TransformedTriMeshCollider,
		)
		if len(contacts) == 0 {
			return nil
		}

		triEntityID := triMeshEntity.GetID()
		capsuleEntityID := capsuleEntity.GetID()

		for _, contact := range contacts {
			contact.EntityID = &capsuleEntityID
			contact.SourceEntityID = &triEntityID
			if *contact.EntityID == 80002 {
				// fmt.Println(s.world.GetSingleton().CommandFrame, utils.PPrintVec(capsuleEntity.GetComponentContainer().TransformComponent.Position), utils.PPrintVec(contact.SeparatingVector))
			}
		}

		return contacts
	}

	if ok := isCapsuleCapsuleCollision(e1, e2); ok {
		contact := collision.CheckCollisionCapsuleCapsule(
			*e1.GetComponentContainer().ColliderComponent.TransformedCapsuleCollider,
			*e2.GetComponentContainer().ColliderComponent.TransformedCapsuleCollider,
		)
		if contact == nil {
			return nil
		}

		e1ID := e1.GetID()
		e2ID := e2.GetID()
		contact.EntityID = &e1ID
		contact.SourceEntityID = &e2ID
		return []*collision.Contact{contact}
	}

	return nil
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

func isCapsuleCapsuleCollision(e1, e2 entities.Entity) bool {
	e1cc := e1.GetComponentContainer()
	e2cc := e2.GetComponentContainer()

	if e1cc.ColliderComponent.CapsuleCollider != nil {
		if e2cc.ColliderComponent.CapsuleCollider != nil {
			return true
		}
	}

	return false
}
