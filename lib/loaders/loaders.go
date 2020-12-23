package loaders

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

type Joint struct {
	ID            int
	Name          string
	BindTransform mgl32.Mat4
	Children      []*Joint
}

type Animation struct {
	KeyFrames []*KeyFrame
	Length    time.Duration
}

// KeyFrame contains a "Pose" which is the mapping from joint name to
// the transformtations that should be applied to the joint for this pose
type KeyFrame struct {
	Pose  map[string]*JointTransform
	Start time.Duration
}

// JointTransform represents the joint-space transformations that should be
// applied to the joint for the KeyFrame it is associated with.
type JointTransform struct {
	Position mgl32.Vec3
	Rotation mgl32.Quat
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
