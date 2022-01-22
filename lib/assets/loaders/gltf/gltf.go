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

type ParsedMesh struct {
	PositionSource         []mgl32.Vec3
	NormalSource           []mgl32.Vec3
	TextureSource          []mgl32.Vec2
	VertexAttributeIndices []int
	JointIDs               [][]int
	JointWeights           [][]float32
}

type ParsedSkin struct {
	RootJoint  *modelspec.JointSpec
	Joints     map[int]*modelspec.JointSpec
	JointOrder []int
}

type ParsedAnimation struct {
	Name          string
	AnimationSpec *modelspec.AnimationSpec
}

func ParseGLTF(documentPath string) (*modelspec.ModelSpecification, error) {
	document, err := gltf.Open(documentPath)
	if err != nil {
		return nil, err
	}

	if len(document.Meshes) > 1 {
		panic("unable to handle > 1 mesh")
	}

	// Parse skinning info to construct the joint hierarchy
	var parsedSkin *ParsedSkin
	for _, skin := range document.Skins {
		parsedSkin, err = parseSkin(document, skin)
		if err != nil {
			return nil, err
		}
	}

	var parsedAnimation *ParsedAnimation
	for _, animation := range document.Animations {
		parsedAnimation, err = parseAnimation(document, animation, parsedSkin)
		if err != nil {
			return nil, err
		}
	}

	// animation info
	// iterate channels
	// only consider nodes that are joints (as determined by skinning info)
	// collect all the samplers
	//     the "input" for the sampler is the key frame time stamp. use this to create each key frame
	//     the "output" for the sampler is the value for the transform
	//     the type of transformation is based off of the channel "path" field
	// each sampler does a create/update on a keyframe. keyframe is added to a map, key: timestamp
	// key frames are sorted by time stamp and placed in the animation spec

	var parsedMesh *ParsedMesh

	// should theoretically only find one node with a mesh/skin
	for _, node := range document.Nodes {
		if node.Mesh != nil && node.Skin != nil {
			// found a node that has a mesh and skinning info
			mesh := document.Meshes[*node.Mesh]

			parsedMesh, err = parseMesh(document, mesh, parsedSkin.JointOrder)
			if err != nil {
				return nil, err
			}

			break
		}
	}

	result := &modelspec.ModelSpecification{
		VertexAttributeIndices: parsedMesh.VertexAttributeIndices,
		// TODO: FIX tHIS NUMBER WHEN WE ADD ANIMATIONS
		VertexAttributesStride: 3,

		PositionSourceData: parsedMesh.PositionSource,
		NormalSourceData:   parsedMesh.NormalSource,
		TextureSourceData:  parsedMesh.TextureSource,
		// EffectSpecData:     effectSpec,

		JointIDs:     parsedMesh.JointIDs,
		JointWeights: parsedMesh.JointWeights,

		RootJoint: parsedSkin.RootJoint,
		Animation: parsedAnimation.AnimationSpec,
		// Animation: &modelspec.AnimationSpec{KeyFrames: keyFrames, Length: keyFrames[len(keyFrames)-1].Start},
	}

	return result, nil
}

func parseAnimation(document *gltf.Document, animation *gltf.Animation, parsedSkin *ParsedSkin) (*ParsedAnimation, error) {
	joints := parsedSkin.Joints
	keyFrames := map[float32]*modelspec.KeyFrame{}

	for _, channel := range animation.Channels {
		nodeID := int(*channel.Target.Node)
		if _, ok := joints[nodeID]; !ok {
			continue
		}
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
			if _, ok := keyFrames[timestamp].Pose[nodeID]; !ok {
				keyFrames[timestamp].Pose[nodeID] = modelspec.NewDefaultJointTransform()
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
				keyFrames[timestamp].Pose[nodeID].Translation = mgl32.Vec3{f32Output[0], f32Output[1], f32Output[2]}
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
				keyFrames[timestamp].Pose[nodeID].Rotation = mgl32.Quat{V: mgl32.Vec3{f32Output[0], f32Output[1], f32Output[2]}, W: f32Output[3]}
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
				keyFrames[timestamp].Pose[nodeID].Scale = mgl32.Vec3{f32Output[0], f32Output[1], f32Output[2]}
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

	return &ParsedAnimation{
		Name: animation.Name,
		AnimationSpec: &modelspec.AnimationSpec{
			KeyFrames: keyFrameSlice,
			Length:    keyFrameSlice[len(keyFrameSlice)-1].Start,
		},
	}, nil
}

func parseSkin(document *gltf.Document, skin *gltf.Skin) (*ParsedSkin, error) {
	jms := map[int]*jointMeta{}
	jointIDs := uint32SliceToIntSlice(skin.Joints)
	for _, id := range jointIDs {
		jms[id] = &jointMeta{}
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
	for i, jointID := range jointIDs {
		matrix := inverseBindMatrices[i]
		inverseBindMatrix := mgl32.Mat4FromRows(
			mgl32.Vec4{matrix[0][0], matrix[0][1], matrix[0][2], matrix[0][3]},
			mgl32.Vec4{matrix[1][0], matrix[1][1], matrix[1][2], matrix[1][3]},
			mgl32.Vec4{matrix[2][0], matrix[2][1], matrix[2][2], matrix[2][3]},
			mgl32.Vec4{matrix[3][0], matrix[3][1], matrix[3][2], matrix[3][3]},
		)
		jms[jointID].inverseBindMatrix = inverseBindMatrix
	}

	joints := map[int]*modelspec.JointSpec{}
	for i, node := range document.Nodes {
		if jm, ok := jms[i]; ok {
			// node is a joint
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

			joints[i] = &modelspec.JointSpec{
				Name:                 fmt.Sprintf("joint_%s_%d", node.Name, i),
				ID:                   i,
				BindTransform:        translationMatrix.Mul4(rotationMatrix.Mul4(scaleMatrix)),
				InverseBindTransform: jm.inverseBindMatrix,
			}
		}
	}

	childIDSet := map[int]int{}
	// setup children slice
	for id, joint := range joints {
		children := uint32SliceToIntSlice(document.Nodes[id].Children)
		for _, childID := range children {
			childIDSet[childID] = childID
			joint.Children = append(joint.Children, joints[childID])
		}
	}

	var root *modelspec.JointSpec
	for id, _ := range joints {
		if _, ok := childIDSet[id]; !ok {
			root = joints[id]
		}
	}

	parsedSkin := &ParsedSkin{
		RootJoint:  root,
		Joints:     joints,
		JointOrder: jointIDs,
	}
	return parsedSkin, nil
}

func parseMesh(document *gltf.Document, mesh *gltf.Mesh, jointOrder []int) (*ParsedMesh, error) {
	parsedMesh := &ParsedMesh{}

	if len(document.Meshes) > 1 {
		panic("unable to handle > 1 mesh")
	}
	for _, primitive := range mesh.Primitives {
		acrIndex := *primitive.Indices
		meshIndices, err := modeler.ReadIndices(document, document.Accessors[int(acrIndex)], nil)
		if err != nil {
			return nil, err
		}

		for _, index := range meshIndices {
			parsedMesh.VertexAttributeIndices = append(parsedMesh.VertexAttributeIndices, []int{int(index), int(index), int(index)}...)
		}

		// attributes := primitive.Attributes
		for attribute, index := range primitive.Attributes {
			if attribute == gltf.POSITION {
				acr := document.Accessors[int(index)]
				positions, err := modeler.ReadPosition(document, acr, nil)
				if err != nil {
					return nil, err
				}
				parsedMesh.PositionSource = loosenFloat32Array3ToVec(positions)
			} else if attribute == gltf.NORMAL {
				acr := document.Accessors[int(index)]
				normals, err := modeler.ReadPosition(document, acr, nil)
				if err != nil {
					return nil, err
				}
				parsedMesh.NormalSource = loosenFloat32Array3ToVec(normals)
			} else if attribute == gltf.TEXCOORD_0 {
				acr := document.Accessors[int(index)]
				textureCoords, err := modeler.ReadTextureCoord(document, acr, nil)
				if err != nil {
					return nil, err
				}
				parsedMesh.TextureSource = loosenFloat32Array2ToVec(textureCoords)
			} else if attribute == gltf.JOINTS_0 {
				acr := document.Accessors[int(index)]
				joints, err := modeler.ReadJoints(document, acr, nil)
				if err != nil {
					return nil, err
				}
				jointIndices := loosenUint16Array(joints)
				for _, jointSet := range jointIndices {
					var allZero bool = true
					for i, _ := range jointSet {
						if jointSet[i] != 0 {
							allZero = false
						}
						// convert from the joint index to the actual joint ID
						jointSet[i] = jointOrder[jointSet[i]]
						// if jointSet[i] == 42 {
						// 	jointSet[i] = 2
						// }
					}

					// all joint indices were zero
					// 42 - hips, 2 - headtop_end, 3 - head, 31 - spine
					if allZero {
						jointSet[0] = 2
						jointSet[1] = 2
						jointSet[2] = 2
						jointSet[3] = 2
					}
				}
				parsedMesh.JointIDs = jointIndices
			} else if attribute == gltf.WEIGHTS_0 {
				acr := document.Accessors[int(index)]
				weights, err := modeler.ReadWeights(document, acr, nil)
				if err != nil {
					return nil, err
				}
				jointWeights := loosenFloat32Array4(weights)
				parsedMesh.JointWeights = jointWeights
			} else {
				fmt.Printf("unhandled attribute %s\n", attribute)
			}
		}
	}
	return parsedMesh, nil
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
