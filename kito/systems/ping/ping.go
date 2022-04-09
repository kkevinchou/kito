package ping

import (
	"fmt"
	"time"

	"github.com/kkevinchou/kito/kito/knetwork"
	"github.com/kkevinchou/kito/kito/managers/player"
	"github.com/kkevinchou/kito/kito/systems/base"
)

type World interface {
	GetPlayer() *player.Player
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

func (s *PingSystem) Update(delta time.Duration) {
	player := s.world.GetPlayer()

	pingMessage := &knetwork.PingMessage{
		SendTime: time.Now(),
	}

	err := player.Client.SendMessage(knetwork.MessageTypePing, pingMessage)
	if err != nil {
		fmt.Printf("error sending ping message: %s\n", err)
	}
}
