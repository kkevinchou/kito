package networkinput

import (
	"time"

	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/lib/network"
)

type World interface {
	GetSingleton() *singleton.Singleton
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
	playerManager := directory.GetDirectory().PlayerManager()

	player := playerManager.GetPlayer(singleton.PlayerID)
	playerInput := singleton.PlayerInput[player.ID]

	inputMessage := &network.InputMessage{
		CommandFrame: singleton.CommandFrame,
		Input:        playerInput,
	}

	// only send the input message if we detected new input
	if playerInput.NewInput {
		// fmt.Println("---------------")
		// fmt.Println("SENT MOVE ON CF", singleton.CommandFrame)
		// for _, e := range s.entities {
		// 	if e.GetID() == singleton.PlayerID {
		// 		t := e.GetComponentContainer().TransformComponent
		// 		fmt.Println(t.Position)
		// 	}
		// }

		player.Client.SendMessage(network.MessageTypeInput, inputMessage)
	}

}
