package grass

import (
	"github.com/kkevinchou/ant/animation"
	"github.com/kkevinchou/ant/components"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/kkevinchou/ant/systems"
)

type Grass struct {
	*RenderComponent
	*components.PositionComponent
}

func New(x, y float64) *Grass {
	entity := &Grass{}

	entity.PositionComponent = &components.PositionComponent{}

	assetManager := systems.GetDirectory().AssetManager()

	entity.RenderComponent = &RenderComponent{
		entity:         entity,
		animationState: animation.CreateStateFromAnimationDef(assetManager.GetAnimation("grass")),
	}
	renderSystem := systems.GetDirectory().RenderSystem()
	renderSystem.Register(entity)

	entity.SetPosition(vector.Vector{X: x, Y: y})

	return entity
}
