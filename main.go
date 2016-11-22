package main

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/kkevinchou/ant/ant"
	"github.com/kkevinchou/ant/directory"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	width  = 800
	height = 600
)

func init() {
	// We want to lock the main thread to this goroutine.  Otherwise,
	// SDL rendering will randomly panic
	//
	// For more details: https://github.com/golang/go/wiki/LockOSThread
	runtime.LockOSThread()
}

func setupDisplay() (*sdl.Window, error) {
	var err error
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, err
	}

	if err := gl.Init(); err != nil {
		return nil, err
	}

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_OPENGL)
	if err != nil {
		return nil, err
	}

	_, err = sdl.GL_CreateContext(window)
	if err != nil {
		return nil, err
	}

	return window, nil
}

func main() {
	window, err := setupDisplay()
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	rand.Seed(time.Now().Unix())

	game := ant.Game{}
	game.Init(window)

	directory := directory.GetDirectory()
	movementSystem := directory.MovementSystem()
	renderSystem := directory.RenderSystem()

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
				game.CameraView(int(e.X), int(e.Y))
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

				game.MoveCamera(cameraMovement)

			}
		}

		game.Update()
		movementSystem.Update(delta)
		renderSystem.Update(delta)
	}
	sdl.Quit()
}
