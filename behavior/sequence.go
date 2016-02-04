package behavior

import "time"

type Sequence struct {
	children []Node
}

func (s *Sequence) Tick(state AIState, delta time.Duration) Status {
	for _, child := range s.children {
		status := child.Tick(state, delta)
		if status != SUCCESS {
			return status
		}
	}

	return SUCCESS
}
