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
	CommandFrames map[int]*CommandFrame
}

func NewCommandFrameHistory() *CommandFrameHistory {
	return &CommandFrameHistory{
		CommandFrames: map[int]*CommandFrame{},
	}
}

func (h *CommandFrameHistory) AddCommandFrame(frameNumber int, frameInput input.Input, player entities.Entity) {
	cf := &CommandFrame{
		FrameNumber: frameNumber,
		FrameInput:  frameInput.Copy(),
	}
	transformComponent := player.GetComponentContainer().TransformComponent

	cf.PostCFState = EntityState{
		ID:          player.GetID(),
		Position:    transformComponent.Position,
		Orientation: transformComponent.Orientation,
	}

	h.CommandFrames[frameNumber] = cf
}
