package kito

import (
	"encoding/json"
	"fmt"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/entities/singleton"
	"github.com/kkevinchou/kito/lib/assets"
	"github.com/kkevinchou/kito/lib/network"
	"github.com/kkevinchou/kito/lib/shaders"
	"github.com/kkevinchou/kito/settings"
	"github.com/kkevinchou/kito/systems/animation"
	camerasys "github.com/kkevinchou/kito/systems/camera"
	"github.com/kkevinchou/kito/systems/charactercontroller"
	"github.com/kkevinchou/kito/systems/networkinput"
	"github.com/kkevinchou/kito/systems/physics"
	"github.com/kkevinchou/kito/systems/render"
	"github.com/kkevinchou/kito/types"
)

func NewClientGame(assetsDirectory string, shaderDirectory string) *Game {
	initSeed()

	g := &Game{
		gameMode:  types.GameModePlaying,
		singleton: singleton.NewSingleton(),
	}

	clientSystemSetup(g, assetsDirectory, shaderDirectory)
	compileShaders()

	// Connect to server
	connectionComponent, err := components.NewConnectionComponent(
		settings.Host,
		settings.Port,
		settings.ConnectionType,
	)
	if err != nil {
		panic(err)
	}

	err = connectionComponent.Client.SendMessage(
		&network.Message{
			SenderID:    connectionComponent.PlayerID,
			MessageType: network.MessageTypeCreatePlayer,
		},
	)
	if err != nil {
		panic(err)
	}

	recvMessage := connectionComponent.Client.SyncReceiveMessage()
	var ack network.AckCreatePlayerMessage
	err = json.Unmarshal(recvMessage.Body, &ack)
	if err != nil {
		panic(err)
	}

	fmt.Println("successfully received ack player creation with id", ack.ID)
	fmt.Println(ack)

	g.singleton.Client = connectionComponent.Client
	g.singleton.PlayerID = connectionComponent.PlayerID

	initialEntities := clientEntitySetup(g, connectionComponent.PlayerID)
	g.RegisterEntities(initialEntities)

	return g
}

func clientEntitySetup(g *Game, playerID int) []entities.Entity {
	block := entities.NewBlock()

	bob := entities.NewBob(mgl64.Vec3{})
	bob.ID = playerID
	camera := entities.NewThirdPersonCamera(cameraStartPosition, cameraStartView, bob.GetID())
	cameraComponentContainer := camera.GetComponentContainer()
	fmt.Println("Camera initialized at position", cameraComponentContainer.TransformComponent.Position)

	bob.GetComponentContainer().ThirdPersonControllerComponent.CameraID = camera.GetID()

	g.SetCamera(camera)
	return []entities.Entity{camera, block, bob}
}

func clientSystemSetup(g *Game, assetsDirectory, shaderDirectory string) {
	renderSystem := render.NewRenderSystem(g)

	// TODO: asset manager creation has to happen after the render system is set up
	// because it depends on GL initializations
	assetManager := assets.NewAssetManager(assetsDirectory, true)
	renderSystem.SetAssetManager(assetManager)

	shaderManager := shaders.NewShaderManager(shaderDirectory)

	networkInputSystem := networkinput.NewNetworkInputSystem(g)
	cameraSystem := camerasys.NewCameraSystem(g)
	animationSystem := animation.NewAnimationSystem(g)
	physicsSystem := physics.NewPhysicsSystem(g)
	characterControllerSystem := charactercontroller.NewCharacterControllerSystem(g)

	d := directory.GetDirectory()
	d.RegisterRenderSystem(renderSystem)
	d.RegisterAssetManager(assetManager)
	d.RegisterShaderManager(shaderManager)

	g.systems = append(g.systems, []System{
		networkInputSystem,
		characterControllerSystem,
		physicsSystem,
		animationSystem,
		cameraSystem,
		renderSystem,
	}...)
}
