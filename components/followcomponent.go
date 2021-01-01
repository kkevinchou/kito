package components

type FollowComponent struct {
	FollowTargetEntityID *int
}

func (c *FollowComponent) AddToComponentContainer(container *ComponentContainer) {
	container.FollowComponent = c
}
