package connectors

import (
	"time"

	"github.com/kkevinchou/ant/behavior"
)

type Memory struct {
	data map[string]interface{}
}

func NewMemory() *Memory {
	return &Memory{
		data: map[string]interface{}{},
	}
}

func (m *Memory) Set(key string) *Set {
	return &Set{
		memory: m,
		key:    key,
	}
}

func (m *Memory) Get(key string) *Get {
	return &Get{
		memory: m,
		key:    key,
	}
}

type Set struct {
	memory *Memory
	key    string
}

func (s *Set) Tick(input interface{}, state behavior.AIState, delta time.Duration) (interface{}, behavior.Status) {
	s.memory.data[s.key] = input
	return input, behavior.SUCCESS
}

func (v *Set) Reset() {}

type Get struct {
	memory *Memory
	key    string
}

func (g *Get) Tick(input interface{}, state behavior.AIState, delta time.Duration) (interface{}, behavior.Status) {
	if value, ok := g.memory.data[g.key]; ok {
		return value, behavior.SUCCESS
	}
	return nil, behavior.FAILURE
}

func (v *Get) Reset() {}
