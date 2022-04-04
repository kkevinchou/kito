package modelspec

import "github.com/go-gl/mathgl/mgl32"

type PBRMetallicRoughness struct {
	BaseColorTexture *uint32
	BaseColorFactor  mgl32.Vec4
	MetalicFactor    float32
	RoughnessFactor  float32
}

type PBRMaterial struct {
	PBRMetallicRoughness *PBRMetallicRoughness
}

type Vertex struct {
	Position mgl32.Vec3
	Normal   mgl32.Vec3
	Texture  mgl32.Vec2

	JointIDs     []int
	JointWeights []float32
}

type MeshChunkSpecification struct {
	VertexIndices []uint32
	// the unique vertices in the mesh chunk. VertexIndices details
	// how the unique vertices are arranged to construct the mesh
	UniqueVertices []Vertex
	Vertices       []Vertex
	// PBR
	PBRMaterial *PBRMaterial
}
type MeshSpecification struct {
	MeshChunks []*MeshChunkSpecification

	// // Geometry
	// // VertexAttributeIndices defines indices that are lookups for individual vertex properties
	// // VertexAttributesStride defines how many contiguous indices within VertexAttributeIndices define a vertex
	// //		Example arrangement:
	// //		[
	// //			triangle1PositionIndex, triangle1NormalIndex, triangle1TextureCoordIndex,
	// //			triangle2PositionIndex, triangle2NormalIndex, triangle2TextureCoordIndex,
	// //		]
	// // VertexAttributesStride would have a value of 3 here
	// // Three contiguous vertices define a triangle, after which the next triangle is defined
	// VertexAttributesStride int

	// PositionSourceData []mgl32.Vec3
	// NormalSourceData   []mgl32.Vec3
	// TextureSourceData  []mgl32.Vec2

	// // sorted by vertex order
	// JointIDs     [][]int
	// JointWeights [][]float32
}

// ModelSpecification is the output of any parsed model files (e.g. from Blender, Maya, etc)
// and acts a the blueprint for the model that contains all the associated vertex and
// animation data. This struct should be agnostic to the 3D modelling tool that produced the data.
type ModelSpecification struct {
	Meshes []*MeshSpecification

	// Joint Hierarchy
	RootJoint *JointSpec

	// Animations
	Animations map[string]*AnimationSpec
}
