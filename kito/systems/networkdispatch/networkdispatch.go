package networkdispatch

import (
	"time"

	"github.com/kkevinchou/kito/kito/commandframe"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/knetwork"
	"github.com/kkevinchou/kito/kito/managers/eventbroker"
	"github.com/kkevinchou/kito/kito/managers/player"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/lib/metrics"
	"github.com/kkevinchou/kito/lib/network"
)

type MessageFetcher func(world World) []*network.Message
type MessageHandler func(world World, message *network.Message)

type World interface {
	RegisterEntities([]entities.Entity)
	GetSingleton() *singleton.Singleton
	GetEventBroker() eventbroker.EventBroker
	GetCommandFrameHistory() *commandframe.CommandFrameHistory
	CommandFrame() int
	GetCamera() entities.Entity
	GetPlayerEntity() entities.Entity
	MetricsRegistry() *metrics.MetricsRegistry
	GetPlayer() *player.Player
	GetPlayerByID(id int) *player.Player
	QueryEntity(componentFlags int) []entities.Entity
	GetEntityByID(id int) entities.Entity
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

func (s *NetworkDispatchSystem) Update(delta time.Duration) {
	var latestGameStateUpdate *network.Message
	messages := s.messageFetcher(s.world)
	for _, message := range messages {
		if message.MessageType == knetwork.MessageTypeGameStateUpdate {
			latestGameStateUpdate = message
		}
	}

	var filteredMessages []*network.Message
	for _, message := range messages {
		// only take the latest gamestate update message
		if message.MessageType == knetwork.MessageTypeGameStateUpdate && message != latestGameStateUpdate {
			continue
		}
		filteredMessages = append(filteredMessages, message)
	}

	sawInputMessage := false
	for _, message := range filteredMessages {
		if utils.IsServer() {
			if message.MessageType == knetwork.MessageTypeInput {
				sawInputMessage = true
			}
		}
		s.messageHandler(s.world, message)
	}
	_ = sawInputMessage
	// if utils.IsServer() && !sawInputMessage {
	// 	fmt.Println("MISSED AN INPUT MESSAGE ON CF", s.world.CommandFrame())
	// }
}

func (s *NetworkDispatchSystem) Name() string {
	return "NetworkDispatchSystem"
}
