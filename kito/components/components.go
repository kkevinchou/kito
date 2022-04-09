package components

type EntityTypeBitFlag int

const (
	ComponentFlagAI                    = 1 << 0
	ComponentFlagAnimation             = 1 << 1
	ComponentFlagCamera                = 1 << 2
	ComponentFlagCollider              = 1 << 3
	ComponentFlagControl               = 1 << 4
	ComponentFlagFollow                = 1 << 5
	ComponentFlagMesh                  = 1 << 6
	ComponentFlagNetwork               = 1 << 7
	ComponentFlagPhysics               = 1 << 8
	ComponentFlagRender                = 1 << 9
	ComponentFlagThirdPersonController = 1 << 10
	ComponentFlagTransform             = 1 << 11
)

type Component interface {
	AddToComponentContainer(container *ComponentContainer)
	ComponentFlag() int
}

type ComponentContainer struct {
	bitflags int

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

func NewComponentContainer(components ...Component) *ComponentContainer {
	container := &ComponentContainer{}
	for _, component := range components {
		component.AddToComponentContainer(container)
	}
	return container
}

func (cc *ComponentContainer) SetBitFlag(b int) {
	cc.bitflags |= b
}

func (cc *ComponentContainer) MatchBitFlags(b int) bool {
	return b&cc.bitflags == b
}
