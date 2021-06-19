package kito

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/entities/singleton"
	"github.com/kkevinchou/kito/managers/item"
	"github.com/kkevinchou/kito/managers/path"
	"github.com/kkevinchou/kito/systems/animation"
	"github.com/kkevinchou/kito/systems/networklistener"
	"github.com/kkevinchou/kito/systems/physics"
	"github.com/kkevinchou/kito/types"
)

func emptyRenderFunction(delta time.Duration) {

}

func NewServerGame(assetsDirectory string, shaderDirectory string) *Game {
	seed := int64(time.Now().Nanosecond())
	fmt.Println(fmt.Sprintf("Initializing server with seed %d ...", seed))
	rand.Seed(seed)

	g := &Game{
		gameMode:  types.GameModePlaying,
		singleton: singleton.New(),
	}

	itemManager := item.NewManager()
	pathManager := path.NewManager()

	// System Setup

	networkListener := networklistener.NewNetworkListenerSystem(g, "localhost", "8080", "tcp")
	physicsSystem := physics.NewPhysicsSystem(g)
	animationSystem := animation.NewAnimationSystem(g)

	d := directory.GetDirectory()
	d.RegisterItemManager(itemManager)
	d.RegisterPathManager(pathManager)

	g.systems = append(g.systems, networkListener)
	g.systems = append(g.systems, physicsSystem)
	g.systems = append(g.systems, animationSystem)

	worldEntities := []entities.Entity{
		entities.NewBlock(),
	}

	g.entities = map[int]entities.Entity{}
	for _, entity := range worldEntities {
		g.entities[entity.GetID()] = entity
	}

	for _, entity := range worldEntities {
		for _, system := range g.systems {
			system.RegisterEntity(entity)
		}
	}

	g.renderFunction = emptyRenderFunction

	return g
}
