package components

type EntityTypeBitFlag int

const (
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
	ComponentFlagAI                    = 1 << 12
	ComponentFlagNotepad               = 1 << 13
	ComponentFlagHealth                = 1 << 14
	ComponentFlagLootDropper           = 1 << 15
	ComponentFlagLoot                  = 1 << 16
	ComponentFlagInventory             = 1 << 17
)

type Component interface {
	AddToComponentContainer(container *ComponentContainer)
	ComponentFlag() int
}

type ComponentContainer struct {
	bitflags int

	AIComponent                    *AIComponent
	AnimationComponent             *AnimationComponent
	RenderComponent                *RenderComponent
	TransformComponent             *TransformComponent
	PhysicsComponent               *PhysicsComponent
	TopDownViewComponent           *TopDownViewComponent
	ThirdPersonControllerComponent *ThirdPersonControllerComponent
	CameraComponent                *CameraComponent
	NetworkComponent               *NetworkComponent
	MeshComponent                  *MeshComponent
	ColliderComponent              *ColliderComponent
	ControlComponent               *ControlComponent
	NotepadComponent               *NotepadComponent
	HealthComponent                *HealthComponent
	LootDropperComponent           *LootDropperComponent
	LootComponent                  *LootComponent
	InventoryComponent             *InventoryComponent
}

func NewComponentContainer(components ...Component) *ComponentContainer {
	container := &ComponentContainer{}
	for _, component := range components {
		component.AddToComponentContainer(container)
		container.bitflags |= component.ComponentFlag()
	}
	return container
}

func (cc *ComponentContainer) SetBitFlag(b int) {
	cc.bitflags |= b
}

func (cc *ComponentContainer) MatchBitFlags(b int) bool {
	return b&cc.bitflags == b
}
