package singleton

import "github.com/kkevinchou/kito/kito/commands"

type InputComponent struct {
	KeyboardInputSet *commands.KeyboardInputSet
}

func New() *InputComponent {
	return &InputComponent{
		KeyboardInputSet: &commands.KeyboardInputSet{},
	}
}

func (c *InputComponent) GetKeyboardInputSet() *commands.KeyboardInputSet {
	return c.KeyboardInputSet
}

func (c *InputComponent) SetKeyboardInputSet(input *commands.KeyboardInputSet) {
	c.KeyboardInputSet = input
}
