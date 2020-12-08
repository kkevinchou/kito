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
}

func NewInputHandler() *InputHandler {
	return &InputHandler{}
}

// TODO: event polling will return no events even though the key is being held down
func (i *InputHandler) CommandPoller() []commands.Command {
	commandList := []commands.Command{}

	keyboardInputSet := commands.KeyboardInputSet{}

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
		case *sdl.KeyboardEvent:
			var repeat bool
			if e.Repeat >= 1 {
				repeat = true
			}
			key := commands.KeyboardKey(sdl.GetKeyName(e.Keysym.Sym))

			var keyboardEvent commands.KeyboardEvent
			if e.Type == sdl.KEYUP {
				keyboardEvent = commands.KeyboardEventUp
			} else if e.Type == sdl.KEYDOWN {
				keyboardEvent = commands.KeyboardEventDown
			} else {
				panic("unexpected keyboard event" + string(e.Type))
			}

			keyboardInputSet[key] = commands.KeyboardInput{
				Key:    key,
				Repeat: repeat,
				Event:  keyboardEvent,
			}
		}
	}

	// // TODO: only check for keys we care about - keyState contains 512 keys
	// sdl.PumpEvents()
	// keyState := sdl.GetKeyboardState()
	// for k, v := range keyState {
	// 	if v <= 0 {
	// 		continue
	// 	}
	// 	key := commands.KeyboardKey(sdl.GetScancodeName(sdl.Scancode(k)))
	// 	keyboardInputSet[key] = commands.KeyboardInput{
	// 		Key:   key,
	// 		Event: commands.KeyboardEventDown,
	// 	}
	// }

	commandList = append(commandList, &keyboardInputSet)

	return commandList
}

func main() {
	game := kito.NewGame()
	inputHandler := NewInputHandler()
	game.Start(inputHandler.CommandPoller)
	sdl.Quit()
}
