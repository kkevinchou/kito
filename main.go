package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/kkevinchou/ant/assets"
	"github.com/kkevinchou/ant/entities/ant"
	"github.com/kkevinchou/ant/entities/food"
	"github.com/kkevinchou/ant/entities/grass"
	"github.com/kkevinchou/ant/lib/geometry"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/kkevinchou/ant/managers/item"
	"github.com/kkevinchou/ant/managers/path"
	"github.com/kkevinchou/ant/pathing"
	"github.com/kkevinchou/ant/systems"
	"github.com/kkevinchou/ant/systems/movement"
	"github.com/kkevinchou/ant/systems/render"
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

func setupGrass() {
	grass.New(366, 450)
	grass.New(386, 450)
	grass.New(406, 450)
	grass.New(406, 350)
	grass.New(436, 350)
}

func setupSystems() *systems.Directory {
	itemManager := item.NewManager()
	pathManager := path.NewManager()
	assetManager := assets.NewAssetManager(renderer, "assets")
	renderSystem := render.NewRenderSystem(renderer, assetManager)
	movementSystem := movement.NewMovementSystem()

	d := systems.GetDirectory()
	d.RegisterRenderSystem(renderSystem)
	d.RegisterMovementSystem(movementSystem)
	d.RegisterAssetManager(assetManager)
	d.RegisterItemManager(itemManager)
	d.RegisterPathManager(pathManager)

	renderSystem.Register(pathManager.NavMesh())

	return d
}

func main() {
	rand.Seed(time.Now().Unix())
	setupDisplay()
	defer window.Destroy()
	defer renderer.Destroy()

	directory := setupSystems()
	renderSystem := directory.RenderSystem()
	movementSystem := directory.MovementSystem()
	pathManager := directory.PathManager()

	ant := ant.New()
	ant.SetPosition(vector.Vector{400, 350})

	setupGrass()

	p := pathing.Planner{}
	p.SetNavMesh(pathManager.NavMesh())

	food.New(150, 100)

	var event sdl.Event
	gameOver := false

	var path []pathing.Node
	pathIndex := 0

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
					position := ant.Position()
					path = p.FindPath(
						geometry.Point{X: position.X, Y: position.Y},
						geometry.Point{X: float64(e.X), Y: float64(e.Y)},
					)
					if path != nil {
						pathIndex = 1
						ant.SetTarget(path[1].Vector())
					}
				}
			case *sdl.KeyUpEvent:
				if e.Keysym.Sym == sdl.K_ESCAPE {
					gameOver = true
				}
			}
		}

		if path != nil {
			if ant.Position().Sub(path[pathIndex].Vector()).Length() <= 2 {
				pathIndex += 1
				if pathIndex == len(path) {
					path = nil
					ant.SetSeekActive(false)
					ant.SetVelocity(vector.Zero())
				} else {
					ant.SetTarget(path[pathIndex].Vector())
				}
			}
		}

		movementSystem.Update(delta)
		renderSystem.Update(delta)
	}
	sdl.Quit()
}
