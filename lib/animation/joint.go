package animation

import (
	"github.com/go-gl/mathgl/mgl32"
)

var (
	jointIDCounter = 0
)

type Joint struct {
	ID       int
	Name     string
	Children []*Joint

	LocalBindTransform   mgl32.Mat4
	InverseBindTransform mgl32.Mat4
}

func NewJoint(id int, name string, localBindTransform mgl32.Mat4) *Joint {
	joint := Joint{
		ID:                 id,
		Name:               name,
		Children:           []*Joint{},
		LocalBindTransform: localBindTransform,
	}

	return &joint
}
