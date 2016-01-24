package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/kkevinchou/ant/assets"
	"github.com/kkevinchou/ant/entities/ant"
	"github.com/kkevinchou/ant/entities/grass"
	"github.com/kkevinchou/ant/lib/geometry"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/kkevinchou/ant/movement"
	"github.com/kkevinchou/ant/pathing"
	"github.com/kkevinchou/ant/render"
	"github.com/veandco/go-sdl2/sdl"
)

func init() {
	// We want to lock the main thread to this goroutine.  Otherwise,
	// SDL rendering will randomly panic
	//
	// For more details: https://github.com/golang/go/wiki/LockOSThread
	runtime.LockOSThread()
}

func setupWindow() *sdl.Window {
	sdl.Init(sdl.INIT_EVERYTHING)

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	return window
}

func sqWithOffset(size, xOffset, yOffset float64) *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{xOffset * size, yOffset * size},
		geometry.Point{xOffset * size, yOffset*size + size},
		geometry.Point{xOffset*size + size, yOffset*size + size},
		geometry.Point{xOffset*size + size, yOffset * size},
	}
	return geometry.NewPolygon(points)
}

func funkyShape1() *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{180, 360},
		geometry.Point{180, 420},
		geometry.Point{600, 560},
		geometry.Point{400, 120},
	}
	return geometry.NewPolygon(points)
}

func funkyShape2() *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{500, 50},
		geometry.Point{300, 100},
		geometry.Point{400, 100},
	}
	return geometry.NewPolygon(points)
}

func setupNavMesh() *pathing.NavMesh {
	polygons := []*geometry.Polygon{
		sqWithOffset(60, 0, 0),
		sqWithOffset(60, 1, 0),
		sqWithOffset(60, 2, 0),
		sqWithOffset(60, 2, 1),
		sqWithOffset(60, 2, 2),
		sqWithOffset(60, 1, 2),
		sqWithOffset(60, 0, 2),
		sqWithOffset(60, 0, 3),
		sqWithOffset(60, 0, 4),
		sqWithOffset(60, 1, 4),
		sqWithOffset(60, 2, 4),
		sqWithOffset(60, 2, 5),
		sqWithOffset(60, 2, 6),
		sqWithOffset(60, 1, 6),
		sqWithOffset(60, 0, 6),
		funkyShape1(),
		funkyShape2(),
	}

	return pathing.ConstructNavMesh(polygons)
}

func setupGrass(assetManager *assets.Manager, renderSystem *render.RenderSystem) {
	grass2 := grass.New(366, 450, assetManager)
	grass1 := grass.New(386, 450, assetManager)
	grass3 := grass.New(406, 450, assetManager)
	grass4 := grass.New(406, 350, assetManager)
	grass5 := grass.New(436, 350, assetManager)
	renderSystem.Register(grass1)
	renderSystem.Register(grass2)
	renderSystem.Register(grass3)
	renderSystem.Register(grass4)
	renderSystem.Register(grass5)
}

func main() {
	window := setupWindow()
	defer window.Destroy()

	rand.Seed(time.Now().Unix())

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(fmt.Sprintf("Failed to create renderer: %s\n", err))
	}
	defer renderer.Destroy()

	assetManager := assets.NewAssetManager(renderer, "assets")
	renderSystem := render.NewRenderSystem(renderer, assetManager)

	ant := ant.New(assetManager)
	ant.SetPosition(vector.Vector{1, 1})

	movementSystem := movement.NewMovementSystem()
	movementSystem.Register(ant)

	navMesh := setupNavMesh()
	renderSystem.Register(ant)
	renderSystem.Register(navMesh)
	setupGrass(assetManager, renderSystem)

	p := pathing.Planner{}
	p.SetNavMesh(navMesh)

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
