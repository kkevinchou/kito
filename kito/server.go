package kito

import (
	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/entities/singleton"
	"github.com/kkevinchou/kito/lib/assets"
	"github.com/kkevinchou/kito/managers/player"
	"github.com/kkevinchou/kito/settings"
	"github.com/kkevinchou/kito/systems/animation"
	"github.com/kkevinchou/kito/systems/networkdispatch"
	"github.com/kkevinchou/kito/systems/networklistener"
	"github.com/kkevinchou/kito/systems/networkupdate"
	"github.com/kkevinchou/kito/systems/physics"
	"github.com/kkevinchou/kito/systems/servercamera"
	"github.com/kkevinchou/kito/systems/servercharactercontroller"
	"github.com/kkevinchou/kito/types"
)

func NewServerGame(assetsDirectory string) *Game {
	initSeed()

	g := &Game{
		gameMode:  types.GameModePlaying,
		singleton: singleton.NewSingleton(),
	}

	serverSystemSetup(g, assetsDirectory)
	initialEntities := serverEntitySetup(g)
	g.RegisterEntities(initialEntities)

	return g
}

func serverEntitySetup(g *Game) []entities.Entity {
	return []entities.Entity{
		entities.NewBlock(),
	}
}

func serverSystemSetup(g *Game, assetsDirectory string) {
	d := directory.GetDirectory()

	playerManager := player.NewPlayerManager()
	d.RegisterPlayerManager(playerManager)

	// asset manager is needed to load animation data. we don't load the meshes themselves to avoid
	// depending on OpenGL on the server
	assetManager := assets.NewAssetManager(assetsDirectory, false)
	d.RegisterAssetManager(assetManager)

	networkListenerSystem := networklistener.NewNetworkListenerSystem(g, settings.Host, settings.Port, settings.ConnectionType)
	networkDispatchSystem := networkdispatch.NewNetworkDispatchSystem(g)
	serverCharacterControllerSystem := servercharactercontroller.NewServerCharacterControllerSystem(g)
	physicsSystem := physics.NewPhysicsSystem(g)
	animationSystem := animation.NewAnimationSystem(g)
	serverCameraSystem := servercamera.NewServerCameraSystem(g)
	networkUpdateSystem := networkupdate.NewNetworkUpdateSystem(g)

	g.systems = append(g.systems, []System{
		networkListenerSystem,
		networkDispatchSystem,
		serverCharacterControllerSystem,
		physicsSystem,
		animationSystem,
		serverCameraSystem,
		networkUpdateSystem,
	}...)
}
