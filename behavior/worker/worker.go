package worker

import (
	"time"

	"github.com/kkevinchou/ant/behavior"
)

// Embed each of the interfaces that this behavior uses
type WorkerI interface {
	AddItem(interface{})
	DropItem(interface{})
	HasItem(interface{})
}

func NewBT(worker WorkerI) *BehaviorTree {
	return &BehaviorTree{
		root:  CreateWorkerBT(worker),
		state: behavior.AIState{},
	}
}

func CreateWorkerBT(worker WorkerI) behavior.Node {
	return &behavior.AddItem{Entity: worker}
}

type BehaviorTree struct {
	root  behavior.Node
	state behavior.AIState
}

func (b *BehaviorTree) Tick(delta time.Duration) {
	b.root.Tick(b.state, delta)
}
