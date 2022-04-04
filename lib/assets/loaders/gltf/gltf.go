package gltf

import (
	"fmt"
	"sort"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/lib/modelspec"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

type jointMeta struct {
	inverseBindMatrix mgl32.Mat4
}

type ParsedJoints struct {
	RootJoint       *modelspec.JointSpec
	NodeIDToJointID map[int]int
}

type TextureCoordStyle int

const (
	TextureCoordStyleOpenGL = 1
)

type ParseConfig struct {
	TextureCoordStyle TextureCoordStyle
}

func ParseGLTF(documentPath string, config *ParseConfig) (*modelspec.ModelSpecification, error) {
	document, err := gltf.Open(documentPath)
	if err != nil {
		return nil, err
	}

	var parsedJoints *ParsedJoints
	for _, skin := range document.Skins {
		parsedJoints, err = parseJoints(document, skin)
		if err != nil {
			return nil, err
		}
	}

	parsedAnimations := map[string]*modelspec.AnimationSpec{}
	for _, animation := range document.Animations {
		parsedAnimation, err := parseAnimation(document, animation, parsedJoints)
		parsedAnimations[animation.Name] = parsedAnimation
		if err != nil {
			return nil, err
		}
	}

	modelSpec := &modelspec.ModelSpecification{}

	for _, mesh := range document.Meshes {
		meshSpec, err := parseMesh(document, mesh, config)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		modelSpec.Meshes = append(modelSpec.Meshes, meshSpec)
	}

	if parsedJoints != nil {
		modelSpec.RootJoint = parsedJoints.RootJoint
	}

	modelSpec.Animations = parsedAnimations

	return modelSpec, nil
}

func parseAnimation(document *gltf.Document, animation *gltf.Animation, parsedJoints *ParsedJoints) (*modelspec.AnimationSpec, error) {
	keyFrames := map[float32]*modelspec.KeyFrame{}

	for _, channel := range animation.Channels {
		nodeID := int(*channel.Target.Node)
		if _, ok := parsedJoints.NodeIDToJointID[nodeID]; !ok {
			continue
		}

		jointID := parsedJoints.NodeIDToJointID[nodeID]
		sampler := animation.Samplers[(*channel.Sampler)]
		inputAccessorIndex := int(*sampler.Input)
		outputAccessorIndex := int(*sampler.Output)

		inputAccessor := document.Accessors[inputAccessorIndex]
		if inputAccessor.ComponentType != gltf.ComponentFloat {
			return nil, fmt.Errorf("unexpected component type %v", inputAccessor.ComponentType)
		}
		if inputAccessor.Type != gltf.AccessorScalar {
			return nil, fmt.Errorf("unexpected accessor type %v", inputAccessor.Type)
		}

		input, err := modeler.ReadAccessor(document, inputAccessor, nil)
		if err != nil {
			panic("WHA")
		}

		timestamps := input.([]float32)
		for _, timestamp := range timestamps {
			if _, ok := keyFrames[timestamp]; !ok {
				// hacky way to keep precision on tiny fractional seconds
				keyFrames[timestamp] = &modelspec.KeyFrame{
					Start: time.Duration(timestamp*1000) * time.Millisecond,
					Pose:  map[int]*modelspec.JointTransform{},
				}
			}
		}

		for _, timestamp := range timestamps {
			if _, ok := keyFrames[timestamp].Pose[jointID]; !ok {
				keyFrames[timestamp].Pose[jointID] = modelspec.NewDefaultJointTransform()
			}
		}

		outputAccessor := document.Accessors[outputAccessorIndex]
		if channel.Target.Path == gltf.TRSTranslation {
			if outputAccessor.ComponentType != gltf.ComponentFloat {
				return nil, fmt.Errorf("unexpected component type %v", outputAccessor.ComponentType)
			}
			if outputAccessor.Type != gltf.AccessorVec3 {
				return nil, fmt.Errorf("unexpected accessor type %v", outputAccessor.Type)
			}
			output, err := modeler.ReadAccessor(document, outputAccessor, nil)
			if err != nil {
				panic("WHA")
			}
			f32OutputValues := output.([][3]float32)
			for i, timestamp := range timestamps {
				f32Output := f32OutputValues[i]
				keyFrames[timestamp].Pose[jointID].Translation = mgl32.Vec3{f32Output[0], f32Output[1], f32Output[2]}
			}
		} else if channel.Target.Path == gltf.TRSRotation {
			if outputAccessor.ComponentType != gltf.ComponentFloat {
				return nil, fmt.Errorf("unexpected component type %v", outputAccessor.ComponentType)
			}
			if outputAccessor.Type != gltf.AccessorVec4 {
				return nil, fmt.Errorf("unexpected accessor type %v", outputAccessor.Type)
			}
			output, err := modeler.ReadAccessor(document, outputAccessor, nil)
			if err != nil {
				panic("WHA")
			}
			f32OutputValues := output.([][4]float32)
			for i, timestamp := range timestamps {
				f32Output := f32OutputValues[i]
				keyFrames[timestamp].Pose[jointID].Rotation = mgl32.Quat{V: mgl32.Vec3{f32Output[0], f32Output[1], f32Output[2]}, W: f32Output[3]}
			}
		} else if channel.Target.Path == gltf.TRSScale {
			if outputAccessor.ComponentType != gltf.ComponentFloat {
				return nil, fmt.Errorf("unexpected component type %v", outputAccessor.ComponentType)
			}
			if outputAccessor.Type != gltf.AccessorVec3 {
				return nil, fmt.Errorf("unexpected accessor type %v", outputAccessor.Type)
			}
			output, err := modeler.ReadAccessor(document, outputAccessor, nil)
			if err != nil {
				panic("WHA")
			}
			f32OutputValues := output.([][3]float32)
			for i, timestamp := range timestamps {
				f32Output := f32OutputValues[i]
				keyFrames[timestamp].Pose[jointID].Scale = mgl32.Vec3{f32Output[0], f32Output[1], f32Output[2]}
			}
		}
	}

	var timestamps []float32
	for timestamp := range keyFrames {
		timestamps = append(timestamps, timestamp)
	}
	sort.Slice(timestamps, func(i, j int) bool { return timestamps[i] < timestamps[j] })

	var keyFrameSlice []*modelspec.KeyFrame
	for _, timestamp := range timestamps {
		keyFrameSlice = append(keyFrameSlice, keyFrames[timestamp])
	}

	return &modelspec.AnimationSpec{
		Name:      animation.Name,
		KeyFrames: keyFrameSlice,
		Length:    keyFrameSlice[len(keyFrameSlice)-1].Start,
	}, nil
}

func parseJoints(document *gltf.Document, skin *gltf.Skin) (*ParsedJoints, error) {
	jms := map[int]*jointMeta{}
	jointNodeIDs := uint32SliceToIntSlice(skin.Joints)
	nodeIDToJointID := map[int]int{}

	for jointID, nodeID := range jointNodeIDs {
		nodeIDToJointID[nodeID] = jointID
	}

	for jointID := 0; jointID < len(jointNodeIDs); jointID++ {
		jms[jointID] = &jointMeta{}
	}

	acr := document.Accessors[int(*skin.InverseBindMatrices)]
	if acr.ComponentType != gltf.ComponentFloat {
		return nil, fmt.Errorf("unexpected component type %v", acr.ComponentType)
	}
	if acr.Type != gltf.AccessorMat4 {
		return nil, fmt.Errorf("unexpected accessor type %v", acr.Type)
	}

	data, err := modeler.ReadAccessor(document, acr, nil)
	if err != nil {
		return nil, err
	}

	inverseBindMatrices := data.([][4][4]float32)
	for jointID, _ := range jms {
		matrix := inverseBindMatrices[jointID]
		inverseBindMatrix := mgl32.Mat4FromRows(
			mgl32.Vec4{matrix[0][0], matrix[0][1], matrix[0][2], matrix[0][3]},
			mgl32.Vec4{matrix[1][0], matrix[1][1], matrix[1][2], matrix[1][3]},
			mgl32.Vec4{matrix[2][0], matrix[2][1], matrix[2][2], matrix[2][3]},
			mgl32.Vec4{matrix[3][0], matrix[3][1], matrix[3][2], matrix[3][3]},
		)
		jms[jointID].inverseBindMatrix = inverseBindMatrix
	}

	joints := map[int]*modelspec.JointSpec{}
	for nodeID, node := range document.Nodes {
		if _, ok := nodeIDToJointID[nodeID]; !ok {
			continue
		}

		jointID := nodeIDToJointID[nodeID]
		translation := node.Translation
		rotation := node.Rotation
		scale := node.Scale

		// from the gltf spec:
		//
		// When a node is targeted for animation (referenced by an animation.channel.target),
		// only TRS properties MAY be present; matrix MUST NOT be present.

		translationMatrix := mgl32.Translate3D(translation[0], translation[1], translation[2])
		rotationMatrix := mgl32.Quat{V: mgl32.Vec3{rotation[0], rotation[1], rotation[2]}, W: rotation[3]}.Mat4()
		scaleMatrix := mgl32.Scale3D(scale[0], scale[1], scale[2])

		joints[jointID] = &modelspec.JointSpec{
			Name:                 fmt.Sprintf("joint_%s_%d", node.Name, jointID),
			ID:                   jointID,
			BindTransform:        translationMatrix.Mul4(rotationMatrix.Mul4(scaleMatrix)),
			InverseBindTransform: jms[jointID].inverseBindMatrix,
		}
	}

	// set up the joint hierarchy
	childIDSet := map[int]bool{}
	for jointID, nodeID := range jointNodeIDs {
		children := uint32SliceToIntSlice(document.Nodes[nodeID].Children)
		for _, childNodeID := range children {
			childJointID := nodeIDToJointID[childNodeID]
			childIDSet[childJointID] = true
			joints[jointID].Children = append(joints[jointID].Children, joints[childJointID])
		}

	}

	// find the root
	var root *modelspec.JointSpec
	for id, _ := range joints {
		if _, ok := childIDSet[id]; !ok {
			joint := joints[id]
			if len(joint.Children) > 0 {
				// sometimes people put joints as control objects that aren't actual parents
				root = joint
			}
		}
	}

	parsedJoints := &ParsedJoints{
		RootJoint:       root,
		NodeIDToJointID: nodeIDToJointID,
	}
	return parsedJoints, nil
}

func parseMesh(document *gltf.Document, mesh *gltf.Mesh, config *ParseConfig) (*modelspec.MeshSpecification, error) {
	// parsedMesh := &ParsedMesh{}
	meshSpec := &modelspec.MeshSpecification{}

	for _, primitive := range mesh.Primitives {
		meshChunkSpec := &modelspec.MeshChunkSpecification{}
		acrIndex := *primitive.Indices
		meshIndices, err := modeler.ReadIndices(document, document.Accessors[int(acrIndex)], nil)
		if err != nil {
			return nil, err
		}
		meshChunkSpec.VertexIndices = meshIndices

		// TODO: not sure when a mesh would have multiple primitives
		// do i need to support multiple materials that come from multiple
		// primitives?
		if primitive.Material != nil {
			materialIndex := int(*primitive.Material)
			material := document.Materials[materialIndex]
			pbr := *material.PBRMetallicRoughness
			meshChunkSpec.PBRMaterial = &modelspec.PBRMaterial{
				PBRMetallicRoughness: &modelspec.PBRMetallicRoughness{
					BaseColorFactor: mgl32.Vec4{pbr.BaseColorFactor[0], pbr.BaseColorFactor[1], pbr.BaseColorFactor[2], pbr.BaseColorFactor[3]},
					MetalicFactor:   *pbr.MetallicFactor,
					RoughnessFactor: *pbr.RoughnessFactor,
				},
			}
			if pbr.BaseColorTexture != nil {
				var tex uint32
				meshChunkSpec.PBRMaterial.PBRMetallicRoughness.BaseColorTexture = &tex
			}
		}

		for attribute, index := range primitive.Attributes {
			acr := document.Accessors[int(index)]
			if meshChunkSpec.UniqueVertices == nil {
				meshChunkSpec.UniqueVertices = make([]modelspec.Vertex, int(acr.Count))
			}

			if attribute == gltf.POSITION {
				positions, err := modeler.ReadPosition(document, acr, nil)
				if err != nil {
					return nil, err
				}

				if len(positions) != len(meshChunkSpec.UniqueVertices) {
					fmt.Println("dafuq")
				}

				for i, position := range positions {
					meshChunkSpec.UniqueVertices[i].Position = position
				}

				// meshSpec.PositionSourceData = loosenFloat32Array3ToVec(positions)
			} else if attribute == gltf.NORMAL {
				normals, err := modeler.ReadNormal(document, acr, nil)
				if err != nil {
					return nil, err
				}
				for i, normal := range normals {
					meshChunkSpec.UniqueVertices[i].Normal = normal
				}
				// meshSpec.NormalSourceData = loosenFloat32Array3ToVec(normals)
			} else if attribute == gltf.TEXCOORD_0 {
				textureCoords, err := modeler.ReadTextureCoord(document, acr, nil)
				if err != nil {
					return nil, err
				}
				for i, textureCoord := range textureCoords {
					if config.TextureCoordStyle == TextureCoordStyleOpenGL {
						textureCoord[1] = 1 - textureCoord[1]
					}
					meshChunkSpec.UniqueVertices[i].Texture = textureCoord
				}
				// meshSpec.TextureSourceData = loosenFloat32Array2ToVec(textureCoords)
			} else if attribute == gltf.JOINTS_0 {
				jointsSlice, err := modeler.ReadJoints(document, acr, nil)
				if err != nil {
					return nil, err
				}
				readJointIDs := loosenUint16Array(jointsSlice)
				for i, jointIDs := range readJointIDs {
					meshChunkSpec.UniqueVertices[i].JointIDs = jointIDs
				}
			} else if attribute == gltf.WEIGHTS_0 {
				weights, err := modeler.ReadWeights(document, acr, nil)
				if err != nil {
					return nil, err
				}
				readJointWeights := loosenFloat32Array4(weights)
				for i, jointWeights := range readJointWeights {
					meshChunkSpec.UniqueVertices[i].JointWeights = jointWeights
				}
			} else {
				fmt.Printf("[%s] unhandled attribute %s\n", mesh.Name, attribute)
			}
		}

		for _, index := range meshChunkSpec.VertexIndices {
			meshChunkSpec.Vertices = append(meshChunkSpec.Vertices, meshChunkSpec.UniqueVertices[index])
		}

		meshSpec.MeshChunks = append(meshSpec.MeshChunks, meshChunkSpec)
	}
	return meshSpec, nil
}

func loosenFloat32Array4(floats [][4]float32) [][]float32 {
	result := make([][]float32, len(floats))
	for i, children := range floats {
		result[i] = make([]float32, len(children))
		for j, float := range children {
			result[i][j] = float
		}
	}
	return result
}

func loosenUint16Array(uints [][4]uint16) [][]int {
	result := make([][]int, len(uints))
	for i, children := range uints {
		result[i] = make([]int, len(children))
		for j, uint := range children {
			result[i][j] = int(uint)
		}
	}
	return result
}

func loosenFloat32Array2ToVec(floats [][2]float32) []mgl32.Vec2 {
	var result []mgl32.Vec2
	for _, props := range floats {
		result = append(result, mgl32.Vec2(props))
	}
	return result
}

func loosenFloat32Array3ToVec(floats [][3]float32) []mgl32.Vec3 {
	var result []mgl32.Vec3
	for _, props := range floats {
		result = append(result, mgl32.Vec3(props))
	}
	return result
}

func uint32SliceToIntSlice(slice []uint32) []int {
	var result []int
	for _, value := range slice {
		result = append(result, int(value))
	}
	return result
}
