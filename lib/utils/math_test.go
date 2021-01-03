package utils_test

import (
	"fmt"
	"testing"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/lib/utils"
)

func TestLookAt(t *testing.T) {
	eye := mgl64.Vec3{0, 0, 0}
	up := mgl64.Vec3{0, 1, 0}
	forward := mgl64.Vec3{0, 0, -1}

	dir := mgl64.Vec3{0, 1, 0}
	lookAt := utils.QuatLookAt(eye, dir, up)
	fmt.Println("up", lookAt.Rotate(forward))

	dir = mgl64.Vec3{0, 0, -1}
	lookAt = utils.QuatLookAt(eye, dir, up)
	fmt.Println("forward", lookAt.Rotate(forward))

	dir = mgl64.Vec3{-1, 0, 0}
	lookAt = utils.QuatLookAt(eye, dir, up)
	fmt.Println("left", lookAt.Rotate(forward))

	dir = mgl64.Vec3{1, 0, 0}
	lookAt = utils.QuatLookAt(eye, dir, up)
	fmt.Println("right", lookAt.Rotate(forward))

	dir = mgl64.Vec3{0, -1, 0}
	lookAt = utils.QuatLookAt(eye, dir, up)
	fmt.Println("down", lookAt.Rotate(forward))

	dir = mgl64.Vec3{0, 0, 1}
	lookAt = utils.QuatLookAt(eye, dir, up)
	fmt.Println("back", lookAt.Rotate(forward))

	t.Fail()
}

func TestQuat(t *testing.T) {
	v1 := mgl64.Vec3{0, 0, -1}
	v2 := mgl64.Vec3{0, 1, 0}

	q := mgl64.QuatBetweenVectors(v1, v2)
	fmt.Println(q)
	fmt.Println(q.Rotate(v1))
	t.Fail()
}
