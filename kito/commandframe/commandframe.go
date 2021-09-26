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
	PostCFState []EntityState
}

type CommandFrameHistory struct {
	CommandFrames map[int]*CommandFrame
}

func NewCommandFrameHistory() *CommandFrameHistory {
	return &CommandFrameHistory{
		CommandFrames: map[int]*CommandFrame{},
	}
}

func (h *CommandFrameHistory) AddCommandFrame(frameNumber int, frameInput input.Input, entities []entities.Entity) {
	cf := &CommandFrame{
		FrameNumber: frameNumber,
		FrameInput:  frameInput.Copy(),
	}
	for _, e := range entities {
		transformComponent := e.GetComponentContainer().TransformComponent
		cf.PostCFState = append(cf.PostCFState, EntityState{
			ID:          e.GetID(),
			Position:    transformComponent.Position,
			Orientation: transformComponent.Orientation,
		})
	}

	h.CommandFrames[frameNumber] = cf
}
