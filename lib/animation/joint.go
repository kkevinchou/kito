package animation

import "github.com/go-gl/mathgl/mgl32"

var (
	jointIDCounter = 0
)

type Joint struct {
	ID       int
	Name     string
	Children []*Joint

	LocalBindTransform   mgl32.Mat4
	InverseBindTransform mgl32.Mat4

	AnimationTransform mgl32.Mat4 // calculated by the animator
}

func NewJoint(name string, localBindTransform mgl32.Mat4) *Joint {
	joint := Joint{
		ID:                 jointIDCounter,
		Name:               name,
		Children:           []*Joint{},
		LocalBindTransform: localBindTransform,
	}
	jointIDCounter++

	return &joint
}

func (j *Joint) AddChild(child *Joint) {
	j.Children = append(j.Children, child)
}

func (j *Joint) GetInverseBindTransform() mgl32.Mat4 {
	return j.InverseBindTransform
}

func (j *Joint) CalculateInverseBindTransform(parentBindTransform mgl32.Mat4) {
	bindTransform := parentBindTransform.Mul4(j.LocalBindTransform) // model-space relative to the origin
	j.InverseBindTransform = bindTransform.Inv()
	for _, child := range j.Children {
		child.CalculateInverseBindTransform(bindTransform)
	}
}
