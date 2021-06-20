package kito

import (
	"fmt"
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/entities/singleton"
	"github.com/kkevinchou/kito/lib/assets"
	"github.com/kkevinchou/kito/lib/shaders"
	"github.com/kkevinchou/kito/managers/item"
	"github.com/kkevinchou/kito/managers/path"
	"github.com/kkevinchou/kito/settings"
	"github.com/kkevinchou/kito/systems/animation"
	camerasys "github.com/kkevinchou/kito/systems/camera"
	"github.com/kkevinchou/kito/systems/charactercontroller"
	"github.com/kkevinchou/kito/systems/physics"
	"github.com/kkevinchou/kito/systems/render"
	"github.com/kkevinchou/kito/types"
)

func NewClientGame(assetsDirectory string, shaderDirectory string) *Game {
	seed := settings.Seed
	fmt.Printf("Initializing client game with seed %d ...\n", seed)
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
	camera := entities.NewThirdPersonCamera(cameraStartPosition, cameraStartView, bob.GetID())

	cameraComponentContainer := camera.GetComponentContainer()
	fmt.Println("Camera initialized at position", cameraComponentContainer.TransformComponent.Position)

	bobComponentContainer := bob.GetComponentContainer()
	if bobComponentContainer.ThirdPersonControllerComponent != nil {
		bobComponentContainer.ThirdPersonControllerComponent.CameraID = camera.GetID()
	}

	g.camera = camera

	worldEntities := []entities.Entity{
		bob,
	}

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
		for _, system := range g.systems {
			system.RegisterEntity(entity)
		}
		renderSystem.RegisterEntity(entity)
	}

	g.renderFunction = renderSystem.Update

	return g
}
