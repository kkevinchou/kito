package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/kkevinchou/kito/kito"
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

type InputPoller struct {
}

func NewInputPoller() *InputPoller {
	return &InputPoller{}
}

func (i *InputPoller) PollInput() []interface{} {
	inputList := []interface{}{}

	keyboardInput := types.KeyboardInput{}
	mouseInput := types.MouseInput{
		MouseWheel: types.MouseWheelDirectionNeutral,
	}

	_, _, mouseState := sdl.GetMouseState()
	if mouseState&sdl.BUTTON_LEFT > 0 {
		mouseInput.LeftButtonDown = true
	}
	if mouseState&sdl.BUTTON_MIDDLE > 0 {
		mouseInput.MiddleButtonDown = true
	}
	if mouseState&sdl.BUTTON_RIGHT > 0 {
		mouseInput.RightButtonDown = true
	}

	var event sdl.Event

	// The same event can be fired multiple times in the same PollEvent loop
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			inputList = append(inputList, &types.QuitCommand{})
		case *sdl.MouseButtonEvent:
			// ?
		case *sdl.MouseMotionEvent:
			if mouseInput.MouseMotionEvent == nil {
				mouseInput.MouseMotionEvent = &types.MouseMotionEvent{}
			}
			mouseInput.MouseMotionEvent.XRel += float64(e.XRel)
			mouseInput.MouseMotionEvent.YRel += float64(e.YRel)
		case *sdl.MouseWheelEvent:
			direction := types.MouseWheelDirectionNeutral
			if e.Y > 0 {
				direction = types.MouseWheelDirectionUp
			} else if e.Y < 0 {
				direction = types.MouseWheelDirectionDown
			}
			mouseInput.MouseWheel = direction
		case *sdl.KeyboardEvent:
			// ?
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

	inputList = append(inputList, &keyboardInput)
	inputList = append(inputList, &mouseInput)

	return inputList
}

const (
	modeLocal  string = "LOCAL"
	modeClient string = "CLIENT"
	modeServer string = "SERVER"

	host           = "localhost"
	port           = "8080"
	connectionType = "tcp"
)

func main() {
	var mode string = modeLocal
	if len(os.Args) > 1 {
		mode = strings.ToUpper(os.Args[1])
		if mode != modeLocal && mode != modeClient && mode != modeServer {
			panic(fmt.Sprintf("unexpected mode %s", mode))
		}
	}

	fmt.Println("starting game on mode:", mode)
	if mode == modeLocal {
		serverGame := kito.NewServerGame("_assets", "shaders")

		go func() {
			serverGame.Start(kito.NullInputPoller)
		}()

		game := kito.NewLocalGame("_assets", "shaders")
		inputPoller := NewInputPoller()
		game.Start(inputPoller.PollInput)
	} else if mode == modeClient {
		game := kito.NewClientGame("_assets", "shaders")
		inputPoller := NewInputPoller()
		game.Start(inputPoller.PollInput)

	} else if mode == modeServer {
		game := kito.NewServerGame("_assets", "shaders")
		game.Start(kito.NullInputPoller)
	}

	sdl.Quit()
}
