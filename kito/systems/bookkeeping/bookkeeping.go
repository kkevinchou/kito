package bookkeeping

import (
	"time"

	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/entities"
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
	for _, entity := range s.world.QueryEntity(components.ComponentFlagNotepad) {
		entity.GetComponentContainer().NotepadComponent.LastAction = components.ActionNone
	}
	for _, entity := range s.world.QueryEntity(components.ComponentFlagCollider) {
		if entity.Type() == types.EntityTypeProjectile {
			contacts := entity.GetComponentContainer().ColliderComponent.Contacts

			for e2ID, _ := range contacts {
				e2 := s.world.GetEntityByID(e2ID)
				if e2.GetComponentContainer().HealthComponent != nil {
					s.world.UnregisterEntity(e2)
				}
			}

			if len(contacts) > 0 {
				s.world.UnregisterEntity(entity)
			}
		}
	}
	for _, entity := range s.world.QueryEntity(components.ComponentFlagCollider) {
		entity.GetComponentContainer().ColliderComponent.Contacts = map[int]*collision.Contact{}
	}
}
