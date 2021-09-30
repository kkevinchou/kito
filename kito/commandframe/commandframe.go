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
	ReadOffset    int
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

	// this handles the initial client startup phase where a command frame goes buy without recording a command frame.
	// this is due to the player not existing yet. there's probably a cleaner way to handle this
	if len(h.CommandFrames) == 1 {
		h.ReadOffset = frameNumber
	}
}

func (h *CommandFrameHistory) GetCommandFrame(frameNumber int) *CommandFrame {
	if frameNumber-h.ReadOffset >= len(h.CommandFrames) {
		return nil
	}
	return &h.CommandFrames[frameNumber-h.ReadOffset]
}

func (h *CommandFrameHistory) ClearUntilFrameNumber(frameNumber int) {
	// fmt.Println(frameNumber-h.ReadOffset, len(h.CommandFrames))
	h.CommandFrames = h.CommandFrames[frameNumber-h.ReadOffset:]
	h.ReadOffset = frameNumber
}
