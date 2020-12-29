package food

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/components/id"
	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/types"
)

type Food interface {
	types.Item
}

type FoodImpl struct {
	*components.RenderComponent
	*components.PositionComponent
	*id.IdComponent
	*components.ItemComponent
}

func New(x, y, z float64) *FoodImpl {
	entity := &FoodImpl{}

	entity.IdComponent = id.NewIdComponent()
	entity.PositionComponent = &components.PositionComponent{}
	entity.ItemComponent = &components.ItemComponent{}

	renderData := &components.ItemRenderData{
		ID:     "mushroom-gills",
		Entity: entity,
	}

	entity.RenderComponent = &components.RenderComponent{
		RenderData: renderData,
	}
	// renderSystem := directory.GetDirectory().RenderSystem()
	// renderSystem.Register(entity)

	itemManager := directory.GetDirectory().ItemManager()
	itemManager.Register(entity)

	entity.SetPosition(mgl64.Vec3{x, y, z})

	return entity
}
