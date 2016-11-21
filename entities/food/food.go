package food

import (
	"github.com/kkevinchou/ant/components"
	"github.com/kkevinchou/ant/components/id"
	"github.com/kkevinchou/ant/directory"
	"github.com/kkevinchou/ant/interfaces"
	"github.com/kkevinchou/ant/lib/math/vector"
)

type Food interface {
	interfaces.Item
}

type FoodImpl struct {
	*RenderComponent
	*components.PositionComponent
	*id.IdComponent
	*components.ItemComponent
}

func New(x, y, z float64) *FoodImpl {
	entity := &FoodImpl{}

	entity.IdComponent = id.NewIdComponent()
	entity.PositionComponent = &components.PositionComponent{}
	entity.ItemComponent = &components.ItemComponent{}

	entity.RenderComponent = &RenderComponent{
		entity: entity,
	}
	renderSystem := directory.GetDirectory().RenderSystem()
	renderSystem.Register(entity)

	itemManager := directory.GetDirectory().ItemManager()
	itemManager.Register(entity)

	entity.SetPosition(vector.Vector3{X: x, Y: y, Z: z})

	return entity
}
