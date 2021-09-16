package input

import (
	"github.com/veandco/go-sdl2/sdl"
)

type SDLInputPoller struct {
}

func NewSDLInputPoller() *SDLInputPoller {
	return &SDLInputPoller{}
}

func (i *SDLInputPoller) PollInput() Input {
	sdl.PumpEvents()

	// Mouse inputs
	mouseInput := MouseInput{}

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

	// Event inputs
	var commands []interface{}
	var event sdl.Event

	wheelCount := 0
	// The same event type can be fired multiple times in the same PollEvent loop
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			commands = append(commands, QuitCommand{})
		case *sdl.MouseButtonEvent:
			// ?
		case *sdl.MouseMotionEvent:
			mouseInput.MouseMotionEvent.XRel += float64(e.XRel)
			mouseInput.MouseMotionEvent.YRel += float64(e.YRel)
		case *sdl.MouseWheelEvent:
			mouseInput.MouseWheelDelta += int(e.Y)
			wheelCount++
		}
	}
	if wheelCount > 0 {
		// fmt.Println(wheelCount)
	}

	// Keyboard inputs
	// TODO: only check for keys we care about - keyState contains 512 keys
	keyboardInput := KeyboardInput{}
	keyState := sdl.GetKeyboardState()
	for k, v := range keyState {
		if v <= 0 {
			continue
		}
		key := KeyboardKey(sdl.GetScancodeName(sdl.Scancode(k)))
		keyboardInput[key] = KeyState{
			Key:   key,
			Event: KeyboardEventDown,
		}
	}

	input := Input{
		KeyboardInput: keyboardInput,
		MouseInput:    mouseInput,
		Commands:      commands,
	}

	return input
}
