package history

import (
	"time"

	"github.com/kkevinchou/kito/kito/commandframe"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/managers/player"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
)

type World interface {
	GetSingleton() *singleton.Singleton
	GetCommandFrameHistory() *commandframe.CommandFrameHistory
	GetPlayerEntity() entities.Entity
	GetPlayer() *player.Player
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

func (s *HistorySystem) Update(delta time.Duration) {
	singleton := s.world.GetSingleton()
	player := s.world.GetPlayer()
	playerEntity := s.world.GetPlayerEntity()

	playerInput := singleton.PlayerInput[player.ID]
	cfHistory := s.world.GetCommandFrameHistory()
	cfHistory.AddCommandFrame(singleton.CommandFrame, playerInput, playerEntity)
}
