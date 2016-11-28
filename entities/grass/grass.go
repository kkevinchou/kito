package grass

import (
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/interfaces"
	"github.com/kkevinchou/kito/lib/math/vector"
)

type Grass interface {
	interfaces.Positionable
}

type GrassImpl struct {
	*components.RenderComponent
	*components.PositionComponent
}

func New(x, y, z float64) *GrassImpl {
	entity := &GrassImpl{}

	entity.PositionComponent = &components.PositionComponent{}

	renderData := &components.TextureRenderData{
		Visible: true,
		ID:      "high-grass",
	}

	entity.RenderComponent = &components.RenderComponent{
		RenderData: renderData,
	}

	renderSystem := directory.GetDirectory().RenderSystem()
	renderSystem.Register(entity)

	entity.SetPosition(vector.Vector3{X: x, Y: y, Z: z})

	return entity
}
