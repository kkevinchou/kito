package kito

import (
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

func (g *Game) GetEntityByID(id int) entities.Entity {
	return g.entityManager.GetEntityByID(id)
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
	return g.GetEntityByID(player.EntityID)

}

func (g *Game) GetCamera() entities.Entity {
	return g.GetEntityByID(g.singleton.CameraID)
}

func (g *Game) RegisterEntities(entityList []entities.Entity) {
	for _, entity := range entityList {
		g.RegisterEntity(entity)
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

func (g *Game) RegisterEntity(e entities.Entity) {
	g.entityManager.RegisterEntity(e)
}

func (g *Game) QueryEntity(componentFlags int) []entities.Entity {
	return g.entityManager.Query(componentFlags)
}

func (g *Game) UnregisterEntity(entity entities.Entity) {
	g.entityManager.UnregisterEntity(entity)
}
