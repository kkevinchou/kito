package kito

import (
	"fmt"
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/managers/player"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/kito/systems/ability"
	"github.com/kkevinchou/kito/kito/systems/ai"
	"github.com/kkevinchou/kito/kito/systems/animation"
	"github.com/kkevinchou/kito/kito/systems/bookkeeping"
	"github.com/kkevinchou/kito/kito/systems/charactercontroller"
	"github.com/kkevinchou/kito/kito/systems/collision"
	"github.com/kkevinchou/kito/kito/systems/combat"
	"github.com/kkevinchou/kito/kito/systems/loot"
	"github.com/kkevinchou/kito/kito/systems/networkdispatch"
	"github.com/kkevinchou/kito/kito/systems/networkupdate"
	"github.com/kkevinchou/kito/kito/systems/physics"
	"github.com/kkevinchou/kito/kito/systems/playerinput"
	"github.com/kkevinchou/kito/kito/systems/playerregistration"
	"github.com/kkevinchou/kito/kito/systems/preframe"
	"github.com/kkevinchou/kito/kito/systems/rpcreceiver"
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
	scene := entities.NewScene()

	enemies := []entities.Entity{}
	for i := 0; i < 5; i++ {
		enemy := entities.NewEnemy()
		x := rand.Intn(1000) - 500
		z := rand.Intn(1000) - 500
		enemy.GetComponentContainer().TransformComponent.Position = mgl64.Vec3{float64(x), 0, float64(z)}
		enemies = append(enemies, enemy)
	}

	lootbox := entities.NewLootbox()

	entities := []entities.Entity{
		scene,
		lootbox,
	}
	entities = append(entities, enemies...)
	return entities
}

func serverSystemSetup(g *Game, assetsDirectory string) {
	d := directory.GetDirectory()

	playerManager := player.NewPlayerManager(g)
	d.RegisterPlayerManager(playerManager)

	// asset manager is needed to load animation data. we don't load the meshes themselves to avoid
	// depending on OpenGL on the server
	assetManager := assets.NewAssetManager(assetsDirectory, false)
	d.RegisterAssetManager(assetManager)

	playerRegistrationSystem := playerregistration.NewPlayerRegistrationSystem(g, settings.ListenAddress, fmt.Sprintf("%d", settings.Port), settings.ConnectionType)
	networkDispatchSystem := networkdispatch.NewNetworkDispatchSystem(g)
	rpcReceiverSystem := rpcreceiver.NewRPCReceiverSystem(g)
	playerInputSystem := playerinput.NewPlayerInputSystem(g)
	aiSystem := ai.NewAnimationSystem(g)
	preframeSystem := preframe.NewPreFrameSystem(g)

	// systems that can manipulate the transform of an entity
	characterControllerSystem := charactercontroller.NewCharacterControllerSystem(g)
	physicsSystem := physics.NewPhysicsSystem(g)
	collisionSystem := collision.NewCollisionSystem(g)

	abilitySystem := ability.NewAbilitySystem(g)
	combatSystem := combat.NewCombatSystem(g)
	lootSystem := loot.NewLootSystem(g)
	animationSystem := animation.NewAnimationSystem(g)
	networkUpdateSystem := networkupdate.NewNetworkUpdateSystem(g)
	bookKeepingSystem := bookkeeping.NewBookKeepingSystem(g)

	g.systems = append(g.systems, []System{
		playerRegistrationSystem,
		networkDispatchSystem,
		rpcReceiverSystem,
		playerInputSystem,
		aiSystem,
		preframeSystem,
		characterControllerSystem,
		physicsSystem,
		collisionSystem,
		abilitySystem,
		combatSystem,
		lootSystem,
		animationSystem,
		bookKeepingSystem,
		networkUpdateSystem,
	}...)
}
