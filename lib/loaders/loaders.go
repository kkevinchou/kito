package loaders

import "github.com/go-gl/mathgl/mgl32"

type Joint struct {
	ID            int
	Name          string
	BindTransform mgl32.Mat4
	Children      []*Joint
}

type ModelSpecification struct {
	// Geometry
	TriIndices []int // vertex indices in triangle order. Each triplet defines a face

	PositionSourceData []mgl32.Vec3
	NormalSourceData   []mgl32.Vec3
	ColorSourceData    []mgl32.Vec3
	TextureSourceData  []mgl32.Vec2

	// Controllers

	JointIDs     [][]int
	JointWeights [][]int

	JointsSourceData       []string // index is the joint id, the string value is the name
	JointWeightsSourceData []float32

	// Joint Hierarchy

	Root *Joint
}
