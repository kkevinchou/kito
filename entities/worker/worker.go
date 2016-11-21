package worker

import (
	"time"

	"github.com/kkevinchou/ant/components"
	"github.com/kkevinchou/ant/components/physics"
	"github.com/kkevinchou/ant/components/steering"
	"github.com/kkevinchou/ant/directory"
	"github.com/kkevinchou/ant/interfaces"
	"github.com/kkevinchou/ant/lib/math/vector"
)

type Worker interface {
	interfaces.ItemGiverReceiver
	SetTarget(vector.Vector3)
	Velocity() vector.Vector3
	Heading() vector.Vector3
}

type WorkerImpl struct {
	*physics.PhysicsComponent
	*steering.SeekComponent
	*RenderComponent
	*components.PositionComponent
	*AIComponent
	*components.InventoryComponent
}

func New() *WorkerImpl {
	entity := &WorkerImpl{}

	entity.PhysicsComponent = &physics.PhysicsComponent{}
	entity.PhysicsComponent.Init(entity, 5, 10)

	entity.PositionComponent = &components.PositionComponent{}
	entity.SeekComponent = &steering.SeekComponent{Entity: entity}

	entity.RenderComponent = &RenderComponent{
		entity: entity,
	}

	renderSystem := directory.GetDirectory().RenderSystem()
	renderSystem.Register(entity)

	movementSystem := directory.GetDirectory().MovementSystem()
	movementSystem.Register(entity)

	entity.AIComponent = NewAIComponent(entity)
	entity.InventoryComponent = components.NewInventoryComponent()

	return entity
}

func (e *WorkerImpl) Update(delta time.Duration) {
	e.PhysicsComponent.Update(delta)
	e.AIComponent.Update(delta)
}
