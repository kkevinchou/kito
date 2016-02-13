package behavior

import "time"

type Selector struct {
	children []Node
}

func (s *Selector) Tick(state AiState, delta time.Duration) Status {
	for _, child := range s.children {
		status := child.Tick(state, delta)
		if status == SUCCESS {
			return SUCCESS
		}
	}

	return FAILURE
}
