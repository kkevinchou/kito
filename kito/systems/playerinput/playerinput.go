package playerinput

import (
	"fmt"
	"time"

	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/inputbuffer"
	"github.com/kkevinchou/kito/kito/managers/player"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
)

type World interface {
	CommandFrame() int
	GetSingleton() *singleton.Singleton
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

func (s *PlayerInputSystem) Update(delta time.Duration) {
	singleton := s.world.GetSingleton()
	playerManager := directory.GetDirectory().PlayerManager()
	players := playerManager.GetPlayers()

	for _, player := range players {
		bufferedInput := singleton.InputBuffer.PullInput(singleton.CommandFrame, player.ID)
		if bufferedInput != nil {
			handlePlayerInput(player, bufferedInput.LocalCommandFrame, bufferedInput, s.world)
		}
	}
}

func handlePlayerInput(player *player.Player, commandFrame int, bufferedInput *inputbuffer.BufferedInput, world World) {
	// This is to somewhat handle out of order messages coming to the server.
	// we take the latest command frame. However the current implementation risks
	// dropping inputs because we simply use only the latest
	if commandFrame > player.LastInputLocalCommandFrame {
		player.LastInputLocalCommandFrame = commandFrame
		player.LastInputGlobalCommandFrame = world.CommandFrame()

		singleton := world.GetSingleton()
		singleton.PlayerInput[player.ID] = bufferedInput.Input
		singleton.PlayerCommands[player.ID] = bufferedInput.PlayerCommands
	} else {
		fmt.Printf("received input out of order, last saw %d but got %d\n", player.LastInputLocalCommandFrame, commandFrame)
	}
}

func (s *PlayerInputSystem) Name() string {
	return "PlayerInputSystem"
}
