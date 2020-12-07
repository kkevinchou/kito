package main

import (
	"runtime"

	"github.com/kkevinchou/kito/kito"
	"github.com/kkevinchou/kito/kito/commands"
	"github.com/kkevinchou/kito/lib/math/vector"
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
func (i *InputHandler) CommandPoller() []commands.Command {
	commandList := []commands.Command{}

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
	var zoom int

	if i.KeyState[sdl.SCANCODE_ESCAPE] > 0 {
		commandList = append(commandList, &commands.QuitCommand{})
	}

	sdl.PumpEvents()
	var event sdl.Event
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			commandList = append(commandList, &commands.QuitCommand{})
		case *sdl.MouseButtonEvent:
			cameraControlled := false
			if e.State == sdl.RELEASED { // Mouse Up
				if e.Button == sdl.BUTTON_LEFT {
				} else if e.Button == sdl.BUTTON_RIGHT {
				} else if e.Button == sdl.BUTTON_MIDDLE {
					// commands = append(commands, &kito.CameraRaycastCommand{X: float64(e.X), Y: float64(e.Y)})
				}
			} else if e.State == sdl.PRESSED {
				if e.Button == sdl.BUTTON_LEFT {
					cameraControlled = true
				}
			}

			commandList = append(commandList, &commands.ToggleCameraControlCommand{Value: cameraControlled})

		case *sdl.MouseMotionEvent:
			x := float64(e.XRel)
			y := float64(e.YRel)

			commandList = append(commandList, &commands.UpdateViewCommand{Value: vector.Vector{x, y}})
		case *sdl.MouseWheelEvent:
			zoom = int(e.Y)
		}
	}
	commandList = append(commandList, &commands.MoveCommand{Value: vector.Vector3{x, y, z}, Zoom: zoom})

	return commandList
}

func main() {
	game := kito.NewGame()
	inputHandler := NewInputHandler()
	game.Start(inputHandler.CommandPoller)
	sdl.Quit()
}
