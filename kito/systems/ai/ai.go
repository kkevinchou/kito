package ai

import (
	"math"
	"math/rand"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/kito/systems/base"
)

const (
	enemyMoveSpeed = 30
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

	if len(playerEntities) <= 0 {
		return
	}

	for _, entity := range s.world.QueryEntity(components.ComponentFlagAI) {
		componentContainer := entity.GetComponentContainer()
		transformComponent := componentContainer.TransformComponent
		aiComponent := componentContainer.AIComponent

		if time.Since(aiComponent.LastUpdate) > 5*time.Second {
			aiComponent.LastUpdate = time.Now()
			aiComponent.MovementDir = mgl64.QuatRotate(rand.Float64()*2*math.Pi, mgl64.Vec3{0, 1, 0})
		}

		aiComponent.Velocity = aiComponent.Velocity.Add(settings.AccelerationDueToGravity.Mul(delta.Seconds()))
		movementVec := aiComponent.MovementDir.Rotate(mgl64.Vec3{0, 0, -1})
		velocity := aiComponent.Velocity.Add(movementVec.Mul(10))
		transformComponent.Position = transformComponent.Position.Add(velocity.Mul(delta.Seconds()))
		transformComponent.Orientation = aiComponent.MovementDir

		// safeguard falling off the map
		if transformComponent.Position[1] < -1000 {
			transformComponent.Position[1] = 25
		}
	}

	// s.spawnTrigger += int(delta.Milliseconds())
	// if s.spawnTrigger > 3000 {
	// 	enemy := entities.NewEnemy()
	// 	x := rand.Intn(600) - 300
	// 	z := rand.Intn(600) - 300
	// 	enemy.GetComponentContainer().TransformComponent.Position = mgl64.Vec3{float64(x), 0, float64(z)}
	// 	s.world.RegisterEntities([]entities.Entity{enemy})
	// 	s.spawnTrigger -= 3000
	// }
}
