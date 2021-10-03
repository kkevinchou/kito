package components

type RenderComponent struct {
	IsVisible bool
}

func (c *RenderComponent) AddToComponentContainer(container *ComponentContainer) {
	container.RenderComponent = c
}
