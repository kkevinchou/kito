package kito

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/inkyblackness/imgui-go/v4"
	"github.com/kkevinchou/kito/kito/commandframe"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/knetwork"
	"github.com/kkevinchou/kito/kito/managers/player"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/kito/systems/animation"
	camerasys "github.com/kkevinchou/kito/kito/systems/camera"
	"github.com/kkevinchou/kito/kito/systems/charactercontroller"
	"github.com/kkevinchou/kito/kito/systems/collision"
	"github.com/kkevinchou/kito/kito/systems/collisionresolver"
	historysys "github.com/kkevinchou/kito/kito/systems/history"
	"github.com/kkevinchou/kito/kito/systems/networkdispatch"
	"github.com/kkevinchou/kito/kito/systems/networkinput"
	"github.com/kkevinchou/kito/kito/systems/physics"
	"github.com/kkevinchou/kito/kito/systems/ping"
	"github.com/kkevinchou/kito/kito/systems/render"
	"github.com/kkevinchou/kito/kito/systems/spawner"
	"github.com/kkevinchou/kito/kito/systems/stateinterpolator"
	"github.com/kkevinchou/kito/lib/assets"
	"github.com/kkevinchou/kito/lib/input"
	"github.com/kkevinchou/kito/lib/network"
	"github.com/kkevinchou/kito/lib/shaders"
	"github.com/veandco/go-sdl2/sdl"
)

type Platform interface {
	NewFrame()
	DisplaySize() [2]float32
	FramebufferSize() [2]float32
}

func NewClientGame(assetsDirectory string, shaderDirectory string) *Game {
	initSeed()
	settings.CurrentGameMode = settings.GameModeClient

	window, err := initializeOpenGL(settings.Width, settings.Height)
	if err != nil {
		panic(err)
	}
	imgui.CreateContext(nil)
	imguiIO := imgui.CurrentIO()
	platform := input.NewSDLPlatform(window, imguiIO)

	g := NewBaseGame()
	g.inputPollingFn = platform.PollInput3
	g.commandFrameHistory = commandframe.NewCommandFrameHistory()

	// Connect to server
	client, _, err := network.Connect(settings.Host, fmt.Sprintf("%d", settings.Port), settings.ConnectionType)
	if err != nil {
		panic(err)
	}
	client.SetCommandFrameFunction(func() int { return g.CommandFrame() })

	err = client.SendMessage(knetwork.MessageTypeCreatePlayer, nil)
	if err != nil {
		panic(err)
	}

	clientSystemSetup(g, window, imguiIO, platform, assetsDirectory, shaderDirectory)
	ackCreatePlayer(g, client)

	initialEntities := clientEntitySetup(g)
	g.RegisterEntities(initialEntities)

	compileShaders()

	return g
}

func clientEntitySetup(g *Game) []entities.Entity {
	return []entities.Entity{}
}

func clientSystemSetup(g *Game, window *sdl.Window, imguiIO imgui.IO, platform Platform, assetsDirectory, shaderDirectory string) {
	d := directory.GetDirectory()

	assetManager := assets.NewAssetManager(assetsDirectory, true)
	renderSystem := render.NewRenderSystem(g, window, platform, imguiIO, settings.Width, settings.Height)

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
	stateInterpolatorSystem := stateinterpolator.NewStateInterpolatorSystem(g)
	spawnerSystem := spawner.NewSpawnerSystem(g)
	pingSystem := ping.NewPingSystem(g)
	collisionSystem := collision.NewCollisionSystem(g)
	controllerResolverSystem := charactercontroller.NewCharacterControllerResolverSystem(g)
	collisionResolverSystem := collisionresolver.NewCollisionResolverSystem(g)

	d.RegisterRenderSystem(renderSystem)
	d.RegisterAssetManager(assetManager)
	d.RegisterShaderManager(shaderManager)
	d.RegisterPlayerManager(playerManager)

	g.systems = append(g.systems, []System{
		networkInputSystem,
		networkDispatchSystem,
		characterControllerSystem,
		spawnerSystem,
		stateInterpolatorSystem,
		physicsSystem,
		collisionSystem,
		controllerResolverSystem,
		collisionResolverSystem,
		animationSystem,
		cameraSystem,
		historySystem,
		pingSystem,
		renderSystem,
	}...)
}

func initializeOpenGL(windowWidth, windowHeight int) (*sdl.Window, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, fmt.Errorf("failed to init SDL %s", err)
	}

	// Enable hints for multisampling which allows opengl to use the default
	// multisampling algorithms implemented by the OpenGL rasterizer
	sdl.GLSetAttribute(sdl.GL_MULTISAMPLEBUFFERS, 1)
	sdl.GLSetAttribute(sdl.GL_MULTISAMPLESAMPLES, 4)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 4)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_FLAGS, sdl.GL_CONTEXT_FORWARD_COMPATIBLE_FLAG)

	// sdl.ShowCursor(sdl.DISABLE)
	// sdl.SetRelativeMouseMode(true)
	// window, err := sdl.CreateWindow("KITO DEMO", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(windowWidth), int32(windowHeight), sdl.WINDOW_OPENGL|sdl.WINDOW_FULLSCREEN)
	window, err := sdl.CreateWindow("KITO DEMO", 400, sdl.WINDOWPOS_UNDEFINED, int32(windowWidth), int32(windowHeight), sdl.WINDOW_OPENGL)
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

func compileShaders() {
	d := directory.GetDirectory()
	shaderManager := d.ShaderManager()
	if err := shaderManager.CompileShaderProgram("basic", "basic", "basic"); err != nil {
		panic(err)
	}
	if err := shaderManager.CompileShaderProgram("basicShadow", "basicshadow", "basicshadow"); err != nil {
		panic(err)
	}
	if err := shaderManager.CompileShaderProgram("skybox", "skybox", "skybox"); err != nil {
		panic(err)
	}
	if err := shaderManager.CompileShaderProgram("model", "model", "model"); err != nil {
		panic(err)
	}
	if err := shaderManager.CompileShaderProgram("model_static", "model_static", "model"); err != nil {
		panic(err)
	}
	if err := shaderManager.CompileShaderProgram("depthDebug", "basictexture", "depthvalue"); err != nil {
		panic(err)
	}
	if err := shaderManager.CompileShaderProgram("material", "model", "material"); err != nil {
		panic(err)
	}
	if err := shaderManager.CompileShaderProgram("quadtex", "quadtex", "quadtex"); err != nil {
		panic(err)
	}
	if err := shaderManager.CompileShaderProgram("basicsolid", "basicsolid", "basicsolid"); err != nil {
		panic(err)
	}
}

func ackCreatePlayer(g *Game, client *network.Client) {
	var messageBody *knetwork.AckCreatePlayerMessage
	for messageBody == nil {
		message := client.SyncReceiveMessage()
		// discard any messages that are not for acking the create player
		if message.MessageType != network.MessageTypeAckCreatePlayer {
			fmt.Printf("during ack create player, discarded message %s\n", string(message.Body))
			continue
		}

		messageBody = &knetwork.AckCreatePlayerMessage{}
		err := network.DeserializeBody(message, messageBody)
		if err != nil {
			fmt.Println(err)
			return
		}
		break
	}

	singleton := g.GetSingleton()
	singleton.PlayerID = messageBody.ID
	singleton.CameraID = messageBody.CameraID

	directory := directory.GetDirectory()
	directory.PlayerManager().RegisterPlayer(messageBody.ID, client)

	bob := entities.NewBob()
	bob.ID = messageBody.ID

	camera := entities.NewThirdPersonCamera(settings.CameraStartPosition, settings.CameraStartView, bob.GetID())
	camera.ID = messageBody.CameraID

	bob.GetComponentContainer().ThirdPersonControllerComponent.CameraID = camera.GetID()

	g.RegisterEntities([]entities.Entity{
		bob,
		camera,
	})
}
