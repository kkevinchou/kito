package model

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
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
	name          string
	rootJoint     *modelspec.JointSpec
	animationSpec *modelspec.AnimationSpec

	// vertexAttributeIndices []int
	// vertexAttributesStride int
	// jointIDs               [][]int
	// jointWeights           [][]float32
}

func (a *Animation) Name() string {
	return a.name
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

func NewAnimations(spec *modelspec.ModelSpecification) map[string]*Animation {
	// TODO: handle animations for multiple meshes
	// mesh := spec.Meshes[0]
	// vertexAttributeIndices := mesh.VertexAttributeIndices
	// vertexAttributesStride := mesh.VertexAttributesStride
	// jointIDs := mesh.JointIDs
	// jointWeights := mesh.JointWeights

	animations := map[string]*Animation{}
	for name, animation := range spec.Animations {
		animations[name] = &Animation{
			name:          name,
			animationSpec: animation,
			rootJoint:     spec.RootJoint,

			// vertexAttributeIndices: vertexAttributeIndices,
			// vertexAttributesStride: vertexAttributesStride,
			// jointIDs:               jointIDs,
			// jointWeights:           jointWeights,
		}
	}

	return animations
}

// lays out the vertex atrributes for:
// 3 - joint indices    vec3
// 4 - joint weights    vec3

// regardless of the number of joints affecting the joint
// we always pad out the full number of joints and weights with zeros
// func (a *Animation) BindVertexAttributes() {
// 	jointIDsAttribute := []int32{}
// 	jointWeightsAttribute := []float32{}

// 	// TODO: it seems like we currently duplicate vertex data in a vertex array rather than using an EBO store indices to the vertices
// 	// this is probably less efficient? we store redundant data for the same vertex as if it were a new vertex. e.g. duplicated positions and
// 	// joint weights

// 	for i := 0; i < len(a.vertexAttributeIndices); i += a.vertexAttributesStride {
// 		// vertex index is the index of the position which we assume to be the first property
// 		vertexIndex := a.vertexAttributeIndices[i]

// 		ids, weights := FillWeights(a.jointIDs[vertexIndex], a.jointWeights[vertexIndex])
// 		for _, id := range ids {
// 			jointIDsAttribute = append(jointIDsAttribute, int32(id))
// 		}
// 		jointWeightsAttribute = append(jointWeightsAttribute, weights...)
// 	}

// 	var vboJointIDs uint32
// 	gl.GenBuffers(1, &vboJointIDs)
// 	gl.BindBuffer(gl.ARRAY_BUFFER, vboJointIDs)
// 	gl.BufferData(gl.ARRAY_BUFFER, len(jointIDsAttribute)*4, gl.Ptr(jointIDsAttribute), gl.STATIC_DRAW)
// 	gl.VertexAttribIPointer(3, int32(settings.AnimationMaxJointWeights), gl.INT, int32(settings.AnimationMaxJointWeights)*4, nil)
// 	gl.EnableVertexAttribArray(3)

// 	var vboJointWeights uint32
// 	gl.GenBuffers(1, &vboJointWeights)
// 	gl.BindBuffer(gl.ARRAY_BUFFER, vboJointWeights)
// 	gl.BufferData(gl.ARRAY_BUFFER, len(jointWeightsAttribute)*4, gl.Ptr(jointWeightsAttribute), gl.STATIC_DRAW)
// 	gl.VertexAttribPointer(4, int32(settings.AnimationMaxJointWeights), gl.FLOAT, false, int32(settings.AnimationMaxJointWeights)*4, nil)
// 	gl.EnableVertexAttribArray(4)
// }
