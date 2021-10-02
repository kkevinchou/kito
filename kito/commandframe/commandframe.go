package commandframe

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/lib/input"
)

type EntityState struct {
	ID          int
	Position    mgl64.Vec3
	Orientation mgl64.Quat
	// TODO scale
}

type CommandFrame struct {
	FrameNumber int
	FrameInput  input.Input
	PostCFState EntityState
}

type CommandFrameHistory struct {
	CommandFrames []CommandFrame
}

func NewCommandFrameHistory() *CommandFrameHistory {
	return &CommandFrameHistory{
		CommandFrames: []CommandFrame{},
	}
}

func (h *CommandFrameHistory) AddCommandFrame(frameNumber int, frameInput input.Input, player entities.Entity) {
	transformComponent := player.GetComponentContainer().TransformComponent

	cf := CommandFrame{
		FrameNumber: frameNumber,
		FrameInput:  frameInput.Copy(),
		PostCFState: EntityState{
			ID:          player.GetID(),
			Position:    transformComponent.Position,
			Orientation: transformComponent.Orientation,
		},
	}

	h.CommandFrames = append(h.CommandFrames, cf)
}

func (h *CommandFrameHistory) GetCommandFrame(frameNumber int) *CommandFrame {
	if len(h.CommandFrames) == 0 {
		return nil
	}

	startFrameNumber := h.CommandFrames[0].FrameNumber
	if frameNumber-startFrameNumber >= len(h.CommandFrames) {
		return nil
	}

	return &h.CommandFrames[frameNumber-startFrameNumber]
}

func (h *CommandFrameHistory) ClearUntilFrameNumber(frameNumber int) {
	if len(h.CommandFrames) == 0 {
		return
	}

	startFrameNumber := h.CommandFrames[0].FrameNumber
	h.CommandFrames = h.CommandFrames[frameNumber-startFrameNumber:]
}

// func (h *CommandFrameHistory) PopUntilFrameNumberInclusive(frameNumber int) {
// 	startFrameNumber := h.CommandFrames[0].FrameNumber
// 	h.CommandFrames = h.CommandFrames[frameNumber-startFrameNumber:]
// }

func (h *CommandFrameHistory) ClearFrames() {
	h.CommandFrames = []CommandFrame{}
}
