package collisionresolver

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/types"
)

type World interface {
	GetSingleton() *singleton.Singleton
	GetEntityByID(id int) (entities.Entity, error)
	GetPlayer() entities.Entity
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

	if componentContainer.ColliderComponent != nil && componentContainer.TransformComponent != nil {
		s.entities = append(s.entities, entity)
	}
}

func (s *CollisionResolverSystem) Update(delta time.Duration) {
	for _, e := range s.entities {
		cc := e.GetComponentContainer()
		colliderComponent := cc.ColliderComponent
		transformComponent := cc.TransformComponent
		physicsComponent := cc.PhysicsComponent
		contactManifold := colliderComponent.ContactManifold
		if contactManifold != nil {
			transformComponent.Position = transformComponent.Position.Add(contactManifold.Contacts[0].SeparatingVector)
			physicsComponent.Impulses = map[string]types.Impulse{}
			physicsComponent.Velocity = mgl64.Vec3{}
		}
	}
}
