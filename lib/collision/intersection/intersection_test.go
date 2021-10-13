package intersection_test

import (
	"testing"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/lib/collision"
	"github.com/kkevinchou/kito/lib/collision/collider"
	"github.com/kkevinchou/kito/lib/collision/intersection"
)

// Lines are perpendicular with separation along origin
func TestClosestPointsLineVSLine(t *testing.T) {
	points, distance := intersection.ClosestPointsLineVSLine(
		collider.Line{
			P1: mgl64.Vec3{-1, 1, 0},
			P2: mgl64.Vec3{1, 1, 0},
		},
		collider.Line{
			P1: mgl64.Vec3{0, -1, -1},
			P2: mgl64.Vec3{0, -1, 1},
		},
	)

	expectedPoint0 := mgl64.Vec3{0, 1, 0}
	if points[0] != expectedPoint0 {
		t.Errorf("expected first point to be %v but got %v\n", expectedPoint0, points[0])
	}
	expectedPoint1 := mgl64.Vec3{0, -1, 0}
	if points[1] != expectedPoint1 {
		t.Errorf("expected second point to be %v but got %v\n", expectedPoint1, points[1])
	}
	var expectedDistance float64 = 2
	if distance != expectedDistance {
		t.Errorf("expected distance to be %f but got %f\n", expectedDistance, distance)
	}
}

// closest point is P1 one of the line segment and the center of the triangle
func TestClosestPointsLineVsTriangle(t *testing.T) {
	line := collider.Line{
		P1: mgl64.Vec3{0, 1, -0.5},
		P2: mgl64.Vec3{0, 2, -1},
	}
	trianglePoints := []mgl64.Vec3{
		{0, 0, 0},
		{1, 0, -1},
		{-1, 0, -1},
	}

	triangle := collider.NewTriangle(trianglePoints)
	points, distance := intersection.ClosestPointsLineVSTriangle(line, triangle)

	expectedPoint0 := mgl64.Vec3{0, 1, -0.5}
	if points[0] != expectedPoint0 {
		t.Errorf("expected first point to be %v but got %v\n", expectedPoint0, points[0])
	}
	expectedPoint1 := mgl64.Vec3{0, 0, -0.5}
	if points[1] != expectedPoint1 {
		t.Errorf("expected second point to be %v but got %v\n", expectedPoint1, points[1])
	}
	var expectedDistance float64 = 1
	if distance != expectedDistance {
		t.Errorf("expected distance to be %f but got %f\n", expectedDistance, distance)
	}
}

// closest point is P2 one of the line segment and the center of the triangle
func TestClosestPointsLineVsTriangle2(t *testing.T) {
	line := collider.Line{
		P1: mgl64.Vec3{0, 2, -1},
		P2: mgl64.Vec3{0, 1, -0.5},
	}
	trianglePoints := []mgl64.Vec3{
		{0, 0, 0},
		{1, 0, -1},
		{-1, 0, -1},
	}

	triangle := collider.NewTriangle(trianglePoints)
	points, distance := intersection.ClosestPointsLineVSTriangle(line, triangle)

	expectedPoint0 := mgl64.Vec3{0, 1, -0.5}
	if points[0] != expectedPoint0 {
		t.Errorf("expected first point to be %v but got %v\n", expectedPoint0, points[0])
	}
	expectedPoint1 := mgl64.Vec3{0, 0, -0.5}
	if points[1] != expectedPoint1 {
		t.Errorf("expected second point to be %v but got %v\n", expectedPoint1, points[1])
	}
	var expectedDistance float64 = 1
	if distance != expectedDistance {
		t.Errorf("expected distance to be %f but got %f\n", expectedDistance, distance)
	}
}

func TestTriangleEdgeClosestToLine(t *testing.T) {
	line := collider.Line{
		P1: mgl64.Vec3{0, -1, -2},
		P2: mgl64.Vec3{0, 1, -2},
	}
	trianglePoints := []mgl64.Vec3{
		{0, 0, 0},
		{1, 0, -1},
		{-1, 0, -1},
	}

	triangle := collider.NewTriangle(trianglePoints)
	points, distance := intersection.ClosestPointsLineVSTriangle(line, triangle)

	expectedPoint0 := mgl64.Vec3{0, 0, -2}
	if points[0] != expectedPoint0 {
		t.Errorf("expected first point to be %v but got %v\n", expectedPoint0, points[0])
	}
	expectedPoint1 := mgl64.Vec3{0, 0, -1}
	if points[1] != expectedPoint1 {
		t.Errorf("expected second point to be %v but got %v\n", expectedPoint1, points[1])
	}
	var expectedDistance float64 = 1
	if distance != expectedDistance {
		t.Errorf("expected distance to be %f but got %f\n", expectedDistance, distance)
	}
}

func TestCheckCollisionCapsuleTriangle(t *testing.T) {
	capsule := collider.Capsule{
		Radius: 1,
		Top:    mgl64.Vec3{0, 10, -0.5},
		Bottom: mgl64.Vec3{0, 0.5, -0.5},
	}

	trianglePoints := []mgl64.Vec3{
		{0, 0, 0},
		{1, 0, -1},
		{-1, 0, -1},
	}

	triangle := collider.NewTriangle(trianglePoints)
	manifold := collision.CheckCollisionCapsuleTriangle(capsule, triangle)

	if len(manifold.Contacts) != 1 {
		t.Errorf("expected %d contact but got %d instead", 1, len(manifold.Contacts))
	}

	contact := manifold.Contacts[0]

	expectedNormal := mgl64.Vec3{0, 1, 0}
	if contact.Normal != expectedNormal {
		t.Errorf("expected contact normal to be %v but got %v", expectedNormal, contact.Normal)
	}

	if contact.SeparatingDistance != 0.5 {
		t.Errorf("expected separating distance to be %f but got %f", 0.5, contact.SeparatingDistance)
	}

	expectedContactPoint := mgl64.Vec3{0, 0, -0.5}
	if contact.Point != expectedContactPoint {
		t.Errorf("expected contact point to be %v but got %v", expectedContactPoint, contact.Point)
	}
}
