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

func Vec3F64ToVec3F32(v mgl64.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{float32(v.X()), float32(v.Y()), float32(v.Z())}
}

func QuatF64ToQuatF32(q mgl64.Quat) mgl32.Quat {
	return mgl32.Quat{W: float32(q.W), V: Vec3F64ToVec3F32(q.V)}
}
