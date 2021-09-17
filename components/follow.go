package components

type FollowComponent struct {
	FollowTargetEntityID int
	FollowDistance       float64
	MaxFollowDistance    float64
	YOffset              float64

	// this zoom stuff probably doesn't belong here
	Zoom      float64
	ZoomSpeed float64
}

func (c *FollowComponent) AddToComponentContainer(container *ComponentContainer) {
	container.FollowComponent = c
}
