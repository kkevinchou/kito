package animation

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

// This is static data that we read from the animation data files
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

func ParseCollada(filename string) {
	// colladaDoc, err := collada.LoadDocument("sample/model.dae")
	// if err != nil {
	// 	t.Error(err)
	// }
}
