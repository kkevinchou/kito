package kito

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kkevinchou/kito/common/enums"
	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/entities/food"
	"github.com/kkevinchou/kito/entities/grass"
	"github.com/kkevinchou/kito/entities/worker"
	"github.com/kkevinchou/kito/lib"
	"github.com/kkevinchou/kito/lib/geometry"
	"github.com/kkevinchou/kito/lib/math/vector"
	"github.com/kkevinchou/kito/managers/item"
	"github.com/kkevinchou/kito/managers/path"
	"github.com/kkevinchou/kito/systems/movement"
	"github.com/kkevinchou/kito/systems/render"
)

const (
	gameUpdateDelta = 10 * time.Millisecond
)

var (
	fps                 = 60.0
	cameraStartPosition = vector.Vector3{X: 0, Y: 0, Z: 5}
	cameraStartView     = vector.Vector{X: 0, Y: 0}
)

func setupGrass() {
	grass.New(5, 0, 4)
	grass.New(2, 0, 2)
	grass.New(6, 0, 1)
	grass.New(6, 0, 7)
	grass.New(4, 0, 2)
	grass.New(0, 0, 0)
}

func (g *Game) setupSystems() *directory.Directory {
	itemManager := item.NewManager()
	pathManager := path.NewManager()
	assetManager := lib.NewAssetManager(nil, "_assets")
	renderSystem := render.NewRenderSystem(g, assetManager, g.camera)
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
	path      []geometry.Point
	worker    *worker.WorkerImpl
	pathIndex int
	gameOver  bool
	camera    *Camera
	gameMode  enums.GameMode
}

func NewGame() *Game {
	seed := int64(time.Now().Nanosecond())
	fmt.Println(fmt.Sprintf("Game Initializing with seed %d ...", seed))
	rand.Seed(seed)

	camera := NewCamera(cameraStartPosition, cameraStartView)
	fmt.Println("Camera initialized at position", camera.Position(), "and view", camera.View())

	g := &Game{
		camera:   camera,
		gameMode: enums.GameModePlaying,
	}

	g.camera.Position()

	g.setupSystems()
	setupGrass()
	food.New(0, 0, 0)
	g.worker = worker.New()
	g.worker.SetPosition(vector.Vector3{X: 19, Y: 12, Z: -10})
	// g.camera.Follow(g.worker)

	return g
}

// func (g *Game) MoveAnt(x, y float64) {
// 	position := g.worker.Position()
// 	pathManager := directory.GetDirectory().PathManager()
// 	g.path = pathManager.FindPath(
// 		geometry.Point{X: position.X, Y: position.Y},
// 		geometry.Point{X: x, Y: y},
// 	)
// 	if g.path != nil {
// 		g.pathIndex = 1
// 		g.worker.SetTarget(g.path[1].Vector3())
// 	}
// }

func (g *Game) PlaceFood(x, y float64) {
	food.New(x, 0, y)
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
	g.camera.Update(delta)
	movementSystem := directory.MovementSystem()
	movementSystem.Update(delta)
}

func (g *Game) Start(commandPoller CommandPoller) {
	rand.Seed(time.Now().Unix())

	previousTime := time.Now()
	var accumulator time.Duration
	var renderAccumulator time.Duration

	msPerFrame := time.Duration(1000000.0/fps) * time.Microsecond
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
			frameCount = 0
		}

		commands := commandPoller(g)
		for _, command := range commands {
			command.Execute(g)
		}

		accumulator += delta
		renderAccumulator += delta

		for accumulator >= gameUpdateDelta {
			g.update(gameUpdateDelta)
			accumulator -= gameUpdateDelta
		}
		if accumulator > 0 { // Temporary update to not lose physics time
			g.update(accumulator)
			accumulator = 0
		}

		if renderAccumulator >= msPerFrame {
			frameCount++
			renderSystem.Update(msPerFrame)
		}

		for renderAccumulator > msPerFrame {
			renderAccumulator -= msPerFrame
		}
	}
}

func (g *Game) CameraViewChange(v vector.Vector) {
	g.camera.ChangeView(v)
}

func (g *Game) SetCameraCommandHeading(v vector.Vector3) {
	g.camera.SetCommandHeading(v)
}

func (g *Game) GetGameMode() enums.GameMode {
	return g.gameMode
}
