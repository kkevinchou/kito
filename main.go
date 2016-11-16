package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/kkevinchou/ant/antz"
	"github.com/kkevinchou/ant/directory"
	"github.com/veandco/go-sdl2/sdl"
)

var window *sdl.Window
var renderer *sdl.Renderer

func init() {
	// We want to lock the main thread to this goroutine.  Otherwise,
	// SDL rendering will randomly panic
	//
	// For more details: https://github.com/golang/go/wiki/LockOSThread
	runtime.LockOSThread()
}

func setupDisplay() {
	sdl.Init(sdl.INIT_EVERYTHING)

	var err error

	window, err = sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(fmt.Sprintf("Failed to create renderer: %s\n", err))
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	setupDisplay()
	defer window.Destroy()
	defer renderer.Destroy()

	game := ant.Game{}
	game.Init(renderer)
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
					game.MoveAnt(float64(e.X), float64(e.Y))
				}
			case *sdl.KeyUpEvent:
				if e.Keysym.Sym == sdl.K_ESCAPE {
					gameOver = true
				}
			}
		}

		game.Update()
		movementSystem.Update(delta)
		renderSystem.Update(delta)
	}
	sdl.Quit()
}
