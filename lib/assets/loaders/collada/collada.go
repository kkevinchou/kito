package collada

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/lib/libutils"
	"github.com/kkevinchou/kito/lib/modelspec"
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

func ParseCollada(documentPath string) (*modelspec.ModelSpecification, error) {
	rawCollada, err := LoadDocument(documentPath)
	if err != nil {
		return nil, err
	}

	// parse materials
	var effectSpec *modelspec.EffectSpec
	var effectId Id

	if len(rawCollada.LibraryEffects[0].Effect) > 0 {
		if len(rawCollada.LibraryEffects[0].Effect) != 1 {
			panic(" more than 1 effect not supported")
		}
		effect := rawCollada.LibraryEffects[0].Effect[0]
		effectId = effect.Id

		if len(effect.ProfileCommon.TechniqueFx) != 1 {
			panic(" more than 1 technique not supported")
		}

		technique := effect.ProfileCommon.TechniqueFx[0]
		var shaderElement string
		var techniqueNode *RenderingTechnique
		if technique.Lambert != nil {
			shaderElement = "lambert"
			techniqueNode = technique.Lambert
		} else if technique.Phong != nil {
			shaderElement = "phong"
			techniqueNode = technique.Phong
		} else if technique.Blinn != nil {
			shaderElement = "blinn"
			techniqueNode = technique.Blinn
		} else if technique.ConstantFx != nil {
			shaderElement = "constant"
			techniqueNode = technique.ConstantFx
		} else {
			panic("Unknown technique for material")
		}

		effectSpec = &modelspec.EffectSpec{
			ShaderElement:          shaderElement,
			EmissionColor:          parseVect3String(techniqueNode.Emission.Floats.V),
			DiffuseColor:           parseVect3String(techniqueNode.Diffuse.Floats.V),
			IndexOfRefractionFloat: float32(techniqueNode.IndexOfRefraction.Value),
			ReflectivityFloat:      float32(techniqueNode.Reflectivity.Value),
			ReflectivityColor:      parseVect3String(techniqueNode.Reflective.Floats.V),
			ShininessFloat:         float32(techniqueNode.Shininess.Value),
			TransparencyFloat:      float32(techniqueNode.Transparency.Value),
			TransparencyColor:      parseVect3String(techniqueNode.Transparent.Floats.V),
		}
	}

	libraryMaterial := rawCollada.LibraryMaterials[0]
	var materialId Id
	if len(libraryMaterial.Material) > 0 {
		if len(libraryMaterial.Material) != 1 {
			panic("more than 1 material not supported")
		}
		material := libraryMaterial.Material[0]
		materialId = material.Id
		// materialName := material.Name
		instanceEffect := material.InstanceEffect
		instanceEffectURL := instanceEffect.Url[1:] // remove leading "#"

		if string(instanceEffectURL) != string(effectId) {
			panic("material effect not found")
		}
	}

	// parse geometry
	mesh := rawCollada.LibraryGeometries[0].Geometry[0].Mesh

	// Blender supports both PolyList and Triangles for storing mesh polygons
	var polyInput []*InputShared
	var polyValues string
	var polyMaterialSymbol string

	if len(mesh.Polylist) > 0 {
		polyInput = mesh.Polylist[0].Input
		polyValues = mesh.Polylist[0].P.V
		polyMaterialSymbol = mesh.Polylist[0].Material
	} else if len(mesh.Triangles) > 0 {
		polyInput = mesh.Triangles[0].Input
		polyValues = mesh.Triangles[0].P.V
		polyMaterialSymbol = mesh.Triangles[0].Material
	}

	if polyMaterialSymbol != string(materialId) {
		panic("Material on mesh not found")
	}

	var normalSourceID, textureSourceID Uri
	for _, input := range polyInput {
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

	triVertices := parseIntArrayString(polyValues)

	if len(rawCollada.LibraryControllers) == 0 || len(rawCollada.LibraryAnimations) == 0 {
		// no animations
		return &modelspec.ModelSpecification{
			TriIndices:       triVertices,
			TriIndicesStride: len(polyInput),

			PositionSourceData: positionSource,
			NormalSourceData:   normalSource,
			TextureSourceData:  textureSource,
			EffectSpecData:     effectSpec,
			ColorSourceData:    nil,
		}, nil
	}

	// parse skinning information, joint weights
	// TODO: handle multiple skins (sometimes the model is broken down into multiple indepdendent meshes)
	skin := rawCollada.LibraryControllers[0].Controller[0].Skin

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

	var jointSourceID string
	var weightSourceID string

	for _, input := range skin.VertexWeights.Input {
		if input.Semantic == string(SemanticJoint) {
			jointSourceID = string(input.Source)[1:] // remove preceding #
		} else if input.Semantic == string(SemanticWeight) {
			weightSourceID = string(input.Source)[1:] // remove preceding #
		}
	}

	var inverseBindMatrixSourceID string
	for _, input := range skin.Joints.Input {
		if input.Semantic == string(SemanticInverseBindMatrix) {
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

	jointNamesToIndex := map[string]int{}
	for i, name := range joints {
		jointNamesToIndex[name] = i
	}

	// parse joint hierarchy
	visualScene := rawCollada.LibraryVisualScenes[0].VisualScene[0]
	var rootNode *Node

	// finding the root
	for _, node := range visualScene.Node {
		if node.Type == "JOINT" {
			rootNode = node
			break
		}

		// TODO: no idea if this works
		if len(node.Node) > 0 && node.Node[0].Type == "JOINT" {
			rootNode = node.Node[0]
			break
		}
	}

	rootJoint := parseJointElement(rootNode, jointNamesToIndex, inverseBindMatrices)

	// parse animations
	timeStampToPose := map[float32]map[int]*modelspec.JointTransform{}

	animations := rawCollada.LibraryAnimations[0].Animations
	if len(animations[0].Animations) > 0 {
		// dirty hack for handling blender exported colladas. not sure why they nest it under
		// yet another animation element. This actually trips up the vscode collada renderer too...
		animations = animations[0].Animations
	}

	for _, animationElement := range animations {
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

		// dirty hack for handling blender collada exports
		if strings.HasPrefix(jointName, "Armature_") {
			jointName = jointName[len("Armature_"):]
		}

		if _, ok := jointNamesToIndex[jointName]; !ok {
			// dirty hack for testing blender collada
			// target actually looks for IDs, and we shouldn't be using jointNamesToIndex
			// todo: update root finding logic and parseJointElement to produce more rich datastructures
			// i'd like to have access to a jointIDsToIndex map that either maps to a name or a a joint index directly.
			if _, ok := jointNamesToIndex["Bone"]; !ok {
				fmt.Println(documentPath)
				panic(fmt.Sprintf("couldn't find joint name \"%s\" in joint listing. available joints: %v", jointName, jointNamesToIndex))
			} else {
				jointName = "Bone"
			}
		}
		jointIndex := jointNamesToIndex[jointName]

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
				timeStampToPose[timeStamp] = map[int]*modelspec.JointTransform{}
			}

			// translation := transform.Col(3).Vec3()
			// rotation := mgl32.Mat4ToQuat(transform)
			translation, scale, rotation := libutils.Decompose(transform)

			timeStampToPose[timeStamp][jointIndex] = &modelspec.JointTransform{
				Translation: translation,
				Rotation:    rotation,
				Scale:       scale,
			}
		}
	}

	timeStamps := []float32{}
	for timeStamp := range timeStampToPose {
		timeStamps = append(timeStamps, timeStamp)
	}

	sort.Sort(byFloat32(timeStamps))

	keyFrames := []*modelspec.KeyFrame{}
	for _, timeStamp := range timeStamps {
		keyFrames = append(keyFrames, &modelspec.KeyFrame{
			Start: time.Duration(int(timeStamp*1000)) * time.Millisecond,
			Pose:  timeStampToPose[timeStamp],
		})
	}

	// TODO: perform assertions on number of joints, verts, etc

	// TODO clean up ModelSpecification. This represents the API between our loader logic
	// and our internal model representation. Theoretically struct shouldn't change
	// if we load a different file format which allows our internal model representation code
	// to not require changes either.
	result := &modelspec.ModelSpecification{
		TriIndices:       triVertices,
		TriIndicesStride: len(polyInput),

		PositionSourceData: positionSource,
		NormalSourceData:   normalSource,
		TextureSourceData:  textureSource,
		EffectSpecData:     effectSpec,

		ColorSourceData: nil,

		JointsSourceData:       joints,
		JointWeightsSourceData: weights,

		JointIDs:     jointIDs,
		JointWeights: jointWeights,

		Root:      rootJoint,
		Animation: &modelspec.AnimationSpec{KeyFrames: keyFrames, Length: keyFrames[len(keyFrames)-1].Start},
	}

	return result, nil
}

func parseJointElement(node *Node, jointNamesToIndex map[string]int, inverseBindMatrices []mgl32.Mat4) *modelspec.JointSpecification {
	children := []*modelspec.JointSpecification{}

	for _, childNode := range node.Node {
		children = append(children, parseJointElement(childNode, jointNamesToIndex, inverseBindMatrices))
	}

	// TODO: cowboy.dae did not have a matrix but instead used translate, rotate scale.
	// too lazy to handle that so just assume identity matrix
	if len(node.Matrix) == 0 {
		fmt.Println("empty node matrix")
	}

	bindTransform := ParseMatrixArrayString(node.Matrix[0].V)
	jointID := jointNamesToIndex[string(node.Id)]

	joint := &modelspec.JointSpecification{
		ID:                   jointID,
		Name:                 string(node.Id),
		BindTransform:        bindTransform,
		InverseBindTransform: inverseBindMatrices[jointID],
		Children:             children,
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

func printHierarchy(j *modelspec.JointSpecification, level int) {
	indentation := ""
	for i := 0; i < level; i++ {
		indentation += "    "
	}
	fmt.Println(indentation + j.Name + fmt.Sprintf(" %d", j.ID))
	for _, c := range j.Children {
		printHierarchy(c, level+1)
	}
}
