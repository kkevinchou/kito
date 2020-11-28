package viewer

import (
	"time"

	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/components/physics"
	"github.com/kkevinchou/kito/lib/math/vector"
)

type ViewerImpl struct {
	*physics.PhysicsComponent
	*components.PositionComponent
	*components.ViewComponent
	*components.CharacterControllerComponent
}

func New(position vector.Vector3, view vector.Vector) *ViewerImpl {
	entity := &ViewerImpl{}

	entity.PhysicsComponent = &physics.PhysicsComponent{}
	entity.PhysicsComponent.Init(entity, 5, 10)

	entity.PositionComponent = &components.PositionComponent{}
	entity.SetPosition(position)
	entity.ViewComponent = &components.ViewComponent{}
	entity.SetView(view)

	entity.CharacterControllerComponent = components.NewCharacterControllerComponent(entity)

	return entity
}

func (e *ViewerImpl) Update(delta time.Duration) {
	e.CharacterControllerComponent.Update(delta)
	e.PhysicsComponent.Update(delta)
}
