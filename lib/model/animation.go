package model

import (
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/lib/modelspec"
)

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

	triIndices       []int
	triIndicesStride int
	jointIDs         [][]int
	jointWeights     [][]float32
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

		triIndices:       spec.TriIndices,
		triIndicesStride: spec.TriIndicesStride,
		jointIDs:         spec.JointIDs,
		jointWeights:     spec.JointWeights,
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

		ids, weights := FillWeights(a.jointIDs[vertexIndex], a.jointWeights[vertexIndex])

		for _, id := range ids {
			jointIDsAttribute = append(jointIDsAttribute, int32(id))
		}
		jointWeightsAttribute = append(jointWeightsAttribute, weights...)
	}

	var vboJointIDs uint32
	gl.GenBuffers(1, &vboJointIDs)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboJointIDs)
	gl.BufferData(gl.ARRAY_BUFFER, len(jointIDsAttribute)*4, gl.Ptr(jointIDsAttribute), gl.STATIC_DRAW)
	gl.VertexAttribIPointer(4, int32(settings.AnimationMaxJointWeights), gl.INT, int32(settings.AnimationMaxJointWeights)*4, nil)
	gl.EnableVertexAttribArray(4)

	var vboJointWeights uint32
	gl.GenBuffers(1, &vboJointWeights)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboJointWeights)
	gl.BufferData(gl.ARRAY_BUFFER, len(jointWeightsAttribute)*4, gl.Ptr(jointWeightsAttribute), gl.STATIC_DRAW)
	gl.VertexAttribPointer(5, int32(settings.AnimationMaxJointWeights), gl.FLOAT, false, int32(settings.AnimationMaxJointWeights)*4, nil)
	gl.EnableVertexAttribArray(5)
}
