package networkdispatch

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/commandframe"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/kito/utils/controllerutils"
	"github.com/kkevinchou/kito/kito/utils/physutils"
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
		singleton.StateBuffer.PushEntityUpdate(world.CommandFrame(), &gameStateUpdate)
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
	player, err := world.GetEntityByID(entitySnapshot.ID)
	if err != nil {
		// this sometimes fails on startup - probably some networking timing
		panic(err)
	}

	if cf != nil {
		historyEntity := cf.PostCFState
		if historyEntity.Position != entitySnapshot.Position || historyEntity.Orientation != entitySnapshot.Orientation {
			fmt.Printf(
				"CLIENT-SIDE PREDICTION MISS %d %d %d\n%v\n%v\n",
				gameStateUpdate.LastInputCommandFrame,
				gameStateUpdate.LastInputGlobalCommandFrame,
				gameStateUpdate.CurrentGlobalCommandFrame,
				historyEntity.Position,
				entitySnapshot.Position,
			)

			// When we miss the client-side prediction, we set the player's state to the snapshot state.
			// what's important to note is that the snapshot state is in the past! At the very least,
			// a whole RTT in the past. So, we want to replay our historical inputs to catch up to the player's
			// present.

			cc := player.GetComponentContainer()
			cc.TransformComponent.Position = entitySnapshot.Position
			cc.TransformComponent.Orientation = entitySnapshot.Orientation
			replayInputs(player, world.GetCamera(), lookupCommandFrame, cfHistory)
		} else {
			// fmt.Println(
			// 	"CLIENT-SIDE PREDICTION HIT",
			// 	gameStateUpdate.LastInputCommandFrame,
			// 	gameStateUpdate.LastInputGlobalCommandFrame,
			// 	gameStateUpdate.CurrentGlobalCommandFrame,
			// )
			cfHistory.ClearUntilFrameNumber(lookupCommandFrame)
		}
	} else {
		// fmt.Println("empty cfHistory", gameStateUpdate.LastInputCommandFrame, deltaGCF, len(world.GetCommandFrameHistory().CommandFrames))
	}
}

func replayInputs(
	entity entities.Entity,
	camera entities.Entity,
	startFrame int,
	cfHistory *commandframe.CommandFrameHistory,
) {
	frameIndex := startFrame + 1
	cf := cfHistory.GetCommandFrame(frameIndex)

	cfs := []*commandframe.CommandFrame{}
	for cf != nil {
		cfs = append(cfs, cf)
		frameIndex += 1
		cf = cfHistory.GetCommandFrame(frameIndex)
	}

	cfHistory.ClearFrames()

	// replay inputs and add the new results to the command frame history
	for i, cf := range cfs {
		controllerutils.UpdateCharacterController(entity, camera, cf.FrameInput)
		physutils.PhysicsStep(time.Duration(settings.MSPerCommandFrame)*time.Millisecond, []entities.Entity{entity}, entity.GetID())
		cfHistory.AddCommandFrame(startFrame+i+1, cf.FrameInput, entity)
	}
}
