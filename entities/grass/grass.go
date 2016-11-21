package grass

import (
	"github.com/kkevinchou/ant/components"
	"github.com/kkevinchou/ant/directory"
	"github.com/kkevinchou/ant/interfaces"
	"github.com/kkevinchou/ant/lib/math/vector"
)

type Grass interface {
	interfaces.Positionable
}

type GrassImpl struct {
	*RenderComponent
	*components.PositionComponent
}

func New(x, y, z float64) *GrassImpl {
	entity := &GrassImpl{}

	entity.PositionComponent = &components.PositionComponent{}

	entity.RenderComponent = &RenderComponent{
		entity: entity,
	}
	renderSystem := directory.GetDirectory().RenderSystem()
	renderSystem.Register(entity)

	entity.SetPosition(vector.Vector3{X: x, Y: y, Z: z})

	return entity
}
