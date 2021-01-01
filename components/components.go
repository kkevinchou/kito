package components

type ComponentContainer struct {
	AnimationComponent   *AnimationComponent
	RenderComponent      *RenderComponent
	TransformComponent   *TransformComponent
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
