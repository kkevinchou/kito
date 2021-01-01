package components

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/types"
)

type RenderData interface {
	IsVisible() bool
}

type TextureRenderData struct {
	ID      string
	Visible bool
}

func (t *TextureRenderData) IsVisible() bool {
	return t.Visible
}

type ItemRenderData struct {
	ID     string
	Entity types.Ownable
}

func (t *ItemRenderData) IsVisible() bool {
	return !t.Entity.Owned()
}

type ModelRenderData struct {
	ID       string
	Visible  bool
	Animated bool
}

func (m *ModelRenderData) IsVisible() bool {
	return m.Visible
}

type BlockRenderData struct {
	Visible bool
	Color   mgl64.Vec3
	Size    mgl64.Vec3
}

func (t *BlockRenderData) IsVisible() bool {
	return t.Visible
}
