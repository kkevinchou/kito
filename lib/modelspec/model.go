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
	// VertexAttributeIndices defines indices that are lookups for individual vertex properties
	// VertexAttributesStride defines how many contiguous indices within VertexAttributeIndices define a vertex
	//		Example arrangement:
	//		[
	//			triangle1PositionIndex, triangle1NormalIndex, triangle1TextureCoordIndex,
	//			triangle2PositionIndex, triangle2NormalIndex, triangle2TextureCoordIndex,
	//		]
	// VertexAttributesStride would have a value of 3 here
	// Three contiguous vertices define a triangle, after which the next triangle is defined
	VertexAttributeIndices []int
	VertexAttributesStride int

	PositionSourceData []mgl32.Vec3
	NormalSourceData   []mgl32.Vec3
	ColorSourceData    []mgl32.Vec3
	TextureSourceData  []mgl32.Vec2

	// sorted by vertex order
	JointIDs     [][]int
	JointWeights [][]float32

	// Effects
	EffectSpecData *EffectSpec

	// Joint Hierarchy
	RootJoint *JointSpec

	// Animations
	Animation *AnimationSpec
}
