package stateinterpolator

import (
	"fmt"
	"time"

	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/statebuffer"
	"github.com/kkevinchou/kito/kito/systems/base"
)

type World interface {
	CommandFrame() int
	GetSingleton() *singleton.Singleton
	GetEntityByID(int) (entities.Entity, error)
	RegisterEntities([]entities.Entity)
}

type StateInterpolatorSystem struct {
	*base.BaseSystem
	world World
}

func NewStateInterpolatorSystem(world World) *StateInterpolatorSystem {
	return &StateInterpolatorSystem{
		world: world,
	}
}

func (s *StateInterpolatorSystem) RegisterEntity(entity entities.Entity) {
}

func (s *StateInterpolatorSystem) Update(delta time.Duration) {
	singleton := s.world.GetSingleton()
	state := singleton.StateBuffer.PullEntityInterpolations(s.world.CommandFrame())
	if state != nil {
		handleGameStateUpdate(state, s.world)
	}
}

func handleGameStateUpdate(bufferedState *statebuffer.BufferedState, world World) {
	singleton := world.GetSingleton()

	for _, entitySnapshot := range bufferedState.InterpolatedEntities {
		if entitySnapshot.ID == world.GetSingleton().PlayerID || entitySnapshot.ID == singleton.CameraID {
			continue
		}

		foundEntity, err := world.GetEntityByID(entitySnapshot.ID)
		if err != nil {
			fmt.Printf("[%d] failed to find entity with id %d type %d to interpolate\n", world.CommandFrame(), entitySnapshot.ID, entitySnapshot.Type)
		} else {
			cc := foundEntity.GetComponentContainer()
			cc.TransformComponent.Position = entitySnapshot.Position
			cc.TransformComponent.Orientation = entitySnapshot.Orientation
		}
	}
}
