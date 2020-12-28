package viewer

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/components/physics"
	"github.com/kkevinchou/kito/lib/math/vector"
)

type ViewerImpl struct {
	*physics.PhysicsComponent
	*components.PositionComponent
	*components.TopDownViewComponent
	*components.ControllerComponent
}

func New(position mgl64.Vec3, view vector.Vector) *ViewerImpl {
	entity := &ViewerImpl{}

	entity.PhysicsComponent = &physics.PhysicsComponent{}
	entity.PhysicsComponent.Init(entity, 50, 10)

	entity.PositionComponent = &components.PositionComponent{}
	entity.SetPosition(position)
	entity.TopDownViewComponent = &components.TopDownViewComponent{}
	entity.SetView(view)

	entity.ControllerComponent = components.NewControllerComponent()
	entity.ControllerComponent.SetControlled(true)

	return entity
}

func (e *ViewerImpl) Update(delta time.Duration) {
	e.PhysicsComponent.Update(delta)
}
