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

	// we use a gcf adjusted command frame lookup because even though an input may happen on only one command
	// frame, the entity can continue to be updated due to that command frame. Therefore we need to make sure
	// we advance the command frame as much as the server has to see if we've mispredicted
	deltaGCF := gameStateupdate.CurrentGlobalCommandFrame - gameStateupdate.LastInputGlobalCommandFrame
	lookupCommandFrame := gameStateupdate.LastInputCommandFrame + deltaGCF

	// TODO: we should use the latest cfHistory if we're not able to find an exact command frame history
	// with the lookup. Standing "still" is still a prediction, and if some outside factor affects the
	// player, we should detect that as a misprediction and move our character accordingly
	cfHistory := world.GetCommandFrameHistory().CommandFrames[lookupCommandFrame]

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
			if entitySnapshot.ID == singleton.CameraID {
				continue
			}

			if entitySnapshot.ID == singleton.PlayerID {
				if cfHistory != nil {
					historyEntity := cfHistory.PostCFState
					// fmt.Println("HIT", gameStateupdate.LastInputCommandFrame, historyEntity.Position)
					if historyEntity.Position != entitySnapshot.Position || historyEntity.Orientation != entitySnapshot.Orientation {
						fmt.Println(
							"CLIENT-SIDE PREDICTION MISS",
							gameStateupdate.LastInputCommandFrame,
							gameStateupdate.LastInputGlobalCommandFrame,
							gameStateupdate.CurrentGlobalCommandFrame,
							historyEntity,
							entitySnapshot,
						)

						// if I was a god programmer I would re-apply historical user inputs to catch
						// up from the latest command frame the server saw. for now I'm just going to
						// snap the user to the server view which will show a bit of a hitch to the user
						fmt.Println("SNAPPING", entitySnapshot.ID)
						cc := foundEntity.GetComponentContainer()
						cc.TransformComponent.Position = entitySnapshot.Position
						cc.TransformComponent.Orientation = entitySnapshot.Orientation
					} else {
						// fmt.Println(
						// 	"CLIENT-SIDE PREDICTION HIT",
						// 	gameStateupdate.LastInputCommandFrame,
						// 	gameStateupdate.LastInputGlobalCommandFrame,
						// 	gameStateupdate.CurrentGlobalCommandFrame,
						// )
					}
				} else {
					// fmt.Println("empty cfHistory", gameStateupdate.LastInputCommandFrame, deltaGCF, len(world.GetCommandFrameHistory().CommandFrames))
				}

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
