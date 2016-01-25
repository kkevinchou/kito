package grass

import (
	"github.com/kkevinchou/ant/animation"
	"github.com/kkevinchou/ant/assets"
	"github.com/kkevinchou/ant/lib/math/vector"
)

type Grass struct {
	*RenderComponent
	*PositionComponent
}

func New(x, y float64, assetManager *assets.Manager) *Grass {
	entity := &Grass{}

	entity.PositionComponent = &PositionComponent{vector.Vector{X: x, Y: y}}

	entity.RenderComponent = &RenderComponent{
		entity:         entity,
		animationState: animation.CreateStateFromAnimationDef(assetManager.GetAnimation("grass")),
	}

	return entity
}
