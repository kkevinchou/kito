package behavior

import (
	"time"
)

type Status int

const (
	RUNNING Status = iota
	SUCCESS Status = iota
	FAILURE Status = iota
)

type Node interface {
	Tick(AIState, time.Duration) Status
}

type BehaviorTree interface {
	Tick(time.Duration)
}

type AIState struct {
	BlackBoard map[string]string
}
