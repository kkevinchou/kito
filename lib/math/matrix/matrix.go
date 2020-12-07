package matrix

import "github.com/go-gl/mathgl/mgl32"

func Mat4FromValues(values []float32) mgl32.Mat4 {
	var col0 mgl32.Vec4
	var col1 mgl32.Vec4
	var col2 mgl32.Vec4
	var col3 mgl32.Vec4
	cols := []mgl32.Vec4{col0, col1, col2, col3}

	for i := 0; i < 16; i++ {
		cols[i/4][i%4] = values[i]
	}

	return mgl32.Mat4FromCols(cols[0], cols[1], cols[2], cols[3])
}
