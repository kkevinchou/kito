package kito

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/managers/eventbroker"
	"github.com/kkevinchou/kito/kito/managers/player"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/animation"
	"github.com/kkevinchou/kito/kito/systems/bookkeeping"
	"github.com/kkevinchou/kito/kito/systems/camera"
	"github.com/kkevinchou/kito/kito/systems/charactercontroller"
	"github.com/kkevinchou/kito/kito/systems/collision"
	"github.com/kkevinchou/kito/kito/systems/collisionresolver"
	"github.com/kkevinchou/kito/kito/systems/networkdispatch"
	"github.com/kkevinchou/kito/kito/systems/networklistener"
	"github.com/kkevinchou/kito/kito/systems/networkupdate"
	"github.com/kkevinchou/kito/kito/systems/physics"
	"github.com/kkevinchou/kito/kito/systems/playerinput"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/lib/assets"
)

func NewServerGame(assetsDirectory string) *Game {
	initSeed()
	settings.CurrentGameMode = settings.GameModeServer

	g := &Game{
		gameMode:    types.GameModePlaying,
		singleton:   singleton.NewSingleton(),
		entities:    map[int]entities.Entity{},
		eventBroker: eventbroker.NewEventBroker(),
	}

	serverSystemSetup(g, assetsDirectory)
	initialEntities := serverEntitySetup(g)
	g.RegisterEntities(initialEntities)

	return g
}

func serverEntitySetup(g *Game) []entities.Entity {
	return []entities.Entity{
		entities.NewScene(mgl64.Vec3{}),
		entities.NewSlime(mgl64.Vec3{-50, 0, -50}),
	}
}

func serverSystemSetup(g *Game, assetsDirectory string) {
	d := directory.GetDirectory()

	playerManager := player.NewPlayerManager(g)
	d.RegisterPlayerManager(playerManager)

	// asset manager is needed to load animation data. we don't load the meshes themselves to avoid
	// depending on OpenGL on the server
	assetManager := assets.NewAssetManager(assetsDirectory, false)
	d.RegisterAssetManager(assetManager)

	networkListenerSystem := networklistener.NewNetworkListenerSystem(g, settings.Host, fmt.Sprintf("%d", settings.Port), settings.ConnectionType)
	networkDispatchSystem := networkdispatch.NewNetworkDispatchSystem(g)
	characterControllerSystem := charactercontroller.NewCharacterControllerSystem(g)
	physicsSystem := physics.NewPhysicsSystem(g)
	animationSystem := animation.NewAnimationSystem(g)
	cameraSystem := camera.NewCameraSystem(g)
	networkUpdateSystem := networkupdate.NewNetworkUpdateSystem(g)
	bookKeepingSystem := bookkeeping.NewBookKeepingSystem(g)
	playerInputSystem := playerinput.NewPlayerInputSystem(g)
	collisionSystem := collision.NewCollisionSystem(g)
	collisionResolverSystem := collisionresolver.NewCollisionResolverSystem(g)

	g.systems = append(g.systems, []System{
		networkListenerSystem,
		networkDispatchSystem,
		playerInputSystem,
		characterControllerSystem,
		physicsSystem,
		collisionSystem,
		collisionResolverSystem,
		animationSystem,
		cameraSystem,
		networkUpdateSystem,
		bookKeepingSystem,
	}...)
}
