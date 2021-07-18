package networkdispatch

import (
	"encoding/json"
	"fmt"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/lib/network"
	"github.com/kkevinchou/kito/managers/player"
)

type MessageHandler func(world World, message *network.Message)

func defaultMessageHandler(world World, message *network.Message) {
}

func ServerMessageHandler(world World, message *network.Message) {
	playerManager := directory.GetDirectory().PlayerManager()
	player := playerManager.GetPlayer(message.SenderID)
	if player == nil {
		fmt.Println(fmt.Errorf("failed to find player with id %d", message.SenderID))
		return
	}

	if message.MessageType == network.MessageTypeCreatePlayer {
		handleCreatePlayer(player, message, world)
	} else if message.MessageType == network.MessageTypeInput {
		handlePlayerInput(player, message, world)
	} else {
		fmt.Println("unknown message type:", message.MessageType, string(message.Body))
	}
}

func ClientMessageHandler(world World, message *network.Message) {
	if message.MessageType == network.MessageTypeGameStateSnapshot {
		handleEntitySnapshot(message, world)
	} else if message.MessageType == network.MessageTypeAckCreatePlayer {
		// This is handled as a part of client initialization and doesn't
		// need to be handled here
	} else {
		fmt.Println("unknown message type:", message.MessageType, string(message.Body))
	}
}

func handleEntitySnapshot(message *network.Message, world World) {
	var gameStateSnapshot network.GameStateSnapshotMessage
	err := json.Unmarshal(message.Body, &gameStateSnapshot)
	if err != nil {
		panic(err)
	}

	for _, entitySnapshot := range gameStateSnapshot.Entities {
		entity, err := world.GetEntityByID(entitySnapshot.ID)
		if err != nil {
			continue
		}

		cc := entity.GetComponentContainer()
		cc.TransformComponent.Position = entitySnapshot.Position
		cc.TransformComponent.Orientation = entitySnapshot.Orientation
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

	world.RegisterEntities([]entities.Entity{bob, camera})
	fmt.Println("Created and registered a new bob with id", bob.ID)

	ack := &network.AckCreatePlayerMessage{
		ID:          playerID,
		CameraID:    camera.ID,
		Position:    cc.TransformComponent.Position,
		Orientation: cc.TransformComponent.Orientation,
	}

	player.Client.SendMessage(network.MessageTypeAckCreatePlayer, ack)
	fmt.Println("Sent entity ack creation message")
}
