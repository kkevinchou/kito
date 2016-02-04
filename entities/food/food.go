package food

import (
	"github.com/kkevinchou/ant/animation"
	"github.com/kkevinchou/ant/components/id"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/kkevinchou/ant/systems"
)

type Food struct {
	*RenderComponent
	*PositionComponent
	*id.IdComponent
	*ItemComponent
}

func New(x, y float64) *Food {
	entity := &Food{}

	entity.IdComponent = id.NewIdComponent()
	entity.PositionComponent = &PositionComponent{vector.Vector{X: x, Y: y}}
	entity.ItemComponent = &ItemComponent{}

	assetManager := systems.GetDirectory().AssetManager()

	entity.RenderComponent = &RenderComponent{
		entity:         entity,
		animationState: animation.CreateStateFromAnimationDef(assetManager.GetAnimation("grass")),
	}
	renderSystem := systems.GetDirectory().RenderSystem()
	renderSystem.Register(entity)

	itemManager := systems.GetDirectory().ItemManager()
	itemManager.Register(entity)

	return entity
}
