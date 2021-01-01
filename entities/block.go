package entities

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/components"
)

func NewBlock() *EntityImpl {
	transformComponent := &components.TransformComponent{
		View: mgl64.Vec3{0, 0, -1},
	}

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
			transformComponent,
			renderComponent,
		),
	)

	return entity
}
