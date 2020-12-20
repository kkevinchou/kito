package kito

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kkevinchou/kito/systems/animation"

	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/entities/food"
	"github.com/kkevinchou/kito/entities/grass"
	"github.com/kkevinchou/kito/entities/singleton"
	"github.com/kkevinchou/kito/entities/viewer"
	"github.com/kkevinchou/kito/entities/worker"
	"github.com/kkevinchou/kito/lib"
	"github.com/kkevinchou/kito/lib/geometry"
	"github.com/kkevinchou/kito/lib/math/vector"
	"github.com/kkevinchou/kito/managers/item"
	"github.com/kkevinchou/kito/managers/path"
	"github.com/kkevinchou/kito/systems/camera"
	"github.com/kkevinchou/kito/systems/movement"
	"github.com/kkevinchou/kito/systems/render"
	"github.com/kkevinchou/kito/types"
)

const (
	gameUpdateDelta = 10 * time.Millisecond
)

var (
	fps                 = 60.0
	viewerStartPosition = vector.Vector3{X: 0, Y: 40, Z: 40}
	viewerStartView     = vector.Vector{X: 10, Y: 0}
)

type System interface {
	Update(delta time.Duration)
}

type Input interface{}
type InputPoller func() []Input

type Game struct {
	path           []geometry.Point
	worker         *worker.WorkerImpl
	pathIndex      int
	gameOver       bool
	viewer         types.Viewer
	gameMode       types.GameMode
	viewControlled bool

	singleton *singleton.Singleton
	systems   []System
}

func NewGame() *Game {
	seed := int64(time.Now().Nanosecond())
	fmt.Println(fmt.Sprintf("Game Initializing with seed %d ...", seed))
	rand.Seed(seed)

	viewer := viewer.New(viewerStartPosition, viewerStartView)
	fmt.Println("Viewer initialized at position", viewer.Position(), "and view", viewer.View())

	g := &Game{
		viewer:    viewer,
		gameMode:  types.GameModePlaying,
		singleton: singleton.New(),
	}

	g.setupSystems()
	// setupGrass()
	// food.New(0, 0, 0)
	// g.worker = worker.New()
	// g.worker.SetPosition(vector.Vector3{X: 19, Y: 12, Z: -10})
	// g.camera.Follow(g.worker)

	return g
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

	g.viewer.Update(delta)
	for _, system := range g.systems {
		system.Update(delta)
	}
}

func (g *Game) Start(pollInputFunc InputPoller) {
	rand.Seed(time.Now().Unix())

	previousTime := time.Now()
	var accumulator time.Duration
	var renderAccumulator time.Duration
	var debugAccumulator time.Duration

	msPerFrame := time.Duration(1000000.0/fps) * time.Microsecond
	directory := directory.GetDirectory()
	renderSystem := directory.RenderSystem()

	for g.gameOver != true {
		now := time.Now()
		delta := time.Since(previousTime)
		if delta > 250*time.Millisecond {
			delta = 250 * time.Millisecond
		}
		previousTime = now

		accumulator += delta
		renderAccumulator += delta
		debugAccumulator += delta

		if debugAccumulator > time.Duration(1*time.Second) {
			// fmt.Println("LOOP START")
		}

		for accumulator >= gameUpdateDelta {
			// input is handled once per frame
			inputList := pollInputFunc()
			for _, input := range inputList {
				g.HandleInput(input)
			}
			g.update(gameUpdateDelta)
			accumulator -= gameUpdateDelta
		}

		// Temporary update to not lose physics time, is this needed? was in a weird
		// case where the game updates weren't running since we would always set accumulation to zero
		// if accumulator > 0 {
		// 	g.update(accumulator)
		// 	accumulator = 0
		// }

		if renderAccumulator >= msPerFrame {
			renderSystem.Update(msPerFrame)
		}
		for renderAccumulator > msPerFrame {
			renderAccumulator -= msPerFrame
		}

		if debugAccumulator > time.Duration(1*time.Second) {
			// fmt.Println("LOOP END")
			debugAccumulator = 0
		}
	}
}

func (g *Game) GetCamera() types.Viewer {
	return g.viewer
}

func (g *Game) GetSingleton() types.Singleton {
	return g.singleton
}

func (g *Game) PlaceFood(x, y float64) {
	food.New(x, 0, y)
}

func (g *Game) setupSystems() *directory.Directory {
	itemManager := item.NewManager()
	pathManager := path.NewManager()
	assetManager := lib.NewAssetManager(nil, "_assets")

	renderSystem := render.NewRenderSystem(g, assetManager, g.viewer)
	movementSystem := movement.NewMovementSystem()
	cameraSystem := camera.NewCameraSystem(g)
	animationSystem := animation.NewAnimationSystem(g)

	d := directory.GetDirectory()
	d.RegisterRenderSystem(renderSystem)
	d.RegisterMovementSystem(movementSystem)
	d.RegisterAssetManager(assetManager)
	d.RegisterItemManager(itemManager)
	d.RegisterPathManager(pathManager)

	renderSystem.Register(pathManager.NavMesh())

	g.systems = append(g.systems, cameraSystem)
	g.systems = append(g.systems, movementSystem)
	g.systems = append(g.systems, animationSystem)

	return d
}

func setupGrass() {
	grass.New(5, 0, 4)
	grass.New(2, 0, 2)
	grass.New(6, 0, 1)
	grass.New(6, 0, 7)
	grass.New(4, 0, 2)
	// grass.New(0, 0, 0)
}
