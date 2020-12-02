package noise

import (
	"fmt"

	"github.com/aquilax/go-perlin"
)

const (
	alpha       = 2.
	beta        = 2.
	n           = 3
	seed  int64 = 100
)

func GenerateNoiseMap(xMax, yMax int) [][]float64 {
	result := make([][]float64, xMax)

	p := perlin.NewPerlin(alpha, beta, n, seed)
	for x := 0; x < xMax; x++ {
		result[x] = make([]float64, yMax)
		for y := 0; y < yMax; y++ {
			result[x][y] = p.Noise2D(float64(x)/0.3, float64(y)/0.3)
		}
	}
	fmt.Println(result)

	return result
}
