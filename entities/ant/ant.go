package ant

import (
	"time"

	"github.com/kkevinchou/ant/animation"
	"github.com/kkevinchou/ant/components/physics"
	"github.com/kkevinchou/ant/components/steering"
	"github.com/kkevinchou/ant/systems"
)

type Ant struct {
	*physics.PhysicsComponent
	*steering.SeekComponent
	*RenderComponent
	*PositionComponent
}

func New() *Ant {
	entity := &Ant{}

	entity.PhysicsComponent = &physics.PhysicsComponent{}
	entity.PhysicsComponent.Init(entity, 100, 10)

	entity.PositionComponent = &PositionComponent{}
	entity.SeekComponent = &steering.SeekComponent{Entity: entity}

	assetManager := systems.GetDirectory().AssetManager()

	entity.RenderComponent = &RenderComponent{
		entity:         entity,
		animationState: animation.CreateStateFromAnimationDef(assetManager.GetAnimation("ant")),
	}

	renderSystem := systems.GetDirectory().RenderSystem()
	renderSystem.Register(entity)

	movementSystem := systems.GetDirectory().MovementSystem()
	movementSystem.Register(entity)

	return entity
}

func (e *Ant) Update(delta time.Duration) {
	e.PhysicsComponent.Update(delta)
}
