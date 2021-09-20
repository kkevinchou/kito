package model

import "github.com/go-gl/mathgl/mgl32"

type JointSpec struct {
	ID            int
	Name          string
	BindTransform mgl32.Mat4
	Children      []*JointSpec
}

type ModelSpec struct {
	// Geometry
	TriIndices       []int // vertex indices in triangle order. Each triplet defines a face
	TriIndicesStride int

	PositionSourceData []mgl32.Vec3
	NormalSourceData   []mgl32.Vec3
	ColorSourceData    []mgl32.Vec3
	TextureSourceData  []mgl32.Vec2

	// Controllers

	// sorted by vertex order
	JointIDs     [][]int
	JointWeights [][]int

	JointsSourceData       []string // index is the joint id, the string value is the name
	JointWeightsSourceData []float32

	// Joint Hierarchy

	Root *JointSpec

	// Animations

	Animation *Animation
}
