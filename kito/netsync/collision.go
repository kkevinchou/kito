package netsync

import (
	"sort"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/lib/collision"
)

type World interface {
	QueryEntity(componentFlags int) []entities.Entity
	GetPlayerEntity() entities.Entity
	GetEntityByID(id int) entities.Entity
}

func ResolveCollisionsForPlayer(player entities.Entity, world World) {
	entityPairs := [][]entities.Entity{}
	for _, e2 := range world.QueryEntity(components.ComponentFlagCollider | components.ComponentFlagTransform) {
		entityPairs = append(entityPairs, []entities.Entity{player, e2})
	}
	resolveCollisionsForEntityPairs(entityPairs, world)
}

func ResolveCollisions(world World) {
	entityPairs := [][]entities.Entity{}
	for _, e1 := range world.QueryEntity(components.ComponentFlagCollider | components.ComponentFlagTransform) {
		for _, e2 := range world.QueryEntity(components.ComponentFlagCollider | components.ComponentFlagTransform) {
			entityPairs = append(entityPairs, []entities.Entity{e1, e2})
		}
	}
	resolveCollisionsForEntityPairs(entityPairs, world)
}

func resolveCollisionsForEntityPairs(entityPairs [][]entities.Entity, world World) {
	// 1. collect pairs of entities that are colliding, sorted by separating vector
	// 2. perform collision resolution for any colliding entities
	// 3. this can cause more collisions, repeat until no more further detected collisions, or we hit the max

	positionalResolutionEntityPairs := [][]entities.Entity{}
	nonPositionalResolutionEntityPairs := [][]entities.Entity{}

	for _, pair := range entityPairs {
		cc1 := pair[0].GetComponentContainer()
		cc2 := pair[1].GetComponentContainer()

		if cc1.ColliderComponent.SkipMovementResolution || cc2.ColliderComponent.SkipMovementResolution {
			nonPositionalResolutionEntityPairs = append(nonPositionalResolutionEntityPairs, pair)
		} else {
			positionalResolutionEntityPairs = append(positionalResolutionEntityPairs, pair)
		}
	}

	resolveCount := map[int]int{}
	reachedMaxCount := false
	for !reachedMaxCount {
		collisionCandidates := collectSortedCollisionCandidates(positionalResolutionEntityPairs, world)
		if len(collisionCandidates) == 0 {
			break
		}

		resolvedEntities := resolveCollisions(collisionCandidates, world)
		for entityID, otherEntityID := range resolvedEntities {
			resolveCount[entityID] += 1
			if resolveCount[entityID] > resolveCountMax {
				reachedMaxCount = true
			}

			e1 := world.GetEntityByID(entityID)
			e2 := world.GetEntityByID(otherEntityID)

			colliderComponent1 := e1.GetComponentContainer().ColliderComponent
			colliderComponent2 := e2.GetComponentContainer().ColliderComponent

			if _, ok := colliderComponent1.Contacts[e2.GetID()]; !ok {
				colliderComponent1.Contacts[e2.GetID()] = nil
			}
			if _, ok := colliderComponent2.Contacts[e1.GetID()]; !ok {
				colliderComponent2.Contacts[e1.GetID()] = nil
			}
		}
	}

	collisionCandidates := collectSortedCollisionCandidates(nonPositionalResolutionEntityPairs, world)
	for _, candidate := range collisionCandidates {
		e1 := world.GetEntityByID(*candidate.EntityID)
		e2 := world.GetEntityByID(*candidate.SourceEntityID)

		colliderComponent1 := e1.GetComponentContainer().ColliderComponent
		colliderComponent2 := e2.GetComponentContainer().ColliderComponent

		if _, ok := colliderComponent1.Contacts[e2.GetID()]; !ok {
			colliderComponent1.Contacts[e2.GetID()] = nil
		}
		if _, ok := colliderComponent2.Contacts[e1.GetID()]; !ok {
			colliderComponent2.Contacts[e1.GetID()] = nil
		}
	}
}

// collectSortedCollisionCandidates collects all potential collisions that can occur in the frame.
// these are "candidates" in that if we resolve some of the collisions in the list, some will be
// invalidated
func collectSortedCollisionCandidates(entityPairs [][]entities.Entity, world World) []*collision.Contact {
	seen := map[int]map[int]bool{}

	// initialize collision state
	for _, e := range world.QueryEntity(components.ComponentFlagCollider | components.ComponentFlagTransform) {
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
		contacts := collide(e1, e2)
		allContacts = append(allContacts, contacts...)

		seen[e1.GetID()][e2.GetID()] = true
		seen[e2.GetID()][e1.GetID()] = true
	}
	sort.Sort(collision.ContactsBySeparatingDistance(allContacts))

	return allContacts
}

func collide(e1 entities.Entity, e2 entities.Entity) []*collision.Contact {
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

func resolveCollisions(contacts []*collision.Contact, world World) map[int]int {
	resolved := map[int]int{}
	for _, contact := range contacts {
		if _, ok := resolved[*contact.EntityID]; ok {
			continue
		}
		entity := world.GetEntityByID(*contact.EntityID)
		sourceEntity := world.GetEntityByID(*contact.SourceEntityID)
		resolveCollision(entity, sourceEntity, contact)

		resolved[*contact.EntityID] = *contact.SourceEntityID
		resolved[*contact.SourceEntityID] = *contact.EntityID
	}

	return resolved
}

func resolveCollision(entity entities.Entity, sourceEntity entities.Entity, contact *collision.Contact) {
	if contact.Type == collision.ContactTypeCapsuleTriMesh {
		cc := entity.GetComponentContainer()
		transformComponent := cc.TransformComponent
		tpcComponent := cc.ThirdPersonControllerComponent
		aiComponent := cc.AIComponent

		if tpcComponent != nil {
			separatingVector := contact.SeparatingVector
			transformComponent.Position = transformComponent.Position.Add(separatingVector)
			if separatingVector.Normalize().Dot(mgl64.Vec3{0, 1, 0}) >= groundedStrictness {
				tpcComponent.Velocity[1] = 0
				tpcComponent.BaseVelocity[1] = 0
				tpcComponent.ZipVelocity = mgl64.Vec3{}
				tpcComponent.Grounded = true
			}
		} else if aiComponent != nil {
			separatingVector := contact.SeparatingVector
			transformComponent.Position = transformComponent.Position.Add(separatingVector)
			aiComponent.Velocity[1] = 0
		}
	} else if contact.Type == collision.ContactTypeCapsuleCapsule {
		cc := entity.GetComponentContainer()
		transformComponent := cc.TransformComponent
		tpcComponent := cc.ThirdPersonControllerComponent

		separatingVector := contact.SeparatingVector.Mul(0.5)
		transformComponent.Position = transformComponent.Position.Add(separatingVector)

		if tpcComponent != nil {
			if separatingVector.Normalize().Dot(mgl64.Vec3{0, 1, 0}) >= groundedStrictness {
				tpcComponent.Grounded = true
				tpcComponent.Velocity[1] = 0
				tpcComponent.BaseVelocity[1] = 0
				tpcComponent.ZipVelocity = mgl64.Vec3{}
			}
		}

		cc2 := sourceEntity.GetComponentContainer()
		transformComponent2 := cc2.TransformComponent
		tpcComponent2 := cc2.ThirdPersonControllerComponent

		separatingVector2 := separatingVector.Mul(-1)
		transformComponent2.Position = transformComponent2.Position.Add(separatingVector2)

		if tpcComponent2 != nil {
			if separatingVector2.Normalize().Dot(mgl64.Vec3{0, 1, 0}) >= groundedStrictness {
				tpcComponent2.Grounded = true
				tpcComponent2.Velocity[1] = 0
				tpcComponent2.BaseVelocity[1] = 0
				tpcComponent2.ZipVelocity = mgl64.Vec3{}
			}
		}
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
