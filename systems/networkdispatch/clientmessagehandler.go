package networkdispatch

import (
	"encoding/json"
	"fmt"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/events"
	"github.com/kkevinchou/kito/lib/network"
	"github.com/kkevinchou/kito/settings"
	"github.com/kkevinchou/kito/types"
)

func ClientMessageHandler(world World, message *network.Message) {
	if message.MessageType == network.MessageTypeGameStateUpdate {
		handleGameStateUpdate(message, world)
	} else if message.MessageType == network.MessageTypeAckCreatePlayer {
		handleAckCreatePlayer(message, world)
	} else {
		fmt.Println("unknown message type:", message.MessageType, string(message.Body))
	}
}

func handleGameStateUpdate(message *network.Message, world World) {
	singleton := world.GetSingleton()
	var gameStateupdate network.GameStateUpdateMessage
	err := json.Unmarshal(message.Body, &gameStateupdate)
	if err != nil {
		panic(err)
	}

	for _, entitySnapshot := range gameStateupdate.Entities {
		var entity entities.Entity
		var err error
		entity, err = world.GetEntityByID(entitySnapshot.ID)
		if err != nil {
			if types.EntityType(entitySnapshot.Type) == types.EntityTypeBob {
				bob := entities.NewBob(mgl64.Vec3{})
				bob.ID = entitySnapshot.ID
				world.RegisterEntities([]entities.Entity{bob})

				entity = bob
			} else {
				return
			}
		}

		cc := entity.GetComponentContainer()
		cc.TransformComponent.Position = entitySnapshot.Position
		cc.TransformComponent.Orientation = entitySnapshot.Orientation
	}

	for _, event := range gameStateupdate.Events {
		if events.EventType(event.Type) == events.EventTypeCreateEntity {
			var realEvent events.CreateEntityEvent
			json.Unmarshal(event.Bytes, &realEvent)

			if realEvent.EntityType == types.EntityTypeBob {
				if realEvent.EntityID == singleton.PlayerID {
					// TODO: skip this is event as its the creation event of the original player.
					// In the future, the server should simply not send this message to the player
					continue
				}

				//TODO: undecided what the best way to handle creation of entities. currently using the entity snapshot
				// to determine the need to create an entity
				// bob := entities.NewBob(mgl64.Vec3{})
				// bob.ID = realEvent.EntityID
				// world.RegisterEntities([]entities.Entity{bob})
			}
		}
	}
}

func handleAckCreatePlayer(message *network.Message, world World) {
	subMessage := &network.AckCreatePlayerMessage{}
	err := json.Unmarshal(message.Body, subMessage)
	if err != nil {
		fmt.Println(err)
		return
	}

	singleton := world.GetSingleton()
	singleton.PlayerID = subMessage.ID
	singleton.CameraID = subMessage.CameraID

	bob := entities.NewBob(mgl64.Vec3{})
	bob.ID = subMessage.ID

	camera := entities.NewThirdPersonCamera(settings.CameraStartPosition, settings.CameraStartView, bob.GetID())
	camera.ID = subMessage.CameraID

	bob.GetComponentContainer().ThirdPersonControllerComponent.CameraID = camera.GetID()

	world.RegisterEntities([]entities.Entity{
		bob,
		camera,
	})
}
