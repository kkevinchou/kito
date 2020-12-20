package collada_test

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/kkevinchou/kito/lib/loaders/collada"
)

func TestCollada(t *testing.T) {
	_, _ = collada.ParseCollada(rawCollada)
	t.Error()
}

func parseNormals(rawCollada *collada.RawCollada) []mgl32.Vec3 {
	return nil
}

func parseAnimations(rawCollada *collada.RawCollada) {
}
