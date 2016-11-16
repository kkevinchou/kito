package behavior

import "time"

type Sequence struct {
	children []Node
	cache    *NodeCache
}

func (s *Sequence) AddChild(node Node) {
	s.children = append(s.children, node)
}

func NewSequence() *Sequence {
	return &Sequence{children: []Node{}, cache: NewNodeCache()}
}

func (s *Sequence) Tick(state AIState, delta time.Duration) Status {
	for _, child := range s.children {
		if s.cache.Contains(child) {
			continue
		}

		status := child.Tick(state, delta)
		if status == SUCCESS || status == FAILURE {
			s.cache.Add(child, status)
		}

		if status != SUCCESS {
			return status
		}
	}

	return SUCCESS
}

func (s *Sequence) Reset() {
	s.cache.Reset()
	for _, child := range s.children {
		child.Reset()
	}
}
