package components

import (
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/lib/behavior"
)

type AIComponent struct {
	// behaviorTree behavior.BehaviorTree
	LastUpdate  time.Time
	MovementDir mgl64.Quat
	Velocity    mgl64.Vec3
}

func NewAIComponent(behaviorTree behavior.BehaviorTree) *AIComponent {
	return &AIComponent{
		LastUpdate:  time.Now(),
		MovementDir: mgl64.QuatRotate(0, mgl64.Vec3{0, 1, 0}),
		// behaviorTree: behaviorTree,
	}
}

func (c *AIComponent) AddToComponentContainer(container *ComponentContainer) {
	container.AIComponent = c
}

func (c *AIComponent) ComponentFlag() int {
	return ComponentFlagAI
}
