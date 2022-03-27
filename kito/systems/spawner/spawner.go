package spawner

import (
	"time"

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

func (s *SpawnerSystem) RegisterEntity(entity entities.Entity) {
}

func (s *SpawnerSystem) Update(delta time.Duration) {
	singleton := s.world.GetSingleton()
	state := singleton.StateBuffer.PeekEntityInterpolations(s.world.CommandFrame())
	if state != nil {
		handleGameStateUpdate(state, s.world)
	}
}

func handleGameStateUpdate(bufferedState *statebuffer.BufferedState, world World) {
	singleton := world.GetSingleton()
	playerEntity := world.GetPlayerEntity()

	var newEntities []entities.Entity
	for _, entitySnapshot := range bufferedState.InterpolatedEntities {
		if entitySnapshot.ID == playerEntity.GetID() || entitySnapshot.ID == singleton.CameraID {
			continue
		}

		_, err := world.GetEntityByID(entitySnapshot.ID)
		if err != nil {
			var newEntity *entities.EntityImpl
			if types.EntityType(entitySnapshot.Type) == types.EntityTypeBob {
				newEntity = entities.NewBob()
				newEntity.ID = entitySnapshot.ID
			} else if types.EntityType(entitySnapshot.Type) == types.EntityTypeScene {
				newEntity = entities.NewScene(entitySnapshot.Position)
				newEntity.ID = entitySnapshot.ID
			} else if types.EntityType(entitySnapshot.Type) == types.EntityTypeStaticSlime {
				newEntity = entities.NewSlime(entitySnapshot.Position)
				newEntity.ID = entitySnapshot.ID
			} else if types.EntityType(entitySnapshot.Type) == types.EntityTypeDynamicRigidBody {
				newEntity = entities.NewDynamicRigidBody(entitySnapshot.Position)
				newEntity.ID = entitySnapshot.ID
			} else if types.EntityType(entitySnapshot.Type) == types.EntityTypeStaticRigidBody {
				newEntity = entities.NewStaticRigidBody(entitySnapshot.Position)
				newEntity.ID = entitySnapshot.ID
			} else if types.EntityType(entitySnapshot.Type) == types.EntityTypeProjectile {
				newEntity = entities.NewProjectile(entitySnapshot.Position)
				newEntity.ID = entitySnapshot.ID
			} else {
				continue
			}
			cc := newEntity.GetComponentContainer()
			cc.TransformComponent.Position = entitySnapshot.Position
			cc.TransformComponent.Orientation = entitySnapshot.Orientation
			newEntities = append(newEntities, newEntity)
		}
	}

	world.RegisterEntities(newEntities)
}
