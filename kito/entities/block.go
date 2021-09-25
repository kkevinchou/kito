package entities

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/kito/components"
	"github.com/kkevinchou/kito/kito/types"
)

func NewBlock() *EntityImpl {
	transformComponent := &components.TransformComponent{
		Position:    mgl64.Vec3{0, 15, 0},
		Orientation: mgl64.QuatIdent(),
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
		types.EntityTypeBlock,
		components.NewComponentContainer(
			transformComponent,
			renderComponent,
		),
	)

	return entity
}
