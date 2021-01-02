package components

type FollowComponent struct {
	FollowTargetEntityID *int
	FollowDistance       float64
}

func (c *FollowComponent) AddToComponentContainer(container *ComponentContainer) {
	container.FollowComponent = c
}
