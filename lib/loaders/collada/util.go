package collada

import (
	"strconv"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
)

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

func ParseVec2Array(source *Source) []mgl32.Vec2 {
	splitString := strings.Split(source.FloatArray.Floats.Values.V, " ")
	result := make([]mgl32.Vec2, len(splitString)/2)
	for i := 0; i < len(splitString); i += 2 {
		x := mustParseFloat32(splitString[i])
		y := mustParseFloat32(splitString[i+1])
		v := mgl32.Vec2{x, y}
		result[i/2] = v
	}
	return result
}

func mustParseFloat32(input string) float32 {
	num, err := strconv.ParseFloat(input, 32)
	if err != nil {
		panic(err)
	}
	return float32(num)
}

func convertToFloatList(v []mgl32.Vec3) []float32 {
	result := make([]float32, len(v)*3)
	for i := range v {
		result[i] = v[i].X()
		result[i+1] = v[i].Y()
		result[i+2] = v[i].Z()
	}
	return result
}

func mustParseInt(input string) int {
	num, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		panic(err)
	}
	return int(num)
}

func parseFloatArrayString(s string) []float32 {
	splitString := strings.Split(strings.TrimSpace(s), " ")
	result := make([]float32, len(splitString))

	for i, f := range splitString {
		result[i] = mustParseFloat32(f)
	}
	return result
}

func parseIntArrayString(s string) []int {
	splitString := strings.Split(strings.TrimSpace(s), " ")
	result := make([]int, len(splitString))

	for i, f := range splitString {
		result[i] = mustParseInt(f)
	}
	return result
}

func parseMatrixArrayString(s string) mgl32.Mat4 {
	splitString := strings.Split(strings.TrimSpace(s), " ")
	data := make([]float32, len(splitString))

	for i, f := range splitString {
		data[i] = mustParseFloat32(f)
	}

	return mgl32.Mat4FromRows(
		mgl32.Vec4{data[0], data[1], data[2], data[3]},
		mgl32.Vec4{data[4], data[5], data[6], data[7]},
		mgl32.Vec4{data[8], data[9], data[10], data[11]},
		mgl32.Vec4{data[12], data[13], data[14], data[15]},
	)
}
