package food

import (
	"github.com/kkevinchou/ant/components"
	"github.com/kkevinchou/ant/components/id"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/kkevinchou/ant/systems"
)

type Food struct {
	*RenderComponent
	*components.PositionComponent
	*id.IdComponent
	*ItemComponent
}

func New(x, y float64) *Food {
	entity := &Food{}

	entity.IdComponent = id.NewIdComponent()
	entity.PositionComponent = &components.PositionComponent{}
	entity.ItemComponent = &ItemComponent{}

	assetManager := systems.GetDirectory().AssetManager()

	entity.RenderComponent = &RenderComponent{
		entity:  entity,
		texture: assetManager.GetTexture("F"),
	}
	renderSystem := systems.GetDirectory().RenderSystem()
	renderSystem.Register(entity)

	itemManager := systems.GetDirectory().ItemManager()
	itemManager.Register(entity)

	entity.SetPosition(vector.Vector{X: x, Y: y})

	return entity
}
