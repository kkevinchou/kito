package networkdispatch

import (
	"encoding/json"
	"fmt"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/lib/network"
)

func clientMessageHandler(world World, message *network.Message) {
	if message.MessageType == network.MessageTypeGameStateUpdate {
		singleton := world.GetSingleton()
		var gameStateUpdate network.GameStateUpdateMessage
		err := json.Unmarshal(message.Body, &gameStateUpdate)
		if err != nil {
			panic(err)
		}

		handleGameStateUpdate(&gameStateUpdate, world)
		singleton.StateBuffer.PushEntityUpdate(world.CommandFrame(), message, &gameStateUpdate)
	} else if message.MessageType == network.MessageTypeAckCreatePlayer {
		handleAckCreatePlayer(message, world)
	} else {
		fmt.Println("unknown message type:", message.MessageType, string(message.Body))
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

func handleGameStateUpdate(gameStateUpdate *network.GameStateUpdateMessage, world World) {
	// we use a gcf adjusted command frame lookup because even though an input may happen on only one command
	// frame, the entity can continue to be updated due to that command frame. Therefore we need to make sure
	// we advance the command frame as much as the server has to see if we've mispredicted
	deltaGCF := gameStateUpdate.CurrentGlobalCommandFrame - gameStateUpdate.LastInputGlobalCommandFrame
	lookupCommandFrame := gameStateUpdate.LastInputCommandFrame + deltaGCF

	// TODO: we should use the latest cfHistory if we're not able to find an exact command frame history
	// with the lookup. Standing "still" is still a prediction, and if some outside factor affects the
	// player, we should detect that as a misprediction and move our character accordingly
	cfHistory := world.GetCommandFrameHistory()
	cf := cfHistory.GetCommandFrame(lookupCommandFrame)

	entitySnapshot := gameStateUpdate.Entities[world.GetSingleton().PlayerID]
	foundEntity, err := world.GetEntityByID(entitySnapshot.ID)
	if err != nil {
		// this sometimes fails on startup - probably some networking timing
		panic(err)
	}

	if cf != nil {
		historyEntity := cf.PostCFState
		// fmt.Println("HIT", gameStateupdate.LastInputCommandFrame, historyEntity.Position)
		if historyEntity.Position != entitySnapshot.Position || historyEntity.Orientation != entitySnapshot.Orientation {
			fmt.Printf(
				"CLIENT-SIDE PREDICTION MISS %d %d %d\n%v\n%v\n",
				gameStateUpdate.LastInputCommandFrame,
				gameStateUpdate.LastInputGlobalCommandFrame,
				gameStateUpdate.CurrentGlobalCommandFrame,
				historyEntity.Position,
				entitySnapshot.Position,
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
			cfHistory.ClearUntilFrameNumber(lookupCommandFrame)
		}
	} else {
		// fmt.Println("empty cfHistory", gameStateupdate.LastInputCommandFrame, deltaGCF, len(world.GetCommandFrameHistory().CommandFrames))
	}
}
