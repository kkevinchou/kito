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
	"github.com/kkevinchou/kito/lib/assets"
	"github.com/kkevinchou/kito/lib/geometry"
	"github.com/kkevinchou/kito/lib/shaders"
	"github.com/kkevinchou/kito/managers/item"
	"github.com/kkevinchou/kito/managers/path"
	camerasys "github.com/kkevinchou/kito/systems/camera"
	"github.com/kkevinchou/kito/systems/charactercontroller"
	"github.com/kkevinchou/kito/systems/render"
	"github.com/kkevinchou/kito/types"
)

const (
	fps               float64 = 60
	msPerCommandFrame float64 = 16
	maxTimeStep       float64 = 250 // in milliseconds
)

var (
	cameraStartPosition = mgl64.Vec3{0, 10, 30}
	cameraStartView     = mgl64.Vec2{0, 0}
)

type System interface {
	Update(delta time.Duration)
}

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

func NewGame(assetsDirectory string, shaderDirectory string) *Game {
	seed := int64(time.Now().Nanosecond())
	fmt.Println(fmt.Sprintf("Game Initializing with seed %d ...", seed))
	rand.Seed(seed)

	g := &Game{
		gameMode:  types.GameModePlaying,
		singleton: singleton.New(),
	}

	itemManager := item.NewManager()
	pathManager := path.NewManager()

	// System Setup

	renderSystem := render.NewRenderSystem(g)

	// TODO: asset manager creation has to happen after the render system is set up
	// because it depends on GL initializations
	assetManager := assets.NewAssetManager(assetsDirectory)
	renderSystem.SetAssetManager(assetManager)

	shaderManager := shaders.NewShaderManager(shaderDirectory)

	cameraSystem := camerasys.NewCameraSystem(g)
	animationSystem := animation.NewAnimationSystem(g)
	physicsSystem := physics.NewPhysicsSystem(g)
	characterControllerSystem := charactercontroller.NewCharacterControllerSystem(g)

	d := directory.GetDirectory()
	d.RegisterRenderSystem(renderSystem)
	d.RegisterAssetManager(assetManager)
	d.RegisterShaderManager(shaderManager)
	d.RegisterItemManager(itemManager)
	d.RegisterPathManager(pathManager)

	g.systems = append(g.systems, characterControllerSystem)
	g.systems = append(g.systems, physicsSystem)
	g.systems = append(g.systems, animationSystem)
	g.systems = append(g.systems, cameraSystem)

	// Shader Compilation

	shaderManager.CompileShaderProgram("basic", "basic", "basic")
	shaderManager.CompileShaderProgram("basicShadow", "basicshadow", "basicshadow")
	shaderManager.CompileShaderProgram("skybox", "skybox", "skybox")
	shaderManager.CompileShaderProgram("model", "model", "model")
	shaderManager.CompileShaderProgram("depth", "depth", "depth")
	shaderManager.CompileShaderProgram("depthDebug", "basictexture", "depthvalue")

	// Entity Setup

	bob := entities.NewBob(mgl64.Vec3{0, 15, 0})
	camera := entities.NewThirdPersonCamera(cameraStartPosition, cameraStartView)

	cameraComponentContainer := camera.GetComponentContainer()
	cameraComponentContainer.FollowComponent.FollowTargetEntityID = bob.GetID()
	fmt.Println("Camera initialized at position", cameraComponentContainer.TransformComponent.Position)

	bobComponentContainer := bob.GetComponentContainer()
	if bobComponentContainer.ThirdPersonControllerComponent != nil {
		bobComponentContainer.ThirdPersonControllerComponent.CameraID = camera.GetID()
	}

	// offset := 5
	// bobDimension := 20
	// var bobs []entities.Entity
	// for i := 0; i < bobDimension; i++ {
	// 	for j := 0; j < bobDimension; j++ {
	// 		x := i - bobDimension/2
	// 		z := j - bobDimension/2
	// 		bobGuy := entities.NewBob(mgl64.Vec3{float64(x * offset), 0, float64(z * offset)})
	// 		container := bobGuy.GetComponentContainer()
	// 		container.ThirdPersonControllerComponent.CameraID = camera.GetID()
	// 		bobs = append(bobs, bobGuy)
	// 	}
	// }

	g.camera = camera

	worldEntities := []entities.Entity{
		bob,
	}
	// worldEntities = append(worldEntities, bobs...)
	worldEntities = append(
		worldEntities,
		camera,
		entities.NewBlock(),
	)

	g.entities = map[int]entities.Entity{}
	for _, entity := range worldEntities {
		g.entities[entity.GetID()] = entity
	}

	for _, entity := range worldEntities {
		characterControllerSystem.RegisterEntity(entity)
		physicsSystem.RegisterEntity(entity)
		animationSystem.RegisterEntity(entity) // animation system should render at the same rate as the render system
		renderSystem.RegisterEntity(entity)
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
	directory := directory.GetDirectory()
	renderSystem := directory.RenderSystem()

	var fpsAccumulator float64

	previousTimeStamp := float64(time.Now().UnixNano()) / 1000000
	frameCount := 0
	for g.gameOver != true {
		now := float64(time.Now().UnixNano()) / 1000000
		delta := math.Min(now-previousTimeStamp, maxTimeStep)
		previousTimeStamp = now

		accumulator += delta
		renderAccumulator += delta

		for accumulator >= msPerCommandFrame {
			// input is handled once per command frame
			inputList := pollInputFunc()
			for _, input := range inputList {
				g.HandleInput(input)
			}
			g.update(time.Duration(msPerCommandFrame) * time.Millisecond)
			accumulator -= msPerCommandFrame
		}

		if renderAccumulator >= msPerFrame {
			frameCount++
			renderSystem.Update(time.Duration(renderAccumulator) * time.Millisecond)
			for renderAccumulator >= msPerFrame {
				renderAccumulator -= msPerFrame
			}
		}

		fpsAccumulator += delta
		if fpsAccumulator > 1000 {
			// fmt.Println(fmt.Sprintf("%d frames rendered last second", frameCount))
			frameCount = 0
			fpsAccumulator -= 1000
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
