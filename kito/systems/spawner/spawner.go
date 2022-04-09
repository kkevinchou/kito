package spawner

import (
	"time"

	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/singleton"
	"github.com/kkevinchou/kito/kito/statebuffer"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/kito/utils/entityutils"
)

type World interface {
	CommandFrame() int
	GetSingleton() *singleton.Singleton
	GetEntityByID(int) (entities.Entity, error)
	RegisterEntities([]entities.Entity)
	GetPlayerEntity() entities.Entity
}

type SpawnerSystem struct {
	*base.BaseSystem
	world World
}

func NewSpawnerSystem(world World) *SpawnerSystem {
	return &SpawnerSystem{
		world: world,
	}
}

func (s *SpawnerSystem) Update(delta time.Duration) {
	singleton := s.world.GetSingleton()
	state := singleton.StateBuffer.PeekEntityInterpolations(s.world.CommandFrame())
	if state != nil {
		handleGameStateUpdate(state, s.world)
	}
}

func handleGameStateUpdate(bufferedState *statebuffer.BufferedState, world World) {
	playerEntity := world.GetPlayerEntity()

	var newEntities []entities.Entity
	for _, snapshot := range bufferedState.InterpolatedEntities {
		if snapshot.ID == playerEntity.GetID() {
			continue
		}

		_, err := world.GetEntityByID(snapshot.ID)
		if err != nil {
			newEntity := entityutils.Spawn(snapshot.ID, types.EntityType(snapshot.Type), snapshot.Position, snapshot.Orientation)
			newEntities = append(newEntities, newEntity)
		}
	}

	world.RegisterEntities(newEntities)
}
