package kito

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/kkevinchou/kito/kito/commandframe"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/managers/eventbroker"
	"github.com/kkevinchou/kito/kito/managers/player"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/animation"
	camerasys "github.com/kkevinchou/kito/kito/systems/camera"
	"github.com/kkevinchou/kito/kito/systems/charactercontroller"
	"github.com/kkevinchou/kito/kito/systems/common"
	historysys "github.com/kkevinchou/kito/kito/systems/history"
	"github.com/kkevinchou/kito/kito/systems/networkdispatch"
	"github.com/kkevinchou/kito/kito/systems/networkinput"
	"github.com/kkevinchou/kito/kito/systems/physics"
	"github.com/kkevinchou/kito/kito/systems/render"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/lib/assets"
	"github.com/kkevinchou/kito/lib/network"
	"github.com/kkevinchou/kito/lib/shaders"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	width  = 1024
	height = 760
)

func NewClientGame(assetsDirectory string, shaderDirectory string) *Game {
	initSeed()
	settings.CurrentGameMode = settings.GameModeClient

	g := &Game{
		gameMode:            types.GameModePlaying,
		singleton:           singleton.NewSingleton(),
		entities:            map[int]entities.Entity{},
		eventBroker:         eventbroker.NewEventBroker(),
		commandFrameHistory: commandframe.NewCommandFrameHistory(),
	}

	clientSystemSetup(g, assetsDirectory, shaderDirectory)
	compileShaders()

	// Connect to server
	nClient, playerID, err := network.Connect(settings.RemoteHost, settings.Port, settings.ConnectionType)
	if err != nil {
		panic(err)
	}

	nClient.SetCommandFrameFunction(func() int { return g.CommandFrame() })

	var client types.NetworkClient = nClient
	if settings.ArtificialClientLatency > 0 {
		client = common.NewArtificallySlowClient(nClient, settings.ArtificialClientLatency)
	}

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

	window, err := initializeOpenGL(width, height)
	if err != nil {
		panic(err)
	}

	assetManager := assets.NewAssetManager(assetsDirectory, true)
	renderSystem := render.NewRenderSystem(g, window, width, height)

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
	historySystem := historysys.NewHistorySystem(g)

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
		historySystem,
		renderSystem,
	}...)
}

func initializeOpenGL(windowWidth, windowHeight int) (*sdl.Window, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, fmt.Errorf("Failed to init SDL %s", err)
	}

	// Enable hints for multisampling which allows opengl to use the default
	// multisampling algorithms implemented by the OpenGL rasterizer
	sdl.GLSetAttribute(sdl.GL_MULTISAMPLEBUFFERS, 1)
	sdl.GLSetAttribute(sdl.GL_MULTISAMPLESAMPLES, 4)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 4)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_FLAGS, sdl.GL_CONTEXT_FORWARD_COMPATIBLE_FLAG)

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(windowWidth), int32(windowHeight), sdl.WINDOW_OPENGL)
	if err != nil {
		return nil, fmt.Errorf("failed to create window %s", err)
	}

	_, err = window.GLCreateContext()
	if err != nil {
		return nil, fmt.Errorf("failed to create context %s", err)
	}

	if err := gl.Init(); err != nil {
		return nil, fmt.Errorf("failed to init OpenGL %s", err)
	}

	return window, nil
}
