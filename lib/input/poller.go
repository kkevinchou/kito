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
	var commands []any
	var event sdl.Event

	// used as a flag for the network input system to determine whether the player
	// has triggered new input which warrants notifying the server
	newInput := false

	// Keyboard inputs
	// TODO: only check for keys we care about - keyState contains 512 keys
	keyboardInput := KeyboardInput{}

	// The same event type can be fired multiple times in the same PollEvent loop
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			commands = append(commands, QuitCommand{})
		case *sdl.MouseButtonEvent:
			newInput = true
			// ?
		case *sdl.MouseMotionEvent:
			mouseInput.MouseMotionEvent.XRel += float64(e.XRel)
			mouseInput.MouseMotionEvent.YRel += float64(e.YRel)
			if mouseInput.LeftButtonDown {
				newInput = true
			}
		case *sdl.MouseWheelEvent:
			newInput = true
			mouseInput.MouseWheelDelta += int(e.Y)
		case *sdl.KeyboardEvent:
			if e.Type == sdl.KEYUP {
				newInput = true
				key := KeyboardKey(sdl.GetScancodeName(e.Keysym.Scancode))
				keyboardInput[key] = KeyState{
					Key:   key,
					Event: KeyboardEventUp,
				}
			}
			// mouseInput.MouseWheelDelta += int(e.Y)
		}
	}

	keyState := sdl.GetKeyboardState()
	for k, v := range keyState {
		if v <= 0 {
			continue
		}
		newInput = true
		key := KeyboardKey(sdl.GetScancodeName(sdl.Scancode(k)))

		// don't overwrite keys we've fetched from sdl.PollEvent()
		if _, ok := keyboardInput[key]; !ok {
			keyboardInput[key] = KeyState{
				Key:   key,
				Event: KeyboardEventDown,
			}
		}
	}

	// TODO: make input return a null input on no new input for safety
	input := Input{
		NewInput:      newInput,
		KeyboardInput: keyboardInput,
		MouseInput:    mouseInput,
		Commands:      commands,
	}

	return input
}
