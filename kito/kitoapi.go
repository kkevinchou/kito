package kito

import (
	"fmt"
	"runtime/debug"

	"github.com/kkevinchou/kito/kito/commandframe"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/managers/eventbroker"
	"github.com/kkevinchou/kito/kito/managers/player"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/lib/metrics"
)

func (g *Game) GetSingleton() *singleton.Singleton {
	return g.singleton
}

func (g *Game) GetEntityByID(id int) (entities.Entity, error) {
	if entity, ok := g.entities[id]; ok {
		return entity, nil
	}

	stack := debug.Stack()

	return nil, fmt.Errorf("%sfailed to find entity with ID %d", string(stack), id)
}

func (g *Game) GetEntities() []entities.Entity {
	var result []entities.Entity

	for _, entity := range g.entities {
		result = append(result, entity)
	}
	return result
}

func (g *Game) GetPlayer() *player.Player {
	if utils.IsServer() {
		panic("invalid call to GetPlayer() as server")
	}

	d := directory.GetDirectory().PlayerManager()
	player := d.GetPlayer(g.singleton.PlayerID)

	return player
}

func (g *Game) GetPlayerByID(id int) *player.Player {
	d := directory.GetDirectory().PlayerManager()
	player := d.GetPlayer(id)
	return player
}

func (g *Game) GetPlayerEntity() entities.Entity {
	if utils.IsServer() {
		panic("invalid call to GetPlayer() as server")
	}
	player := g.GetPlayer()

	if entity, ok := g.entities[player.EntityID]; ok {
		return entity
	}
	return nil
}

func (g *Game) GetCamera() entities.Entity {
	if entity, ok := g.entities[g.singleton.CameraID]; ok {
		return entity
	}
	return nil
}

func (g *Game) RegisterEntities(entityList []entities.Entity) {
	for _, entity := range entityList {
		g.entities[entity.GetID()] = entity
	}

	for _, entity := range entityList {
		for _, system := range g.systems {
			system.RegisterEntity(entity)
		}
	}
}

func (g *Game) CommandFrame() int {
	return g.singleton.CommandFrame
}

func (g *Game) GetEventBroker() eventbroker.EventBroker {
	return g.eventBroker
}

func (g *Game) GetCommandFrameHistory() *commandframe.CommandFrameHistory {
	return g.commandFrameHistory
}

func (g *Game) MetricsRegistry() *metrics.MetricsRegistry {
	return g.metricsRegistry
}
