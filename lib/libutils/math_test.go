package libutils_test

import (
	"testing"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/lib/libutils"
)

func TestLookAt(t *testing.T) {
	eye := mgl64.Vec3{0, 0, 0}
	up := mgl64.Vec3{0, 1, 0}

	dir := mgl64.Vec3{0, 1, 0}
	libutils.QuatLookAt(eye, dir, up)

	dir = mgl64.Vec3{0, 0, -1}
	libutils.QuatLookAt(eye, dir, up)

	dir = mgl64.Vec3{-1, 0, 0}
	libutils.QuatLookAt(eye, dir, up)

	dir = mgl64.Vec3{1, 0, 0}
	libutils.QuatLookAt(eye, dir, up)

	dir = mgl64.Vec3{0, -1, 0}
	libutils.QuatLookAt(eye, dir, up)

	dir = mgl64.Vec3{0, 0, 1}
	libutils.QuatLookAt(eye, dir, up)
}

func TestQuat(t *testing.T) {
	v1 := mgl64.Vec3{0, 0, -1}
	v2 := mgl64.Vec3{0, 1, 0}

	mgl64.QuatBetweenVectors(v1, v2)
}
