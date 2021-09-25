package components

import (
	"time"

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

func (c *AIComponent) Update(delta time.Duration) {
	c.behaviorTree.Tick(delta)
}
