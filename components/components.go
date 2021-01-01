package components

type ComponentContainer struct {
	AnimationComponent   *AnimationComponent
	RenderComponent      *RenderComponent
	PositionComponent    *PositionComponent
	PhysicsComponent     *PhysicsComponent
	TopDownViewComponent *TopDownViewComponent
	ControllerComponent  *ControllerComponent
	FollowComponent      *FollowComponent
}

type Component interface {
	AddToComponentContainer(container *ComponentContainer)
}

func NewComponentContainer(components ...Component) *ComponentContainer {
	container := &ComponentContainer{}
	for _, component := range components {
		component.AddToComponentContainer(container)
	}
	return container
}
