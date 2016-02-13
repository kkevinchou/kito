package behavior

import "time"

type Status int

const (
	RUNNING Status = iota
	SUCCESS Status = iota
	FAILURE Status = iota
)

type Node interface {
	Tick(AiState, time.Duration) Status
}

type BehaviorTree interface {
	Tick(time.Duration)
}

type AiState struct {
	BlackBoard map[string]string
}

type NodeCache struct {
	cache map[Node]Status
}

func NewNodeCache() *NodeCache {
	return &NodeCache{cache: map[Node]Status{}}
}
