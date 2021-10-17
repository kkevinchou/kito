package playerinput

import (
	"fmt"
	"time"

	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/managers/player"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/lib/input"
)

type World interface {
	CommandFrame() int
	GetSingleton() *singleton.Singleton
	GetEntityByID(int) (entities.Entity, error)
}

type PlayerInputSystem struct {
	*base.BaseSystem
	world World
}

func NewPlayerInputSystem(world World) *PlayerInputSystem {
	return &PlayerInputSystem{
		world: world,
	}
}

func (s *PlayerInputSystem) RegisterEntity(entity entities.Entity) {
}

func (s *PlayerInputSystem) Update(delta time.Duration) {
	singleton := s.world.GetSingleton()
	playerManager := directory.GetDirectory().PlayerManager()
	players := playerManager.GetPlayers()

	for _, player := range players {
		bufferedInput := singleton.InputBuffer.PullInput(singleton.CommandFrame, player.ID)
		if bufferedInput != nil {
			// fmt.Println("TIME IN INPUT BUFFER", time.Since(bufferedInput.ReceivedTimestamp))
			handlePlayerInput(player, bufferedInput.PlayerCommandFrame, bufferedInput.Input, s.world)
		}
	}
}

func handlePlayerInput(player *player.Player, commandFrame int, input input.Input, world World) {
	singleton := world.GetSingleton()
	// This is to somewhat handle out of order messages coming to the server.
	// we take the latest command frame. However the current implementation risks
	// dropping inputs because we simply use only the latest
	if commandFrame > player.LastInputCommandFrame {
		player.LastInputCommandFrame = commandFrame
		player.LastInputGlobalCommandFrame = world.CommandFrame()
		singleton.PlayerInput[player.ID] = input
	} else {
		fmt.Printf("received input out of order, last saw %d but got %d\n", player.LastInputCommandFrame, commandFrame)
	}
}
