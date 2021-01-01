package entities

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/components"
)

func NewBlock() *EntityImpl {
	positionComponent := &components.PositionComponent{}

	renderData := &components.BlockRenderData{
		Visible: true,
		Size:    mgl64.Vec3{100, 100, 10},
	}
	renderComponent := &components.RenderComponent{
		RenderData: renderData,
	}

	entity := NewEntity(
		"block",
		components.NewComponentContainer(
			positionComponent,
			renderComponent,
		),
	)

	return entity
}
