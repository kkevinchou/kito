package entities

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/types"
)

func NewCamera(position mgl64.Vec3, view mgl64.Vec2) *EntityImpl {
	physicsComponent := &components.PhysicsComponent{
		Impulses: map[string]types.Impulse{},
	}

	transformComponent := &components.TransformComponent{
		Position: position,
	}

	topDownViewComponent := &components.TopDownViewComponent{}
	topDownViewComponent.SetView(view)

	entity := NewEntity(
		"camera",
		components.NewComponentContainer(
			transformComponent,
			topDownViewComponent,
			physicsComponent,
		),
	)

	return entity
}

func NewThirdPersonCamera(positionOffset mgl64.Vec3, view mgl64.Vec2) *EntityImpl {
	// TODO: sync initial position from transform compnoent from follow component

	transformComponent := &components.TransformComponent{
		ViewQuaternion: mgl64.QuatIdent(),
		Position:       mgl64.Vec3{0, 0, 1},
		UpVector:       mgl64.Vec3{0, 1, 0},
		ForwardVector:  mgl64.Vec3{0, 0, -1},
	}

	followComponent := &components.FollowComponent{
		FollowDistance: 40,
	}

	entity := NewEntity(
		"camera",
		components.NewComponentContainer(
			transformComponent,
			followComponent,
		),
	)

	return entity
}
