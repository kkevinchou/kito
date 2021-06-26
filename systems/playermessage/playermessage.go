package playermessage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/lib/network"
	"github.com/kkevinchou/kito/managers/player"
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
				handleCreatePlayer(player, message, s.world)
			}
		}
	}
}

func handleCreatePlayer(player *player.Player, message *network.Message, world World) {
	fmt.Println("start new bob creation for id", message.SenderID)
	playerId := message.SenderID

	bob := entities.NewServerBob(mgl64.Vec3{})
	bob.ID = playerId

	world.RegisterEntities([]entities.Entity{bob})
	fmt.Println("Created and registered a new bob with id", bob.ID)

	cc := bob.ComponentContainer

	ack := &network.AckCreatePlayerMessage{
		Position:    cc.TransformComponent.Position,
		Orientation: cc.TransformComponent.Orientation,
	}
	ackBytes, err := json.Marshal(ack)
	if err != nil {
		panic(err)
	}

	response := &network.Message{
		MessageType: network.MessageTypeAckCreatePlayer,
		Body:        ackBytes,
	}
	player.Client.SendMessage(response)
	fmt.Println("Sent entity ack creation message")
}
