package networkdispatch

import (
	"encoding/json"
	"fmt"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/events"
	"github.com/kkevinchou/kito/kito/managers/player"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/lib/network"
)

type MessageHandler func(world World, message *network.Message)

func serverMessageHandler(world World, message *network.Message) {
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

func handlePlayerInput(player *player.Player, message *network.Message, world World) {
	singleton := world.GetSingleton()

	inputMessage := network.InputMessage{}
	err := json.Unmarshal(message.Body, &inputMessage)
	if err != nil {
		panic(err)
	}

	// fmt.Println("---------------")
	// fmt.Println("RECEIVED MOVE ON GCF", singleton.CommandFrame, "CF", message.CommandFrame)

	// TODO: This is to somewhat handle out of order messages coming to the server.
	// we take the latest command frame. however the current implementation risks
	// dropping inputs because we simply use only the latest
	if inputMessage.CommandFrame > player.LastInputCommandFrame {
		player.LastInputCommandFrame = inputMessage.CommandFrame
		player.LastInputGlobalCommandFrame = world.GetSingleton().CommandFrame
		singleton.PlayerInput[player.ID] = inputMessage.Input
	} else {
		fmt.Printf("received input out of order, last saw %d but got %d\n", player.LastInputCommandFrame, inputMessage.CommandFrame)
	}
}

// TODO: in the future this should be handled by some other system via an event
func handleCreatePlayer(player *player.Player, message *network.Message, world World) {
	playerID := message.SenderID

	bob := entities.NewBob(mgl64.Vec3{})
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

	event := &events.CreateEntityEvent{
		EntityType: types.EntityTypeBob,
		EntityID:   bob.GetID(),
	}
	world.GetEventBroker().Broadcast(event)
}