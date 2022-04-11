package components

import (
	"github.com/kkevinchou/kito/lib/behavior"
)

type AIComponent struct {
	behaviorTree behavior.BehaviorTree
}

func NewAIComponent(behaviorTree behavior.BehaviorTree) *AIComponent {
	return &AIComponent{
		behaviorTree: behaviorTree,
	}
}

func (c *AIComponent) AddToComponentContainer(container *ComponentContainer) {
	container.AIComponent = c
}

func (c *AIComponent) ComponentFlag() int {
	return ComponentFlagAI
}
