package combat

import (
	"time"

	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/events"
	"github.com/kkevinchou/kito/kito/managers/eventbroker"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/kito/utils"
)

type World interface {
	QueryEntity(componentFlags int) []entities.Entity
	GetEntityByID(id int) entities.Entity
	CommandFrame() int
	UnregisterEntity(entity entities.Entity)
	GetEventBroker() eventbroker.EventBroker
}

type CombatSystem struct {
	*base.BaseSystem

	world World
}

func NewCombatSystem(world World) *CombatSystem {
	return &CombatSystem{
		world: world,
	}
}

func (s *CombatSystem) Update(delta time.Duration) {
	if utils.IsClient() {
		return
	}

	// handle fireball collisions
	for _, entity := range s.world.QueryEntity(components.ComponentFlagCollider) {
		if entity.Type() == types.EntityTypeProjectile {
			contacts := entity.GetComponentContainer().ColliderComponent.Contacts
			if len(contacts) == 0 {
				continue
			}

			for e2ID, _ := range contacts {
				e2 := s.world.GetEntityByID(e2ID)
				health := e2.GetComponentContainer().HealthComponent
				if health != nil {
					health.Data.Value -= 50
				}
			}

			event := &events.UnregisterEntityEvent{
				GlobalCommandFrame: s.world.CommandFrame(),
				EntityID:           entity.GetID(),
			}
			s.world.GetEventBroker().Broadcast(event)
		}
	}

	// handle death events
	for _, entity := range s.world.QueryEntity(components.ComponentFlagHealth) {
		if entity.GetComponentContainer().HealthComponent.Data.Value <= 0 {
			event := &events.UnregisterEntityEvent{
				GlobalCommandFrame: s.world.CommandFrame(),
				EntityID:           entity.GetID(),
			}
			s.world.GetEventBroker().Broadcast(event)
		}
	}
}

func (s *CombatSystem) Name() string {
	return "CombatSystem"
}
