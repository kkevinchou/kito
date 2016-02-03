package grass

import (
	"github.com/kkevinchou/ant/animation"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/kkevinchou/ant/systems"
)

type Grass struct {
	*RenderComponent
	*PositionComponent
}

func New(x, y float64) *Grass {
	entity := &Grass{}

	entity.PositionComponent = &PositionComponent{vector.Vector{X: x, Y: y}}

	assetManager := systems.GetDirectory().AssetManager()

	entity.RenderComponent = &RenderComponent{
		entity:         entity,
		animationState: animation.CreateStateFromAnimationDef(assetManager.GetAnimation("grass")),
	}
	renderSystem := systems.GetDirectory().RenderSystem()
	renderSystem.Register(entity)

	return entity
}
