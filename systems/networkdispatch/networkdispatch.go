package networkdispatch

import (
	"time"

	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/entities/singleton"
	"github.com/kkevinchou/kito/lib/input"
	"github.com/kkevinchou/kito/managers/eventbroker"
	"github.com/kkevinchou/kito/systems/base"
)

type InputBuffer struct {
	bufferSize int
	index      int
	buffer     []input.Input
}

type World interface {
	RegisterEntities([]entities.Entity)
	GetEntityByID(id int) (entities.Entity, error)
	GetSingleton() *singleton.Singleton
	GetEventBroker() eventbroker.EventBroker
}

type NetworkDispatchSystem struct {
	*base.BaseSystem
	world          World
	messageFetcher MessageFetcher
	messageHandler MessageHandler
}

func NewNetworkDispatchSystem(world World) *NetworkDispatchSystem {
	networkDispatchSystem := &NetworkDispatchSystem{
		BaseSystem:     base.NewBaseSystem(),
		world:          world,
		messageFetcher: defaultMessageFetcher,
		messageHandler: defaultMessageHandler,
	}
	eventBroker := world.GetEventBroker()
	eventBroker.AddObserver(networkDispatchSystem, []eventbroker.EventType{
		eventbroker.EventCreatePlayer,
	})

	return networkDispatchSystem
}

// todo: just check game mode and directly set the function
func (s *NetworkDispatchSystem) SetMessageFetcher(f MessageFetcher) {
	s.messageFetcher = f
}

// todo: just check game mode and directly set the function
func (s *NetworkDispatchSystem) SetMessageHandler(f MessageHandler) {
	s.messageHandler = f
}

func (s *NetworkDispatchSystem) RegisterEntity(entity entities.Entity) {
}

func (s *NetworkDispatchSystem) Update(delta time.Duration) {
	for _, message := range s.messageFetcher(s.world) {
		s.messageHandler(s.world, message)
	}
}
