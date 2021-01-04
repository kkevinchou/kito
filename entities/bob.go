package entities

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/lib/animation"
	"github.com/kkevinchou/kito/lib/loaders/collada"
	"github.com/kkevinchou/kito/types"
)

func NewBob(position mgl64.Vec3) *EntityImpl {
	transformComponent := &components.TransformComponent{
		Position:       position,
		ViewQuaternion: mgl64.QuatIdent(),
		ForwardVector:  mgl64.Vec3{0, 0, -1},
		UpVector:       mgl64.Vec3{0, 1, 0},
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

	physicsComponent := &components.PhysicsComponent{
		Impulses: map[string]types.Impulse{},
	}

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
