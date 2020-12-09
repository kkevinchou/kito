package singleton

import (
	"github.com/kkevinchou/kito/components/singleton"
)

type Singleton struct {
	*singleton.KeyboardInputComponent
	*singleton.MouseInputComponent
}

func New() *Singleton {
	return &Singleton{
		KeyboardInputComponent: singleton.NewKeyboardInputComponent(),
		MouseInputComponent:    singleton.NewMouseInputComponent(),
	}
}
