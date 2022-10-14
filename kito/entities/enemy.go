package entities

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/lib/animation"
	"github.com/kkevinchou/kito/lib/collision"
	"github.com/kkevinchou/kito/lib/collision/collider"
	"github.com/kkevinchou/kito/lib/model"
)

func NewEnemy() *EntityImpl {
	modelName := "mutant"
	assetManager := directory.GetDirectory().AssetManager()

	transformComponent := &components.TransformComponent{
		Position:    mgl64.Vec3{78, 78, -73},
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
		Scale:       mgl64.Ident4(),
		Orientation: yr,
		Model:       m,
	}

	// capsule := collider.NewCapsule(mgl64.Vec3{0, 18, 0}, mgl64.Vec3{0, 6, 0}, 6)
	capsule := collider.NewCapsuleFromModel(m)
	boundingBox := collider.BoundingBoxFromCapsule(capsule)

	colliderComponent := &components.ColliderComponent{
		BoundingBoxCollider: boundingBox,
		CapsuleCollider:     &capsule,
		Contacts:            map[int]*collision.Contact{},
	}

	entityComponents := []components.Component{
		&components.NetworkComponent{},
		transformComponent,
		animationComponent,
		meshComponent,
		colliderComponent,
		renderComponent,
		components.NewAIComponent(nil),
		&components.HealthComponent{Value: 100},
		&components.LootDropperComponent{},
	}

	entity := NewEntity(
		"enemy",
		types.EntityTypeEnemy,
		components.NewComponentContainer(entityComponents...),
	)

	return entity
}
