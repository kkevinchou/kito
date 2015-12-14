package main

import (
	"fmt"
	"time"

	"github.com/kkevinchou/ant/assets"
	"github.com/kkevinchou/ant/entity"
	"github.com/kkevinchou/ant/math/vector"
	"github.com/kkevinchou/ant/movement"
	"github.com/kkevinchou/ant/render"
	"github.com/veandco/go-sdl2/sdl"
)

func setupWindow() *sdl.Window {
	sdl.Init(sdl.INIT_EVERYTHING)

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	return window
}

func main() {
	window := setupWindow()
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(fmt.Sprintf("Failed to create renderer: %s\n", err))
	}
	defer renderer.Destroy()

	entity := entity.New()
	entity.SetTarget(vector.Vector{100, 0})

	movementSystem := movement.NewMovementSystem()
	movementSystem.Register(entity)

	assetManager := assets.NewAssetManager(renderer, "assets/icons")
	renderSystem := render.NewRenderSystem(renderer, assetManager)
	renderSystem.Register(entity)

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
			case *sdl.MouseMotionEvent:
				entity.SetTarget(vector.Vector{float64(e.X), float64(e.Y)})
			case *sdl.KeyUpEvent:
				if e.Keysym.Sym == sdl.K_ESCAPE {
					gameOver = true
				}
			}

		}

		movementSystem.Update(delta)
		renderSystem.Update(delta)
	}
	sdl.Quit()

	// node1 := pathing.CreateNode(0, 1)
	// node2 := pathing.CreateNode(1, 1)

	// nodes := []*pathing.Node{
	// 	node1,
	// 	node2,
	// }

	// edges := []*pathing.Edge{
	// 	pathing.CreateEdge(node1, node2),
	// }

	// planner := pathing.CreatePlanner(nodes, edges)
	// fmt.Println(planner)
	// fmt.Println(node1)
	// fmt.Println(node2)

	// surface, err := window.GetSurface()
	// if err != nil {
	// 	panic(err)
	// }

	// var i int32
	// for i = 0; i < 10; i++ {
	// 	surface.FillRect(&sdl.Rect{0, 0, 800, 600}, 0x0)
	// 	rect := sdl.Rect{i * 75, 0, 200, 200}
	// 	surface.FillRect(&rect, 0xffff0000)
	// 	window.UpdateSurface()
	// 	sdl.Delay(500)
	// }

}
