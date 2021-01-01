package kito

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/systems/animation"
	"github.com/kkevinchou/kito/systems/physics"

	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/entities/singleton"
	"github.com/kkevinchou/kito/lib"
	"github.com/kkevinchou/kito/lib/geometry"
	"github.com/kkevinchou/kito/managers/item"
	"github.com/kkevinchou/kito/managers/path"
	camerasys "github.com/kkevinchou/kito/systems/camera"
	"github.com/kkevinchou/kito/systems/charactercontroller"
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
	pathIndex      int
	gameOver       bool
	camera         entities.Entity
	gameMode       types.GameMode
	viewControlled bool

	singleton *singleton.Singleton
	systems   []System
	entities  map[int]entities.Entity
}

func NewGame() *Game {
	seed := int64(time.Now().Nanosecond())
	fmt.Println(fmt.Sprintf("Game Initializing with seed %d ...", seed))
	rand.Seed(seed)

	g := &Game{
		gameMode:  types.GameModePlaying,
		singleton: singleton.New(),
	}

	itemManager := item.NewManager()
	pathManager := path.NewManager()
	assetManager := lib.NewAssetManager(nil, "_assets")

	// System Setup

	renderSystem := render.NewRenderSystem(g, assetManager)
	cameraSystem := camerasys.NewCameraSystem(g)
	animationSystem := animation.NewAnimationSystem(g)
	physicsSystem := physics.NewPhysicsSystem(g)
	characterControllerSystem := charactercontroller.NewCharacterControllerSystem(g)

	d := directory.GetDirectory()
	d.RegisterRenderSystem(renderSystem)
	d.RegisterAssetManager(assetManager)
	d.RegisterItemManager(itemManager)
	d.RegisterPathManager(pathManager)

	g.systems = append(g.systems, characterControllerSystem)
	g.systems = append(g.systems, physicsSystem)
	g.systems = append(g.systems, cameraSystem)
	g.systems = append(g.systems, animationSystem)

	// Entity Setup

	bob := entities.NewBob()

	// camera := entities.NewCamera(cameraStartPosition, cameraStartView)
	camera := entities.NewThirdPersonCamera(bob.GetID(), cameraStartPosition, cameraStartView)
	componentContainer := camera.GetComponentContainer()
	// fmt.Println("Camera initialized at position", componentContainer.PositionComponent.Position, "and view", componentContainer.TopDownViewComponent.View())
	fmt.Println("Camera initialized at position", componentContainer.PositionComponent.Position)
	g.camera = camera
	renderSystem.SetCamera(camera)

	worldEntities := []entities.Entity{
		bob,
		camera,
		entities.NewBlock(),
	}

	g.entities = map[int]entities.Entity{}
	for _, entity := range worldEntities {
		g.entities[entity.GetID()] = entity
	}

	for _, entity := range worldEntities {
		animationSystem.RegisterEntity(entity)
		renderSystem.RegisterEntity(entity)
		physicsSystem.RegisterEntity(entity)
		characterControllerSystem.RegisterEntity(entity)
	}

	return g
}

func (g *Game) update(delta time.Duration) {
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

func (g *Game) GetCamera() entities.Entity {
	return g.camera
}

func (g *Game) GetSingleton() types.Singleton {
	return g.singleton
}

func (g *Game) GetEntityByID(id int) (entities.Entity, error) {
	if entity, ok := g.entities[id]; ok {
		return entity, nil
	}

	return nil, fmt.Errorf("failed to find entity with ID %d", id)
}
