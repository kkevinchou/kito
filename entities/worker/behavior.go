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
	seq.AddChild(memory.Set("item_position"))
	seq.AddChild(&behavior.Move{Entity: worker})
	seq.AddChild(memory.Get("item"))
	seq.AddChild(&behavior.PickupItem{Entity: worker})

	seq2 := behavior.NewSequence()
	seq2.AddChild(&behavior.Value{Value: vector.Vector3{X: 0, Y: 0, Z: 0}})
	seq2.AddChild(&behavior.Move{Entity: worker})
	seq2.AddChild(memory.Get("item"))
	seq2.AddChild(&behavior.DropItem{Entity: worker})

	seq3 := behavior.NewSequence()
	seq3.AddChild(memory.Get("item_position"))
	seq3.AddChild(&behavior.Move{Entity: worker})

	seq4 := behavior.NewSequence()
	seq4.AddChild(memory.Get("item"))
	seq4.AddChild(&connectors.Position{})
	seq4.AddChild(&behavior.Move{Entity: worker})
	seq4.AddChild(memory.Get("item"))
	seq4.AddChild(&behavior.PickupItem{Entity: worker})

	seq5 := behavior.NewSequence()
	seq5.AddChild(memory.Get("item_position"))
	seq5.AddChild(&behavior.Move{Entity: worker})
	seq5.AddChild(memory.Get("item"))
	seq5.AddChild(&behavior.DropItem{Entity: worker})

	seq6 := behavior.NewSequence()
	seq6.AddChild(&behavior.Value{Value: vector.Vector3{X: 0, Y: 0, Z: 0}})
	seq6.AddChild(&behavior.Move{Entity: worker})

	findFood := behavior.NewSequence()
	findFood.AddChild(seq)
	findFood.AddChild(seq2)
	findFood.AddChild(seq3)
	findFood.AddChild(seq4)
	findFood.AddChild(seq5)
	findFood.AddChild(seq6)

	// final := behavior.NewSelector()
	// final.AddChild(findFood)
	// final.AddChild(seq4)
	return findFood
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
