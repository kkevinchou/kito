package components

type FollowComponent struct {
	FollowTargetEntityID int
	FollowDistance       float64
	MaxFollowDistance    float64

	Zoom          float64
	ZoomDirection int
	ZoomVelocity  float64
}

func (c *FollowComponent) AddToComponentContainer(container *ComponentContainer) {
	container.FollowComponent = c
}
