package matrix

import (
	"fmt"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestMultiply(t *testing.T) {
	m1 := []float32{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, -5, 1,
	}
	m := mat4FromValues(m1)
	fmt.Println(m.Mul4(m.Inv()))

	t.Fail()
}

func mat4FromValues(values []float32) mgl32.Mat4 {
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
