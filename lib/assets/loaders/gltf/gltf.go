package gltf

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/lib/modelspec"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

func ParseGLTF(documentPath string) (*modelspec.ModelSpecification, error) {
	document, err := gltf.Open(documentPath)
	if err != nil {
		return nil, err
	}

	var positionSource []mgl32.Vec3
	var normalSource []mgl32.Vec3
	var textureSource []mgl32.Vec2
	var triIndices []int
	var jointIDs [][]int
	var jointWeights [][]float32

	if len(document.Meshes) > 1 {
		panic("unable to handle > 1 mesh")
	}

	for _, mesh := range document.Meshes {
		for _, primitive := range mesh.Primitives {
			acrIndex := *primitive.Indices
			meshIndices, err := modeler.ReadIndices(document, document.Accessors[int(acrIndex)], nil)
			if err != nil {
				return nil, err
			}
			// triIndices = uint32SliceToIntSlice(meshIndices)
			for _, index := range meshIndices {
				triIndices = append(triIndices, int(index))
				triIndices = append(triIndices, int(index))
				triIndices = append(triIndices, int(index))
			}

			// attributes := primitive.Attributes
			for attribute, index := range primitive.Attributes {
				if attribute == gltf.POSITION {
					acr := document.Accessors[int(index)]
					positions, err := modeler.ReadPosition(document, acr, nil)
					if err != nil {
						return nil, err
					}
					positionSource = float3SliceToVec3Slice(positions)
				} else if attribute == gltf.NORMAL {
					acr := document.Accessors[int(index)]
					normals, err := modeler.ReadPosition(document, acr, nil)
					if err != nil {
						return nil, err
					}
					normalSource = float3SliceToVec3Slice(normals)
				} else if attribute == gltf.TEXCOORD_0 {
					acr := document.Accessors[int(index)]
					textureCoords, err := modeler.ReadTextureCoord(document, acr, nil)
					if err != nil {
						return nil, err
					}
					textureSource = float2SliceToVec2Slice(textureCoords)
				} else if attribute == gltf.JOINTS_0 {
					acr := document.Accessors[int(index)]
					joints, err := modeler.ReadJoints(document, acr, nil)
					if err != nil {
						return nil, err
					}
					jointIDs = uint16SliceToIntSlice(joints)
				} else if attribute == gltf.WEIGHTS_0 {
					acr := document.Accessors[int(index)]
					weights, err := modeler.ReadWeights(document, acr, nil)
					if err != nil {
						return nil, err
					}
					jointWeights = float32ArraySliceTofloat32SliceSlice(weights)
				} else {
					panic(fmt.Sprintf("unexpected attribute %s\n", attribute))
				}
			}
		}
	}

	result := &modelspec.ModelSpecification{
		TriIndices: triIndices,
		// TODO: FIX tHIS NUMBER
		TriIndicesStride: 3,

		PositionSourceData: positionSource,
		NormalSourceData:   normalSource,
		TextureSourceData:  textureSource,
		// EffectSpecData:     effectSpec,

		JointIDs:     jointIDs,
		JointWeights: jointWeights,

		// RootJoint: rootJoint,
		// Animation: &modelspec.AnimationSpec{KeyFrames: keyFrames, Length: keyFrames[len(keyFrames)-1].Start},
	}

	return result, nil
}

func float32ArraySliceTofloat32SliceSlice(floats [][4]float32) [][]float32 {
	result := make([][]float32, len(floats))
	for i, children := range floats {
		result[i] = make([]float32, len(children))
		for j, float := range children {
			result[i][j] = float
		}
	}
	return result
}

func uint16SliceToIntSlice(uints [][4]uint16) [][]int {
	result := make([][]int, len(uints))
	for i, children := range uints {
		result[i] = make([]int, len(children))
		for j, uint := range children {
			result[i][j] = int(uint)
		}
	}
	return result
}

func float2SliceToVec2Slice(floats [][2]float32) []mgl32.Vec2 {
	var result []mgl32.Vec2
	for _, props := range floats {
		result = append(result, mgl32.Vec2(props))
	}
	return result
}

func float3SliceToVec3Slice(floats [][3]float32) []mgl32.Vec3 {
	var result []mgl32.Vec3
	for _, props := range floats {
		result = append(result, mgl32.Vec3(props))
	}
	return result
}

func uint32SliceToIntSlice(uints []uint32) []int {
	var result []int
	for _, uint := range uints {
		result = append(result, int(uint))
	}
	return result
}
