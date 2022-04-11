package ai

import (
	"time"

	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/entities"
	"github.com/kkevinchou/kito/kito/systems/base"
)

type World interface {
	QueryEntity(componentFlags int) []entities.Entity
	GetEntityByID(id int) entities.Entity
}

type AISystem struct {
	*base.BaseSystem
	world World
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
		_ = target

		for _, p := range playerEntities {
			cc := p.GetComponentContainer()
			if cc.TransformComponent.Position.Sub(transform.Position).LenSqr() < targetDist {
				target = p
			}
		}

	}
}
