package geometry

import (
	"testing"

	"github.com/kkevinchou/ant/geometry"
)

func defaultPolygon() *geometry.Polygon {
	points := []geometry.Point{
		geometry.Point{0, 0},
		geometry.Point{0, 6},
		geometry.Point{6, 6},
		geometry.Point{6, 0},
	}
	return geometry.NewPolygon(points)
}

func assert(t *testing.T, errorMessage string, actual, expected bool) {
	if expected != actual {
		t.Fatalf("%s - Expected [%v] but got [%v]", errorMessage, expected, actual)
	}
}

// We consider the borders to be inclusive, may be subject to change in the future
func TestContainsPointOnBorder(t *testing.T) {
	polygon := defaultPolygon()

	assert(t, "Should contain point when it lies on the top border", polygon.ContainsPoint(geometry.Point{3, 0}), true)
	assert(t, "Should contain point when it lies on the bottom border", polygon.ContainsPoint(geometry.Point{3, 6}), true)
	assert(t, "Should contain point when it lies on the left border", polygon.ContainsPoint(geometry.Point{0, 3}), true)
	assert(t, "Should contain point when it lies on the right border", polygon.ContainsPoint(geometry.Point{0, 3}), true)

	assert(t, "Should contain point when it overlaps a point", polygon.ContainsPoint(geometry.Point{0, 0}), true)
	assert(t, "Should contain point when it overlaps a point", polygon.ContainsPoint(geometry.Point{0, 6}), true)
	assert(t, "Should contain point when it overlaps a point", polygon.ContainsPoint(geometry.Point{6, 6}), true)
	assert(t, "Should contain point when it overlaps a point", polygon.ContainsPoint(geometry.Point{6, 0}), true)
}

func TestContainsPointWithinBorder(t *testing.T) {
	polygon := defaultPolygon()
	assert(t, "Should contain point when it lies within the borders", polygon.ContainsPoint(geometry.Point{3, 3}), true)
	assert(t, "Should contain point when it lies within the borders", polygon.ContainsPoint(geometry.Point{1, 2}), true)
	assert(t, "Should contain point when it lies within the borders", polygon.ContainsPoint(geometry.Point{5, 4}), true)
	assert(t, "Should contain point when it lies within the borders", polygon.ContainsPoint(geometry.Point{3, 1}), true)
}

func TestDoesNotContainPoint(t *testing.T) {
	polygon := defaultPolygon()
	assert(t, "Should contain point when it lies above the polygon", polygon.ContainsPoint(geometry.Point{3, -10}), false)
	assert(t, "Should contain point when it lies below the polygon", polygon.ContainsPoint(geometry.Point{3, 10}), false)
	assert(t, "Should contain point when it lies left of the polygon", polygon.ContainsPoint(geometry.Point{10, 3}), false)
	assert(t, "Should contain point when it lies right of the polygon", polygon.ContainsPoint(geometry.Point{-10, 3}), false)
}
