package modelspec

import "github.com/go-gl/mathgl/mgl32"

type EffectSpec struct { // todo(kevin): rename to MaterialSpec
	ID                     string
	ShaderElement          string
	EmissionColor          *mgl32.Vec3
	DiffuseColor           *mgl32.Vec3
	IndexOfRefractionFloat float32
	ReflectivityFloat      float32
	ReflectivityColor      *mgl32.Vec3
	ShininessFloat         float32
	TransparencyFloat      float32
	TransparencyColor      *mgl32.Vec3
}

// ModelSpecification is the output of any parsed model files (e.g. from Blender, Maya, etc)
// and acts a the blueprint for the model that contains all the associated vertex and
// animation data. This struct should be agnostic to the 3D modelling tool that produced the data.
type ModelSpecification struct {
	// Geometry
	TriIndices       []int // vertex indices in triangle order. Each triplet defines a face
	TriIndicesStride int

	PositionSourceData []mgl32.Vec3
	NormalSourceData   []mgl32.Vec3
	ColorSourceData    []mgl32.Vec3
	TextureSourceData  []mgl32.Vec2

	EffectSpecData *EffectSpec

	// Controllers

	// sorted by vertex order
	JointIDs     [][]int
	JointWeights [][]int

	JointsSourceData       []string // index is the joint id, the string value is the name
	JointWeightsSourceData []float32

	// Joint Hierarchy

	Root *JointSpecification

	// Animations

	Animation *AnimationSpec
}
