package components

type FollowComponent struct {
	FollowTargetEntityID int
	FollowDistance       float64
	MaxFollowDistance    float64
	YOffset              float64

	// this zoom stuff probably doesn't belong here
	ZoomSpeed float64
}

func (c *FollowComponent) AddToComponentContainer(container *ComponentContainer) {
	container.FollowComponent = c
}
