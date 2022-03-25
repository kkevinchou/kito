package kito

import (
	"fmt"
	"runtime/debug"

	"github.com/kkevinchou/kito/kito/commandframe"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/managers/eventbroker"
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

func (g *Game) GetPlayer() entities.Entity {
	if utils.IsServer() {
		panic("invalid call to GetPlayer() as server")
	}

	if entity, ok := g.entities[g.singleton.PlayerID]; ok {
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