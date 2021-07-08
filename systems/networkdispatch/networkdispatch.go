package networkdispatch

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/entities/singleton"
	"github.com/kkevinchou/kito/lib/network"
	"github.com/kkevinchou/kito/managers/player"
	"github.com/kkevinchou/kito/systems/base"
)

type World interface {
	RegisterEntities([]entities.Entity)
	GetEntityByID(id int) (entities.Entity, error)
	GetSingleton() *singleton.Singleton
	SetCamera(camera entities.Entity)
}

type NetworkDispatchSystem struct {
	*base.BaseSystem
	world World
}

func NewNetworkDispatchSystem(world World) *NetworkDispatchSystem {
	return &NetworkDispatchSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}
}

func (s *NetworkDispatchSystem) RegisterEntity(entity entities.Entity) {
}

func (s *NetworkDispatchSystem) Update(delta time.Duration) {
	d := directory.GetDirectory()
	playerManager := d.PlayerManager()

	for _, player := range playerManager.GetPlayers() {
		messages := player.Client.PullIncomingMessages()
		for _, message := range messages {
			if message.MessageType == network.MessageTypeCreatePlayer {
				handleCreatePlayer(player, message, s.world)
			} else if message.MessageType == network.MessageTypeInput {
				handlePlayerInput(player, message, s.world)
			} else {
				fmt.Println("unknown message type:", message.MessageType, string(message.Body))
			}
		}
	}
}

func handlePlayerInput(player *player.Player, message *network.Message, world World) {
	singleton := world.GetSingleton()

	inputMessage := network.InputMessage{}
	err := json.Unmarshal(message.Body, &inputMessage)
	if err != nil {
		panic(err)
	}

	singleton.PlayerInput[player.ID] = inputMessage.Input
}

// todo: in the future this should be handled by some other system via an event
func handleCreatePlayer(player *player.Player, message *network.Message, world World) {
	playerID := message.SenderID

	bob := entities.NewServerBob(mgl64.Vec3{})
	bob.ID = playerID

	cc := bob.ComponentContainer

	camera := entities.NewThirdPersonCamera(mgl64.Vec3{}, mgl64.Vec2{0, 0}, bob.GetID())
	cameraComponentContainer := camera.GetComponentContainer()
	fmt.Println("Server camera initialized at position", cameraComponentContainer.TransformComponent.Position)

	cc.ThirdPersonControllerComponent.CameraID = camera.GetID()

	world.SetCamera(camera)

	world.RegisterEntities([]entities.Entity{bob, camera})
	fmt.Println("Created and registered a new bob with id", bob.ID)

	ack := &network.AckCreatePlayerMessage{
		ID:          playerID,
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
