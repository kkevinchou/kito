package model

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/kkevinchou/kito/kito/settings"
	"github.com/kkevinchou/kito/lib/modelspec"
)

type MeshChunk struct {
	vao  uint32
	spec *modelspec.MeshChunkSpecification
}

type Mesh struct {
	meshChunks []*MeshChunk
}

func (m *MeshChunk) VAO() uint32 {
	return m.vao
}

func (m *MeshChunk) Vertices() []modelspec.Vertex {
	return m.spec.Vertices
}

func (m *MeshChunk) VertexCount() int {
	return len(m.spec.Vertices)
}

func (m *MeshChunk) PBRMaterial() *modelspec.PBRMaterial {
	return m.spec.PBRMaterial
}

func NewMesh(spec *modelspec.MeshSpecification) *Mesh {
	var meshChunks []*MeshChunk
	for _, mc := range spec.MeshChunks {
		meshChunks = append(meshChunks, &MeshChunk{
			spec: mc,
		})
	}
	return &Mesh{
		meshChunks: meshChunks,
		// meshChunks: spec.,
	}
}

func (m *Mesh) MeshChunks() []*MeshChunk {
	return m.meshChunks
}

func (m *Mesh) Prepare() {
	for _, chunk := range m.meshChunks {
		chunk.Prepare()
	}
}

func (m *MeshChunk) Prepare() {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	m.vao = vao

	// vertexIndices := []uint32{}
	// for i := 0; i < len(m.spec.VertexIndices); i++ {
	// 	vertexIndices = append(vertexIndices, uint32(i))
	// }

	var vertexAttributes []float32
	var jointIDsAttribute []int32
	var jointWeightsAttribute []float32

	// fmt.Println("-------------")
	// fmt.Println(len(m.spec.VertexIndices))
	// fmt.Println(m.spec.VertexIndices)
	// fmt.Println(len(m.spec.Vertices))

	for _, vertex := range m.spec.UniqueVertices {
		position := vertex.Position
		normal := vertex.Normal
		texture := vertex.Texture
		jointIDs := vertex.JointIDs
		jointWeights := vertex.JointWeights

		vertexAttributes = append(vertexAttributes,
			position.X(), position.Y(), position.Z(),
			normal.X(), normal.Y(), normal.X(),
			texture.X(), texture.Y(),
		)

		ids, weights := FillWeights(jointIDs, jointWeights)
		for _, id := range ids {
			jointIDsAttribute = append(jointIDsAttribute, int32(id))
		}
		jointWeightsAttribute = append(jointWeightsAttribute, weights...)
	}

	totalAttributeSize := len(vertexAttributes) / len(m.spec.UniqueVertices)

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

	// TODO: it seems like we currently duplicate vertex data in a vertex array rather than using an EBO store indices to the vertices
	// this is probably less efficient? we store redundant data for the same vertex as if it were a new vertex. e.g. duplicated positions and
	// joint weights

	var vboJointIDs uint32
	gl.GenBuffers(1, &vboJointIDs)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboJointIDs)
	gl.BufferData(gl.ARRAY_BUFFER, len(jointIDsAttribute)*4, gl.Ptr(jointIDsAttribute), gl.STATIC_DRAW)
	gl.VertexAttribIPointer(3, int32(settings.AnimationMaxJointWeights), gl.INT, int32(settings.AnimationMaxJointWeights)*4, nil)
	gl.EnableVertexAttribArray(3)

	var vboJointWeights uint32
	gl.GenBuffers(1, &vboJointWeights)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboJointWeights)
	gl.BufferData(gl.ARRAY_BUFFER, len(jointWeightsAttribute)*4, gl.Ptr(jointWeightsAttribute), gl.STATIC_DRAW)
	gl.VertexAttribPointer(4, int32(settings.AnimationMaxJointWeights), gl.FLOAT, false, int32(settings.AnimationMaxJointWeights)*4, nil)
	gl.EnableVertexAttribArray(4)

	// // TODO: see if this works, but this should probably be instead
	// // using m.spec.VertexIndices that can be reused and repointed to already seen
	// // vertices to save on memory. Right now we can potentially readd duplicated
	// // vertex attributes to `vertexAttributes`.
	// indices := []uint32{}
	// for i := 0; i < len(vertexAttributes)/8; i++ {
	// 	indices = append(indices, uint32(i))
	// }

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(m.spec.VertexIndices)*4, gl.Ptr(m.spec.VertexIndices), gl.STATIC_DRAW)
}

// func constructMeshVertexAttributes(
// 	spec *modelspec.ModelSpecification,
// ) ([]float32, int, []mgl64.Vec3) {
// 	var vertices []mgl64.Vec3
// 	vertexAttributes := []float32{}

// 	for _, mesh := range spec.Meshes {
// 		positionSourceData := mesh.PositionSourceData
// 		normalSourceData := mesh.NormalSourceData
// 		textureSourceData := mesh.TextureSourceData
// 		vertexAttributeIndices := mesh.VertexAttributeIndices
// 		vertexAttributesStride := mesh.VertexAttributesStride

// 		if mesh.VertexAttributesStride <= 0 {
// 			panic(fmt.Sprintf("unexpected stride value %d", mesh.VertexAttributesStride))
// 		}

// 		// triIndicies format: position, normal, texture, color
// 		for i := 0; i < len(vertexAttributeIndices); i += vertexAttributesStride {
// 			// TODO: we are assuming this ordering of position, normal, texture but this is not
// 			// necessarily the case. it depends on the <input> elements are ordered in the collada file
// 			position := positionSourceData[vertexAttributeIndices[i]]
// 			normal := normalSourceData[vertexAttributeIndices[i+1]]
// 			texture := textureSourceData[vertexAttributeIndices[i+2]]

// 			vertexAttributes = append(vertexAttributes, position.X(), position.Y(), position.Z())
// 			vertexAttributes = append(vertexAttributes, normal.X(), normal.Y(), normal.Z())
// 			vertexAttributes = append(vertexAttributes, texture.X(), texture.Y())

// 			vertices = append(vertices, mgl64.Vec3{float64(position.X()), float64(position.Y()), float64(position.Z())})
// 		}
// 	}

// 	totalAttributeSize := len(spec.Meshes[0].PositionSourceData[0]) + len(spec.Meshes[0].NormalSourceData[0]) + len(spec.Meshes[0].TextureSourceData[0])
// 	return vertexAttributes, totalAttributeSize, vertices
// }

// lays out the vertex atrributes for:
// 0 - position         vec3
// 1 - normal           vec3
// 2 - texture coord    vec2
// 3 - color            vec3
// func (m *Mesh) BindVertexAttributes() {
// 	var vbo uint32
// 	gl.GenBuffers(1, &vbo)
// 	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
// 	gl.BufferData(gl.ARRAY_BUFFER, len(m.vertexAttributes)*4, gl.Ptr(m.vertexAttributes), gl.STATIC_DRAW)

// 	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(m.totalAttributeSize)*4, nil)
// 	gl.EnableVertexAttribArray(0)

// 	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(m.totalAttributeSize)*4, gl.PtrOffset(3*4))
// 	gl.EnableVertexAttribArray(1)

// 	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, int32(m.totalAttributeSize)*4, gl.PtrOffset(6*4))
// 	gl.EnableVertexAttribArray(2)
// }
