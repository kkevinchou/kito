package entities

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/components"
)

func NewThirdPersonCamera(positionOffset mgl64.Vec3, view mgl64.Vec2, followTargetEntityID int) *EntityImpl {
	// TODO: sync initial position from transform compnoent from follow component

	transformComponent := &components.TransformComponent{
		Orientation:   mgl64.QuatIdent(),
		Position:      mgl64.Vec3{0, 0, 1},
		UpVector:      mgl64.Vec3{0, 1, 0},
		ForwardVector: mgl64.Vec3{0, 0, -1},
	}

	followComponent := &components.FollowComponent{
		FollowTargetEntityID: followTargetEntityID,
		FollowDistance:       40,
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
