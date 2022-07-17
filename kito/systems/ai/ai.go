package ai

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/kito/systems/base"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/lib/libutils"
)

const (
	enemyMoveSpeed = 40
)

type World interface {
	QueryEntity(componentFlags int) []entities.Entity
	GetEntityByID(id int) entities.Entity
	RegisterEntities(es []entities.Entity)
}

type AISystem struct {
	*base.BaseSystem
	world        World
	spawnTrigger int
}

func NewAnimationSystem(world World) *AISystem {
	return &AISystem{
		BaseSystem: &base.BaseSystem{},
		world:      world,
	}
}

func (s *AISystem) Update(delta time.Duration) {
	playerManager := directory.GetDirectory().PlayerManager()
	players := playerManager.GetPlayers()
	var playerEntities []entities.Entity

	for _, p := range players {
		e := s.world.GetEntityByID(p.EntityID)
		if e == nil {
			continue
		}
		playerEntities = append(playerEntities, e)
	}

	if len(playerEntities) == 0 {
		return
	}
	playerPosition := playerEntities[0].GetComponentContainer().TransformComponent.Position

	for _, entity := range s.world.QueryEntity(components.ComponentFlagAI) {
		componentContainer := entity.GetComponentContainer()
		transformComponent := componentContainer.TransformComponent
		aiComponent := componentContainer.AIComponent

		if entity.Type() == types.EntityTypeEnemy {
			if time.Since(aiComponent.LastUpdate) > time.Duration(rand.Intn(5)+2)*time.Second {
				aiComponent.AIState = components.AIStateWalk
				aiComponent.LastUpdate = time.Now()
				aiToPlayer := playerPosition.Sub(transformComponent.Position)
				aiToPlayer[1] = 0

				dir := mgl64.Vec3{}
				if aiToPlayer.Len() < 200 {
					dir = aiToPlayer.Normalize()
					aiComponent.AIState = components.AIStateAttack
				} else {
					dir = mgl64.Vec3{rand.Float64()*2 - 1, 0, rand.Float64()*2 - 1}.Normalize()
				}

				aiComponent.MovementDir = libutils.Vec3ToQuat(dir)
			}
		} else {
			fmt.Println("unhandled ai entity type")
			continue
		}

		aiComponent.Velocity = aiComponent.Velocity.Add(settings.AccelerationDueToGravity.Mul(delta.Seconds()))
		movementVec := aiComponent.MovementDir.Rotate(mgl64.Vec3{0, 0, -1})
		velocity := aiComponent.Velocity.Add(movementVec.Mul(enemyMoveSpeed))
		transformComponent.Position = transformComponent.Position.Add(velocity.Mul(delta.Seconds()))
		transformComponent.Orientation = aiComponent.MovementDir

		// safeguard falling off the map
		if transformComponent.Position[1] < -1000 {
			transformComponent.Position[1] = 25
		}
	}

	aiCount := len(s.world.QueryEntity(components.ComponentFlagAI))
	if aiCount < 5 {
		s.spawnTrigger += int(delta.Milliseconds())
		if s.spawnTrigger > 2000 {
			enemy := entities.NewEnemy()
			x := rand.Intn(1000) - 500
			z := rand.Intn(1000) - 500
			enemy.GetComponentContainer().TransformComponent.Position = mgl64.Vec3{float64(x), 0, float64(z)}
			s.world.RegisterEntities([]entities.Entity{enemy})
			s.spawnTrigger -= 3000
		}
	}
}

func (s *AISystem) Name() string {
	return "AISystem"
}
