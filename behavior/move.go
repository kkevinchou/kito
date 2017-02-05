package behavior

import (
	"time"

	"github.com/kkevinchou/kito/directory"
	"github.com/kkevinchou/kito/interfaces"
	libbehavior "github.com/kkevinchou/kito/lib/behavior"
	"github.com/kkevinchou/kito/lib/geometry"
	"github.com/kkevinchou/kito/lib/math/vector"
	"github.com/kkevinchou/kito/logger"
)

type Mover interface {
	interfaces.Positionable
	SetTarget(target vector.Vector3)
}

type Move struct {
	Entity    Mover
	path      []geometry.Point
	pathIndex int
}

func (m *Move) Tick(input interface{}, state libbehavior.AIState, delta time.Duration) (interface{}, libbehavior.Status) {
	logger.Debug("Move - ENTER")

	if m.path == nil {
		var target vector.Vector3
		var ok bool

		if target, ok = input.(vector.Vector3); !ok {
			logger.Debug("Move - FAIL")
			return nil, libbehavior.FAILURE
		}

		pathManager := directory.GetDirectory().PathManager()
		position := m.Entity.Position()

		path := pathManager.FindPath(
			geometry.Point{X: position.X, Y: position.Y, Z: position.Z},
			geometry.Point{X: target.X, Y: target.Y, Z: target.Z},
		)

		if path != nil {
			m.path = path
			m.pathIndex = 1
			m.Entity.SetTarget(m.path[m.pathIndex].Vector3())
		}
	}

	if m.path == nil {
		logger.Debug("Move - FAIL")
		return nil, libbehavior.FAILURE
	}

	if m.pathIndex == len(m.path) {
		logger.Debug("Move - SUCCESS")
		return nil, libbehavior.SUCCESS
	}

	if m.Entity.Position().Sub(m.path[m.pathIndex].Vector3()).Length() <= 0.1 {
		m.pathIndex++
		if m.pathIndex < len(m.path) {
			m.Entity.SetTarget(m.path[m.pathIndex].Vector3())
		}
	}

	if m.pathIndex == len(m.path) {
		logger.Debug("Move - SUCCESS")
		return nil, libbehavior.SUCCESS
	}

	logger.Debug("Move - RUNNING")
	return nil, libbehavior.RUNNING
}

func (m *Move) Reset() {
	m.path = nil
	m.pathIndex = 0
}
