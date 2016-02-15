package worker

import (
	"time"

	"github.com/kkevinchou/ant/behavior"
	"github.com/kkevinchou/ant/lib/math/vector"
)

// Embed each of the interfaces that this behavior uses
type WorkerI interface {
	AddItem(interface{})
	DropItem(interface{})
	HasItem(interface{})
	SetTarget(vector.Vector)
	Position() vector.Vector
}

func NewBT(worker WorkerI) *BehaviorTree {
	return &BehaviorTree{
		root:  CreateWorkerBT(worker),
		state: behavior.AiState{BlackBoard: map[string]string{}},
	}
}

func CreateWorkerBT(worker WorkerI) behavior.Node {
	seq := behavior.NewSequence()
	seq.AddChild(&behavior.LocateItem{Entity: worker})
	seq.AddChild(&behavior.Move{Entity: worker})

	seq2 := behavior.NewSequence()
	seq2.AddChild(behavior.NewAiStateModifier(func(s behavior.AiState) { s.BlackBoard["output"] = "406_450" }))
	seq2.AddChild(&behavior.Move{Entity: worker})

	final := behavior.NewSequence()
	final.AddChild(seq)
	final.AddChild(seq2)
	return final
}

type BehaviorTree struct {
	root  behavior.Node
	state behavior.AiState
}

func (b *BehaviorTree) Tick(delta time.Duration) {
	tickResult := b.root.Tick(b.state, delta)
	if tickResult == behavior.SUCCESS {
		b.root.Reset()
	}
}
