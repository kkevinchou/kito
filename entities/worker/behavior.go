package worker

import (
	"time"

	"github.com/kkevinchou/kito/behavior"
	libbehavior "github.com/kkevinchou/kito/lib/behavior"
)

func NewBT(worker Worker) *BehaviorTree {
	return &BehaviorTree{
		root:  CreateWorkerBT(worker),
		state: libbehavior.AIState{BlackBoard: map[string]string{}},
	}
}

func CreateWorkerBT(worker Worker) libbehavior.Node {
	memory := libbehavior.NewMemory()
	seq := libbehavior.NewSequence()

	// SWH: bug when adding first three nodes
	seq.AddChild(&libbehavior.Value{Value: worker})
	seq.AddChild(&libbehavior.Position{})
	seq.AddChild(memory.Set("entity_position"))
	seq.AddChild(&behavior.RandomItem{})
	seq.AddChild(memory.Set("item"))
	seq.AddChild(&libbehavior.Position{})
	seq.AddChild(memory.Set("item_position"))
	seq.AddChild(&behavior.Move{Entity: worker})
	seq.AddChild(memory.Get("item"))
	seq.AddChild(&behavior.PickupItem{Entity: worker})

	seq2 := libbehavior.NewSequence()
	seq2.AddChild(memory.Get("entity_position"))
	seq2.AddChild(&behavior.Move{Entity: worker})
	seq2.AddChild(memory.Get("item"))
	seq2.AddChild(&behavior.DropItem{Entity: worker})

	seq3 := libbehavior.NewSequence()
	seq3.AddChild(memory.Get("item_position"))
	seq3.AddChild(&behavior.Move{Entity: worker})

	seq4 := libbehavior.NewSequence()
	seq4.AddChild(memory.Get("item"))
	seq4.AddChild(&libbehavior.Position{})
	seq4.AddChild(&behavior.Move{Entity: worker})
	seq4.AddChild(memory.Get("item"))
	seq4.AddChild(&behavior.PickupItem{Entity: worker})

	seq5 := libbehavior.NewSequence()
	seq5.AddChild(memory.Get("item_position"))
	seq5.AddChild(&behavior.Move{Entity: worker})
	seq5.AddChild(memory.Get("item"))
	seq5.AddChild(&behavior.DropItem{Entity: worker})

	seq6 := libbehavior.NewSequence()
	seq6.AddChild(memory.Get("entity_position"))
	seq6.AddChild(&behavior.Move{Entity: worker})

	findFood := libbehavior.NewSequence()
	findFood.AddChild(seq)
	findFood.AddChild(seq2)
	findFood.AddChild(seq3)
	findFood.AddChild(seq4)
	findFood.AddChild(seq5)
	findFood.AddChild(seq6)

	final := libbehavior.NewSelector()
	final.AddChild(findFood)
	final.AddChild(seq4)
	return findFood
}

type BehaviorTree struct {
	root  libbehavior.Node
	state libbehavior.AIState
}

func (b *BehaviorTree) Tick(delta time.Duration) {
	_, tickResult := b.root.Tick(nil, b.state, delta)
	if tickResult == libbehavior.SUCCESS {
		b.root.Reset()
	}
}
