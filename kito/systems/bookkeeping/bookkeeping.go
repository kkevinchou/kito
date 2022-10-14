package bookkeeping

import (
	"time"

	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/events"
	"github.com/kkevinchou/kito/kito/managers/eventbroker"
	"github.com/kkevinchou/kito/kito/netsync"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/lib/input"
)

type World interface {
	GetSingleton() *singleton.Singleton
	QueryEntity(componentFlags int) []entities.Entity
	GetEventBroker() eventbroker.EventBroker
	UnregisterEntityByID(id int)
}

type BookKeepingSystem struct {
	*base.BaseSystem
	events []events.Event

	world World
}

func NewBookKeepingSystem(world World) *BookKeepingSystem {
	s := &BookKeepingSystem{
		world: world,
	}

	eventBroker := world.GetEventBroker()
	eventBroker.AddObserver(s, []events.EventType{
		events.EventTypeUnregisterEntity,
	})

	return s
}
func (s *BookKeepingSystem) clearEvents() {
	s.events = nil
}

func (s *BookKeepingSystem) Observe(event events.Event) {
	if event.Type() == events.EventTypeUnregisterEntity {
		s.events = append(s.events, event)
	}
}

func (s *BookKeepingSystem) Update(delta time.Duration) {
	defer s.clearEvents()

	if utils.IsServer() {
		singleton := s.world.GetSingleton()
		for i, _ := range singleton.PlayerInput {
			singleton.PlayerInput[i] = input.Input{}
		}

		for _, event := range s.events {
			if e, ok := event.(*events.UnregisterEntityEvent); ok {
				s.world.UnregisterEntityByID(e.EntityID)
			}
		}
	}

	// reset collision contacts
	for _, entity := range s.world.QueryEntity(components.ComponentFlagCollider) {
		netsync.CollisionBookKeeping(entity)
	}

	// reset notepad
	for _, entity := range s.world.QueryEntity(components.ComponentFlagNotepad) {
		entity.GetComponentContainer().NotepadComponent.LastAction = components.ActionNone
	}
}

func (s *BookKeepingSystem) Name() string {
	return "BookKeepingSystem"
}
