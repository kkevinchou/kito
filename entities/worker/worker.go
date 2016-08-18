package worker

import (
	"time"

	"github.com/kkevinchou/ant/animation"
	"github.com/kkevinchou/ant/components"
	"github.com/kkevinchou/ant/components/physics"
	"github.com/kkevinchou/ant/components/steering"
	"github.com/kkevinchou/ant/systems"
)

type Worker struct {
	*physics.PhysicsComponent
	*steering.SeekComponent
	*RenderComponent
	*PositionComponent
	*AIComponent
	*components.InventoryComponent
}

func New() *Worker {
	entity := &Worker{}

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

	entity.AIComponent = NewAIComponent(entity)
	entity.InventoryComponent = components.NewInventoryComponent()

	return entity
}

func (e *Worker) Update(delta time.Duration) {
	e.PhysicsComponent.Update(delta)
	e.AIComponent.Update(delta)
}
