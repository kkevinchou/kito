package networkdispatch

import (
	"time"

	"github.com/kkevinchou/kito/kito/commandframe"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/managers/eventbroker"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/lib/metrics"
	"github.com/kkevinchou/kito/lib/network"
)

type MessageFetcher func(world World) []*network.Message
type MessageHandler func(world World, message *network.Message)
type MessageHandlerInit func(world World)

type World interface {
	RegisterEntities([]entities.Entity)
	GetEntityByID(id int) (entities.Entity, error)
	GetSingleton() *singleton.Singleton
	GetEventBroker() eventbroker.EventBroker
	GetCommandFrameHistory() *commandframe.CommandFrameHistory
	CommandFrame() int
	GetCamera() entities.Entity
	GetPlayer() entities.Entity
	MetricsRegistry() *metrics.MetricsRegistry
}

type NetworkDispatchSystem struct {
	*base.BaseSystem
	world              World
	messageFetcher     MessageFetcher
	messageHandler     MessageHandler
	messageHandlerInit MessageHandlerInit
}

func NewNetworkDispatchSystem(world World) *NetworkDispatchSystem {
	networkDispatchSystem := &NetworkDispatchSystem{
		BaseSystem: base.NewBaseSystem(),
		world:      world,
	}

	if utils.IsClient() {
		networkDispatchSystem.messageFetcher = clientMessageFetcher
		networkDispatchSystem.messageHandler = clientMessageHandler
		networkDispatchSystem.messageHandlerInit = func(world World) { return }
	} else if utils.IsServer() {
		networkDispatchSystem.messageFetcher = connectedPlayersMessageFetcher
		networkDispatchSystem.messageHandler = serverMessageHandler
		networkDispatchSystem.messageHandlerInit = serverMessageHandlerInit
	}

	return networkDispatchSystem
}

func (s *NetworkDispatchSystem) RegisterEntity(entity entities.Entity) {
}

func (s *NetworkDispatchSystem) Update(delta time.Duration) {
	s.messageHandlerInit(s.world)
	for _, message := range s.messageFetcher(s.world) {
		s.messageHandler(s.world, message)
	}
}
