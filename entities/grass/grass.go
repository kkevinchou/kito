package grass

import (
	"github.com/kkevinchou/ant/assets"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/kkevinchou/ant/render"
)

type Grass struct {
	*RenderComponent
	*PositionComponent
}

func New(x, y float64, assetManager *assets.Manager) *Grass {
	entity := &Grass{}

	entity.PositionComponent = &PositionComponent{vector.Vector{X: x, Y: y}}

	entity.RenderComponent = &RenderComponent{
		entity: entity,
		animationState: render.AnimationState{
			MetaData: assetManager.GetAnimationMetaData("grass"),
		},
	}

	return entity
}
