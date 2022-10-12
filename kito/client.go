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
	"github.com/kkevinchou/kito/kito/systems/ability"
	"github.com/kkevinchou/kito/kito/systems/animation"
	"github.com/kkevinchou/kito/kito/systems/bookkeeping"
	camerasys "github.com/kkevinchou/kito/kito/systems/camera"
	"github.com/kkevinchou/kito/kito/systems/charactercontroller"
	"github.com/kkevinchou/kito/kito/systems/collision"
	historysys "github.com/kkevinchou/kito/kito/systems/history"
	"github.com/kkevinchou/kito/kito/systems/networkdispatch"
	"github.com/kkevinchou/kito/kito/systems/networkinput"
	"github.com/kkevinchou/kito/kito/systems/physics"
	"github.com/kkevinchou/kito/kito/systems/ping"
	"github.com/kkevinchou/kito/kito/systems/preframe"
	"github.com/kkevinchou/kito/kito/systems/render"
	"github.com/kkevinchou/kito/kito/systems/rpcsender"
	"github.com/kkevinchou/kito/kito/systems/spawner"
	"github.com/kkevinchou/kito/kito/systems/stateinterpolator"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/kito/utils/entityutils"
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
	g.inputPollingFn = platform.PollInput
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
	cameraSystem := camerasys.NewCameraSystem(g)
	networkInputSystem := networkinput.NewNetworkInputSystem(g)
	networkDispatchSystem := networkdispatch.NewNetworkDispatchSystem(g)
	spawnerSystem := spawner.NewSpawnerSystem(g)
	stateInterpolatorSystem := stateinterpolator.NewStateInterpolatorSystem(g)
	preframeSystem := preframe.NewPreFrameSystem(g)

	// systems that can manipulate the transform of an entity
	characterControllerSystem := charactercontroller.NewCharacterControllerSystem(g)
	physicsSystem := physics.NewPhysicsSystem(g)
	collisionSystem := collision.NewCollisionSystem(g)

	abilitySystem := ability.NewAbilitySystem(g)
	animationSystem := animation.NewAnimationSystem(g)
	historySystem := historysys.NewHistorySystem(g)
	pingSystem := ping.NewPingSystem(g)
	rpcSenderSystem := rpcsender.NewRPCSenderSystem(g)
	bookKeepingSystem := bookkeeping.NewBookKeepingSystem(g)

	d.RegisterRenderSystem(renderSystem)
	d.RegisterAssetManager(assetManager)
	d.RegisterShaderManager(shaderManager)
	d.RegisterPlayerManager(playerManager)

	g.systems = append(g.systems, []System{
		cameraSystem,
		networkInputSystem,
		networkDispatchSystem,
		spawnerSystem,
		stateInterpolatorSystem,
		preframeSystem,
		characterControllerSystem,
		physicsSystem,
		collisionSystem,
		abilitySystem,
		animationSystem,
		historySystem,
		pingSystem,
		renderSystem,
		rpcSenderSystem,
		bookKeepingSystem,
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
	sdl.SetRelativeMouseMode(false)

	windowFlags := sdl.WINDOW_OPENGL
	if settings.Fullscreen {
		windowFlags |= sdl.WINDOW_FULLSCREEN
	}
	window, err := sdl.CreateWindow("KITO DEMO", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(windowWidth), int32(windowHeight), uint32(windowFlags))
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

	fmt.Println("Open GL Version:", gl.GoStr(gl.GetString(gl.VERSION)))

	return window, nil
}

func compileShaders() {
	d := directory.GetDirectory()
	shaderManager := d.ShaderManager()
	if err := shaderManager.CompileShaderProgram("skybox", "skybox", "skybox"); err != nil {
		panic(err)
	}
	if err := shaderManager.CompileShaderProgram("modelpbr", "model", "pbr"); err != nil {
		panic(err)
	}
	if err := shaderManager.CompileShaderProgram("model_debug", "model_debug", "pbr_debug"); err != nil {
		panic(err)
	}
	if err := shaderManager.CompileShaderProgram("model_static", "model_static", "pbr"); err != nil {
		panic(err)
	}
	if err := shaderManager.CompileShaderProgram("depthDebug", "basictexture", "depthvalue"); err != nil {
		panic(err)
	}
	if err := shaderManager.CompileShaderProgram("flat", "flat", "flat"); err != nil {
		panic(err)
	}
	if err := shaderManager.CompileShaderProgram("quadtex", "quadtex", "quadtex"); err != nil {
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
	singleton.PlayerID = messageBody.PlayerID
	singleton.CameraID = messageBody.CameraID

	bob := entities.NewBob()
	bob.ID = messageBody.EntityID

	playerManager := directory.GetDirectory().PlayerManager()
	playerManager.RegisterPlayer(messageBody.PlayerID, client)
	player := playerManager.GetPlayer(messageBody.PlayerID)
	player.EntityID = bob.ID

	camera := entities.NewThirdPersonCamera(settings.CameraStartPosition, settings.CameraStartView, player.ID, player.EntityID)
	camera.ID = messageBody.CameraID
	fmt.Println("set camera id", camera.ID)

	tpcComponent := bob.GetComponentContainer().ThirdPersonControllerComponent
	tpcComponent.CameraID = camera.GetID()

	initialEntities := []entities.Entity{bob, camera}
	for _, snapshot := range messageBody.Entities {
		entity := entityutils.SpawnWithID(snapshot.ID, types.EntityType(snapshot.Type), snapshot.Position, snapshot.Orientation)
		initialEntities = append(initialEntities, entity)
	}

	g.RegisterEntities(initialEntities)
}
