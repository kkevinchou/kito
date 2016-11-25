package ant

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kkevinchou/ant/directory"
	"github.com/kkevinchou/ant/entities/food"
	"github.com/kkevinchou/ant/entities/grass"
	"github.com/kkevinchou/ant/entities/worker"
	"github.com/kkevinchou/ant/lib"
	"github.com/kkevinchou/ant/lib/geometry"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/kkevinchou/ant/lib/pathing"
	"github.com/kkevinchou/ant/managers/item"
	"github.com/kkevinchou/ant/managers/path"
	"github.com/kkevinchou/ant/systems/movement"
	"github.com/kkevinchou/ant/systems/render"
)

const (
	gameUpdateDelta = 10 * time.Millisecond
)

var (
	fps = 60.0
)

func setupGrass() {
	grass.New(5, 0, 4)
	grass.New(2, 0, 2)
	grass.New(6, 0, 1)
	grass.New(6, 0, 7)
	grass.New(4, 0, 2)
	grass.New(0, 0, 0)
}

func setupSystems() *directory.Directory {
	itemManager := item.NewManager()
	pathManager := path.NewManager()
	assetManager := lib.NewAssetManager(nil, "_assets")
	renderSystem := render.NewRenderSystem(assetManager)
	movementSystem := movement.NewMovementSystem()

	d := directory.GetDirectory()
	d.RegisterRenderSystem(renderSystem)
	d.RegisterMovementSystem(movementSystem)
	d.RegisterAssetManager(assetManager)
	d.RegisterItemManager(itemManager)
	d.RegisterPathManager(pathManager)

	renderSystem.Register(pathManager.NavMesh())

	return d
}

type Game struct {
	path      []pathing.Node
	worker    *worker.WorkerImpl
	pathIndex int
	gameOver  bool
}

func NewGame() *Game {
	g := &Game{}
	rand.Seed(int64(time.Now().Nanosecond()))
	setupSystems()
	setupGrass()
	food.New(0, 0, 0)
	g.worker = worker.New()
	g.worker.SetPosition(vector.Vector3{X: 19, Y: 12, Z: -10})

	return g
}

func (g *Game) MoveAnt(x, y float64) {
	position := g.worker.Position()
	pathManager := directory.GetDirectory().PathManager()
	g.path = pathManager.FindPath(
		geometry.Point{X: position.X, Y: position.Y},
		geometry.Point{X: x, Y: y},
	)
	if g.path != nil {
		g.pathIndex = 1
		g.worker.SetTarget(g.path[1].Vector3())
	}
}

func (g *Game) PlaceFood(x, y float64) {
	food.New(x, 0, y)
}

func (g *Game) CameraViewChange(x, y int) {
	renderSystem := directory.GetDirectory().RenderSystem()
	renderSystem.CameraViewChange(x, y)
}

func (g *Game) MoveCamera(v vector.Vector3) {
	renderSystem := directory.GetDirectory().RenderSystem()
	renderSystem.MoveCamera(v)
}

func (g *Game) GameOver() {
	g.gameOver = true
}

func (g *Game) update(delta time.Duration) {
	if g.path != nil {
		if g.worker.Position().Sub(g.path[g.pathIndex].Vector3()).Length() <= 2 {
			g.pathIndex++
			if g.pathIndex == len(g.path) {
				g.path = nil
				g.worker.SetSeekActive(false)
				g.worker.SetVelocity(vector.Vector3{})
			} else {
				g.worker.SetTarget(g.path[g.pathIndex].Vector3())
			}
		}
	}

	directory := directory.GetDirectory()
	movementSystem := directory.MovementSystem()
	movementSystem.Update(delta)
}

func (g *Game) Start(commandPoller CommandPoller) {
	rand.Seed(time.Now().Unix())

	previousTime := time.Now()
	// var accumulator time.Duration
	// var renderAccumulator time.Duration

	// msPerFrame := time.Duration(1000.0/fps) * time.Millisecond
	directory := directory.GetDirectory()
	renderSystem := directory.RenderSystem()

	var fpsAccumulator time.Duration
	frameCount := 0

	for g.gameOver != true {
		now := time.Now()
		delta := time.Since(previousTime)
		if delta > 250*time.Millisecond {
			delta = 250 * time.Millisecond
		}
		previousTime = now

		fpsAccumulator += delta
		numWholeSeconds := 0
		for fpsAccumulator > time.Second {
			numWholeSeconds++
			fpsAccumulator -= time.Second
		}
		if numWholeSeconds > 0 {
			fmt.Println(frameCount)
			frameCount = 0
		}

		commandPoller(g)

		g.update(delta)
		renderSystem.Update(delta)

		// accumulator += delta
		// renderAccumulator += delta

		// for accumulator >= gameUpdateDelta {
		// 	fmt.Println("UPDATE")
		// 	g.update(delta)
		// 	accumulator -= gameUpdateDelta
		// }

		// if renderAccumulator >= msPerFrame {
		// 	fmt.Println("RENDER")
		// 	renderSystem.Update(delta)
		// }

		// for renderAccumulator > msPerFrame {
		// 	renderAccumulator -= msPerFrame
		// }

		frameCount++
	}
}

type CommandPoller func(game *Game)
