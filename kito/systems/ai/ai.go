package ai

import (
	"time"

	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
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
		transform := componentContainer.TransformComponent

		target := playerEntities[0]
		targetDist := playerEntities[0].GetComponentContainer().TransformComponent.Position.Sub(transform.Position).LenSqr()

		for _, p := range playerEntities {
			cc := p.GetComponentContainer()
			if cc.TransformComponent.Position.Sub(transform.Position).LenSqr() < targetDist {
				target = p
			}
		}

		vecToTarget := target.GetComponentContainer().TransformComponent.Position.Sub(transform.Position)

		if vecToTarget.Len() < 50 {
			continue
		}

		transform.Position = transform.Position.Add(vecToTarget.Normalize().Mul(enemyMoveSpeed * delta.Seconds()))
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
