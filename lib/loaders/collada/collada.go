package collada

import (
	"strconv"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
)

type Collada struct {
	vertexSourceData  []mgl32.Vec3
	normalSourceData  []mgl32.Vec3
	textureSourceData []mgl32.Vec2
}

type SemanticType string

const (
	SemanticVertex   SemanticType = "VERTEX"
	SemanticNormal   SemanticType = "NORMAL"
	SemanticTexCoord SemanticType = "TEXCOORD"
	SemanticColor    SemanticType = "COLOR"
	SemanticPosition SemanticType = "POSITION"
)

func ParseCollada(rawCollada *RawCollada) *Collada {
	mesh := rawCollada.LibraryGeometries[0].Geometry[0].Mesh

	var normalSourceID Uri
	for _, polyList := range mesh.Polylist {
		for _, input := range polyList.HasSharedInput.Input {
			if input.Semantic == string(SemanticNormal) {
				normalSourceID = input.Source[1:] // remove leading "#"
			}
		}
	}

	if mesh.Vertices.Input[0].Semantic != string(SemanticPosition) {
		panic("vertices element expected to have semantic=position")
	}
	vertexSourceID := mesh.Vertices.Input[0].Source[1:] // remove leading "#"

	var vertexPositionSource *Source
	var normalSource *Source

	for _, source := range mesh.Source {
		if string(source.Id) == string(vertexSourceID) {
			vertexPositionSource = source
		} else if string(source.Id) == string(normalSourceID) {
			normalSource = source
		}
	}

	if vertexPositionSource == nil {
		panic("could not find vertex source")
	}

	vertices := ParseVec3Array(vertexPositionSource) // looks at <geometries>
	normals := ParseVec3Array(normalSource)          // looks at <geometries>
	result := &Collada{
		vertexSourceData: vertices,
		normalSourceData: normals,
	}
	// parseAnimations(rawCollada) // looks at <animations> and <controllers> <scenes>(contains hierarchy)

	return result
}

func ParseVec3Array(source *Source) []mgl32.Vec3 {
	splitString := strings.Split(source.FloatArray.Floats.Values.V, " ")
	result := make([]mgl32.Vec3, len(splitString)/3)
	for i := 0; i < len(splitString); i += 3 {
		x := mustParseFloat32(splitString[i])
		y := mustParseFloat32(splitString[i+1])
		z := mustParseFloat32(splitString[i+2])
		v := mgl32.Vec3{x, y, z}
		result[i/3] = v
	}
	return result
}

func mustParseFloat32(input string) float32 {
	out, err := strconv.ParseFloat(input, 64)
	if err != nil {
		panic(err)
	}
	return float32(out)
}
