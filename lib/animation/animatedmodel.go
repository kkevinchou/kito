package animation

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/lib/loaders/collada"
)

type Mesh struct {
	vao         uint32
	vertexCount int
}

type AnimatedModel struct {
	RootJoint *Joint
	Mesh      *Mesh
}

func NewAnimatedModel(c *collada.Collada, maxJoints, maxWeights int) *AnimatedModel {
	mesh := NewMesh(c)
	return &AnimatedModel{
		Mesh: mesh,
	}
}

// func NewMeshFromCollada(c *collada.Collada) *Mesh {
func NewMesh(c *collada.Collada) *Mesh {
	// maxJoints := 50
	// maxWeights := 3

	vertexAttributes, totalAttributeSize := constructVertexAttributes(
		c.TriIndices,
		c.PositionSourceData,
		c.NormalSourceData,
		c.ColorSourceData,
		c.TextureSourceData,
	)
	vertexCount := len(vertexAttributes) / totalAttributeSize

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	configureGeometryVertexAttributes(vao, vertexAttributes, totalAttributeSize, vertexCount)

	return &Mesh{
		vao:         vao,
		vertexCount: vertexCount,
	}
}

func (m *Mesh) VAO() uint32 {
	return m.vao
}

func (m *Mesh) VertexCount() int {
	return m.vertexCount
}

func constructVertexAttributes(
	triIndices []int,
	positionSourceData []mgl32.Vec3,
	normalSourceData []mgl32.Vec3,
	colorSourceData []mgl32.Vec3,
	textureSourceData []mgl32.Vec2,
	// JointWeightsSourceData []float32,
	// jointIDs [][]int,
	// jointWeights [][]int,
	// maxWeights int,
	// vertexCount, dataSize int,
) ([]float32, int) {
	vertexAttributes := []float32{}
	totalAttributeSize := len(positionSourceData[0]) + len(normalSourceData[0]) + len(textureSourceData[0]) + len(colorSourceData[0])
	// TODO: i'm still ordering vertex attributes by the face order, rather than keeping the original exported source order
	// this current way will repeat data since i explicity store data for every vertex, rather than using indicies for lookup
	// in the future, i should refactor this to store the data in source data order then use an index buffer for VAO creation

	// triIndicies format:
	// position, normal, texture, color
	for i := 0; i < len(triIndices); i += 4 {
		position := positionSourceData[triIndices[i]]
		normal := normalSourceData[triIndices[i+1]]
		texture := textureSourceData[triIndices[i+2]]
		// color := colorSourceData[i]

		color := mgl32.Vec3{0, 0, 0}

		// vertJointIDs := jointIDs[i]
		// jointWeights := jointWeights[i]

		// // each weight is a joint ID and weight, hence we multiply by 2
		// totalAttributeSize := len(position) + len(normal) + len(color) + len(texture) + (maxWeights * 2)

		vertexAttributes = append(vertexAttributes, position.X(), position.Y(), position.Z())
		vertexAttributes = append(vertexAttributes, normal.X(), normal.Y(), normal.Z())
		vertexAttributes = append(vertexAttributes, texture.X(), texture.Y())
		vertexAttributes = append(vertexAttributes, color.X(), color.Y(), color.Z())

		// for j := 0; j < maxWeights; j++ {
		// 	vertexAttributes = append(vertexAttributes, vertJointIDs[j])
		// }
		// for j := 0; j < maxWeights; j++ {
		// 	vertexAttributes = append(vertexAttributes, JointWeightsSourceData[jointWeights[j]])
		// }
	}

	return vertexAttributes, totalAttributeSize
}

func configureGeometryVertexAttributes(vao uint32, vertexAttributes []float32, totalAttributeSize int, vertexCount int) {
	var vbo, ebo uint32
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexAttributes)*4, gl.Ptr(vertexAttributes), gl.STATIC_DRAW)

	indices := []uint32{}
	for i := 0; i < vertexCount; i++ {
		indices = append(indices, uint32(i))
	}

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(totalAttributeSize)*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(totalAttributeSize)*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, int32(totalAttributeSize)*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	gl.VertexAttribPointer(3, 3, gl.FLOAT, false, int32(totalAttributeSize)*4, gl.PtrOffset(8*4))
	gl.EnableVertexAttribArray(3)
}

// func configureJointVertexAttributes(vao uint32) {
// 	var vbo, ebo uint32
// 	gl.GenBuffers(1, &vbo)
// 	gl.GenBuffers(1, &ebo)

// 	gl.BindVertexArray(vao)
// 	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
// 	gl.BufferData(gl.ARRAY_BUFFER, len(vertexAttributes)*4, gl.Ptr(vertexAttributes), gl.STATIC_DRAW)

// 	indices := []uint32{}
// 	for i := 0; i < vertexCount; i++ {
// 		indices = append(indices, uint32(i))
// 	}

// 	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
// 	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

// 	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(totalAttributeSize)*4, nil)
// 	gl.EnableVertexAttribArray(0)

// 	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(totalAttributeSize)*4, gl.PtrOffset(3*4))
// 	gl.EnableVertexAttribArray(1)

// 	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, int32(totalAttributeSize)*4, gl.PtrOffset(6*4))
// 	gl.EnableVertexAttribArray(2)

// 	gl.VertexAttribPointer(3, 3, gl.FLOAT, false, int32(totalAttributeSize)*4, gl.PtrOffset(8*4))
// 	gl.EnableVertexAttribArray(3)
// }

func copySliceSliceInt(data [][]int) [][]int {
	result := [][]int{}
	for _, slice := range data {
		result = append(result, slice[:])
	}
	return result
}
