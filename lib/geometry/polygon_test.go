package geometry

import "testing"

func defaultPolygon() *Polygon {
	points := []Point{
		Point{X: 0, Y: 0, Z: 0},
		Point{X: 0, Y: 0, Z: 6},
		Point{X: 6, Y: 0, Z: 6},
		Point{X: 6, Y: 0, Z: 0},
	}
	return NewPolygon(points)
}

func assert(t *testing.T, errorMessage string, actual, expected bool) {
	if expected != actual {
		t.Fatalf("%s - Expected [%v] but got [%v]", errorMessage, expected, actual)
	}
}

// We consider the borders to be inclusive, may be subject to change in the future
func TestContainsPointOnBorder(t *testing.T) {
	polygon := defaultPolygon()

	assert(t, "Should contain point when it lies on the top border", polygon.ContainsPoint(Point{X: 3, Y: 0, Z: 0}), true)
	assert(t, "Should contain point when it lies on the bottom border", polygon.ContainsPoint(Point{X: 3, Y: 0, Z: 6}), true)
	assert(t, "Should contain point when it lies on the left border", polygon.ContainsPoint(Point{X: 0, Y: 0, Z: 3}), true)
	assert(t, "Should contain point when it lies on the right border", polygon.ContainsPoint(Point{X: 0, Y: 0, Z: 3}), true)

	assert(t, "Should contain point when it overlaps a point", polygon.ContainsPoint(Point{X: 0, Y: 0, Z: 0}), true)
	assert(t, "Should contain point when it overlaps a point", polygon.ContainsPoint(Point{X: 0, Y: 0, Z: 6}), true)
	assert(t, "Should contain point when it overlaps a point", polygon.ContainsPoint(Point{X: 6, Y: 0, Z: 6}), true)
	assert(t, "Should contain point when it overlaps a point", polygon.ContainsPoint(Point{X: 6, Y: 0, Z: 0}), true)
}

func TestContainsPointWithinBorder(t *testing.T) {
	polygon := defaultPolygon()
	assert(t, "Should contain point when it lies within the borders", polygon.ContainsPoint(Point{X: 3, Y: 0, Z: 3}), true)
	assert(t, "Should contain point when it lies within the borders", polygon.ContainsPoint(Point{X: 1, Y: 0, Z: 2}), true)
	assert(t, "Should contain point when it lies within the borders", polygon.ContainsPoint(Point{X: 5, Y: 0, Z: 4}), true)
	assert(t, "Should contain point when it lies within the borders", polygon.ContainsPoint(Point{X: 3, Y: 0, Z: 1}), true)
}

func TestDoesNotContainPoint(t *testing.T) {
	polygon := defaultPolygon()
	assert(t, "Should contain point when it lies above the polygon", polygon.ContainsPoint(Point{X: 3, Y: 0, Z: -10}), false)
	assert(t, "Should contain point when it lies below the polygon", polygon.ContainsPoint(Point{X: 3, Y: 0, Z: 10}), false)
	assert(t, "Should contain point when it lies left of the polygon", polygon.ContainsPoint(Point{X: 10, Y: 0, Z: 3}), false)
	assert(t, "Should contain point when it lies right of the polygon", polygon.ContainsPoint(Point{X: -10, Y: 0, Z: 3}), false)
}
