package entities

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/directory"
	"github.com/kkevinchou/kito/kito/types"
	"github.com/kkevinchou/kito/kito/utils"
	"github.com/kkevinchou/kito/lib/animation"
	"github.com/kkevinchou/kito/lib/collision/collider"
	"github.com/kkevinchou/kito/lib/model"
	"github.com/kkevinchou/kito/lib/textures"
)

func NewBob() *EntityImpl {
	modelName := "human"
	textureName := "color_grid"
	assetManager := directory.GetDirectory().AssetManager()

	transformComponent := &components.TransformComponent{
		Position:    mgl64.Vec3{0, 0, 70},
		Orientation: mgl64.QuatIdent(),
	}

	renderComponent := &components.RenderComponent{
		IsVisible: true,
	}

	var texture *textures.Texture
	if utils.IsClient() {
		texture = assetManager.GetTexture(textureName)
	}

	modelSpec := assetManager.GetAnimatedModel(modelName)
	m := model.NewModel(modelSpec)
	if utils.IsClient() {
		m.Prepare()
	}

	animationPlayer := animation.NewAnimationPlayer(m.Animations)
	animationPlayer.PlayAnimation("Idle")
	animationComponent := &components.AnimationComponent{
		Player: animationPlayer,
	}
	_ = animationComponent

	yr := mgl64.QuatRotate(mgl64.DegToRad(180), mgl64.Vec3{0, 1, 0}).Mat4()
	meshComponent := &components.MeshComponent{
		Texture: texture,
		// Scale:            mgl64.Scale3D(1, 1, 1),
		Scale: mgl64.Scale3D(10, 10, 10),
		// Orientation: mgl64.Ident4(),
		Orientation: yr,

		Model: m,
	}

	capsule := collider.NewCapsule(mgl64.Vec3{0, 12, 0}, mgl64.Vec3{0, 3, 0}, 3)
	colliderComponent := &components.ColliderComponent{
		CapsuleCollider: &capsule,
	}

	thirdPersonControllerComponent := &components.ThirdPersonControllerComponent{
		Controlled: true,
	}

	entityComponents := []components.Component{
		&components.NetworkComponent{},
		transformComponent,
		// animationComponent,
		thirdPersonControllerComponent,
		meshComponent,
		colliderComponent,
		renderComponent,
	}

	entity := NewEntity(
		"bob",
		types.EntityTypeBob,
		components.NewComponentContainer(entityComponents...),
	)

	return entity
}
