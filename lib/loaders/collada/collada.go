package collada

import (
	"strings"

	"github.com/kkevinchou/kito/lib/loaders"
)

type TechniqueType string

func (t TechniqueType) String() string {
	return string(t)
}

func (t NodeType) String() string {
	return string(t)
}

type SemanticType string
type NodeType string

const (
	SemanticVertex   SemanticType = "VERTEX"
	SemanticNormal   SemanticType = "NORMAL"
	SemanticTexCoord SemanticType = "TEXCOORD"
	SemanticColor    SemanticType = "COLOR"
	SemanticPosition SemanticType = "POSITION"

	TechniqueJoint     TechniqueType = "JOINT"
	TechniqueTransform TechniqueType = "TRANSFORM"
	TechniqueWeight    TechniqueType = "WEIGHT"

	NodeJoint NodeType = "JOINT"

	ArmatureNodeID = "Armature"
)

func ParseCollada(documentPath string) (*loaders.ModelSpecification, error) {
	rawCollada, err := LoadDocument(documentPath)
	if err != nil {
		return nil, err
	}

	mesh := rawCollada.LibraryGeometries[0].Geometry[0].Mesh

	var normalSourceID, textureSourceID Uri
	for _, input := range mesh.Polylist[0].Input {
		if input.Semantic == string(SemanticNormal) {
			normalSourceID = input.Source[1:] // remove leading "#"
		} else if input.Semantic == string(SemanticTexCoord) {
			textureSourceID = input.Source[1:] // remove leading "#"
		}
	}

	if mesh.Vertices.Input[0].Semantic != string(SemanticPosition) {
		panic("vertices element expected to have semantic=position")
	}
	vertexSourceID := mesh.Vertices.Input[0].Source[1:] // remove leading "#"

	var positionSourceElement *Source
	var normalSourceElement *Source
	var textureSourceElement *Source

	for _, source := range mesh.Source {
		if string(source.Id) == string(vertexSourceID) {
			positionSourceElement = source
		} else if string(source.Id) == string(normalSourceID) {
			normalSourceElement = source
		} else if string(source.Id) == string(textureSourceID) {
			textureSourceElement = source
		}
	}

	if normalSourceElement == nil {
		panic("could not find position source")
	}

	if normalSourceElement == nil {
		panic("could not find normal source")
	}

	if textureSourceElement == nil {
		panic("could not find texture source")
	}

	positionSource := ParseVec3Array(positionSourceElement) // looks at <geometries>
	normalSource := ParseVec3Array(normalSourceElement)     // looks at <geometries>
	textureSource := ParseVec2Array(textureSourceElement)   // looks at <geometries>

	triVertices := parseIntArrayString(mesh.Polylist[0].P.V)

	skin := rawCollada.LibraryControllers[0].Controller.Skin

	var joints []string
	var weights []float32
	// parse controller sources
	for _, source := range skin.Source {
		if source.TechniqueCommon.Accessor.Param.Name == TechniqueJoint.String() {
			joints = strings.Split(source.NameArray.V, " ")
		} else if source.TechniqueCommon.Accessor.Param.Name == TechniqueTransform.String() {
		} else if source.TechniqueCommon.Accessor.Param.Name == TechniqueWeight.String() {
			weights = parseFloatArrayString(source.FloatArray.Floats.V)
		}
	}

	jointsToIndex := map[string]int{}
	for i, name := range joints {
		jointsToIndex[name] = i
	}

	// parse joint weights
	vcount := parseIntArrayString(skin.VertexWeights.VCount)
	v := parseIntArrayString(skin.VertexWeights.V)

	jointIDs := [][]int{}
	jointWeights := [][]int{}
	vIndex := 0
	for _, numWeights := range vcount {
		jointIDsList := []int{}
		jointWeightsList := []int{}

		for i := 0; i < numWeights; i++ {
			jointID := v[vIndex+(i*2)]       // each weight takes up two spots (joint index, weight index)
			weightIndex := v[vIndex+(i*2)+1] // each weight takes up two spots (joint index, weight index)
			jointIDsList = append(jointIDsList, jointID)
			jointWeightsList = append(jointWeightsList, weightIndex)
		}
		jointIDs = append(jointIDs, jointIDsList)
		jointWeights = append(jointWeights, jointWeightsList)
		vIndex += (numWeights * 2)
	}

	// parse joint hierarchy
	var rootJoint *loaders.Joint
	for _, node := range rawCollada.LibraryVisualScenes[0].VisualScene[0].Node {
		if string(node.Id) == ArmatureNodeID {
			rootNode := node.Node[0]
			rootJoint = parseJointElement(rootNode, jointsToIndex)
		}
	}

	result := &loaders.ModelSpecification{
		// VertexFloatsSourceData: vertexFloatData,
		// NormalFloatsSourceData: normalFloatData,
		TriIndices: triVertices,

		PositionSourceData: positionSource,
		NormalSourceData:   normalSource,
		TextureSourceData:  textureSource,
		ColorSourceData:    nil,

		JointsSourceData:       joints,
		JointWeightsSourceData: weights,

		JointIDs:     jointIDs,
		JointWeights: jointWeights,

		Root: rootJoint,
	}

	return result, nil
}

func parseJointElement(node *Node, jointsToIndex map[string]int) *loaders.Joint {
	children := []*loaders.Joint{}
	for _, childNode := range node.Node {
		children = append(children, parseJointElement(childNode, jointsToIndex))
	}

	bindTransform := parseMatrixArrayString(node.Matrix[0].V)
	joint := &loaders.Joint{
		ID:            jointsToIndex[node.Name],
		Name:          node.Name,
		BindTransform: bindTransform,
		Children:      children,
	}
	return joint
}
