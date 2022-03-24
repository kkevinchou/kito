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
	incomingEntityUpdates       []IncomingEntityUpdate
}

func NewStateBuffer(maxStateBufferCommandFrames int) *StateBuffer {
	return &StateBuffer{
		maxStateBufferCommandFrames: maxStateBufferCommandFrames,
		timeline:                    map[int]BufferedState{},
	}
}

func (s *StateBuffer) PushEntityUpdate(localCommandFrame int, gameStateUpdateMessage *knetwork.GameStateUpdateMessage) {
	if len(s.incomingEntityUpdates) == 0 {
		targetCommandFrame := localCommandFrame + s.maxStateBufferCommandFrames
		s.incomingEntityUpdates = append(
			s.incomingEntityUpdates,
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

	lastEntityUpdate := s.incomingEntityUpdates[len(s.incomingEntityUpdates)-1]
	gcfDelta := gameStateUpdateMessage.CurrentGlobalCommandFrame - lastEntityUpdate.gameStateUpdateMessage.CurrentGlobalCommandFrame

	currentEntityUpdate := IncomingEntityUpdate{
		gameStateUpdateMessage: gameStateUpdateMessage,
		targetCommandFrame:     lastEntityUpdate.targetCommandFrame + gcfDelta,
	}

	s.incomingEntityUpdates = append(
		s.incomingEntityUpdates,
		currentEntityUpdate,
	)

	s.generateIntermediateStateUpdates(lastEntityUpdate, currentEntityUpdate)
	s.incomingEntityUpdates = s.incomingEntityUpdates[1:]
}

// TODO: move interpolation logic in stateinterpolator system?
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
				Velocity:    endSnapshot.Velocity.Sub(startSnapshot.Velocity).Mul(float64(i) * cfStep).Add(startSnapshot.Velocity),
			}
		}

		s.timeline[start.targetCommandFrame+i] = BufferedState{
			InterpolatedEntities: interpolatedEntities,
		}
	}
}

func (s *StateBuffer) PeekEntityInterpolations(cf int) *BufferedState {
	if b, ok := s.timeline[cf]; ok {
		return &b
	}
	return nil
}

func (s *StateBuffer) PullEntityInterpolations(cf int) *BufferedState {
	if b := s.PeekEntityInterpolations(cf); b != nil {
		delete(s.timeline, cf)
		return b
	}
	return nil
}
