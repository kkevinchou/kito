package charactercontroller

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/lib/collision"

	"github.com/kkevinchou/kito/kito/entities"
)

const (
	// a value of 1 means the normal vector of what you're on must be exactly Vec3{0, 1, 0}
	groundedStrictness = 0.85
)

type CharacterControllerResolverSystem struct {
	*base.BaseSystem
	world    World
	entities []entities.Entity
}

func NewCharacterControllerResolverSystem(world World) *CharacterControllerResolverSystem {
	return &CharacterControllerResolverSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
		entities:   []entities.Entity{},
	}
}

func (s *CharacterControllerResolverSystem) RegisterEntity(entity entities.Entity) {
	componentContainer := entity.GetComponentContainer()

	if componentContainer.ThirdPersonControllerComponent != nil && componentContainer.TransformComponent != nil {
		s.entities = append(s.entities, entity)
	}

}

func (s *CharacterControllerResolverSystem) Update(delta time.Duration) {
	// collision resolution is synchronized from the server to the client
	if utils.IsClient() {
		player := s.world.GetPlayer()
		if player != nil {
			s.resolve(player)
		}
	} else {
		for _, entity := range s.entities {
			s.resolve(entity)
		}
	}
}

func (s *CharacterControllerResolverSystem) resolve(entity entities.Entity) {
	cc := entity.GetComponentContainer()
	colliderComponent := cc.ColliderComponent
	transformComponent := cc.TransformComponent
	contactManifolds := colliderComponent.ContactManifolds
	if contactManifolds != nil {
		separatingVector := combineSeparatingVectors(contactManifolds)
		transformComponent.Position = transformComponent.Position.Add(separatingVector)
	} else {
		// no collisions were detected (i.e. the ground)
		// physicsComponent.Grounded = false
	}
}

func combineSeparatingVectors(contactManifolds []*collision.ContactManifold) mgl64.Vec3 {
	// only add separating vectors which we haven't seen before. ideally
	// this should handle cases where separating vectors are a basis of another
	// and avoid "overcounting" separating vectors
	seenSeparatingVectors := []mgl64.Vec3{}
	var separatingVector mgl64.Vec3
	for _, contactManifold := range contactManifolds {
		for _, contact := range contactManifold.Contacts {
			seen := false
			for _, v := range seenSeparatingVectors {
				if contact.SeparatingVector.ApproxEqual(v) {
					seen = true
					break
				}
			}
			if !seen {
				separatingVector = separatingVector.Add(contact.SeparatingVector)
				seenSeparatingVectors = append(seenSeparatingVectors, contact.SeparatingVector)
			}
		}
	}
	return separatingVector
}
