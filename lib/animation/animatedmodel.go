package animation

import (
	"fmt"
	"sort"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Mesh struct {
	vao         uint32
	vertexCount int
}

type AnimatedModel struct {
	RootJoint *Joint
	Mesh      *Mesh
}

func JointSpecToJoint(js *JointSpecification) *Joint {
	j := NewJoint(js.ID, js.Name, js.BindTransform)
	for _, c := range js.Children {
		j.Children = append(j.Children, JointSpecToJoint(c))
	}
	return j
}

func printHierarchy(j *Joint, level int) {
	indentation := ""
	for i := 0; i < level; i++ {
		indentation += "    "
	}
	fmt.Println(indentation + j.Name + fmt.Sprintf(" %d", j.ID))
	for _, c := range j.Children {
		printHierarchy(c, level+1)
	}
}

func NewAnimatedModel(c *ModelSpecification, maxJoints, maxWeights int) *AnimatedModel {
	joint := JointSpecToJoint(c.Root)
	printHierarchy(joint, 0)

	mesh := NewMesh(c, maxWeights)
	return &AnimatedModel{
		Mesh:      mesh,
		RootJoint: joint,
	}
}

//  1420 faces
//  4260 individual vertices (vertices can be counted multiple times)
//  740 distinct vertices
func NewMesh(c *ModelSpecification, maxWeights int) *Mesh {
	// maxJoints := 50

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	vertexAttributes, totalAttributeSize := constructGeometryVertexAttributes(
		c.TriIndices,
		c.PositionSourceData,
		c.NormalSourceData,
		c.ColorSourceData,
		c.TextureSourceData,
	)
	vertexCount := len(vertexAttributes) / totalAttributeSize
	configureGeometryVertexAttributes(vertexAttributes, totalAttributeSize)
	jointIDsAttribute := configureJointVertexAttributes(c.TriIndices, c.JointWeightsSourceData, c.JointIDs, c.JointWeights, maxWeights)
	configureIndexBuffer(vertexCount, vertexAttributes, jointIDsAttribute)
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
		// color := colorSourceData[triIndices[i+3]]

		color := mgl32.Vec3{0, 0, 0}

		vertexAttributes = append(vertexAttributes, position.X(), position.Y(), position.Z())
		vertexAttributes = append(vertexAttributes, normal.X(), normal.Y(), normal.Z())
		vertexAttributes = append(vertexAttributes, texture.X(), texture.Y())
		vertexAttributes = append(vertexAttributes, color.X(), color.Y(), color.Z())
	}

	return vertexAttributes, totalAttributeSize
}

func configureGeometryVertexAttributes(vertexAttributes []float32, totalAttributeSize int) {
	var vbo uint32
	gl.GenBuffers(1, &vbo)

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

func configureJointVertexAttributes(triIndices []int, jointWeightsSourceData []float32, jointIDs [][]int, jointWeights [][]int, maxWeights int) []int32 {
	jointIDsAttribute := []int32{}
	jointWeightsAttribute := []float32{}

	for i := 0; i < len(triIndices); i += 4 {
		vertexIndex := triIndices[i]

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

	return jointIDsAttribute
}

func configureIndexBuffer(vertexCount int, vertexAttributes []float32, jointIDsAttributes []int32) {
	var ebo uint32
	gl.GenBuffers(1, &ebo)

	indices := []uint32{}
	for i := 0; i < vertexCount; i++ {
		indices = append(indices, uint32(i))
	}

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
}

// if we exceed maxWeights, drop the weakest weights and normalize
// if we're below maxWeights, fill in dummy weights so we always have "maxWeights" number of weights
func FillWeights(jointIDs []int, weights []int, jointWeightsSourceData []float32, maxWeights int) ([]int, []float32) {
	j := []int{}
	w := []float32{}

	if len(jointIDs) <= maxWeights {
		j = append(j, jointIDs...)
		for _, weightIndex := range weights {
			w = append(w, jointWeightsSourceData[weightIndex])
		}
		// fill in empty jointIDs and weights
		for i := 0; i < maxWeights-len(jointIDs); i++ {
			j = append(j, 0)
			w = append(w, 0)
		}
	} else if len(jointIDs) > maxWeights {
		jointWeights := []JointWeight{}
		for i := range jointIDs {
			jointWeights = append(jointWeights, JointWeight{JointID: jointIDs[i], Weight: jointWeightsSourceData[weights[i]]})
		}
		sort.Sort(sort.Reverse(byWeights(jointWeights)))

		// take top 3 weights
		jointWeights = jointWeights[:maxWeights]
		NormalizeWeights(jointWeights)
		for _, jw := range jointWeights {
			j = append(j, jw.JointID)
			w = append(w, jw.Weight)
		}
	}

	return j, w
}

func NormalizeWeights(jointWeights []JointWeight) {
	var totalWeight float32
	for _, jw := range jointWeights {
		totalWeight += jw.Weight
	}

	for i := range jointWeights {
		jointWeights[i].Weight /= totalWeight
	}
}

type byWeights []JointWeight

type JointWeight struct {
	JointID int
	Weight  float32
}

func (s byWeights) Len() int {
	return len(s)
}
func (s byWeights) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byWeights) Less(i, j int) bool {
	return s[i].Weight < s[j].Weight
}
