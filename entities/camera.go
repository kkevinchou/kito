package entities

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/components"
)

func NewCamera(position mgl64.Vec3, view mgl64.Vec2) *EntityImpl {
	physicsComponent := &components.PhysicsComponent{}
	physicsComponent.Init()

	transformComponent := &components.TransformComponent{
		Position: position,
		View:     mgl64.Vec3{0, 0, -1},
	}

	topDownViewComponent := &components.TopDownViewComponent{}
	topDownViewComponent.SetView(view)

	controllerComponent := &components.ControllerComponent{Controlled: true}

	entity := NewEntity(
		"camera",
		components.NewComponentContainer(
			transformComponent,
			topDownViewComponent,
			physicsComponent,
			controllerComponent,
		),
	)

	return entity
}

func NewThirdPersonCamera(followTarget int, positionOffset mgl64.Vec3, view mgl64.Vec2) *EntityImpl {
	transformComponent := &components.TransformComponent{
		View: mgl64.Vec3{0, 0, -10},
		// View: mgl64.Vec3{0, -1, -10},
	}
	controllerComponent := &components.ControllerComponent{Controlled: false}

	followComponent := &components.FollowComponent{FollowTargetEntityID: &followTarget}

	entity := NewEntity(
		"camera",
		components.NewComponentContainer(
			transformComponent,
			// physicsComponent,
			controllerComponent,
			followComponent,
		),
	)

	return entity
}
