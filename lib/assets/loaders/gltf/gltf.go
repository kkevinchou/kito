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
	var jointIDs [][4]int
	var jointWeights [][4]float32

	for _, mesh := range document.Meshes {
		for _, primitive := range mesh.Primitives {
			acrIndex := *primitive.Indices
			indices, err := modeler.ReadIndices(document, document.Accessors[int(acrIndex)], nil)
			if err != nil {
				return nil, err
			}
			triIndices = uint32SliceToIntSlice(indices)

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
					jointWeights = weights
				} else {
					panic(fmt.Sprintf("unexpected attribute %s\n", attribute))
				}
			}
		}
	}

	_ = jointIDs
	_ = jointWeights

	result := &modelspec.ModelSpecification{
		TriIndices: triIndices,
		// TriIndicesStride: len(polyInput),

		PositionSourceData: positionSource,
		NormalSourceData:   normalSource,
		TextureSourceData:  textureSource,
		// EffectSpecData:     effectSpec,

		// ColorSourceData: nil,

		// JointWeightsSourceData: weights,

		// JointIDs:     jointIDs,
		// JointWeights: jointWeights,

		// RootJoint: rootJoint,
		// Animation: &modelspec.AnimationSpec{KeyFrames: keyFrames, Length: keyFrames[len(keyFrames)-1].Start},
	}

	return result, nil
}

func uint16SliceToIntSlice(uints [][4]uint16) [][4]int {
	var result [][4]int
	for _, props := range uints {
		var casted [4]int
		for i, uint := range props {
			casted[i] = int(uint)
		}
		result = append(result, casted)
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
