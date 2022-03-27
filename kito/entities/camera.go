package entities

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/types"
)

const (
	maxFollowDistance     float64 = 300
	defaultFollowDistance float64 = 60
	defaultFollowY        float64 = 15
)

func NewThirdPersonCamera(positionOffset mgl64.Vec3, view mgl64.Vec2, playerID int, followTargetEntityID int) *EntityImpl {
	followComponent := &components.FollowComponent{
		FollowTargetEntityID: followTargetEntityID,
		FollowDistance:       defaultFollowDistance,
		MaxFollowDistance:    maxFollowDistance,
		YOffset:              defaultFollowY,
	}

	transformComponent := &components.TransformComponent{
		Orientation: mgl64.QuatIdent(),
		Position:    mgl64.Vec3{0, followComponent.YOffset, followComponent.FollowDistance},
	}

	entity := NewEntity(
		"camera",
		types.EntityTypeCamera,
		components.NewComponentContainer(
			&components.NetworkComponent{},
			&components.CameraComponent{},
			transformComponent,
			followComponent,
			&components.ControlComponent{PlayerID: playerID},
		),
	)

	return entity
}
