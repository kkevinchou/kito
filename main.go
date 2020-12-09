package main

import (
	"runtime"

	"github.com/kkevinchou/kito/kito"
	"github.com/kkevinchou/kito/kito/commands"
	"github.com/kkevinchou/kito/lib/math/vector"
	"github.com/kkevinchou/kito/types"
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
}

func NewInputHandler() *InputHandler {
	return &InputHandler{}
}

// TODO: event polling will return no events even though the key is being held down
func (i *InputHandler) CommandPoller() []commands.Command {
	commandList := []commands.Command{}

	keyboardInput := types.KeyboardInput{}
	mouseInput := types.MouseInput{
		MouseWheel: types.MouseWheelDirectionNeutral,
	}

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
			// zoom = int(e.Y)
			direction := types.MouseWheelDirectionNeutral
			if e.Y > 0 {
				direction = types.MouseWheelDirectionUp
			} else if e.Y < 0 {
				direction = types.MouseWheelDirectionDown
			}
			mouseInput.MouseWheel = direction
		case *sdl.KeyboardEvent:
			// var repeat bool
			// if e.Repeat >= 1 {
			// 	repeat = true
			// }
			// key := types.KeyboardKey(sdl.GetKeyName(e.Keysym.Sym))

			// var keyboardEvent types.KeyboardEvent
			// if e.Type == sdl.KEYUP {
			// 	keyboardEvent = types.KeyboardEventUp
			// } else if e.Type == sdl.KEYDOWN {
			// 	keyboardEvent = types.KeyboardEventDown
			// } else {
			// 	panic("unexpected keyboard event" + string(e.Type))
			// }

			// keyboardInput[key] = types.KeyboardInput{
			// 	Key:    key,
			// 	Repeat: repeat,
			// 	Event:  keyboardEvent,
			// }
		}
	}

	// TODO: only check for keys we care about - keyState contains 512 keys
	sdl.PumpEvents()
	keyState := sdl.GetKeyboardState()
	for k, v := range keyState {
		if v <= 0 {
			continue
		}
		key := types.KeyboardKey(sdl.GetScancodeName(sdl.Scancode(k)))
		keyboardInput[key] = types.KeyState{
			Key:   key,
			Event: types.KeyboardEventDown,
		}
	}

	commandList = append(commandList, &keyboardInput)
	commandList = append(commandList, &mouseInput)

	return commandList
}

func main() {
	game := kito.NewGame()
	inputHandler := NewInputHandler()
	game.Start(inputHandler.CommandPoller)
	sdl.Quit()
}
