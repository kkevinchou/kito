package grass

import (
	"github.com/kkevinchou/ant/components"
	"github.com/kkevinchou/ant/directory"
	"github.com/kkevinchou/ant/lib"
	"github.com/kkevinchou/ant/lib/math/vector"
)

type Grass struct {
	*RenderComponent
	*components.PositionComponent
}

func New(x, y float64) *Grass {
	entity := &Grass{}

	entity.PositionComponent = &components.PositionComponent{}

	assetManager := directory.GetDirectory().AssetManager()

	entity.RenderComponent = &RenderComponent{
		entity:         entity,
		animationState: lib.CreateStateFromAnimationDef(assetManager.GetAnimation("grass")),
	}
	renderSystem := directory.GetDirectory().RenderSystem()
	renderSystem.Register(entity)

	entity.SetPosition(vector.Vector{X: x, Y: y})

	return entity
}
