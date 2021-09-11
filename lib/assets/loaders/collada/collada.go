package collada

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/lib/animation"
)

type SemanticType string
type NodeType string

const (
	SemanticVertex            SemanticType = "VERTEX"
	SemanticNormal            SemanticType = "NORMAL"
	SemanticTexCoord          SemanticType = "TEXCOORD"
	SemanticColor             SemanticType = "COLOR"
	SemanticPosition          SemanticType = "POSITION"
	SemanticInput             SemanticType = "INPUT"
	SemanticOutput            SemanticType = "OUTPUT"
	SemanticJoint             SemanticType = "JOINT"
	SemanticWeight            SemanticType = "WEIGHT"
	SemanticInverseBindMatrix SemanticType = "INV_BIND_MATRIX"

	NodeJoint NodeType = "JOINT"

	ArmatureNodeID = "Armature"
	ArmatureName   = "Armature"
)

func ParseCollada(documentPath string) (*animation.ModelSpecification, error) {
	rawCollada, err := LoadDocument(documentPath)
	if err != nil {
		return nil, err
	}

	// parse geometry
	mesh := rawCollada.LibraryGeometries[0].Geometry[0].Mesh

	// for some reason, my personal blender exports will place
	// polygon information in a "triangles" element instead of "polylist"

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

	// parse skinning information, joint weights
	skin := rawCollada.LibraryControllers[0].Controller.Skin

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
	visualScene := rawCollada.LibraryVisualScenes[0].VisualScene[0]
	// armatureID := visualScene.Node[0].Id // was used for bob.dae

	var jointSourceID string
	var weightSourceID string
	var inverseBindMatrixSourceID string

	for _, input := range skin.VertexWeights.Input {
		if input.Semantic == string(SemanticJoint) {
			jointSourceID = string(input.Source)[1:] // remove preceding #
		} else if input.Semantic == string(SemanticWeight) {
			weightSourceID = string(input.Source)[1:] // remove preceding #
		} else if input.Semantic == string(SemanticInverseBindMatrix) {
			inverseBindMatrixSourceID = string(input.Source)[1:]
		}
	}

	var joints []string
	var weights []float32
	var inverseBindMatrices []mgl32.Mat4

	for _, source := range skin.Source {
		if string(source.Id) == jointSourceID {
			for _, j := range strings.Fields(source.NameArray.V) {
				// joints = append(joints, fmt.Sprintf("%s_%s", string(armatureID), j)) // was used for bob.dae
				joints = append(joints, j)
			}
		} else if string(source.Id) == weightSourceID {
			weights = parseFloatArrayString(source.FloatArray.Floats.V)
		} else if string(source.Id) == inverseBindMatrixSourceID {
			inverseBindMatrices = parseMultiMatrixArrayString(source.FloatArray.Floats.V)
		}
	}

	jointsToIndex := map[string]int{}
	for i, name := range joints {
		jointsToIndex[name] = i
	}

	var rootNode *Node
	for _, node := range visualScene.Node {
		if node.Type == "JOINT" {
			rootNode = node
			break
		}

		// TODO: no idea if this works
		if len(node.Node) == 1 && node.Node[0].Type == "JOINT" {
			rootNode = node.Node[0]
			break
		}
	}

	rootJoint := parseJointElement(rootNode, jointsToIndex, inverseBindMatrices)

	// parse animations
	timeStampToPose := map[float32]map[int]*animation.JointTransform{}

	for _, animationElement := range rawCollada.LibraryAnimations[0].Animations {
		// get the input/output sources
		var inputSource Uri
		var outputSource Uri
		for _, input := range animationElement.Sampler.Inputs {
			if input.Semantic == string(SemanticInput) {
				inputSource = input.Source
			} else if input.Semantic == string(SemanticOutput) {
				outputSource = input.Source
			}
		}

		if string(inputSource) == "" {
			panic("could not find input source")
		}

		if string(outputSource) == "" {
			panic("could not find output source")
		}

		target := animationElement.Channel.Target
		jointName := strings.Split(target, "/")[0] // guessing the sample i'm looking at looks like: "Torso/transform"

		if _, ok := jointsToIndex[jointName]; !ok {
			panic(fmt.Sprintf("couldn't find joint name \"%s\" in joint listing. available joints: %v", jointName, jointsToIndex))
		}
		jointID := jointsToIndex[jointName]

		var timeStamps []float32
		var poseMatrices []mgl32.Mat4
		for _, source := range animationElement.Source {
			if string(source.Id) == string(inputSource)[1:] {
				timeStamps = parseFloatArrayString(source.FloatArray.V)
			} else if string(source.Id) == string(outputSource)[1:] {
				poseMatrices = parseMultiMatrixArrayString(source.FloatArray.V)
			}
		}

		if len(timeStamps) != len(poseMatrices) {
			panic("number of timestamps doesn't line up with number of matrices")
		}

		for i := 0; i < len(timeStamps); i++ {
			timeStamp := timeStamps[i]
			transform := poseMatrices[i]
			if timeStampToPose[timeStamp] == nil {
				timeStampToPose[timeStamp] = map[int]*animation.JointTransform{}
			}

			timeStampToPose[timeStamp][jointID] = &animation.JointTransform{
				Translation: transform.Col(3).Vec3(),
				Rotation:    mgl32.Mat4ToQuat(transform),
			}
		}
	}

	timeStamps := []float32{}
	for timeStamp := range timeStampToPose {
		timeStamps = append(timeStamps, timeStamp)
	}

	sort.Sort(byFloat32(timeStamps))

	keyFrames := []*animation.KeyFrame{}
	for _, timeStamp := range timeStamps {
		keyFrames = append(keyFrames, &animation.KeyFrame{
			Start: time.Duration(int(timeStamp*1000)) * time.Millisecond,
			Pose:  timeStampToPose[timeStamp],
		})
	}

	// TODO: perform assertions on number of joints, verts, etc

	result := &animation.ModelSpecification{
		TriIndices:       triVertices,
		TriIndicesStride: len(mesh.Polylist[0].Input),

		PositionSourceData: positionSource,
		NormalSourceData:   normalSource,
		TextureSourceData:  textureSource,
		ColorSourceData:    nil,

		JointsSourceData:       joints,
		JointWeightsSourceData: weights,

		JointIDs:     jointIDs,
		JointWeights: jointWeights,

		Root:      rootJoint,
		Animation: &animation.Animation{KeyFrames: keyFrames, Length: keyFrames[len(keyFrames)-1].Start},
	}

	return result, nil
}

func parseJointElement(node *Node, jointsToIndex map[string]int, inverseBindMatrices []mgl32.Mat4) *animation.JointSpecification {
	children := []*animation.JointSpecification{}

	for _, childNode := range node.Node {
		children = append(children, parseJointElement(childNode, jointsToIndex, inverseBindMatrices))
	}

	// TODO: cowboy.dae did not have a matrix but instead used translate, rotate scale.
	// too lazy to handle that so just assume identity matrix
	if len(node.Matrix) == 0 {
		fmt.Println("empty node matrix")
	}

	bindTransform := parseMatrixArrayString(node.Matrix[0].V)

	joint := &animation.JointSpecification{
		ID:            jointsToIndex[string(node.Id)],
		Name:          string(node.Id),
		BindTransform: bindTransform,
		Children:      children,
	}
	return joint
}

type byFloat32 []float32

func (s byFloat32) Len() int {
	return len(s)
}
func (s byFloat32) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byFloat32) Less(i, j int) bool {
	return s[i] < s[j]
}
