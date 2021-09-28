package history

import (
	"fmt"
	"time"

	"github.com/kkevinchou/kito/kito/commandframe"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
)

type World interface {
	GetSingleton() *singleton.Singleton
	GetCommandFrameHistory() *commandframe.CommandFrameHistory
	GetEntityByID(id int) (entities.Entity, error)
}

type HistorySystem struct {
	*base.BaseSystem
	world    World
	entities []entities.Entity
}

func NewHistorySystem(world World) *HistorySystem {
	return &HistorySystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}
}

func (s *HistorySystem) RegisterEntity(entity entities.Entity) {
	componentContainer := entity.GetComponentContainer()

	// TODO: emulating the filtering that physics systems do... pretty brittle
	if componentContainer.PhysicsComponent != nil && componentContainer.TransformComponent != nil {
		s.entities = append(s.entities, entity)
	}
}

func (s *HistorySystem) Update(delta time.Duration) {
	singleton := s.world.GetSingleton()
	playerInput := singleton.PlayerInput[singleton.PlayerID]

	cfHistory := s.world.GetCommandFrameHistory()
	player, err := s.world.GetEntityByID(singleton.PlayerID)
	if err != nil {
		fmt.Println("history update failed to find player", err)
		return
	}
	cfHistory.AddCommandFrame(singleton.CommandFrame, playerInput, player)
}
