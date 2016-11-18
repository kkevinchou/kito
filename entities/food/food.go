package food

import (
	"github.com/kkevinchou/ant/components"
	"github.com/kkevinchou/ant/components/id"
	"github.com/kkevinchou/ant/directory"
	"github.com/kkevinchou/ant/lib/math/vector"
)

type Food interface {
	Position() vector.Vector
}

type FoodImpl struct {
	*RenderComponent
	*components.PositionComponent
	*id.IdComponent
	*ItemComponent
}

func New(x, y float64) *FoodImpl {
	entity := &FoodImpl{}

	entity.IdComponent = id.NewIdComponent()
	entity.PositionComponent = &components.PositionComponent{}
	entity.ItemComponent = &ItemComponent{}

	assetManager := directory.GetDirectory().AssetManager()

	entity.RenderComponent = &RenderComponent{
		entity:  entity,
		texture: assetManager.GetTexture("F"),
	}
	renderSystem := directory.GetDirectory().RenderSystem()
	renderSystem.Register(entity)

	itemManager := directory.GetDirectory().ItemManager()
	itemManager.Register(entity)

	entity.SetPosition(vector.Vector{X: x, Y: y})

	return entity
}
