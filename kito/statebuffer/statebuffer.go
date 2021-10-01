package statebuffer

import (
	"fmt"

	"github.com/kkevinchou/kito/lib/libutils"
	"github.com/kkevinchou/kito/lib/network"
)

type BufferedState struct {
	InterpolatedEntities map[int]network.EntitySnapshot
}

type IncomingEntityUpdate struct {
	targetCommandFrame     int
	gameStateUpdateMessage *network.GameStateUpdateMessage
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

func (s *StateBuffer) PushEntityUpdate(playerCommandFrame int, gameStateUpdateMessage *network.GameStateUpdateMessage) {
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
		interpolatedEntities := map[int]network.EntitySnapshot{}

		for id, startSnapshot := range start.gameStateUpdateMessage.Entities {
			if _, ok := end.gameStateUpdateMessage.Entities[id]; !ok {
				fmt.Printf("warning, entity from start update (%d) did not exist in the next one\n", id)
			}

			endSnapshot := end.gameStateUpdateMessage.Entities[id]
			interpolatedEntities[id] = network.EntitySnapshot{
				ID:          startSnapshot.ID,
				Type:        startSnapshot.Type,
				Position:    endSnapshot.Position.Sub(startSnapshot.Position).Mul(float64(i) * cfStep).Add(startSnapshot.Position),
				Orientation: libutils.QInterpolate64(startSnapshot.Orientation, endSnapshot.Orientation, float64(i)*cfStep),
				// TODO: Orientation
			}
		}

		s.timeline[start.targetCommandFrame+i] = BufferedState{
			InterpolatedEntities: interpolatedEntities,
		}
		// fmt.Printf("cf %d size %d\n", start.targetCommandFrame+i, i)
	}

	// s.timeline[end.targetCommandFrame] = BufferedState{
	// 	interpolatedEntities: end.gameStateUpdateMessage.Entities,
	// }
}

func (s *StateBuffer) PullEntityInterpolations(cf int) *BufferedState {
	if b, ok := s.timeline[cf]; ok {
		// delete(s.timeline, cf)
		// fmt.Printf("pulling %d\n", cf)
		return &b
	}
	// fmt.Printf("pulling miss %d\n", cf)
	return nil
}

// // AppendCF assumes updates are appended in order
// func (s *StateBuffer) AppendCF(cf int, update *network.GameStateUpdateMessage) {
// 	s.stateBuffer[cf] = &CommandFrame{
// 		cf:     cf,
// 		update: update,
// 	}
// 	s.cfBuffer = append(s.cfBuffer, cf)
// }

// func (s *StateBuffer) Interpolate(cf int) *network.GameStateUpdateMessage {
// 	if len(s.cfBuffer) < s.bufferSize {
// 		if len(s.cfBuffer) > 0 {
// 			// if we're slow in receiving updates we shouldn't return nil
// 			// TODO: have a less aggressive solution than halting the world
// 			return s.stateBuffer[s.cfBuffer[0]].update
// 		}
// 		return nil
// 	}

// 	if cf < s.cfBuffer[0] {
// 		// command frame is in the past
// 		return nil
// 	}

// 	if cf == s.cfBuffer[0] {
// 		return s.stateBuffer[cf].update
// 	}

// 	// TODO: handle cf that is beyond position 1
// 	if cf == s.cfBuffer[1] {
// 		delete(s.stateBuffer, s.cfBuffer[0])
// 		s.cfBuffer = s.cfBuffer[1:]
// 		return s.stateBuffer[cf].update
// 	}

// 	stateBundle := interpolate(cf, s.stateBuffer[s.cfBuffer[0]], s.stateBuffer[s.cfBuffer[1]])

// 	return stateBundle.update
// }
