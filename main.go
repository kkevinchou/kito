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

// TODO: event polling will return no events even though the key is being held down
func CommandPoller(game *ant.Game) {
	var event sdl.Event
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			game.GameOver()
		case *sdl.MouseButtonEvent:
			if e.State == 0 { // Mouse Up
				// game.MoveAnt(float64(e.X), float64(e.Y))
				game.PlaceFood(float64(e.X), float64(e.Y))
			}
		case *sdl.MouseMotionEvent:
			game.CameraViewChange(float64(e.XRel), float64(e.YRel))
		case *sdl.KeyUpEvent:
			if e.Keysym.Sym == sdl.K_ESCAPE {
				game.GameOver()
			}
		case *sdl.KeyDownEvent:
			var x, y, z float64

			if e.Keysym.Sym == sdl.K_w {
				z--
			} else if e.Keysym.Sym == sdl.K_s {
				z++
			} else if e.Keysym.Sym == sdl.K_a {
				x--
			} else if e.Keysym.Sym == sdl.K_d {
				x++
			} else if e.Keysym.Sym == sdl.K_SPACE {
				y++
			}

			game.MoveCamera(x, y, z)
		}
	}
}

func main() {
	game := ant.NewGame()
	game.Start(CommandPoller)
	sdl.Quit()
}
