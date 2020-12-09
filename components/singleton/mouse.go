package singleton

import (
	"github.com/kkevinchou/kito/types"
)

type MouseInputComponent struct {
	MouseInput *types.MouseInput
}

func NewMouseInputComponent() *MouseInputComponent {
	return &MouseInputComponent{
		MouseInput: nil,
	}
}

func (c *MouseInputComponent) GetMouseInput() *types.MouseInput {
	return c.MouseInput
}

func (c *MouseInputComponent) SetMouseInput(input *types.MouseInput) {
	c.MouseInput = input
}
