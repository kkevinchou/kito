package playermessage

import (
	"fmt"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/lib/network"
	"github.com/kkevinchou/kito/systems/base"
)

type World interface {
	RegisterEntities([]entities.Entity)
}

type PlayerMessageSystem struct {
	*base.BaseSystem
	world World
}

func NewPlayerMessageSystem(world World) *PlayerMessageSystem {
	return &PlayerMessageSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}
}

func (s *PlayerMessageSystem) RegisterEntity(entity entities.Entity) {
}

func (s *PlayerMessageSystem) Update(delta time.Duration) {
	d := directory.GetDirectory()
	playerManager := d.PlayerManager()

	for _, player := range playerManager.GetPlayers() {
		messages := player.Client.PullIncomingMessages()
		for _, message := range messages {
			if message.MessageType == network.MessageTypeCreatePlayer {
				handleCreatePlayer(message, s.world)
			}
		}
	}
}

func handleCreatePlayer(message *network.Message, world World) {
	fmt.Println("start new bob creation for id", message.SenderID)
	playerId := message.SenderID

	bob := entities.NewServerBob(mgl64.Vec3{})
	bob.ID = playerId

	world.RegisterEntities([]entities.Entity{bob})
	fmt.Println("Created and registered a new bob with id", bob.ID)
}
