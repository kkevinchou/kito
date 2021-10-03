package statebuffer

import (
	"fmt"

	"github.com/kkevinchou/kito/kito/knetwork"
	"github.com/kkevinchou/kito/lib/libutils"
)

type BufferedState struct {
	InterpolatedEntities map[int]knetwork.EntitySnapshot
}

type IncomingEntityUpdate struct {
	targetCommandFrame     int
	gameStateUpdateMessage *knetwork.GameStateUpdateMessage
}

type StateBuffer struct {
	maxStateBufferCommandFrames int
	timeline                    map[int]BufferedState
	incomingEntityUpdate        []IncomingEntityUpdate
}

func NewStateBuffer(maxStateBufferCommandFrames int) *StateBuffer {
	return &StateBuffer{
		maxStateBufferCommandFrames: maxStateBufferCommandFrames,
		timeline:                    map[int]BufferedState{},
	}
}

func (s *StateBuffer) PushEntityUpdate(playerCommandFrame int, gameStateUpdateMessage *knetwork.GameStateUpdateMessage) {
	if len(s.incomingEntityUpdate) == 0 {
		targetCommandFrame := playerCommandFrame + s.maxStateBufferCommandFrames
		s.incomingEntityUpdate = append(
			s.incomingEntityUpdate,
			IncomingEntityUpdate{
				gameStateUpdateMessage: gameStateUpdateMessage,
				targetCommandFrame:     targetCommandFrame,
			},
		)

		s.timeline[targetCommandFrame] = BufferedState{
			InterpolatedEntities: gameStateUpdateMessage.Entities,
		}

		return
	}

	lastEntityUpdate := s.incomingEntityUpdate[len(s.incomingEntityUpdate)-1]
	gcfDelta := gameStateUpdateMessage.CurrentGlobalCommandFrame - lastEntityUpdate.gameStateUpdateMessage.CurrentGlobalCommandFrame

	currentEntityUpdate := IncomingEntityUpdate{
		gameStateUpdateMessage: gameStateUpdateMessage,
		targetCommandFrame:     lastEntityUpdate.targetCommandFrame + gcfDelta,
	}

	s.incomingEntityUpdate = append(
		s.incomingEntityUpdate,
		currentEntityUpdate,
	)

	s.generateIntermediateStateUpdates(lastEntityUpdate, currentEntityUpdate)
	s.incomingEntityUpdate = s.incomingEntityUpdate[1:]
	// fmt.Println("----------------------")
	// fmt.Println(playerCommandFrame)
	// max := 0
	// for k, _ := range s.timeline {
	// 	if k > max {
	// 		max = k
	// 	}
	// }
	// fmt.Println(max)
}

func (s *StateBuffer) generateIntermediateStateUpdates(start IncomingEntityUpdate, end IncomingEntityUpdate) {
	startGCF := start.gameStateUpdateMessage.CurrentGlobalCommandFrame
	endGCF := end.gameStateUpdateMessage.CurrentGlobalCommandFrame
	gcfDelta := endGCF - startGCF
	cfStep := float64(1) / float64(gcfDelta)

	for i := 1; i <= gcfDelta; i++ {
		interpolatedEntities := map[int]knetwork.EntitySnapshot{}

		for id, startSnapshot := range start.gameStateUpdateMessage.Entities {
			if _, ok := end.gameStateUpdateMessage.Entities[id]; !ok {
				fmt.Printf("warning, entity from start update (%d) did not exist in the next one\n", id)
			}

			endSnapshot := end.gameStateUpdateMessage.Entities[id]
			interpolatedEntities[id] = knetwork.EntitySnapshot{
				ID:          startSnapshot.ID,
				Type:        startSnapshot.Type,
				Position:    endSnapshot.Position.Sub(startSnapshot.Position).Mul(float64(i) * cfStep).Add(startSnapshot.Position),
				Orientation: libutils.QInterpolate64(startSnapshot.Orientation, endSnapshot.Orientation, float64(i)*cfStep),
			}
		}

		s.timeline[start.targetCommandFrame+i] = BufferedState{
			InterpolatedEntities: interpolatedEntities,
		}
	}
}

func (s *StateBuffer) PullEntityInterpolations(cf int) *BufferedState {
	if b, ok := s.timeline[cf]; ok {
		delete(s.timeline, cf)
		// fmt.Println(len(s.timeline))
		return &b
	}
	return nil
}
