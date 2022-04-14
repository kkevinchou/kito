package kito

import (
	"fmt"

	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/managers/player"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/kito/systems/ability"
	"github.com/kkevinchou/kito/kito/systems/ai"
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
	"github.com/kkevinchou/kito/lib/assets"
)

func NewServerGame(assetsDirectory string) *Game {
	initSeed()
	settings.CurrentGameMode = settings.GameModeServer

	g := NewBaseGame()

	serverSystemSetup(g, assetsDirectory)
	initialEntities := serverEntitySetup(g)
	g.RegisterEntities(initialEntities)

	return g
}

func serverEntitySetup(g *Game) []entities.Entity {
	return []entities.Entity{
		entities.NewScene(),
		entities.NewEnemy(),
		// entities.NewSlime(mgl64.Vec3{-100, 0, -50}),
		// entities.NewStaticRigidBody(mgl64.Vec3{-5, 10, 0}),
		// entities.NewDynamicRigidBody(mgl64.Vec3{-5, 10, 0}),
		// entities.NewSlime(mgl64.Vec3{-50, 0, -50}),
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

	networkListenerSystem := networklistener.NewNetworkListenerSystem(g, "localhost", fmt.Sprintf("%d", settings.Port), settings.ConnectionType)
	networkDispatchSystem := networkdispatch.NewNetworkDispatchSystem(g)
	characterControllerSystem := charactercontroller.NewCharacterControllerSystem(g)
	abilitySystem := ability.NewAbilitySystem(g)
	physicsSystem := physics.NewPhysicsSystem(g)
	animationSystem := animation.NewAnimationSystem(g)
	cameraSystem := camera.NewCameraSystem(g)
	networkUpdateSystem := networkupdate.NewNetworkUpdateSystem(g)
	bookKeepingSystem := bookkeeping.NewBookKeepingSystem(g)
	playerInputSystem := playerinput.NewPlayerInputSystem(g)
	collisionSystem := collision.NewCollisionSystem(g)
	controllerResolverSystem := charactercontroller.NewCharacterControllerResolverSystem(g)
	collisionResolverSystem := collisionresolver.NewCollisionResolverSystem(g)
	aiSystem := ai.NewAnimationSystem(g)

	g.systems = append(g.systems, []System{
		networkListenerSystem,
		networkDispatchSystem,
		playerInputSystem,
		aiSystem,
		characterControllerSystem,
		abilitySystem,
		physicsSystem,
		collisionSystem,
		controllerResolverSystem,
		collisionResolverSystem,
		animationSystem,
		cameraSystem,
		networkUpdateSystem,
		bookKeepingSystem,
	}...)
}
