package networkdispatch

import (
	"time"

	"github.com/kkevinchou/kito/entities"
	"github.com/kkevinchou/kito/entities/singleton"
	"github.com/kkevinchou/kito/lib/input"
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
	SetCamera(camera entities.Entity)
}

type NetworkDispatchSystem struct {
	*base.BaseSystem
	world          World
	messageFetcher MessageFetcher
	messageHandler MessageHandler
}

func NewNetworkDispatchSystem(world World) *NetworkDispatchSystem {
	return &NetworkDispatchSystem{
		BaseSystem:     &base.BaseSystem{},
		world:          world,
		messageFetcher: defaultMessageFetcher,
		messageHandler: defaultMessageHandler,
	}
}

func (s *NetworkDispatchSystem) SetMessageFetcher(f MessageFetcher) {
	s.messageFetcher = f
}

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
