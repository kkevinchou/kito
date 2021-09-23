package kito

import (
	"fmt"

	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/lib/assets"
	"github.com/kkevinchou/kito/lib/network"
	"github.com/kkevinchou/kito/lib/shaders"
	"github.com/kkevinchou/kito/managers/eventbroker"
	"github.com/kkevinchou/kito/managers/player"
	"github.com/kkevinchou/kito/settings"
	"github.com/kkevinchou/kito/singleton"
	"github.com/kkevinchou/kito/systems/animation"
	camerasys "github.com/kkevinchou/kito/systems/camera"
	"github.com/kkevinchou/kito/systems/charactercontroller"
	"github.com/kkevinchou/kito/systems/networkdispatch"
	"github.com/kkevinchou/kito/systems/networkinput"
	"github.com/kkevinchou/kito/systems/physics"
	"github.com/kkevinchou/kito/systems/render"
	"github.com/kkevinchou/kito/types"
)

func NewClientGame(assetsDirectory string, shaderDirectory string) *Game {
	initSeed()
	settings.CurrentGameMode = settings.GameModeClient

	g := &Game{
		gameMode:    types.GameModePlaying,
		singleton:   singleton.NewSingleton(),
		entities:    map[int]entities.Entity{},
		eventBroker: eventbroker.NewEventBroker(),
	}

	clientSystemSetup(g, assetsDirectory, shaderDirectory)
	compileShaders()

	// Connect to server
	client, playerID, err := network.Connect(settings.RemoteHost, settings.Port, settings.ConnectionType)
	if err != nil {
		panic(err)
	}

	client.SetCommandFrameFunction(func() int { return g.CommandFrame() })

	err = client.SendMessage(network.MessageTypeCreatePlayer, nil)
	if err != nil {
		panic(err)
	}

	directory := directory.GetDirectory()
	directory.PlayerManager().RegisterPlayer(playerID, client)
	g.GetSingleton().PlayerID = playerID

	initialEntities := clientEntitySetup(g)
	g.RegisterEntities(initialEntities)

	fmt.Println("successfully received ack player creation with id", playerID)

	return g
}

func clientEntitySetup(g *Game) []entities.Entity {
	return []entities.Entity{}
}

func clientSystemSetup(g *Game, assetsDirectory, shaderDirectory string) {
	d := directory.GetDirectory()

	renderSystem := render.NewRenderSystem(g)

	// TODO: asset manager creation has to happen after the render system is set up
	// because it depends on GL initializations. Should probably decouple GL initializations from the rendering system
	assetManager := assets.NewAssetManager(assetsDirectory, true)
	renderSystem.SetAssetManager(assetManager)

	// Managers
	shaderManager := shaders.NewShaderManager(shaderDirectory)
	playerManager := player.NewPlayerManager(g)

	// Systems
	networkInputSystem := networkinput.NewNetworkInputSystem(g)
	networkDispatchSystem := networkdispatch.NewNetworkDispatchSystem(g)
	cameraSystem := camerasys.NewCameraSystem(g)
	animationSystem := animation.NewAnimationSystem(g)
	physicsSystem := physics.NewPhysicsSystem(g)
	characterControllerSystem := charactercontroller.NewCharacterControllerSystem(g)

	d.RegisterRenderSystem(renderSystem)
	d.RegisterAssetManager(assetManager)
	d.RegisterShaderManager(shaderManager)
	d.RegisterPlayerManager(playerManager)

	g.systems = append(g.systems, []System{
		networkInputSystem,
		networkDispatchSystem,
		characterControllerSystem,
		physicsSystem,
		animationSystem,
		cameraSystem,
		renderSystem,
	}...)
}
