package networkdispatch

import (
	"encoding/json"
	"fmt"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/events"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/lib/network"
)

func clientMessageHandler(world World, message *network.Message) {
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

	cfHistory := world.GetCommandFrameHistory().CommandFrames[gameStateupdate.LatestReceivedCommandFrame]

	//TODO: terrible n^2
	if cfHistory != nil {
		for _, historyEntity := range cfHistory.PostCFState {
			for _, entitySnapshot := range gameStateupdate.Entities {
				if historyEntity.ID == entitySnapshot.ID {
					if historyEntity.Position == entitySnapshot.Position {
						fmt.Println("HIT", gameStateupdate.LatestReceivedCommandFrame)
					} else {
						fmt.Println("MiSS", gameStateupdate.LatestReceivedCommandFrame, historyEntity.Position, entitySnapshot.Position)
					}
				}
			}
		}
	}

	var newEntities []entities.Entity
	for _, entitySnapshot := range gameStateupdate.Entities {
		foundEntity, err := world.GetEntityByID(entitySnapshot.ID)

		if err != nil {
			if types.EntityType(entitySnapshot.Type) == types.EntityTypeBob {
				newEntity := entities.NewBob(mgl64.Vec3{})
				newEntity.ID = entitySnapshot.ID

				cc := newEntity.GetComponentContainer()
				cc.TransformComponent.Position = entitySnapshot.Position
				cc.TransformComponent.Orientation = entitySnapshot.Orientation

				newEntities = append(newEntities, newEntity)
			} else if types.EntityType(entitySnapshot.Type) == types.EntityTypeRigidBody {
				newEntity := entities.NewRigidBody(entitySnapshot.Position)
				newEntity.ID = entitySnapshot.ID

				cc := newEntity.GetComponentContainer()
				cc.TransformComponent.Position = entitySnapshot.Position
				cc.TransformComponent.Orientation = entitySnapshot.Orientation

				newEntities = append(newEntities, newEntity)
			} else {
				continue
			}
		} else {
			// TODO: the rewinding here is forcibly making us MISS when comparing past command frames
			// once we more intelligently apply these updates this should cause us to miss less frequently.
			// currently using this hack
			if entitySnapshot.ID == singleton.PlayerID {
				continue
			}
			cc := foundEntity.GetComponentContainer()
			cc.TransformComponent.Position = entitySnapshot.Position
			cc.TransformComponent.Orientation = entitySnapshot.Orientation
		}
	}
	world.RegisterEntities(newEntities)

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
