package singleton

import "github.com/kkevinchou/kito/kito/commands"

type KeyboardInputComponent struct {
	KeyboardInputSet *commands.KeyboardInputSet
}

func NewKeyboardInputComponent() *KeyboardInputComponent {
	return &KeyboardInputComponent{
		KeyboardInputSet: &commands.KeyboardInputSet{},
	}
}

func (c *KeyboardInputComponent) GetKeyboardInputSet() *commands.KeyboardInputSet {
	return c.KeyboardInputSet
}

func (c *KeyboardInputComponent) SetKeyboardInputSet(input *commands.KeyboardInputSet) {
	c.KeyboardInputSet = input
}
