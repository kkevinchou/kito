package ping

import (
	"fmt"
	"time"

	"github.com/kkevinchou/kito/kito/commandframe"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/knetwork"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
)

type World interface {
	GetSingleton() *singleton.Singleton
	GetCommandFrameHistory() *commandframe.CommandFrameHistory
	GetEntityByID(id int) (entities.Entity, error)
	GetPlayer() entities.Entity
}

type PingSystem struct {
	*base.BaseSystem
	world World
}

func NewPingSystem(world World) *PingSystem {
	return &PingSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}
}

func (s *PingSystem) RegisterEntity(entity entities.Entity) {
}

func (s *PingSystem) Update(delta time.Duration) {
	singleton := s.world.GetSingleton()
	playerManager := directory.GetDirectory().PlayerManager()

	player := playerManager.GetPlayer(singleton.PlayerID)

	pingMessage := &knetwork.PingMessage{
		SendTime: time.Now(),
	}

	err := player.Client.SendMessage(knetwork.MessageTypePing, pingMessage)
	if err != nil {
		fmt.Printf("error sending ping message: %s\n", err)
	}
}
