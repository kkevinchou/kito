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
	mesh := NewMesh(c, maxWeights)
	return &AnimatedModel{
		Mesh: mesh,
	}
}

// func NewMeshFromCollada(c *collada.Collada) *Mesh {
func NewMesh(c *collada.Collada, maxWeights int) *Mesh {
	// maxJoints := 50

	var vao uint32
	gl.GenVertexArrays(1, &vao)

	vertexAttributes, totalAttributeSize := constructGeometryVertexAttributes(
		c.TriIndices,
		c.PositionSourceData,
		c.NormalSourceData,
		c.ColorSourceData,
		c.TextureSourceData,
	)
	vertexCount := len(vertexAttributes) / totalAttributeSize
	configureGeometryVertexAttributes(vao, vertexAttributes, totalAttributeSize, vertexCount)
	configureIndexBuffer(vertexCount)
	configureJointVertexAttributes(vao, c.JointWeightsSourceData, c.JointIDs, c.JointWeights, vertexCount, maxWeights)

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

func constructGeometryVertexAttributes(
	triIndices []int,
	positionSourceData []mgl32.Vec3,
	normalSourceData []mgl32.Vec3,
	colorSourceData []mgl32.Vec3,
	textureSourceData []mgl32.Vec2,
) ([]float32, int) {
	vertexAttributes := []float32{}
	totalAttributeSize := len(positionSourceData[0]) + len(normalSourceData[0]) + len(textureSourceData[0]) + len(colorSourceData[0])
	// TODO: i'm still ordering vertex attributes by the face order, rather than keeping the original exported source order
	// this current way will repeat data since i explicity store data for every vertex, rather than using indicies for lookup
	// in the future, i should refactor this to store the data in source data order then use an index buffer for VAO creation

	// triIndicies format: position, normal, texture, color
	for i := 0; i < len(triIndices); i += 4 {
		position := positionSourceData[triIndices[i]]
		normal := normalSourceData[triIndices[i+1]]
		texture := textureSourceData[triIndices[i+2]]
		// color := colorSourceData[i]

		color := mgl32.Vec3{0, 0, 0}

		vertexAttributes = append(vertexAttributes, position.X(), position.Y(), position.Z())
		vertexAttributes = append(vertexAttributes, normal.X(), normal.Y(), normal.Z())
		vertexAttributes = append(vertexAttributes, texture.X(), texture.Y())
		vertexAttributes = append(vertexAttributes, color.X(), color.Y(), color.Z())
	}

	return vertexAttributes, totalAttributeSize
}

func configureIndexBuffer(vertexCount int) {
	var ebo uint32
	gl.GenBuffers(1, &ebo)

	indices := []uint32{}
	for i := 0; i < vertexCount; i++ {
		indices = append(indices, uint32(i))
	}

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
}

func configureGeometryVertexAttributes(vao uint32, vertexAttributes []float32, totalAttributeSize int, vertexCount int) {
	var vbo uint32
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexAttributes)*4, gl.Ptr(vertexAttributes), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(totalAttributeSize)*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(totalAttributeSize)*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, int32(totalAttributeSize)*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	gl.VertexAttribPointer(3, 3, gl.FLOAT, false, int32(totalAttributeSize)*4, gl.PtrOffset(8*4))
	gl.EnableVertexAttribArray(3)
}

func configureJointVertexAttributes(vao uint32, JointWeightsSourceData []float32, jointIDs [][]int, jointWeights [][]int, vertexCount, maxWeights int) {
	jointIDsAttribute := []int{}
	jointWeightsAttribute := []float32{}

	for i := 0; i < len(jointIDs); i++ {
		// fill in empty weights if we're overboard

		ids, weights := normalizeWeights(jointIDs[i], jointWeights[i], JointWeightsSourceData, maxWeights)
		for j := 0; j < maxWeights; j++ {
			jointIDsAttribute = append(jointIDsAttribute, ids...)
			jointWeightsAttribute = append(jointWeightsAttribute, weights...)
		}
	}

	var vboJointIDs, vboJointWeights uint32

	gl.GenBuffers(1, &vboJointIDs)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboJointIDs)
	gl.BufferData(gl.ARRAY_BUFFER, vertexCount*maxWeights*4, gl.Ptr(jointIDsAttribute), gl.STATIC_DRAW)
	gl.VertexAttribPointer(4, int32(maxWeights), gl.INT, false, int32(maxWeights)*4, nil)
	gl.EnableVertexAttribArray(4)

	gl.GenBuffers(1, &vboJointWeights)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboJointWeights)
	gl.BufferData(gl.ARRAY_BUFFER, vertexCount*maxWeights*4, gl.Ptr(jointWeightsAttribute), gl.STATIC_DRAW)
	gl.VertexAttribPointer(5, int32(maxWeights), gl.FLOAT, false, int32(maxWeights)*4, nil)
	gl.EnableVertexAttribArray(5)
}

func copySliceSliceInt(data [][]int) [][]int {
	result := [][]int{}
	for _, slice := range data {
		result = append(result, slice[:])
	}
	return result
}

// if we exceed maxWeights, drop the weakest weights and normalize
// if we're below maxWeights, fill in dummy weights so we always have "maxWeights" number of weights
func normalizeWeights(jointIDs []int, weights []int, JointWeightsSourceData []float32, maxWeights int) ([]int, []float32) {
	// TODO
	j := []int{}
	w := []float32{}

	for i := 0; i < maxWeights; i++ {
		j = append(j, 0)
		w = append(w, 0)
	}

	return j, w
}
