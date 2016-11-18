package worker

import (
	"time"

	"github.com/kkevinchou/ant/behavior"
	"github.com/kkevinchou/ant/behavior/connectors"
	"github.com/kkevinchou/ant/lib/math/vector"
)

func NewBT(worker Worker) *BehaviorTree {
	return &BehaviorTree{
		root:  CreateWorkerBT(worker),
		state: behavior.AIState{BlackBoard: map[string]string{}},
	}
}

func CreateWorkerBT(worker Worker) behavior.Node {
	memory := connectors.NewMemory()
	seq := behavior.NewSequence()

	seq.AddChild(&behavior.RandomItem{})
	seq.AddChild(memory.Set("item"))
	seq.AddChild(&connectors.Position{})
	seq.AddChild(&behavior.Move{Entity: worker})
	seq.AddChild(memory.Get("item"))
	seq.AddChild(&behavior.PickupItem{Entity: worker})

	seq2 := behavior.NewSequence()
	seq2.AddChild(&behavior.Value{Value: vector.Vector{X: 406, Y: 350}})
	seq2.AddChild(&behavior.Move{Entity: worker})

	final := behavior.NewSequence()
	final.AddChild(seq)
	final.AddChild(seq2)
	return final
}

type BehaviorTree struct {
	root  behavior.Node
	state behavior.AIState
}

func (b *BehaviorTree) Tick(delta time.Duration) {
	_, tickResult := b.root.Tick(nil, b.state, delta)
	if tickResult == behavior.SUCCESS {
		b.root.Reset()
	}
}
