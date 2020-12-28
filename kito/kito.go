package kito

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/systems/animation"

	"github.com/kkevinchou/kito/directory"
	cameraEntity "github.com/kkevinchou/kito/entities/camera"
	"github.com/kkevinchou/kito/entities/food"
	"github.com/kkevinchou/kito/entities/grass"
	"github.com/kkevinchou/kito/entities/singleton"
	"github.com/kkevinchou/kito/entities/worker"
	"github.com/kkevinchou/kito/lib"
	"github.com/kkevinchou/kito/lib/geometry"
	"github.com/kkevinchou/kito/managers/item"
	"github.com/kkevinchou/kito/managers/path"
	"github.com/kkevinchou/kito/systems/camera"
	"github.com/kkevinchou/kito/systems/movement"
	"github.com/kkevinchou/kito/systems/render"
	"github.com/kkevinchou/kito/types"
)

const (
	fps                  float64 = 60
	simulationsPerSecond float64 = 60
	maxTimeStep          float64 = 250 // in milliseconds
)

var (
	cameraStartPosition = mgl64.Vec3{0, 10, 30}
	cameraStartView     = mgl64.Vec2{0, 0}
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
	camera         types.Camera
	gameMode       types.GameMode
	viewControlled bool

	singleton *singleton.Singleton
	systems   []System
}

func NewGame() *Game {
	seed := int64(time.Now().Nanosecond())
	fmt.Println(fmt.Sprintf("Game Initializing with seed %d ...", seed))
	rand.Seed(seed)

	camera := cameraEntity.New(cameraStartPosition, cameraStartView)
	fmt.Println("Camera initialized at position", camera.Position(), "and view", camera.View())

	g := &Game{
		camera:    camera,
		gameMode:  types.GameModePlaying,
		singleton: singleton.New(),
	}

	g.setupSystems()
	// setupGrass()
	// food.New(0, 0, 0)
	// g.worker = worker.New()
	// g.worker.SetPosition(mgl64.Vec3{X: 19, Y: 12, Z: -10})
	// g.camera.Follow(g.worker)

	return g
}

func (g *Game) update(delta time.Duration) {
	if g.path != nil {
		if g.worker.Position().Sub(g.path[g.pathIndex].MglVector3()).Len() <= 2 {
			g.pathIndex++
			if g.pathIndex == len(g.path) {
				g.path = nil
				g.worker.SetSeekActive(false)
				g.worker.SetVelocity(mgl64.Vec3{})
			} else {
				g.worker.SetTarget(g.path[g.pathIndex].MglVector3())
			}
		}
	}

	g.camera.Update(delta)
	for _, system := range g.systems {
		system.Update(delta)
	}
}

func (g *Game) Start(pollInputFunc InputPoller) {
	rand.Seed(time.Now().Unix())

	var accumulator float64
	var renderAccumulator float64

	msPerFrame := float64(1000) / fps
	msPerSimulation := float64(1000) / simulationsPerSecond
	directory := directory.GetDirectory()
	renderSystem := directory.RenderSystem()

	previousTimeStamp := float64(time.Now().UnixNano()) / 1000000
	for g.gameOver != true {
		now := float64(time.Now().UnixNano()) / 1000000
		delta := math.Min(now-previousTimeStamp, maxTimeStep)
		previousTimeStamp = now

		accumulator += delta
		renderAccumulator += delta

		for accumulator >= msPerSimulation {
			// input is handled once per simulation frame
			inputList := pollInputFunc()
			for _, input := range inputList {
				g.HandleInput(input)
			}
			g.update(time.Duration(msPerSimulation) * time.Millisecond)
			accumulator -= msPerSimulation
		}

		if renderAccumulator >= msPerFrame {
			renderSystem.Update(time.Duration(renderAccumulator) * time.Millisecond)
			for renderAccumulator >= msPerFrame {
				renderAccumulator -= msPerFrame
			}
		}
	}
}

func (g *Game) GetCamera() types.Camera {
	return g.camera
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

	renderSystem := render.NewRenderSystem(g, assetManager, g.camera)
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
