package components

type ComponentContainer struct {
	AnimationComponent             *AnimationComponent
	RenderComponent                *RenderComponent
	TransformComponent             *TransformComponent
	PhysicsComponent               *PhysicsComponent
	TopDownViewComponent           *TopDownViewComponent
	ThirdPersonControllerComponent *ThirdPersonControllerComponent
	FollowComponent                *FollowComponent
	CameraComponent                *CameraComponent
	NetworkComponent               *NetworkComponent
	MeshComponent                  *MeshComponent
	ColliderComponent              *ColliderComponent
	ControlComponent               *ControlComponent
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
