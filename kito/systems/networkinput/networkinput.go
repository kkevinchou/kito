package networkinput

import (
	"time"

	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/knetwork"
	"github.com/kkevinchou/kito/kito/managers/player"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/lib/metrics"
)

type World interface {
	GetSingleton() *singleton.Singleton
	MetricsRegistry() *metrics.MetricsRegistry
	CommandFrame() int
	GetPlayer() *player.Player
}

type NetworkInputSystem struct {
	*base.BaseSystem
	world    World
	entities []entities.Entity
}

func NewNetworkInputSystem(world World) *NetworkInputSystem {
	return &NetworkInputSystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}
}

func (s *NetworkInputSystem) RegisterEntity(entity entities.Entity) {
	componentContainer := entity.GetComponentContainer()

	if componentContainer.PhysicsComponent != nil && componentContainer.TransformComponent != nil && componentContainer.ThirdPersonControllerComponent != nil {
		s.entities = append(s.entities, entity)
	}
}

func (s *NetworkInputSystem) Update(delta time.Duration) {
	singleton := s.world.GetSingleton()

	player := s.world.GetPlayer()
	playerInput := singleton.PlayerInput[player.ID]

	inputMessage := &knetwork.InputMessage{
		CommandFrame: singleton.CommandFrame,
		Input:        playerInput,
	}

	s.world.MetricsRegistry().Inc("newinput", 1)
	player.Client.SendMessage(knetwork.MessageTypeInput, inputMessage)
}
