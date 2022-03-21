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

type MeshSpecification struct {
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
	TextureSourceData  []mgl32.Vec2

	// sorted by vertex order
	JointIDs     [][]int
	JointWeights [][]float32
}

func (m *ModelSpecification) ConvertTexCoordsFromGLTFToOpenGL() {
	for _, mesh := range m.Meshes {
		mesh.ConvertTexCoordsFromGLTFToOpenGL()
	}
}

func (m *MeshSpecification) ConvertTexCoordsFromGLTFToOpenGL() {
	for i, v := range m.TextureSourceData {
		m.TextureSourceData[i] = mgl32.Vec2{v.X(), 1 - v.Y()}
	}
}

// ModelSpecification is the output of any parsed model files (e.g. from Blender, Maya, etc)
// and acts a the blueprint for the model that contains all the associated vertex and
// animation data. This struct should be agnostic to the 3D modelling tool that produced the data.
type ModelSpecification struct {
	Meshes []*MeshSpecification

	// Effects
	EffectSpecData *EffectSpec

	// Joint Hierarchy
	RootJoint *JointSpec

	// Animations
	Animations map[string]*AnimationSpec
}
