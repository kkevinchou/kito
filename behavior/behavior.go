package behavior

import (
	"time"
)

const (
	RUNNING = iota
	SUCCESS = iota
	FAILURE = iota
)

type Node interface {
	OnEnter()
	OnExit()
	Tick(delta time.Duration) int
}

type Sequence struct {
	children []Node
}

type Selector struct {
	children []Node
}

type BehaviorTree struct {
	root Node
}

type AIState struct {
	PreviousRunningNode Node
	BlackBoard          map[string]string
}

func (s *Sequence) Tick(delta time.Duration) int {
	for _, child := range s.children {
		status := child.Tick(delta)
		if status != SUCCESS {
			return status
		}
	}

	return SUCCESS
}

func (s *Selector) Tick(delta time.Duration) int {
	for _, child := range s.children {
		status := child.Tick(delta)
		if status == SUCCESS {
			return SUCCESS
		}
	}

	return FAILURE
}
