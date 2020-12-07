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
	*components.TopDownViewComponent
	*components.CharacterControllerComponent
}

func New(position vector.Vector3, view vector.Vector) *ViewerImpl {
	entity := &ViewerImpl{}

	entity.PhysicsComponent = &physics.PhysicsComponent{}
	entity.PhysicsComponent.Init(entity, 50, 10)

	entity.PositionComponent = &components.PositionComponent{}
	entity.SetPosition(position)
	entity.TopDownViewComponent = &components.TopDownViewComponent{}
	entity.SetView(view)

	entity.CharacterControllerComponent = components.NewCharacterControllerComponent(entity)

	return entity
}

func (e *ViewerImpl) Update(delta time.Duration) {
	// this code should be in a system
	e.CharacterControllerComponent.Update(delta)
	e.PhysicsComponent.Update(delta)
}
