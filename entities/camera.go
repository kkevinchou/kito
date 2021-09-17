package entities

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/types"
)

const (
	maxFollowDistance     float64 = 500
	defaultFollowDistance float64 = 40
)

func NewThirdPersonCamera(positionOffset mgl64.Vec3, view mgl64.Vec2, followTargetEntityID int) *EntityImpl {
	// TODO: sync initial position from transform compnoent from follow component

	transformComponent := &components.TransformComponent{
		Orientation: mgl64.QuatIdent(),
		Position:    mgl64.Vec3{0, 0, 1},
	}

	followComponent := &components.FollowComponent{
		FollowTargetEntityID: followTargetEntityID,
		FollowDistance:       defaultFollowDistance,
		MaxFollowDistance:    maxFollowDistance,
		YOffset:              5,
	}

	entity := NewEntity(
		"camera",
		types.EntityTypeCamera,
		components.NewComponentContainer(
			&components.NetworkComponent{},
			&components.CameraComponent{},
			transformComponent,
			followComponent,
		),
	)

	return entity
}
