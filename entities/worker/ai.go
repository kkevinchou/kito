package worker

import (
	"time"

	"github.com/kkevinchou/kito/behavior"
)

type AIComponent struct {
	bt behavior.BehaviorTree
}

func NewAIComponent(entity Worker) *AIComponent {
	return &AIComponent{
		bt: NewBT(entity),
	}
}

func (c *AIComponent) Update(delta time.Duration) {
	c.bt.Tick(delta)
}
