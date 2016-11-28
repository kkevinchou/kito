package worker

import (
	"time"

	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/components/physics"
	"github.com/kkevinchou/kito/components/steering"
	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/interfaces"
	"github.com/kkevinchou/kito/lib/math/vector"
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
	*components.RenderComponent
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

	renderData := &components.TextureRenderData{
		Visible: true,
		ID:      "worker",
	}

	entity.RenderComponent = &components.RenderComponent{
		RenderData: renderData,
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
