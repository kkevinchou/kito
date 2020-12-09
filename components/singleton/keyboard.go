package singleton

import "github.com/kkevinchou/kito/types"

type KeyboardInputComponent struct {
	KeyboardInput *types.KeyboardInput
}

func NewKeyboardInputComponent() *KeyboardInputComponent {
	return &KeyboardInputComponent{
		KeyboardInput: &types.KeyboardInput{},
	}
}

func (c *KeyboardInputComponent) GetKeyboardInputSet() *types.KeyboardInput {
	return c.KeyboardInput
}

func (c *KeyboardInputComponent) SetKeyboardInputSet(input *types.KeyboardInput) {
	c.KeyboardInput = input
}
