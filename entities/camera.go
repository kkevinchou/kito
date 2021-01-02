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
	// TODO: sync initial position from transform compnoent from follow component
	initialPosition := mgl64.Vec3{0, 0, 40}

	transformComponent := &components.TransformComponent{
		View:           mgl64.Vec3{0, 0, -10},
		ViewQuaternion: mgl64.QuatIdent(),
		Position:       initialPosition,
	}
	controllerComponent := &components.ControllerComponent{Controlled: false}

	followComponent := &components.FollowComponent{
		FollowTargetEntityID: &followTarget,
		InitialOffset:        initialPosition,
		Rotation:             mgl64.QuatIdent(),
	}

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
