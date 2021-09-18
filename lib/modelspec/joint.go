package modelspec

import "github.com/go-gl/mathgl/mgl32"

type JointSpecification struct {
	ID            int
	Name          string
	BindTransform mgl32.Mat4
	Children      []*JointSpecification
}
