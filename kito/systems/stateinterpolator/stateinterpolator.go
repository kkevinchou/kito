package stateinterpolator

import (
	"fmt"
	"time"

	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/managers/player"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/statebuffer"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/lib/metrics"
)

type World interface {
	CommandFrame() int
	GetSingleton() *singleton.Singleton
	GetEntityByID(int) entities.Entity
	RegisterEntities([]entities.Entity)
	GetPlayerEntity() entities.Entity
	GetPlayer() *player.Player
	MetricsRegistry() *metrics.MetricsRegistry
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

func (s *StateInterpolatorSystem) Update(delta time.Duration) {
	singleton := s.world.GetSingleton()
	state := singleton.StateBuffer.PullEntityInterpolations(s.world.CommandFrame())
	if state != nil {
		handleGameStateUpdate(state, s.world)
	}
}

func handleGameStateUpdate(bufferedState *statebuffer.BufferedState, world World) {
	playerEntity := world.GetPlayerEntity()
	for _, entitySnapshot := range bufferedState.InterpolatedEntities {
		if entitySnapshot.ID == playerEntity.GetID() {
			continue
		}

		// if entitySnapshot.ID > 80001 && entitySnapshot.Velocity.Y() > 0 {
		// 	fmt.Printf("%v cf %d\n", time.Now(), world.CommandFrame())
		// }

		foundEntity := world.GetEntityByID(entitySnapshot.ID)
		if foundEntity == nil {
			fmt.Printf("[%d] failed to find entity with id %d type %d to interpolate\n", world.CommandFrame(), entitySnapshot.ID, entitySnapshot.Type)
		} else {
			cc := foundEntity.GetComponentContainer()
			cc.TransformComponent.Position = entitySnapshot.Position
			cc.TransformComponent.Orientation = entitySnapshot.Orientation
			if cc.ThirdPersonControllerComponent != nil {
				cc.ThirdPersonControllerComponent.Velocity = entitySnapshot.Velocity
			}
		}
	}
}
