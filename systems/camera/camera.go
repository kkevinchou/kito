package camera

import (
	"time"

	"github.com/kkevinchou/kito/interfaces"
	"github.com/kkevinchou/kito/kito/commands"
	"github.com/kkevinchou/kito/lib/math/vector"
)

type Singleton interface {
	GetKeyboardInputSet() *commands.KeyboardInputSet
}

type World interface {
	GetSingleton() interfaces.Singleton
	GetCamera() interfaces.Viewer
}

type CameraSystem struct {
	world World
}

func NewCameraSystem(world World) *CameraSystem {
	s := CameraSystem{
		world: world,
	}
	return &s
}

// if i.KeyState[sdl.SCANCODE_W] > 0 {
// 	z--
// }
// if i.KeyState[sdl.SCANCODE_S] > 0 {
// 	z++
// }
// if i.KeyState[sdl.SCANCODE_A] > 0 {
// 	x--
// }
// if i.KeyState[sdl.SCANCODE_D] > 0 {
// 	x++
// }
// if i.KeyState[sdl.SCANCODE_SPACE] > 0 {
// 	y++
// }
// if i.KeyState[sdl.SCANCODE_LSHIFT] > 0 {
// 	y--
// }

func (s *CameraSystem) Update(delta time.Duration) {
	camera := s.world.GetCamera()
	keyboardInputSet := *s.world.GetSingleton().GetKeyboardInputSet()

	var controlVector vector.Vector3
	if key, ok := keyboardInputSet[commands.KeyboardKeyW]; ok && key.Event == commands.KeyboardEventDown {
		controlVector.Z--
	}
	if key, ok := keyboardInputSet[commands.KeyboardKeyS]; ok && key.Event == commands.KeyboardEventDown {
		controlVector.Z++
	}
	if key, ok := keyboardInputSet[commands.KeyboardKeyA]; ok && key.Event == commands.KeyboardEventDown {
		controlVector.X--
	}
	if key, ok := keyboardInputSet[commands.KeyboardKeyD]; ok && key.Event == commands.KeyboardEventDown {
		controlVector.X++
	}
	if key, ok := keyboardInputSet[commands.KeyboardKeyLShift]; ok && key.Event == commands.KeyboardEventDown {
		controlVector.Y--
	}
	if key, ok := keyboardInputSet[commands.KeyboardKeySpace]; ok && key.Event == commands.KeyboardEventDown {
		controlVector.Y++
	}

	camera.SetControlDirection(controlVector, 0)
}
