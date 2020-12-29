package worker

import (
	"time"

	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/components/physics"
	"github.com/kkevinchou/kito/components/steering"
	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/types"
)

type WorkerImpl struct {
	*physics.PhysicsComponent
	*steering.SeekComponent
	*components.RenderComponent
	*components.PositionComponent
	*components.AIComponent
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

	// renderSystem := directory.GetDirectory().RenderSystem()
	// renderSystem.Register(entity)

	movementSystem := directory.GetDirectory().MovementSystem()
	movementSystem.Register(entity)

	entity.AIComponent = components.NewAIComponent(NewBT(entity))
	entity.InventoryComponent = components.NewInventoryComponent()

	return entity
}

func (e *WorkerImpl) Update(delta time.Duration) {
	e.PhysicsComponent.Update(delta)
	e.AIComponent.Update(delta)
}

func (e *WorkerImpl) MovementType() types.MovementType {
	return types.MovementTypeSteering
}
