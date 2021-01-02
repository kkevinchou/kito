package entities

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/lib/animation"
	"github.com/kkevinchou/kito/lib/loaders/collada"
)

func NewBob() *EntityImpl {
	transformComponent := &components.TransformComponent{
		ViewQuaternion: mgl64.QuatIdent(),
	}

	renderData := &components.ModelRenderData{
		Visible:  true,
		Animated: true,
	}
	renderComponent := &components.RenderComponent{
		RenderData: renderData,
	}

	// TODO: get this garbage out of here
	parsedCollada, err := collada.ParseCollada("_assets/collada/bob.dae")
	if err != nil {
		panic(err)
	}
	animatedModel := animation.NewAnimatedModel(parsedCollada, 50, 3)

	animationComponent := &components.AnimationComponent{
		AnimatedModel: animatedModel, // potentially shared across many entities
		Animation:     parsedCollada.Animation,
	}

	physicsComponent := &components.PhysicsComponent{}
	physicsComponent.Init()

	thirdPersonControllerComponent := &components.ThirdPersonControllerComponent{
		Controlled: true,
	}

	entity := NewEntity(
		"bob",
		components.NewComponentContainer(
			transformComponent,
			renderComponent,
			animationComponent,
			physicsComponent,
			thirdPersonControllerComponent,
		),
	)

	return entity
}
