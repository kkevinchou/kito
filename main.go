package main

import (
	"runtime"

	"github.com/kkevinchou/kito/kito"
	"github.com/veandco/go-sdl2/sdl"
)

func init() {
	// We want to lock the main thread to this goroutine.  Otherwise,
	// SDL rendering will randomly panic
	//
	// For more details: https://github.com/golang/go/wiki/LockOSThread
	runtime.LockOSThread()
}

type InputHandler struct {
	KeyState []uint8
}

func NewInputHandler() *InputHandler {
	return &InputHandler{KeyState: sdl.GetKeyboardState()}
}

// TODO: event polling will return no events even though the key is being held down
func (i *InputHandler) CommandPoller(game *kito.Game) []kito.Command {
	commands := []kito.Command{}

	sdl.PumpEvents()
	var x, y, z float64

	if i.KeyState[sdl.SCANCODE_W] > 0 {
		z--
	}
	if i.KeyState[sdl.SCANCODE_S] > 0 {
		z++
	}
	if i.KeyState[sdl.SCANCODE_A] > 0 {
		x--
	}
	if i.KeyState[sdl.SCANCODE_D] > 0 {
		x++
	}
	if i.KeyState[sdl.SCANCODE_SPACE] > 0 {
		y++
	}
	if i.KeyState[sdl.SCANCODE_LSHIFT] > 0 {
		y--
	}
	if i.KeyState[sdl.SCANCODE_ESCAPE] > 0 {
		commands = append(commands, &kito.QuitCommand{})
	}

	var event sdl.Event
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			commands = append(commands, &kito.QuitCommand{})
		case *sdl.MouseButtonEvent:
			cameraControlled := false
			if e.State == sdl.RELEASED { // Mouse Up
				if e.Button == sdl.BUTTON_LEFT {
				} else if e.Button == sdl.BUTTON_RIGHT {
				} else if e.Button == sdl.BUTTON_MIDDLE {
					commands = append(commands, &kito.CameraRaycastCommand{X: float64(e.X), Y: float64(e.Y)})
				}
			} else if e.State == sdl.PRESSED {
				if e.Button == sdl.BUTTON_LEFT {
					cameraControlled = true
				}
			}

			commands = append(commands, &kito.SetCameraControlCommand{Value: cameraControlled})

		case *sdl.MouseMotionEvent:
			x := float64(e.XRel)
			y := float64(e.YRel)

			commands = append(commands, &kito.CameraViewCommand{X: x, Y: y})
		}
	}

	commands = append(commands, &kito.SetCameraSpeed{X: x, Y: y, Z: z})

	return commands
}

func main() {
	game := kito.NewGame()
	inputHandler := NewInputHandler()
	game.Start(inputHandler.CommandPoller)
	sdl.Quit()
}
