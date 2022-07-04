package entities

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/lib/animation"
	"github.com/kkevinchou/kito/lib/collision/collider"
	"github.com/kkevinchou/kito/lib/model"
)

func NewEnemy() *EntityImpl {
	modelName := "big_cube"
	assetManager := directory.GetDirectory().AssetManager()

	transformComponent := &components.TransformComponent{
		Position:    mgl64.Vec3{150, 0, 70},
		Orientation: mgl64.QuatIdent(),
	}

	renderComponent := &components.RenderComponent{
		IsVisible: true,
	}

	modelSpec := assetManager.GetModel(modelName)
	m := model.NewModel(modelSpec)

	animationPlayer := animation.NewAnimationPlayer(m)
	animationPlayer.PlayAnimation("Idle")
	animationComponent := &components.AnimationComponent{
		Player: animationPlayer,
	}
	_ = animationComponent

	yr := mgl64.QuatRotate(mgl64.DegToRad(180), mgl64.Vec3{0, 1, 0}).Mat4()
	meshComponent := &components.MeshComponent{
		Scale:       mgl64.Scale3D(1, 1, 1),
		Orientation: yr,
		Model:       m,
	}

	capsule := collider.NewCapsule(mgl64.Vec3{0, 12, 0}, mgl64.Vec3{0, 3, 0}, 3)
	colliderComponent := &components.ColliderComponent{
		CapsuleCollider: &capsule,
	}

	entityComponents := []components.Component{
		&components.NetworkComponent{},
		transformComponent,
		// animationComponent,
		meshComponent,
		colliderComponent,
		renderComponent,
		// &components.AIComponent{},
	}

	entity := NewEntity(
		"enemy",
		types.EntityTypeEnemy,
		components.NewComponentContainer(entityComponents...),
	)

	return entity
}
