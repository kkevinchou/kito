package stateinterpolator

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/commandframe"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/statebuffer"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/types"
)

type World interface {
	CommandFrame() int
	GetSingleton() *singleton.Singleton
	GetEntityByID(int) (entities.Entity, error)
	GetCommandFrameHistory() *commandframe.CommandFrameHistory
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

	// TODO: move new entities logic to a separate system / handler
	var newEntities []entities.Entity
	for _, entitySnapshot := range bufferedState.InterpolatedEntities {
		if entitySnapshot.ID == world.GetSingleton().PlayerID || entitySnapshot.ID == singleton.CameraID {
			continue
		}

		foundEntity, err := world.GetEntityByID(entitySnapshot.ID)
		if err != nil {
			if types.EntityType(entitySnapshot.Type) == types.EntityTypeBob {
				newEntity := entities.NewBob(mgl64.Vec3{})
				newEntity.ID = entitySnapshot.ID

				cc := newEntity.GetComponentContainer()
				cc.TransformComponent.Position = entitySnapshot.Position
				cc.TransformComponent.Orientation = entitySnapshot.Orientation

				newEntities = append(newEntities, newEntity)
			} else if types.EntityType(entitySnapshot.Type) == types.EntityTypeScene {
				newEntity := entities.NewScene(entitySnapshot.Position)
				newEntity.ID = entitySnapshot.ID

				cc := newEntity.GetComponentContainer()
				cc.TransformComponent.Position = entitySnapshot.Position
				cc.TransformComponent.Orientation = entitySnapshot.Orientation

				newEntities = append(newEntities, newEntity)
			} else if types.EntityType(entitySnapshot.Type) == types.EntityTypeStaticSlime {
				newEntity := entities.NewSlime(entitySnapshot.Position)
				newEntity.ID = entitySnapshot.ID

				cc := newEntity.GetComponentContainer()
				cc.TransformComponent.Position = entitySnapshot.Position
				cc.TransformComponent.Orientation = entitySnapshot.Orientation

				newEntities = append(newEntities, newEntity)
			} else {
				continue
			}
		} else {
			cc := foundEntity.GetComponentContainer()
			cc.TransformComponent.Position = entitySnapshot.Position
			cc.TransformComponent.Orientation = entitySnapshot.Orientation
		}
	}

	world.RegisterEntities(newEntities)
}
