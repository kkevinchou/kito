package components

import "github.com/kkevinchou/kito/types"

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

type ModelRenderData struct {
	ID      string
	Visible bool
}

func (m *ModelRenderData) IsVisible() bool {
	return m.Visible
}

type ItemRenderData struct {
	ID     string
	Entity types.Ownable
}

func (t *ItemRenderData) IsVisible() bool {
	return !t.Entity.Owned()
}

type RenderComponent struct {
	RenderData RenderData
}

func (r *RenderComponent) GetRenderData() RenderData {
	return r.RenderData
}
