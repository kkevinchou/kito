package behavior

import (
	"time"

	"github.com/kkevinchou/ant/directory"
	"github.com/kkevinchou/ant/lib/geometry"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/kkevinchou/ant/lib/pathing"
)

type Mover interface {
	Position() vector.Vector
	SetTarget(target vector.Vector)
}

type Move struct {
	Entity    Mover
	path      []pathing.Node
	pathIndex int
}

func (m *Move) Tick(input interface{}, state AIState, delta time.Duration) (interface{}, Status) {
	if m.path == nil {
		pathManager := directory.GetDirectory().PathManager()
		position := m.Entity.Position()

		var target vector.Vector
		var ok bool

		if target, ok = input.(vector.Vector); !ok {
			return nil, FAILURE
		}

		path := pathManager.FindPath(
			geometry.Point{X: position.X, Y: position.Y},
			geometry.Point{X: target.X, Y: target.Y},
		)

		if path != nil {
			m.path = path
			m.pathIndex = 1
			m.Entity.SetTarget(m.path[m.pathIndex].Vector())
		}
	}

	if m.path == nil {
		return nil, FAILURE
	}

	if m.pathIndex == len(m.path) {
		return nil, SUCCESS
	}

	if m.Entity.Position().Sub(m.path[m.pathIndex].Vector()).Length() <= 2 {
		m.pathIndex += 1
		if m.pathIndex < len(m.path) {
			m.Entity.SetTarget(m.path[m.pathIndex].Vector())
		}
	}

	if m.pathIndex == len(m.path) {
		return nil, SUCCESS
	} else {
		return nil, RUNNING
	}
}

func (m *Move) Reset() {
	m.path = nil
	m.pathIndex = 0
}
