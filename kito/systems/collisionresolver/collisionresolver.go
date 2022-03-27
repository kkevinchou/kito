package collisionresolver

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/kito/utils"
)

const (
	// a value of 1 means the normal vector of what you're on must be exactly Vec3{0, 1, 0}
	groundedStrictness = 0.85
)

type World interface {
	GetSingleton() *singleton.Singleton
	GetEntityByID(id int) (entities.Entity, error)
	GetPlayerEntity() entities.Entity
}

type CollisionResolverSystem struct {
	*base.BaseSystem
	world    World
	entities []entities.Entity
}

func NewCollisionResolverSystem(world World) *CollisionResolverSystem {
	return &CollisionResolverSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}
}

func (s *CollisionResolverSystem) RegisterEntity(entity entities.Entity) {
	componentContainer := entity.GetComponentContainer()

	if componentContainer.ColliderComponent != nil && componentContainer.TransformComponent != nil && componentContainer.PhysicsComponent != nil {
		s.entities = append(s.entities, entity)
	}
}

func (s *CollisionResolverSystem) Update(delta time.Duration) {
	// collision resolution is synchronized from the server to the client
	if utils.IsClient() {
		return
	}

	for _, entity := range s.entities {
		s.resolve(entity)
	}
}

func (s *CollisionResolverSystem) resolve(entity entities.Entity) {
	cc := entity.GetComponentContainer()
	colliderComponent := cc.ColliderComponent
	transformComponent := cc.TransformComponent
	physicsComponent := cc.PhysicsComponent

	if colliderComponent.CollisionInstances != nil {
		contactManifolds := colliderComponent.CollisionInstances[0].ContactManifolds

		// naively add all separating vectors together
		var separatingVector mgl64.Vec3
		for _, contactManifold := range contactManifolds {
			separatingVector = separatingVector.Add(contactManifold.Contacts[0].SeparatingVector)
		}

		if separatingVector.Normalize().Dot(mgl64.Vec3{0, 1, 0}) >= groundedStrictness {
			physicsComponent.Grounded = true
			physicsComponent.Velocity = mgl64.Vec3{}
		} else {
			physicsComponent.Grounded = false
		}

		transformComponent.Position = transformComponent.Position.Add(separatingVector)
		physicsComponent.Impulses = map[string]types.Impulse{}
	} else {
		// no collisions were detected (i.e. the ground)
		physicsComponent.Grounded = false
	}
}
