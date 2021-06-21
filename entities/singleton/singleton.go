package singleton

import (
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/components/singleton"
	"github.com/kkevinchou/kito/settings"
)

type Singleton struct {
	*singleton.KeyboardInputComponent
	*singleton.MouseInputComponent
	*components.ConnectionComponent
}

func NewServerSingleton() *Singleton {
	return &Singleton{
		KeyboardInputComponent: singleton.NewKeyboardInputComponent(),
		MouseInputComponent:    singleton.NewMouseInputComponent(),
	}
}

func NewClientSingleton() *Singleton {
	connectionComponent, err := components.NewConnectionComponent(
		settings.Host,
		settings.Port,
		settings.ConnectionType,
	)
	if err != nil {
		panic(err)
	}

	return &Singleton{
		KeyboardInputComponent: singleton.NewKeyboardInputComponent(),
		MouseInputComponent:    singleton.NewMouseInputComponent(),
		ConnectionComponent:    connectionComponent,
	}
}
