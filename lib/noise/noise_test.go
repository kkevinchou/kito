package noise_test

import (
	"fmt"
	"testing"

	"github.com/kkevinchou/kito/lib/noise"
)

func TestNoise(t *testing.T) {
	noiseMap := noise.GenerateNoiseMap(10, 10)
	fmt.Println(noiseMap)
}
