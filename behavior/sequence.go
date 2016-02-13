package behavior

import "time"

type Sequence struct {
	children []Node
}

func (s *Sequence) AddChild(node Node) {
	s.children = append(s.children, node)
}

func NewSequence() *Sequence {
	return &Sequence{children: []Node{}}
}

func (s *Sequence) Tick(state AiState, delta time.Duration) Status {
	for _, child := range s.children {
		status := child.Tick(state, delta)
		if status != SUCCESS {
			return status
		}
	}

	return SUCCESS
}
