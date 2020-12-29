package bob

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/lib/animation"
	"github.com/kkevinchou/kito/lib/loaders/collada"
)

type BobImpl struct {
	*components.RenderComponent
	*components.PositionComponent
	*components.AnimationComponent
}

func NewBob() *BobImpl {
	entity := &BobImpl{}

	entity.PositionComponent = &components.PositionComponent{}

	renderData := &components.ModelRenderData{
		Visible: true,
		ID:      "bob",
	}

	entity.RenderComponent = &components.RenderComponent{
		RenderData: renderData,
	}

	renderSystem := directory.GetDirectory().RenderSystem()
	renderSystem.Register(entity)

	// TODO: get this garbage out of here
	parsedCollada, err := collada.ParseCollada("_assets/collada/cube2.dae")
	if err != nil {
		panic(err)
	}
	animatedModel := animation.NewAnimatedModel(parsedCollada, 50, 3)
	animatedModel.RootJoint.CalculateInverseBindTransform(mgl32.Ident4())

	entity.AnimationComponent = &components.AnimationComponent{
		AnimatedModel: animatedModel, // potentially shared across many entities
		Animation:     parsedCollada.Animation,
	}

	return entity
}
