package behavior

import (
	"strconv"
	"strings"
	"time"

	"github.com/kkevinchou/ant/directory"
	"github.com/kkevinchou/ant/lib/geometry"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/kkevinchou/ant/lib/pathing"
)

type MoveI interface {
	Position() vector.Vector
	SetTarget(target vector.Vector)
}

type Move struct {
	Entity    MoveI
	path      []pathing.Node
	pathIndex int
}

func (m *Move) Tick(state AIState, delta time.Duration) Status {
	if m.path == nil {
		pathManager := directory.GetDirectory().PathManager()
		position := m.Entity.Position()
		targetStr := strings.Split(state.BlackBoard["output"], "_")
		targetX, _ := strconv.ParseFloat(targetStr[0], 64)
		targetY, _ := strconv.ParseFloat(targetStr[1], 64)

		path := pathManager.FindPath(
			geometry.Point{X: position.X, Y: position.Y},
			geometry.Point{X: targetX, Y: targetY},
		)

		if path != nil {
			m.path = path
			m.pathIndex = 1
			m.Entity.SetTarget(m.path[m.pathIndex].Vector())
		}
	}

	if m.path == nil {
		return FAILURE
	}

	if m.pathIndex == len(m.path) {
		return SUCCESS
	}

	if m.Entity.Position().Sub(m.path[m.pathIndex].Vector()).Length() <= 2 {
		m.pathIndex += 1
		if m.pathIndex < len(m.path) {
			m.Entity.SetTarget(m.path[m.pathIndex].Vector())
		}
	}

	if m.pathIndex == len(m.path) {
		return SUCCESS
	} else {
		return RUNNING
	}
}

func (m *Move) Reset() {
	m.path = nil
	m.pathIndex = 0
}
