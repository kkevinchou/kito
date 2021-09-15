package render_test

import (
	"fmt"
	"testing"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kkevinchou/kito/systems/render"
)

func TestCalc(t *testing.T) {
	// points := []mgl64.Vec3{
	// 	mgl64.Vec3{-1, 1, -1},
	// 	mgl64.Vec3{1, 1, -1},
	// 	mgl64.Vec3{1, -1, -1},
	// 	mgl64.Vec3{-1, -1, -1},
	// }
	points := []mgl64.Vec3{
		mgl64.Vec3{0, 0, -1},
	}
	lightOrientation := mgl64.QuatRotate(mgl64.DegToRad(-90), mgl64.Vec3{1, 0, 0})
	fmt.Println("TEST VECTOR", lightOrientation.Mat4().Mul4x1(mgl64.Vec3{0, 0, -1}.Vec4(1)).Vec3())

	position, projMatrix := render.ComputeDirectionalLightProps(lightOrientation.Mat4(), points)
	fmt.Println(position)
	fmt.Println(projMatrix)
	t.Fail()
}
