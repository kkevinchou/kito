package entities

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/lib/animation"
	"github.com/kkevinchou/kito/types"
	"github.com/kkevinchou/kito/utils"
)

func NewBob(position mgl64.Vec3) *EntityImpl {
	assetManager := directory.GetDirectory().AssetManager()

	transformComponent := &components.TransformComponent{
		Position:      position,
		Orientation:   mgl64.QuatIdent(),
		ForwardVector: mgl64.Vec3{0, 0, -1},
		UpVector:      mgl64.Vec3{0, 1, 0},
	}

	renderData := &components.ModelRenderData{
		Visible:  true,
		Animated: true,
	}
	renderComponent := &components.RenderComponent{
		RenderData: renderData,
	}

	modelSpec := assetManager.GetAnimatedModel("bob")

	var animatedModel *animation.AnimatedModel
	if utils.IsClient() {
		animatedModel = animation.NewAnimatedModel(modelSpec, 50, 3)
	} else {
		animatedModel = animation.NewAnimatedModelWithoutMesh(modelSpec, 50, 3)
	}

	animationComponent := &components.AnimationComponent{
		AnimatedModel: animatedModel, // potentially shared across many entities
		Animation:     modelSpec.Animation,
	}

	physicsComponent := &components.PhysicsComponent{
		Impulses: map[string]types.Impulse{},
	}

	thirdPersonControllerComponent := &components.ThirdPersonControllerComponent{
		Controlled: true,
	}

	entityComponents := []components.Component{
		&components.NetworkComponent{},
		transformComponent,
		animationComponent,
		physicsComponent,
		thirdPersonControllerComponent,
	}

	if utils.IsClient() {
		entityComponents = append(entityComponents, renderComponent)
	}

	entity := NewEntity(
		"bob",
		components.NewComponentContainer(entityComponents...),
	)

	return entity
}
