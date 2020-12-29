package entities

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/components"
)

func NewCamera(position mgl64.Vec3, view mgl64.Vec2) *EntityImpl {
	physicsComponent := &components.PhysicsComponent{}
	physicsComponent.Init(10, 50)

	positionComponent := &components.PositionComponent{Position: position}

	topDownViewComponent := &components.TopDownViewComponent{}
	topDownViewComponent.SetView(view)

	controllerComponent := components.NewControllerComponent()
	controllerComponent.SetControlled(true)

	entity := &EntityImpl{
		ComponentContainer: components.NewComponentContainer(
			positionComponent,
			topDownViewComponent,
			physicsComponent,
			controllerComponent,
		),
	}

	return entity
}
