package main

import (
	"math/rand"
	"runtime"
	"time"

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

func main() {
	rand.Seed(time.Now().Unix())
	game := ant.NewGame()

	var event sdl.Event
	gameOver := false

	previousTime := time.Now()
	for gameOver != true {
		now := time.Now()
		delta := time.Since(previousTime)
		previousTime = now

		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				gameOver = true
			case *sdl.MouseButtonEvent:
				if e.State == 0 { // Mouse Up
					// game.MoveAnt(float64(e.X), float64(e.Y))
					game.PlaceFood(float64(e.X), float64(e.Y))
				}
			case *sdl.MouseMotionEvent:
				game.CameraViewChange(int(e.XRel), int(e.YRel))
			case *sdl.KeyUpEvent:
				if e.Keysym.Sym == sdl.K_ESCAPE {
					gameOver = true
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

		game.Update(delta)
	}
	sdl.Quit()
}
