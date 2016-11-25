package main

import (
	"runtime"

	"github.com/kkevinchou/ant/ant"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/veandco/go-sdl2/sdl"
)

func init() {
	// We want to lock the main thread to this goroutine.  Otherwise,
	// SDL rendering will randomly panic
	//
	// For more details: https://github.com/golang/go/wiki/LockOSThread
	runtime.LockOSThread()
}

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
			game.CameraViewChange(int(e.XRel), int(e.YRel))
		case *sdl.KeyUpEvent:
			if e.Keysym.Sym == sdl.K_ESCAPE {
				game.GameOver()
			}
		case *sdl.KeyDownEvent:
			cameraMovement := vector.Vector3{}
			if e.Keysym.Sym == sdl.K_w {
				cameraMovement.Z -= 1
			}

			if e.Keysym.Sym == sdl.K_s {
				cameraMovement.Z += 1
			}

			if e.Keysym.Sym == sdl.K_a {
				cameraMovement.X -= 1
			}

			if e.Keysym.Sym == sdl.K_d {
				cameraMovement.X += 1
			}

			if e.Keysym.Sym == sdl.K_SPACE {
				cameraMovement.Y += 1
			}

			game.MoveCamera(cameraMovement)
		}
	}
}

func main() {
	game := ant.NewGame()
	game.Start(CommandPoller)
	sdl.Quit()
}
