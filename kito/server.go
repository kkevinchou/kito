package kito

import (
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/entities/singleton"
	"github.com/kkevinchou/kito/settings"
	"github.com/kkevinchou/kito/systems/animation"
	"github.com/kkevinchou/kito/systems/networklistener"
	"github.com/kkevinchou/kito/systems/physics"
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
	networkListener := networklistener.NewNetworkListenerSystem(g, settings.Host, settings.Port, settings.ConnectionType)
	physicsSystem := physics.NewPhysicsSystem(g)
	animationSystem := animation.NewAnimationSystem(g)

	g.systems = append(g.systems, networkListener)
	g.systems = append(g.systems, physicsSystem)
	g.systems = append(g.systems, animationSystem)
}
