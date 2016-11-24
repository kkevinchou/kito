package components

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

type RenderComponent struct {
	RenderData RenderData
}

func (r *RenderComponent) GetRenderData() RenderData {
	return r.RenderData
}
