package singleton

import (
	"github.com/kkevinchou/kito/components"
	"github.com/kkevinchou/kito/lib/input"
)

type Singleton struct {
	// client connection
	*components.ConnectionComponent

	frameInput input.Input
}

func NewSingleton() *Singleton {
	return &Singleton{}
}

func (s *Singleton) SetInput(i input.Input) {
	s.frameInput = i
}

func (s *Singleton) GetInput() input.Input {
	return s.frameInput
}
