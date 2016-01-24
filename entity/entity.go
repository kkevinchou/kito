package entity

import (
	"time"

	"github.com/kkevinchou/ant/assets"
	"github.com/kkevinchou/ant/physics"
	"github.com/kkevinchou/ant/render"
	"github.com/kkevinchou/ant/steering"
)

type Entity struct {
	*physics.PhysicsComponent
	*steering.SeekComponent
	*RenderComponent
	*PositionComponent
}

func New(assetManager *assets.Manager) *Entity {
	entity := &Entity{}

	entity.PhysicsComponent = &physics.PhysicsComponent{}
	entity.PhysicsComponent.Init(entity, 100, 10)

	entity.PositionComponent = &PositionComponent{}
	entity.SeekComponent = &steering.SeekComponent{Entity: entity}

	entity.RenderComponent = &RenderComponent{
		entity: entity,
		animationState: render.AnimationState{
			MetaData: assetManager.GetAnimationMetaData("ant"),
		},
	}

	return entity
}

func (e *Entity) Update(delta time.Duration) {
	e.PhysicsComponent.Update(delta)
}
