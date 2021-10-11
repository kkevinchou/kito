package networkdispatch

import (
	"fmt"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/commandframe"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/knetwork"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/kito/utils/controllerutils"
	"github.com/kkevinchou/kito/kito/utils/physutils"
	"github.com/kkevinchou/kito/lib/network"
)

func clientMessageHandler(world World, message *network.Message) {
	singleton := world.GetSingleton()
	metricsRegistry := singleton.MetricsRegistry
	if message.MessageType == knetwork.MessageTypeGameStateUpdate {
		var gameStateUpdate knetwork.GameStateUpdateMessage
		err := network.DeserializeBody(message, &gameStateUpdate)
		if err != nil {
			panic(err)
		}

		validateClientPrediction(&gameStateUpdate, world)
		singleton.StateBuffer.PushEntityUpdate(world.CommandFrame(), &gameStateUpdate)
	} else if message.MessageType == knetwork.MessageTypeAckCreatePlayer {
		handleAckCreatePlayer(message, world)
	} else if message.MessageType == knetwork.MessageTypeAckPing {
		var ackPingMessage knetwork.AckPingMessage
		err := network.DeserializeBody(message, &ackPingMessage)
		if err != nil {
			fmt.Printf("error deserializing ackping message %s\n", err)
		}

		metricsRegistry.Inc("ping", float64(time.Since(ackPingMessage.PingSendTime).Milliseconds()))
	} else {
		fmt.Println("unknown message type:", message.MessageType, string(message.Body))
	}
}

func handleAckCreatePlayer(message *network.Message, world World) {
	messageBody := &knetwork.AckCreatePlayerMessage{}
	err := network.DeserializeBody(message, messageBody)
	if err != nil {
		fmt.Println(err)
		return
	}

	singleton := world.GetSingleton()
	singleton.PlayerID = messageBody.ID
	singleton.CameraID = messageBody.CameraID

	bob := entities.NewBob(mgl64.Vec3{})
	bob.ID = messageBody.ID

	camera := entities.NewThirdPersonCamera(settings.CameraStartPosition, settings.CameraStartView, bob.GetID())
	camera.ID = messageBody.CameraID

	bob.GetComponentContainer().ThirdPersonControllerComponent.CameraID = camera.GetID()

	world.RegisterEntities([]entities.Entity{
		bob,
		camera,
	})
}

func validateClientPrediction(gameStateUpdate *knetwork.GameStateUpdateMessage, world World) {
	singleton := world.GetSingleton()
	metricsRegistry := singleton.MetricsRegistry

	// We use a gcf adjusted command frame lookup because even though an input may happen on only one command
	// frame, the entity can continue to be updated due to that input. Therefore we need to make sure
	// we advance the command frame as much as the server has to see if we've mispredicted
	deltaGCF := gameStateUpdate.CurrentGlobalCommandFrame - gameStateUpdate.LastInputGlobalCommandFrame
	lookupCommandFrame := gameStateUpdate.LastInputCommandFrame + deltaGCF

	cfHistory := world.GetCommandFrameHistory()
	cf := cfHistory.GetCommandFrame(lookupCommandFrame)
	if cf == nil {
		// We should use the latest cfHistory if we're not able to find an exact command frame history
		// with the lookup. Standing "still" is still a prediction, and if some outside factor affects the
		// player, we should detect that as a misprediction and move our character accordingly

		// Sometimes the server is a single tick ahead
		cf = cfHistory.GetCommandFrame(lookupCommandFrame - 1)
	}

	player := world.GetPlayer()
	if player == nil {
		fmt.Println("handleGameStateUpdate - could not find player")
		return
	}

	entitySnapshot := gameStateUpdate.Entities[world.GetSingleton().PlayerID]

	if cf != nil {
		historyEntity := cf.PostCFState
		if historyEntity.Position != entitySnapshot.Position || historyEntity.Orientation != entitySnapshot.Orientation {
			metricsRegistry.Inc("predictionMiss", 1)
			fmt.Printf(
				"--------------------------------------\n[CF:%d] CLIENT-SIDE PREDICTION MISS\nlastCF: %d\nlastGlobalCF: %d\ncurrentGlobalCF: %d\n%v\n%v\n",
				world.CommandFrame(),
				gameStateUpdate.LastInputCommandFrame,
				gameStateUpdate.LastInputGlobalCommandFrame,
				gameStateUpdate.CurrentGlobalCommandFrame,
				historyEntity.Position,
				entitySnapshot.Position,
			)

			prevHistoryEntity := cfHistory.GetCommandFrame(lookupCommandFrame - 1)
			nextHistoryEntity := cfHistory.GetCommandFrame(lookupCommandFrame + 1)
			prevHit := 0
			nextHit := 0
			if prevHistoryEntity.PostCFState.Position == entitySnapshot.Position && prevHistoryEntity.PostCFState.Orientation == entitySnapshot.Orientation {
				prevHit = 1
			}
			if nextHistoryEntity != nil {
				if nextHistoryEntity.PostCFState.Position == entitySnapshot.Position && nextHistoryEntity.PostCFState.Orientation == entitySnapshot.Orientation {
					nextHit = 1
				}
			}
			fmt.Printf("prevHit %d nextHit %d\n", prevHit, nextHit)

			// When we miss the client-side prediction, we set the player's state to the snapshot state.
			// what's important to note is that the snapshot state is in the past! At the very least,
			// a whole RTT in the past. So, we want to replay our historical inputs to catch up to the player's
			// present.

			cc := player.GetComponentContainer()
			cc.TransformComponent.Position = entitySnapshot.Position
			cc.TransformComponent.Orientation = entitySnapshot.Orientation
			cc.PhysicsComponent.Velocity = entitySnapshot.Velocity
			cc.PhysicsComponent.Impulses = entitySnapshot.Impulses
			replayInputs(player, world.GetCamera(), lookupCommandFrame, cfHistory)
		} else {
			metricsRegistry.Inc("predictionHit", 1)
			// fmt.Println(
			// 	"CLIENT-SIDE PREDICTION HIT",
			// 	gameStateUpdate.LastInputCommandFrame,
			// 	gameStateUpdate.LastInputGlobalCommandFrame,
			// 	gameStateUpdate.CurrentGlobalCommandFrame,
			// )
			cfHistory.ClearUntilFrameNumber(lookupCommandFrame)
		}
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
		physutils.PhysicsStep(time.Duration(settings.MSPerCommandFrame)*time.Millisecond, entity)
		cfHistory.AddCommandFrame(startFrame+i+1, cf.FrameInput, entity)
	}
	fmt.Printf("replayed %d inputs\n", len(cfs))
}
