package model

import (
	"github.com/go-gl/mathgl/mgl32"
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
