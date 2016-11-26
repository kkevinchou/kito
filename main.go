package main

import (
	"runtime"

	"github.com/kkevinchou/ant/ant"
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
func (i *InputHandler) CommandPoller(game *ant.Game) []ant.Command {
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

	commands := []ant.Command{}

	var event sdl.Event
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			commands = append(commands, &ant.QuitCommand{})
		case *sdl.MouseButtonEvent:
			if e.State == 0 { // Mouse Up
				// game.MoveAnt(float64(e.X), float64(e.Y))
				game.PlaceFood(float64(e.X), float64(e.Y))
			}
		case *sdl.MouseMotionEvent:
			commands = append(commands, &ant.CameraViewCommand{X: float64(e.XRel), Y: float64(e.YRel)})
		case *sdl.KeyUpEvent:
			if e.Keysym.Sym == sdl.K_ESCAPE {
				commands = append(commands, &ant.QuitCommand{})
			}
		}
	}

	commands = append(commands, &ant.SetCameraSpeed{X: x, Y: y, Z: z})

	return commands
}

func main() {
	game := ant.NewGame()
	inputHandler := NewInputHandler()
	game.Start(inputHandler.CommandPoller)
	sdl.Quit()
}
