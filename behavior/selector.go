package behavior

import "time"

type Selector struct {
	children []Node
	cache    NodeCache
}

func (s *Selector) Tick(state AIState, delta time.Duration) Status {
	for _, child := range s.children {
		status := child.Tick(state, delta)
		if status == SUCCESS {
			return SUCCESS
		}
	}

	return FAILURE
}

func (s *Selector) Reset() {
	for _, child := range s.children {
		child.Reset()
	}
}
