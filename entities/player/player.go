package player

import (
	"time"

	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/components/physics"
	"github.com/kkevinchou/kito/components/steering"
	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/interfaces"
	"github.com/kkevinchou/kito/lib/math/vector"
)

type Player interface {
	interfaces.ItemGiverReceiver
	Velocity() vector.Vector3
}

type PlayerImpl struct {
	*physics.PhysicsComponent
	*components.RenderComponent
	*components.PositionComponent
	*components.InventoryComponent
	*components.CharacterControllerComponent
}

func New() *PlayerImpl {
	entity := &PlayerImpl{}

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

	entity.InventoryComponent = components.NewInventoryComponent()

	return entity
}

func (e *PlayerImpl) Update(delta time.Duration) {
	e.PhysicsComponent.Update(delta)
	e.CharacterControllerComponent.Update(delta)
}
