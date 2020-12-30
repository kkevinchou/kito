package entities

import (
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/lib/animation"
	"github.com/kkevinchou/kito/lib/loaders/collada"
)

func NewBob() *EntityImpl {
	positionComponent := &components.PositionComponent{}

	renderData := &components.ModelRenderData{
		Visible:  true,
		ID:       "bob",
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

	entity := &EntityImpl{
		ComponentContainer: components.NewComponentContainer(
			positionComponent,
			renderComponent,
			animationComponent,
		),
	}

	return entity
}
