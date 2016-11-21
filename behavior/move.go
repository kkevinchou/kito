package behavior

import (
	"time"

	"github.com/kkevinchou/ant/directory"
	"github.com/kkevinchou/ant/interfaces"
	"github.com/kkevinchou/ant/lib/geometry"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/kkevinchou/ant/lib/pathing"
)

type Mover interface {
	interfaces.Positionable
	SetTarget(target vector.Vector3)
}

type Move struct {
	Entity    Mover
	path      []pathing.Node
	pathIndex int
}

func (m *Move) Tick(input interface{}, state AIState, delta time.Duration) (interface{}, Status) {
	if m.path == nil {
		var target vector.Vector3
		var ok bool

		if target, ok = input.(vector.Vector3); !ok {
			return nil, FAILURE
		}

		pathManager := directory.GetDirectory().PathManager()
		position := m.Entity.Position()

		path := pathManager.FindPath(
			geometry.Point{X: position.X, Y: position.Z},
			geometry.Point{X: target.X, Y: target.Z},
		)

		if path != nil {
			m.path = path
			m.pathIndex = 1
			m.Entity.SetTarget(m.path[m.pathIndex].Vector3())
		}
	}

	if m.path == nil {
		return nil, FAILURE
	}

	if m.pathIndex == len(m.path) {
		return nil, SUCCESS
	}

	if m.Entity.Position().Sub(m.path[m.pathIndex].Vector3()).Length() <= 2 {
		m.pathIndex++
		if m.pathIndex < len(m.path) {
			m.Entity.SetTarget(m.path[m.pathIndex].Vector3())
		}
	}

	if m.pathIndex == len(m.path) {
		return nil, SUCCESS
	}
	return nil, RUNNING
}

func (m *Move) Reset() {
	m.path = nil
	m.pathIndex = 0
}
