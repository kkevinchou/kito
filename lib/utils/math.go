package utils

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

func Vec3IsZero(v mgl64.Vec3) bool {
	return v[0] == 0 && v[1] == 0 && v[2] == 0
}

func Cross2D(v1, v2 mgl64.Vec3) float64 {
	return (v1.X() * v2.Z()) - (v1.Z() * v2.X())
}

func Vec3F64ToF32(v mgl64.Vec3) mgl32.Vec3 {
	var result mgl32.Vec3
	for i := 0; i < len(v); i++ {
		result[i] = float32(v[i])
	}
	return result
}

func Vec4F64To4F32(v mgl64.Vec4) mgl32.Vec4 {
	return mgl32.Vec4{float32(v.X()), float32(v.Y()), float32(v.Z()), float32(v.W())}
}

func QuatF64ToF32(q mgl64.Quat) mgl32.Quat {
	return mgl32.Quat{W: float32(q.W), V: Vec3F64ToF32(q.V)}
}

func Mat4F64ToF32(m mgl64.Mat4) mgl32.Mat4 {
	var result mgl32.Mat4
	for i := 0; i < len(m); i++ {
		result[i] = float32(m[i])
	}
	return result
}

func SameSign(a, b int) bool {
	return (a > 0 && b > 0) || (a < 0 && b < 0) || (a == 0 && b == 0)
}

func QuatLookAt(eye, center, up mgl64.Vec3) mgl64.Quat {
	// http://www.opengl-tutorial.org/intermediate-tutorials/tutorial-17-quaternions/#I_need_an_equivalent_of_gluLookAt__How_do_I_orient_an_object_towards_a_point__
	// https://bitbucket.org/sinbad/ogre/src/d2ef494c4a2f5d6e2f0f17d3bfb9fd936d5423bb/OgreMain/src/OgreCamera.cpp?at=default#cl-161

	direction := center.Sub(eye).Normalize()

	// Find the rotation between the front of the object (that we assume towards Z-,
	// but this depends on your model) and the desired direction
	rotDir := mgl64.QuatBetweenVectors(mgl64.Vec3{0, 0, -1}, direction)

	// Because of the 1st rotation, the up is probably completely screwed up.
	// Find the rotation between the "up" of the rotated object, and the desired up
	newNormal := rotDir.Rotate(mgl64.Vec3{0, 1, 0})
	rotUp := mgl64.QuatBetweenVectors(newNormal, up)

	rotTarget := rotUp.Mul(rotDir) // remember, in reverse order.
	return rotTarget
}
