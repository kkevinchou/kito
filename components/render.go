package components

type RenderComponent struct {
	RenderData RenderData
}

func (r *RenderComponent) GetRenderData() RenderData {
	return r.RenderData
}

func (c *RenderComponent) AddToComponentContainer(container *ComponentContainer) {
	container.RenderComponent = c
}
