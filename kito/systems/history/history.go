package history

import (
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
	GetPlayer() entities.Entity
}

type HistorySystem struct {
	*base.BaseSystem
	world World
}

func NewHistorySystem(world World) *HistorySystem {
	return &HistorySystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}
}

func (s *HistorySystem) RegisterEntity(entity entities.Entity) {
}

func (s *HistorySystem) Update(delta time.Duration) {
	singleton := s.world.GetSingleton()
	playerInput := singleton.PlayerInput[singleton.PlayerID]

	cfHistory := s.world.GetCommandFrameHistory()
	player := s.world.GetPlayer()
	cfHistory.AddCommandFrame(singleton.CommandFrame, playerInput, player)
}
