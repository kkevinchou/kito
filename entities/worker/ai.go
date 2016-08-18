package worker

import (
	"time"

	"github.com/kkevinchou/ant/behavior"
	"github.com/kkevinchou/ant/behavior/worker"
)

type AIComponent struct {
	bt behavior.BehaviorTree
}

func NewAIComponent(entity worker.WorkerI) *AIComponent {
	return &AIComponent{
		bt: worker.NewBT(entity),
	}
}

func (c *AIComponent) Update(delta time.Duration) {
	c.bt.Tick(delta)
}
