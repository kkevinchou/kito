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
	Scale       mgl32.Vec3
	Rotation    mgl32.Quat
}

type Animation struct {
	rootJoint     *modelspec.JointSpec
	animationSpec *modelspec.AnimationSpec

	triIndices             []int
	triIndicesStride       int
	jointWeightsSourceData []float32
	jointIDs               [][]int
	jointWeights           [][]int
	maxWeights             int
}

func (a *Animation) RootJoint() *modelspec.JointSpec {
	return a.rootJoint
}

func (a *Animation) KeyFrames() []*modelspec.KeyFrame {
	return a.animationSpec.KeyFrames
}

func (a *Animation) Length() time.Duration {
	return a.animationSpec.Length
}

func NewAnimation(spec *modelspec.ModelSpecification) *Animation {
	return &Animation{
		animationSpec: spec.Animation,
		rootJoint:     spec.RootJoint,

		triIndices:             spec.TriIndices,
		triIndicesStride:       spec.TriIndicesStride,
		jointWeightsSourceData: spec.JointWeightsSourceData,
		jointIDs:               spec.JointIDs,
		jointWeights:           spec.JointWeights,
		maxWeights:             maxWeights,
	}
}

// lays out the vertex atrributes for:
// 4 - joint indices    vec3
// 5 - joint weights    vec3

// regardless of the number of joints affecting the joint
// we always pad out the full number of joints and weights with zeros
func (a *Animation) BindVertexAttributes() {
	jointIDsAttribute := []int32{}
	jointWeightsAttribute := []float32{}

	for i := 0; i < len(a.triIndices); i += a.triIndicesStride {
		vertexIndex := a.triIndices[i]

		for j, id := range a.jointIDs[vertexIndex] {
			if id == 38 || id == 50 || id == 51 {
				a.jointWeights[vertexIndex][j] = 0
			}
		}

		ids, weights := FillWeights(a.jointIDs[vertexIndex], a.jointWeights[vertexIndex], a.jointWeightsSourceData, maxWeights)

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
				Scale:       jointTransform.Scale,
			}
		}

		keyFrames = append(keyFrames, keyFrame)
	}
	return keyFrames
}
