package networkdispatch

import (
	"time"

	"github.com/kkevinchou/kito/kito/commandframe"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/managers/eventbroker"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/utils"
)

type World interface {
	RegisterEntities([]entities.Entity)
	GetEntityByID(id int) (entities.Entity, error)
	GetSingleton() *singleton.Singleton
	GetEventBroker() eventbroker.EventBroker
	GetCommandFrameHistory() *commandframe.CommandFrameHistory
	CommandFrame() int
}

type NetworkDispatchSystem struct {
	*base.BaseSystem
	world          World
	messageFetcher MessageFetcher
	messageHandler MessageHandler
}

func NewNetworkDispatchSystem(world World) *NetworkDispatchSystem {
	networkDispatchSystem := &NetworkDispatchSystem{
		BaseSystem: base.NewBaseSystem(),
		world:      world,
	}

	if utils.IsClient() {
		networkDispatchSystem.messageFetcher = clientMessageFetcher
		networkDispatchSystem.messageHandler = clientMessageHandler
	} else if utils.IsServer() {
		networkDispatchSystem.messageFetcher = connectedPlayersMessageFetcher
		networkDispatchSystem.messageHandler = serverMessageHandler
	}

	return networkDispatchSystem
}

func (s *NetworkDispatchSystem) RegisterEntity(entity entities.Entity) {
}

func (s *NetworkDispatchSystem) Update(delta time.Duration) {
	for _, message := range s.messageFetcher(s.world) {
		s.messageHandler(s.world, message)
	}
}
