package model

import (
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/lib/modelspec"
)

// KeyFrame contains a "Pose" which is the mapping from joint name to
// the transformtations that should be applied to the joint for this pose
type KeyFrame struct {
	Pose  map[int]*JointTransform
	Start time.Duration
}

// JointTransform represents the joint-space transformations that should be
// applied to the joint for the KeyFrame it is associated with.
type JointTransform struct {
	Translation mgl32.Vec3
	Rotation    mgl32.Quat
}

type Animation struct {
	RootJoint *Joint
	JointMap  map[int]*Joint

	KeyFrames []*KeyFrame
	Length    time.Duration
}

// this is used by the server. ideally we don't even need to have joints set up
// not sure how important it is for the server to sim the animation state.
func NewJointOnlyAnimation(spec *modelspec.ModelSpecification) *Animation {
	joint := JointSpecToJoint(spec.Root)
	jointMap := map[int]*Joint{}
	return &Animation{
		RootJoint: joint,
		JointMap:  getJointMap(joint, jointMap),

		Length:    spec.Animation.Length,
		KeyFrames: copyKeyFrames(spec.Animation),
	}
}

func NewAnimation(spec *modelspec.ModelSpecification) *Animation {
	configureJointVertexAttributes(spec.TriIndices, spec.TriIndicesStride, spec.JointWeightsSourceData, spec.JointIDs, spec.JointWeights, maxWeights)
	joint := JointSpecToJoint(spec.Root)
	jointMap := map[int]*Joint{}

	return &Animation{
		RootJoint: joint,
		JointMap:  getJointMap(joint, jointMap),

		Length:    spec.Animation.Length,
		KeyFrames: copyKeyFrames(spec.Animation),
	}
}

// lays out the vertex atrributes for:
// 4 - joint indices    vec3
// 5 - joint weights    vec3

// regardless of the number of joints affecting the joint
// we always pad out the full number of joints and weights with zeros
func configureJointVertexAttributes(triIndices []int, triIndicesStride int, jointWeightsSourceData []float32, jointIDs [][]int, jointWeights [][]int, maxWeights int) {
	jointIDsAttribute := []int32{}
	jointWeightsAttribute := []float32{}

	for i := 0; i < len(triIndices); i += triIndicesStride {
		vertexIndex := triIndices[i]

		for j, id := range jointIDs[vertexIndex] {
			if id == 38 || id == 50 || id == 51 {
				jointWeights[vertexIndex][j] = 0
			}
		}

		ids, weights := FillWeights(jointIDs[vertexIndex], jointWeights[vertexIndex], jointWeightsSourceData, maxWeights)

		for _, id := range ids {
			jointIDsAttribute = append(jointIDsAttribute, int32(id))
		}

		jointWeightsAttribute = append(jointWeightsAttribute, weights...)
	}

	var vboJointIDs uint32
	gl.GenBuffers(1, &vboJointIDs)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboJointIDs)
	gl.BufferData(gl.ARRAY_BUFFER, len(jointIDsAttribute)*4, gl.Ptr(jointIDsAttribute), gl.STATIC_DRAW)
	gl.VertexAttribIPointer(4, int32(maxWeights), gl.INT, int32(maxWeights)*4, nil)
	gl.EnableVertexAttribArray(4)

	var vboJointWeights uint32
	gl.GenBuffers(1, &vboJointWeights)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboJointWeights)
	gl.BufferData(gl.ARRAY_BUFFER, len(jointWeightsAttribute)*4, gl.Ptr(jointWeightsAttribute), gl.STATIC_DRAW)
	gl.VertexAttribPointer(5, int32(maxWeights), gl.FLOAT, false, int32(maxWeights)*4, nil)
	gl.EnableVertexAttribArray(5)
}

func copyKeyFrames(spec *modelspec.AnimationSpec) []*KeyFrame {
	var keyFrames []*KeyFrame

	for _, kf := range spec.KeyFrames {
		keyFrame := &KeyFrame{
			Start: kf.Start,
			Pose:  map[int]*JointTransform{},
		}

		for idx, jointTransform := range kf.Pose {
			keyFrame.Pose[idx] = &JointTransform{
				Translation: jointTransform.Translation,
				Rotation:    jointTransform.Rotation,
			}
		}

		keyFrames = append(keyFrames, keyFrame)
	}
	return keyFrames
}
