package entities

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/components"
)

func NewCamera(position mgl64.Vec3, view mgl64.Vec2) *EntityImpl {
	physicsComponent := &components.PhysicsComponent{}
	physicsComponent.Init()

	positionComponent := &components.PositionComponent{Position: position}

	topDownViewComponent := &components.TopDownViewComponent{}
	topDownViewComponent.SetView(view)

	controllerComponent := &components.ControllerComponent{Controlled: true}

	entity := NewEntity(
		"camera",
		components.NewComponentContainer(
			positionComponent,
			topDownViewComponent,
			physicsComponent,
			controllerComponent,
		),
	)

	return entity
}

func NewThirdPersonCamera(followTarget int, positionOffset mgl64.Vec3, view mgl64.Vec2) *EntityImpl {
	// physicsComponent := &components.PhysicsComponent{}
	// physicsComponent.Init()

	positionComponent := &components.PositionComponent{}
	controllerComponent := &components.ControllerComponent{Controlled: false}

	followComponent := &components.FollowComponent{FollowTargetEntityID: &followTarget}

	entity := NewEntity(
		"camera",
		components.NewComponentContainer(
			positionComponent,
			// physicsComponent,
			controllerComponent,
			followComponent,
		),
	)

	return entity
}
