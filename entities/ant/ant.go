package ant

import (
	"time"

	"github.com/kkevinchou/ant/assets"
	"github.com/kkevinchou/ant/physics"
	"github.com/kkevinchou/ant/render"
	"github.com/kkevinchou/ant/steering"
)

type Ant struct {
	*physics.PhysicsComponent
	*steering.SeekComponent
	*RenderComponent
	*PositionComponent
}

func New(assetManager *assets.Manager) *Ant {
	entity := &Ant{}

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

func (e *Ant) Update(delta time.Duration) {
	e.PhysicsComponent.Update(delta)
}
