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

func (s *Sequence) Tick(input interface{}, state AIState, delta time.Duration) (interface{}, Status) {
	var status Status

	for _, child := range s.children {
		if s.cache.Contains(child) {
			continue
		}

		input, status = child.Tick(input, state, delta)
		if status == SUCCESS || status == FAILURE {
			s.cache.Add(child, status)
		}

		if status != SUCCESS {
			return nil, status
		}
	}

	return nil, SUCCESS
}

func (s *Sequence) Reset() {
	s.cache.Reset()
	for _, child := range s.children {
		child.Reset()
	}
}
