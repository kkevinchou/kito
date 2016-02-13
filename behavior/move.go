package behavior

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kkevinchou/ant/lib/geometry"
	"github.com/kkevinchou/ant/lib/math/vector"
	"github.com/kkevinchou/ant/pathing"
	"github.com/kkevinchou/ant/systems"
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

func (m *Move) Tick(state AiState, delta time.Duration) Status {
	if m.path == nil {
		pathManager := systems.GetDirectory().PathManager()
		position := m.Entity.Position()
		targetStr := strings.Split(state.BlackBoard["output"], "_")
		targetX, _ := strconv.ParseFloat(targetStr[0], 64)
		targetY, _ := strconv.ParseFloat(targetStr[1], 64)
		fmt.Println(targetStr)

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
		return SUCCESS
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

	return SUCCESS
}
