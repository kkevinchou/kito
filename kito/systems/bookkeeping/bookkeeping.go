package bookkeeping

import (
	"time"

	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/events"
	"github.com/kkevinchou/kito/kito/managers/eventbroker"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/lib/collision"
	"github.com/kkevinchou/kito/lib/input"
)

type World interface {
	CommandFrame() int
	GetSingleton() *singleton.Singleton
	QueryEntity(componentFlags int) []entities.Entity
	UnregisterEntity(entity entities.Entity)
	GetEntityByID(id int) entities.Entity
	GetEventBroker() eventbroker.EventBroker
}

type BookKeepingSystem struct {
	*base.BaseSystem

	world World
}

func NewBookKeepingSystem(world World) *BookKeepingSystem {
	return &BookKeepingSystem{
		world: world,
	}
}

func (s *BookKeepingSystem) Update(delta time.Duration) {
	if utils.IsServer() {
		singleton := s.world.GetSingleton()
		for i, _ := range singleton.PlayerInput {
			singleton.PlayerInput[i] = input.Input{}
		}
	}

	// handle fireball collisions
	for _, entity := range s.world.QueryEntity(components.ComponentFlagCollider) {
		if entity.Type() == types.EntityTypeProjectile {
			contacts := entity.GetComponentContainer().ColliderComponent.Contacts

			for e2ID, _ := range contacts {
				e2 := s.world.GetEntityByID(e2ID)
				health := e2.GetComponentContainer().HealthComponent
				if health != nil {
					health.Value -= 50
				}
			}

			if len(contacts) > 0 {
				s.world.UnregisterEntity(entity)
				event := &events.UnregisterEntityEvent{
					GlobalCommandFrame: s.world.CommandFrame(),
					EntityID:           entity.GetID(),
				}
				s.world.GetEventBroker().Broadcast(event)
			}
		}
	}

	// handle death events
	for _, entity := range s.world.QueryEntity(components.ComponentFlagHealth) {
		if entity.GetComponentContainer().HealthComponent.Value <= 0 {
			s.world.UnregisterEntity(entity)
			event := &events.UnregisterEntityEvent{
				GlobalCommandFrame: s.world.CommandFrame(),
				EntityID:           entity.GetID(),
			}
			s.world.GetEventBroker().Broadcast(event)
		}
	}

	// reset collision contacts
	for _, entity := range s.world.QueryEntity(components.ComponentFlagCollider) {
		cc := entity.GetComponentContainer()
		if len(cc.ColliderComponent.Contacts) == 0 {
			if cc.ThirdPersonControllerComponent != nil {
				cc.ThirdPersonControllerComponent.Grounded = false
			}
		}
		cc.ColliderComponent.Contacts = map[int]*collision.Contact{}
	}

	// reset notepad
	for _, entity := range s.world.QueryEntity(components.ComponentFlagNotepad) {
		entity.GetComponentContainer().NotepadComponent.LastAction = components.ActionNone
	}
}

func (s *BookKeepingSystem) Name() string {
	return "BookKeepingSystem"
}
